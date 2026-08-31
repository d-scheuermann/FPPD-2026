package main

import (
	"fmt"
	"sync"
)

var contador int
var mu sync.Mutex

func incrementar(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	for i := 0; i < n; i++ {
		contador++
	}
	mu.Unlock()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go incrementar(2000000, &wg)
	go incrementar(2000000, &wg)
	wg.Wait()
	fmt.Printf("Contador: %d (esperado: 4000000)\n", contador)
}
