// Exercício 2B — Largada Sincronizada ("Ready-Set-Go")
//
// Este programa simula uma corrida: 6 corredores devem esperar
// até que o juiz dê o sinal de largada. O juiz prepara a pista
// (simula com sleep) e então libera todos de uma vez.
//
// Atualmente NÃO há sincronização: os corredores partem
// imediatamente, sem esperar o juiz.
//
// Sua tarefa: usar sync.Cond com Broadcast() para que todos
// os corredores esperem o sinal do juiz antes de começar.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numCorredores = 6

// Estado compartilhado
var (
	largou bool
	mu     sync.Mutex
	cond   = sync.NewCond(&mu)
)

func corredor(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[Corredor %d] na linha de largada\n", id)

	// Aguarda o sinal do juiz
	cond.L.Lock()
	for !largou {
		cond.Wait() // Libera a trava temporariamente e bloqueia até receber o Broadcast
	}
	cond.L.Unlock()

	duracao := time.Duration(500+rand.Intn(2000)) * time.Millisecond
	fmt.Printf("[Corredor %d] LARGOU!\n", id)
	time.Sleep(duracao)
	fmt.Printf("[Corredor %d] chegou (tempo: %v)\n", id, duracao)
}

func main2() {
	var wg sync.WaitGroup

	for i := 1; i <= numCorredores; i++ {
		wg.Add(1)
		go corredor(i, &wg)
	}

	// Juiz prepara a pista
	fmt.Println("\n[Juiz] Preparando a pista...")
	time.Sleep(2 * time.Second)

	// Altera o estado e avisa todos os corredores
	cond.L.Lock()
	largou = true
	fmt.Println("[Juiz] VAI!\n")
	cond.Broadcast() // Acorda todas as goroutines em espera
	cond.L.Unlock()

	wg.Wait()
	fmt.Println("\nCorrida encerrada.")
}