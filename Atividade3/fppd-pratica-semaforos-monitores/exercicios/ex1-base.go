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
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	capacidade   = 3
	totalPessoas = 10
)


func usarBanheiro(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[Pessoa %2d] quer usar o banheiro\n", id)

	// TODO: adquirir uma permissão do semáforo antes de entrar

	fmt.Printf("[Pessoa %2d] >>> ENTROU no banheiro\n", id)
	duracao := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(duracao)
	fmt.Printf("[Pessoa %2d] <<< SAIU do banheiro (usou %v)\n", id, duracao)

	// TODO: liberar a permissão do semáforo ao sair
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= totalPessoas; i++ {
		wg.Add(1)
		go usarBanheiro(i, &wg)
	}

	wg.Wait()
	fmt.Println("\nTodos usaram o banheiro.")
}
