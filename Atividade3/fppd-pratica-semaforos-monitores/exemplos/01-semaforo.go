// Exemplo 1 — Semáforo contador com golang.org/x/sync/semaphore
//
// Limita a 2 goroutines simultâneas acessando um recurso.
// Execute: go run exemplos/01-semaforo.go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	sem := semaphore.NewWeighted(2) // permite até 2 acessos simultâneos
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("[%d] esperando...\n", id)
			sem.Acquire(context.Background(), 1) // P (wait) — bloqueia se já há 2 dentro

			fmt.Printf("[%d] >>> entrou\n", id)
			time.Sleep(1 * time.Second)
			fmt.Printf("[%d] <<< saiu\n", id)

			sem.Release(1) // V (signal) — libera uma permissão
		}(i)
	}

	wg.Wait()
}
