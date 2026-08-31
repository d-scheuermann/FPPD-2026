package main

import (
	"fmt"
	"sync"
	"time"
)

func crescente(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		fmt.Println("Crescente:", i)
		time.Sleep(1 * time.Second) }
}

func decrescente(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 10; i >= 1; i-- {
		fmt.Println("Descrescente:", i) 
		time.Sleep(1 * time.Second) }
	}

func main() {
	var wg sync.WaitGroup
	
	wg.Add(2)
	go crescente(&wg)
	go decrescente(&wg)

	wg.Wait()
	fmt.Println("Ambas goroutines finalizaram.")
}