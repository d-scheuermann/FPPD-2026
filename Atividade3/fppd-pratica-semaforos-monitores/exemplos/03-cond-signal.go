// Exemplo 3 — sync.Cond com Signal: acordar uma goroutine por vez
//
// Goroutines esperam por um recurso. A main libera um de cada vez.
// Execute: go run exemplos/03-cond-signal.go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	disponivel := 0

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			mu.Lock()
			for disponivel == 0 {
				cond.Wait()
			}
			disponivel--
			fmt.Printf("[%d] obteve recurso (restam: %d)\n", id, disponivel)
			mu.Unlock()
		}(i)
	}

	for i := 1; i <= 3; i++ {
		time.Sleep(500 * time.Millisecond)
		mu.Lock()
		disponivel++
		fmt.Printf("--- liberou 1 recurso (total: %d)\n", disponivel)
		cond.Signal() // acorda UMA goroutine esperando
		mu.Unlock()
	}

	wg.Wait()
}
