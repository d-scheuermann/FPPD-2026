// Exercício 1 — Ponto de partida
//
// Este programa usa sync.WaitGroup para a main() esperar
// duas goroutines. Sua tarefa: substituir o WaitGroup por
// canais como mecanismo de sinalização.

package main

import (
	"fmt"
	"time"
)

// Recebe um canal para sinalizar quando terminar a execução
func crescente(done chan<- bool) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("[Crescente] %d\n", i)
		time.Sleep(1 * time.Second)
	}
	done <- true // Sinaliza que a goroutine finalizou
}

// Recebe um canal para sinalizar quando terminar a execução
func decrescente(done chan<- bool) {
	for i := 10; i >= 1; i-- {
		fmt.Printf("[Decrescente] %d\n", i)
		time.Sleep(1 * time.Second)
	}
	done <- true // Sinaliza que a goroutine finalizou
}

func main1() {
	// Canal usado para sinalização de término
	done := make(chan bool)

	go crescente(done)

	go decrescente(done)

	for i := 0; i < 2; i++ {
		<-done 
	}

	fmt.Println("Fim!")
}