// Exercício 1 — Controle de Acesso ao Banheiro com Semáforos
//
// Este programa simula 10 pessoas tentando usar um banheiro
// com capacidade para 3. Atualmente NÃO há controle de acesso:
// todas as pessoas entram ao mesmo tempo.
//
// Sua tarefa: adicionar um semáforo contador usando o pacote
// golang.org/x/sync/semaphore para limitar a 3 acessos simultâneos.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	capacidade   = 3
	totalPessoas = 10
)

var sem = semaphore.NewWeighted(int64(capacidade))

func usarBanheiro(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	limiteEspera := time.Duration(1000+rand.Intn(2000)) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), limiteEspera)
	defer cancel() 

	fmt.Printf("[Pessoa %2d] quer usar o banheiro (tolera esperar até %v)\n", id, limiteEspera)

	if err := sem.Acquire(ctx, 1); err != nil {
		fmt.Printf("[Pessoa %2d] DESISTIU e foi embora! (esperou demais)\n", id)
		return
	}

	defer sem.Release(1)

	fmt.Printf("[Pessoa %2d] >>> ENTROU no banheiro\n", id)
	duracaoUso := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(duracaoUso)
	fmt.Printf("[Pessoa %2d] <<< SAIU do banheiro (usou por %v)\n", id, duracaoUso)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= totalPessoas; i++ {
		wg.Add(1)
		go usarBanheiro(i, &wg)
	}

	wg.Wait()
	fmt.Println("\nProcesso finalizado.")
}