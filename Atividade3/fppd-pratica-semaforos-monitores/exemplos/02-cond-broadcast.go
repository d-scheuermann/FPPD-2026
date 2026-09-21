// Exemplo 2 — sync.Cond com Broadcast: acordar todas as goroutines
//
// Três goroutines esperam um evento; a main sinaliza todas de uma vez.
// Execute: go run exemplos/02-cond-broadcast.go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	pronto := false

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			mu.Lock()
			for !pronto { // sempre usar FOR, nunca IF
				cond.Wait() // libera o lock, espera, readquire o lock
			}
			mu.Unlock()

			fmt.Printf("[%d] acordou!\n", id)
		}(i)
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Sinalizando todas as goroutines...")

	mu.Lock()
	pronto = true
	cond.Broadcast() // acorda TODAS as goroutines esperando
	mu.Unlock()

	wg.Wait()
}
