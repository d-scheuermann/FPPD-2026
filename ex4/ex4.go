package main

import (
	"fmt"
	"sync"
	"time"
)

var mutexA, mutexB sync.Mutex

func tarefa1(wg *sync.WaitGroup) {
	defer wg.Done()
	mutexA.Lock()
	fmt.Println("tarefa1: adquiriu A")
	time.Sleep(100 * time.Millisecond)
	mutexB.Lock() // espera por B, que está com tarefa2
	fmt.Println("tarefa1: adquiriu B")
	mutexB.Unlock()
	mutexA.Unlock()
}
func tarefa2(wg *sync.WaitGroup) {
	defer wg.Done()
	mutexB.Lock()
	fmt.Println("tarefa2: adquiriu B")
	time.Sleep(100 * time.Millisecond)
	mutexA.Lock() // espera por A, que está com tarefa1
	fmt.Println("tarefa2: adquiriu A")
	mutexA.Unlock()
	mutexB.Unlock()
}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go tarefa1(&wg)
	go tarefa2(&wg)
	wg.Wait()
	fmt.Println("Fim!")
}
