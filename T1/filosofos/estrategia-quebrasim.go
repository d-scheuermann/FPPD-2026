package main

import (
	"fmt"
	"sync"
	"time"
)

func executarQuebraSimetria(n, r int) {
	fmt.Println("\n--- QUEBRA DE SIMETRIA ---")
	fmt.Println("O ultimo filosofo pega os garfos na ordem inversa.")

	if n < 2 {
		fmt.Println("Sao necessarios pelo menos 2 filosofos.")
		return
	}

	garfos := make([]chan struct{}, n)

	for i := 0; i < n; i++ {
		garfos[i] = make(chan struct{}, 1)
		garfos[i] <- struct{}{}
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

				// Simula o filosofo pensando.
				time.Sleep(10 * time.Millisecond)

				inicioEspera := time.Now()

				var primeiro chan struct{}
				var segundo chan struct{}

				// Todos pegam esquerda -> direita.
				// O ultimo faz o contrario.
				if id == n-1 {
					primeiro = direito
					segundo = esquerdo
				} else {
					primeiro = esquerdo
					segundo = direito
				}

				<-primeiro
				<-segundo

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
				segundo <- struct{}{}
				primeiro <- struct{}{}
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("\nTodos os filosofos terminaram.")
	imprimirMetricas(metricas)
}
