// Problema 2: Produtor/Consumidor com Buffer Limitado (versão com channel)
//
// Uso:
//
//	go run v-channel.go                      # varre K = 1, 10, 100
//	go run v-channel.go -k 10 -p 4 -c 4      # um único K
//	go run v-channel.go -csv > resultados.csv
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type result struct {
	K           int
	Elapsed     time.Duration
	Produced    int64
	Consumed    int64
	PerConsumer []int64
	AvgOcc      float64
	Timeouts    int64
}

// dorme um tempo aleatório em [0, 2*mean] (média = mean), simulando trabalho com jitter.
func sleepJitter(r *rand.Rand, mean time.Duration) {
	if mean <= 0 {
		return
	}
	time.Sleep(time.Duration(r.Float64() * 2 * float64(mean)))
}

func run(P, C, K, N int, prodMean, consMean, timeout time.Duration, quiet bool) (result, error) {
	buf := make(chan int, K) // buffer limitado = channel com capacidade K

	total := P * N
	seen := make([]int32, total) // para checar que cada item foi consumido exatamente 1 vez

	var produced int64
	perConsumer := make([]int64, C)
	var timeouts int64

	// ---------- Produtores ----------
	var prodWG sync.WaitGroup
	start := time.Now()
	for p := 0; p < P; p++ {
		prodWG.Add(1)
		go func(p int) {
			defer prodWG.Done()
			r := rand.New(rand.NewSource(int64(1000 + p)))
			for i := 0; i < N; i++ {
				sleepJitter(r, prodMean)
				buf <- p*N + i // bloqueia se o buffer estiver cheio
				atomic.AddInt64(&produced, 1)
			}
		}(p)
	}

	// ---------- Encerramento ordenado ----------
	// Só quando TODOS os produtores terminam o channel é fechado.
	// Fechar um channel não descarta o que já está no buffer: os consumidores
	// continuam recebendo os itens restantes e só veem ok == false quando esvaziar.
	go func() {
		prodWG.Wait()
		close(buf)
	}()

	// ---------- Consumidores ----------
	var consWG sync.WaitGroup
	for c := 0; c < C; c++ {
		consWG.Add(1)
		go func(c int) {
			defer consWG.Done()
			r := rand.New(rand.NewSource(int64(2000 + c)))

			consume := func(item int) {
				atomic.AddInt32(&seen[item], 1)
				perConsumer[c]++ // cada posição é escrita por uma única goroutine
				sleepJitter(r, consMean)
			}

			if c == 0 {
				// Consumidor com timeout: select + time.After.
				for {
					select {
					case item, ok := <-buf:
						if !ok { // fechado e drenado
							return
						}
						consume(item)
					case <-time.After(timeout):
						// comportamento alternativo: nenhum item chegou no intervalo
						n := atomic.AddInt64(&timeouts, 1)
						if !quiet && n <= 3 {
							fmt.Fprintf(os.Stderr, "[consumidor 0] timeout #%d: sem itens em %v, fazendo trabalho alternativo\n", n, timeout)
						}
					}
				}
			}

			// Consumidores comuns: range termina quando o channel está fechado e vazio.
			for item := range buf {
				consume(item)
			}
		}(c)
	}

	// ---------- Amostrador de ocupação do buffer ----------
	stop := make(chan struct{})
	sampDone := make(chan struct{})
	var occSum, occN int64
	go func() {
		defer close(sampDone)
		t := time.NewTicker(100 * time.Microsecond)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				occSum += int64(len(buf))
				occN++
			case <-stop:
				return
			}
		}
	}()

	consWG.Wait()
	elapsed := time.Since(start)
	close(stop)
	<-sampDone

	var consumed int64
	for _, v := range perConsumer {
		consumed += v
	}
	res := result{
		K: K, Elapsed: elapsed, Produced: atomic.LoadInt64(&produced), Consumed: consumed,
		PerConsumer: perConsumer, Timeouts: atomic.LoadInt64(&timeouts),
	}
	if occN > 0 {
		res.AvgOcc = float64(occSum) / float64(occN)
	}

	// ---------- Verificações finais ----------
	if res.Produced != res.Consumed || res.Produced != int64(total) {
		return res, fmt.Errorf("K=%d: produzido=%d consumido=%d esperado=%d", K, res.Produced, res.Consumed, total)
	}
	for id, v := range seen {
		if v != 1 {
			return res, fmt.Errorf("K=%d: item %d consumido %d vezes", K, id, v)
		}
	}
	return res, nil
}

func main() {
	P := flag.Int("p", 4, "número de produtores")
	C := flag.Int("c", 4, "número de consumidores")
	K := flag.Int("k", 0, "capacidade do buffer (0 = varre 1, 10 e 100)")
	N := flag.Int("n", 500, "itens por produtor")
	prod := flag.Duration("prod", time.Millisecond, "tempo médio de produção por item")
	cons := flag.Duration("cons", time.Millisecond, "tempo médio de consumo por item")
	timeout := flag.Duration("timeout", 2*time.Millisecond, "timeout do consumidor 0")
	reps := flag.Int("reps", 3, "repetições por K (reporta a média)")
	csv := flag.Bool("csv", false, "saída em CSV")
	flag.Parse()

	ks := []int{1, 10, 100}
	if *K > 0 {
		ks = []int{*K}
	}

	if *csv {
		fmt.Println("K,P,C,itens,tempo_s,throughput_itens_s,ocupacao_media,timeouts,por_consumidor")
	} else {
		fmt.Printf("P=%d C=%d N=%d/produtor prod=%v cons=%v timeout=%v reps=%d\n\n", *P, *C, *N, *prod, *cons, *timeout, *reps)
		fmt.Printf("%-5s %-10s %-16s %-12s %-9s %s\n", "K", "tempo", "throughput(it/s)", "ocup. média", "timeouts", "itens por consumidor")
	}

	for _, k := range ks {
		var sumElapsed time.Duration
		var sumOcc float64
		var sumTO int64
		sumPer := make([]int64, *C)
		var items int64

		for r := 0; r < *reps; r++ {
			res, err := run(*P, *C, k, *N, *prod, *cons, *timeout, *csv || r > 0)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERRO:", err)
				os.Exit(1)
			}
			sumElapsed += res.Elapsed
			sumOcc += res.AvgOcc
			sumTO += res.Timeouts
			items = res.Consumed
			for i, v := range res.PerConsumer {
				sumPer[i] += v
			}
		}
		n := int64(*reps)
		avgElapsed := sumElapsed / time.Duration(n)
		thr := float64(items) / avgElapsed.Seconds()
		per := make([]int64, *C)
		for i := range per {
			per[i] = sumPer[i] / n
		}
		if *csv {
			fmt.Printf("%d,%d,%d,%d,%.4f,%.1f,%.2f,%d,%v\n", k, *P, *C, items, avgElapsed.Seconds(), thr, sumOcc/float64(n), sumTO/n, per)
		} else {
			fmt.Printf("%-5d %-10v %-16.0f %-12.2f %-9d %v\n", k, avgElapsed.Round(time.Millisecond), thr, sumOcc/float64(n), sumTO/n, per)
		}
	}
	if !*csv {
		fmt.Println("\nOK: produzido == consumido e cada item foi consumido exatamente uma vez.")
	}
}
