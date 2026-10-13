package main

import (
	"fmt"
	"sync"
)

func executarBase(n, r int) {
	fmt.Println("\n--- VERSAO BASE ---")
	fmt.Println("Todos os filosofos pegam primeiro o garfo da esquerda.")

	if n < 2 {
		fmt.Println("Erro: sao necessarios pelo menos 2 filosofos.")
		return
	}

	// Cada garfo e representado por um channel de capacidade 1.
	// O token dentro do channel significa que o garfo esta disponivel.
	garfos := make([]chan struct{}, n)

	for i := 0; i < n; i++ {
		garfos[i] = make(chan struct{}, 1)
		garfos[i] <- struct{}{}
	}

	var wg sync.WaitGroup
	wg.Add(n)

	// Usado para saber quando todos ja pegaram o garfo esquerdo.
	prontos := make(chan struct{}, n)

	// Libera todos para tentarem pegar o garfo direito.
	tentarDireita := make(chan struct{})

	// Confirma que todos chegaram ao ponto de tentar o direito.
	tentandoDireita := make(chan struct{}, n)

	// Usado para cancelar a demonstracao depois que
	// o deadlock foi identificado.
	cancelar := make(chan struct{})

	// Cada filosofo avisa que percebeu o cancelamento.
	cancelados := make(chan struct{}, n)

	// So depois que todos foram cancelados os garfos
	// esquerdos sao devolvidos.
	liberarGarfos := make(chan struct{})

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()

			esquerdo := garfos[id]
			direito := garfos[(id+1)%n]

			// Todos pegam primeiro o garfo esquerdo.
			<-esquerdo

			fmt.Printf(
				"Filosofo %d pegou o garfo esquerdo\n",
				id,
			)

			// Avisa que conseguiu o esquerdo.
			prontos <- struct{}{}

			// Espera todos os outros tambem pegarem
			// seus garfos esquerdos.
			<-tentarDireita

			fmt.Printf(
				"Filosofo %d tentando pegar o garfo direito\n",
				id,
			)

			// Avisa ao coordenador que chegou na tentativa.
			tentandoDireita <- struct{}{}

			// O garfo direito esta com o vizinho.
			// Portanto, todos ficam bloqueados aqui.
			select {

			case <-direito:
				// Este trecho nao deve acontecer na versao base,
				// pois todos os garfos ja estao ocupados.
				fmt.Printf(
					"Filosofo %d conseguiu os dois garfos\n",
					id,
				)

				direito <- struct{}{}
				esquerdo <- struct{}{}

			case <-cancelar:
				// Avisa que saiu da espera pelo garfo direito.
				cancelados <- struct{}{}

				// Espera os outros tambem serem cancelados
				// antes de devolver o garfo esquerdo.
				<-liberarGarfos

				esquerdo <- struct{}{}

				return
			}
		}(i)
	}

	// Espera todos pegarem o garfo esquerdo.
	for i := 0; i < n; i++ {
		<-prontos
	}

	fmt.Println("\nTodos pegaram o garfo esquerdo.")
	fmt.Println("Liberando tentativa de pegar o direito...")

	// Todos podem agora tentar pegar o direito.
	close(tentarDireita)

	// Espera todos chegarem ao ponto de tentar pegar
	// o segundo garfo.
	for i := 0; i < n; i++ {
		<-tentandoDireita
	}

	/*
		Neste momento:

		F0 possui G0 e espera G1
		F1 possui G1 e espera G2
		F2 possui G2 e espera G3
		...
		Fn possui Gn e espera G0

		Forma-se uma espera circular.
	*/

	fmt.Println("\nDEADLOCK detectado:")
	fmt.Println(
		"cada filosofo possui um garfo e espera pelo garfo do vizinho.",
	)

	// Cancela a espera pelo segundo garfo.
	close(cancelar)

	// Espera todos perceberem o cancelamento.
	for i := 0; i < n; i++ {
		<-cancelados
	}

	// Agora e seguro permitir que devolvam os
	// garfos esquerdos.
	close(liberarGarfos)

	wg.Wait()

	fmt.Println(
		"Demonstracao encerrada sem deixar goroutines bloqueadas.",
	)

	// Na versao base nenhuma refeicao e concluida,
	// pois o deadlock acontece na primeira tentativa.
	fmt.Println("\n--- METRICAS ---")

	for i := 0; i < n; i++ {
		fmt.Printf(
			"Filosofo %d: refeicoes=0 | espera media=indefinida (deadlock)\n",
			i,
		)
	}

	fmt.Printf(
		"\nForam solicitadas %d iteracoes, mas a execucao entrou em deadlock na primeira.\n",
		r,
	)
}
