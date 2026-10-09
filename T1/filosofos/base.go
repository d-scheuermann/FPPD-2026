package main

import (
	"fmt"
	"sync"
	"time"
)

func executarBase(n, r int) {
	fmt.Println("\n--- VERSAO BASE ---")
	fmt.Println("Todos os filosofos pegam primeiro o garfo da esquerda.")

	garfos := make([]chan struct{}, n)

	for i := 0; i < n; i++ {
		garfos[i] = make(chan struct{}, 1)
		garfos[i] <- struct{}{}
	}

	var wg sync.WaitGroup

	prontos := make(chan struct{}, n)
	tentarDireita := make(chan struct{})
	cancelar := make(chan struct{})

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()

			esquerdo := garfos[id]
			direito := garfos[(id+1)%n]

			// Cada filosofo pega primeiro seu garfo esquerdo.
			<-esquerdo
			fmt.Printf("Filosofo %d pegou o garfo esquerdo\n", id)

			// Avisa ao coordenador que ja pegou o esquerdo.
			prontos <- struct{}{}

			// Espera todos fazerem a mesma coisa.
			<-tentarDireita

			fmt.Printf("Filosofo %d tentando pegar o garfo direito\n", id)

			select {
			case <-direito:
				// Em uma execucao normal desta versao,
				// este ponto nao deve ser alcancado.
				fmt.Printf("Filosofo %d conseguiu os dois garfos\n", id)

				direito <- struct{}{}
				esquerdo <- struct{}{}

			case <-cancelar:
				// Permite encerrar a demonstracao sem
				// deixar goroutines presas.
				esquerdo <- struct{}{}
				return
			}
		}(i)
	}

	// Espera todos terem adquirido o garfo esquerdo.
	for i := 0; i < n; i++ {
		<-prontos
	}

	fmt.Println("\nTodos pegaram o garfo esquerdo.")
	fmt.Println("Liberando tentativa de pegar o direito...")

	close(tentarDireita)

	// Dá tempo para observar a espera circular.
	time.Sleep(1 * time.Second)

	fmt.Println("\nDEADLOCK detectado:")
	fmt.Println("cada filosofo possui um garfo e espera pelo garfo do vizinho.")

	close(cancelar)

	wg.Wait()

	fmt.Println("Demonstracao encerrada sem deixar goroutines bloqueadas.")

	_ = r // sera usado depois nas metricas
}