package main

//P produtores, C consumidores e capacidade K de buffer
//Consultas: Ben-Ari, Algorithm 6.8 - produtor e consumidor e IA para sintaxe, correções e sugestões

//COMO RODAR:
//Detec. de data race - go run -race v-semaforos.go 
//Default: go run v-semaforos.go (opcional:) -p=? -c=? -k=? -s=5


import (
	"flag"
	"fmt"
	"sync"
	"time"
)


//notEmpty: quantos tem
//notFull: quantas vagas

//comeca com B maiusculo para ser publico

type BufferProdutos struct {
	buffer []struct{}        
	notFull  chan struct{} //cada struct{}{} é uma vaga disponivel
	notEmpty chan struct{} //cada struct{}{} é uma vaga ocupada
	mutex chan struct{} 
	k int          
}

var totalProduzido int = 0;
var totalConsumido int = 0;
var encerramentoProducao bool = false; //encerramento ordenado

//canal pra nao dar data race
var ocupMediaBuffer chan struct{};
var somaOcupacao int64 = 0;
var totalAlteracoes int64 = 0;

var mutexContadorProducao chan struct{} 
var mutexContadorConsumo chan struct{}


func produtor(nProdutor int, buffer *BufferProdutos, wg *sync.WaitGroup, segundos int) {
	//o defer vai dar done no wait group dos produtores como ultima coisa ou uma das ultimas coisas antes de sair do escopo
	defer wg.Done()
	//timer, canal que vai receber apos dar o tempo
	timerFim := time.After(time.Duration(segundos) * time.Second)
	totalProduzidoProdutor := 0;
	for { //loop forever ate acabar o tempo
		select {
			case <- timerFim:
				fmt.Printf("Produtor n. %d esgotou seu tempo, tendo produzido %d.\n", nProdutor, totalProduzidoProdutor)
				<- mutexContadorProducao
				totalProduzido += totalProduzidoProdutor
				mutexContadorProducao <- struct{}{}
				return
			case <- buffer.notFull:
				//mutex pra nao ter data race no buffer
				<-buffer.mutex
				buffer.buffer = append(buffer.buffer, struct{}{})

				//soma 1 na media de ocup
				<-ocupMediaBuffer                       
				somaOcupacao += int64(len(buffer.buffer)) 
				totalAlteracoes++                         
				ocupMediaBuffer <- struct{}{}  

				totalProduzidoProdutor++
				buffer.mutex <- struct{}{}
				buffer.notEmpty <- struct{}{}
		}
	}
}

func consumidor(nConsumidor int, buffer *BufferProdutos, wg *sync.WaitGroup) {
	defer wg.Done()
	totalConsumidoConsumidor := 0;

	for {
		select {
			case <- buffer.notEmpty:
				<- buffer.mutex
				//tira o primeiro do buffer
				buffer.buffer = buffer.buffer[1:]

				//decrementa a media de ocup
				<-ocupMediaBuffer                       
				somaOcupacao += int64(len(buffer.buffer)) 
				totalAlteracoes++                         
				ocupMediaBuffer <- struct{}{}  

				totalConsumidoConsumidor++
				buffer.mutex <- struct{}{}
				buffer.notFull <- struct{}{}
			//deve vir depois do notEmpty
			case <- time.After(1000 * time.Millisecond):
				<- buffer.mutex
				if (encerramentoProducao && len(buffer.buffer) == 0) {
					fmt.Printf("Consumidor n. %d terminou, consumindo %d itens. \n", nConsumidor, totalConsumidoConsumidor)
					<- mutexContadorConsumo 
					totalConsumido += totalConsumidoConsumidor
					mutexContadorConsumo <- struct{}{}
					buffer.mutex <- struct{}{}
					return
				}
				buffer.mutex <- struct{}{}

				fmt.Printf("Consumidor n. %d: comportamento alternativo. \n", nConsumidor)
		}
	}
}

func main () {
	//comeca contador pro throughput
	inicio := time.Now()

	//argumentos: nome, padrao, descricao

	p := flag.Int("p", 100, "n de produtores")
	c := flag.Int("c", 140, "n de consumidores")
	k := flag.Int("k", 1000, "capacidade do buffer")
	s := flag.Int("s", 20, "depois de quantos segundos os produtores devem parar")
	
	//parse serve para processa-los

	flag.Parse()

	//instancia o buffer

	buffer := BufferProdutos{
		buffer:   make([]struct{}, 0, *k),
		notFull:  make(chan struct{}, *k),
		notEmpty: make(chan struct{}, *k),
		mutex:    make(chan struct{}, 1),
		k:        *k,
	}

	//preenche o mutex dos contadores

	mutexContadorProducao = make(chan struct{}, 1)
	mutexContadorConsumo = make(chan struct{}, 1)

	mutexContadorProducao <- struct{}{}
	mutexContadorConsumo <- struct{}{}

	//preenche as permissões
	for i := 0; i < *k; i++ {
		buffer.notFull <- struct{}{}
	}
	buffer.mutex <- struct{}{}

	//o wait group producao serve para o encerramento ordenado
	//quando todos produtores encerrarem dado o timer um boolean sera true
	//e os consumidores passarao a poder a dar done no seu wait group quando 
	//o buffer estiver vazio E o boolean ser true
	var wgProducao sync.WaitGroup
	var wgConsumidores sync.WaitGroup

	ocupMediaBuffer = make(chan struct{}, 1)
	ocupMediaBuffer <- struct{}{}

	//lanca as go routines e adiciona o tamanho no wait group
	for i := 1; i <= *p; i++ {
		wgProducao.Add(1)

		go produtor(i, &buffer, &wgProducao, *s)
	}

	for i := 1; i <= *c; i++ {
		wgConsumidores.Add(1)
		
		go consumidor(i, &buffer, &wgConsumidores)
	}
	
	//espera os produtores terminarem
	wgProducao.Wait()

	//tem esse mutex pois teria um data race da go routine da main com rotinas dos consumidores
	<- buffer.mutex
	//boolean que os consumidores verificam pra ver se param caso o buffer esteja vazio
	encerramentoProducao = true
	buffer.mutex <- struct{}{}

	wgConsumidores.Wait()
	

	//ve tempo desde entao
	tempoTotal := time.Since(inicio).Seconds()

	throughput := float64(totalConsumido) / tempoTotal

	ocupMedia := float64(somaOcupacao) / float64(totalAlteracoes)

	fmt.Printf("Termino do programa: Total produzido: %d, total consumido: %d \n. Consumo/s de %fs, ocupação média no buffer de %f.", totalProduzido, totalConsumido, throughput, ocupMedia);

}