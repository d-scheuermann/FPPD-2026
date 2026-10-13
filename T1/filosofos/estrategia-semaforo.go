package main

import (
	"fmt"
	"sync"
	"time"
)

func executarSemaforo(n, r int) {
	fmt.Println("\n--- SEMAFORO N-1 ---")
	fmt.Println("No maximo N-1 filosofos podem tentar pegar garfos ao mesmo tempo.")

	if n < 2 {
		fmt.Println("Sao necessarios pelo menos 2 filosofos.")
		return
	}

	// Cria os garfos.
	garfos := make([]chan struct{}, n)

	for i := 0; i < n; i++ {
		garfos[i] = make(chan struct{}, 1)
		garfos[i] <- struct{}{}
	}

	// Semaforo contador com N-1 permissoes.
	lugares := make(chan struct{}, n-1)

	for i := 0; i < n-1; i++ {
		lugares <- struct{}{}
	}

	metricas := make([]MetricaFilosofo, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()

			esquerdo := garfos[id]
			direito := garfos[(id+1)%n]

			for refeicao := 0; refeicao < r; refeicao++ {

				time.Sleep(10 * time.Millisecond)

				inicioEspera := time.Now()

				// Wait(P): pega uma permissao para sentar.
				<-lugares

				// Todos continuam pegando esquerda -> direita.
				<-esquerdo
				<-direito

				tempoEspera := time.Since(inicioEspera)

				metricas[id].EsperaTotal += tempoEspera
				metricas[id].Refeicoes++

				fmt.Printf(
					"Filosofo %d comendo (%d/%d)\n",
					id,
					refeicao+1,
					r,
				)

				time.Sleep(5 * time.Millisecond)

				// Devolve os garfos.
				direito <- struct{}{}
				esquerdo <- struct{}{}

				// Signal(V): devolve a permissao.
				lugares <- struct{}{}
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("\nTodos os filosofos terminaram.")
	imprimirMetricas(metricas)
}
