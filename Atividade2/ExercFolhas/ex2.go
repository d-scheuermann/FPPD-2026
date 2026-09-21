// Produtor-consumidor simples com duas goroutines.
//
// Uma goroutine produz itens e envia pelo canal.
// Outra goroutine consome e imprime cada item.
// O produtor fecha o canal ao terminar, sinalizando
// ao consumidor que não haverá mais dados.

package main

import (
	"fmt"
	"time"
)

func produtor(ch chan<- string) {
	for i := 1; i <= 5; i++ {
		item := fmt.Sprintf("item-%d", i)
		fmt.Printf("[Produtor]   Produzindo %s\n", item)
		ch <- item
		time.Sleep(200 * time.Millisecond)
	}
	close(ch)
	fmt.Println("[Produtor] Canal fechado")
}

func consumidor(ch <-chan string, done chan<- bool) {
    for item := range ch {
        hora := time.Now().Format("15:04:05")
        fmt.Printf("[Consumidor] [%s] Recebeu %s\n", hora, item)
        time.Sleep(1 * time.Second)
    }
    done <- true
}

func main() {

	//a)
	ch := make(chan string)

	//c)
	//ch := make(chan string, 5)

	//e)
	//ch := make(chan string, 2)

	done := make(chan bool)

	go produtor(ch)
	go consumidor(ch, done)

	fmt.Println("[Main] Esperando consumidor...")
	<-done
	fmt.Println("Fim!")
}

//b) Responda: quem dita o ritmo da comunicação, o produtor ou o consumidor? Por quê?

// Quem dita o ritmo da comunicacao e o consumidor, pois o produtor apenas ira produzir
// o proximo item do canal a partir do consumo do item anterior

//d) Responda: o comportamento do produtor mudou? Por quê?

// Sim, o comportamento mudou pois agora o produtor pode enviar no canal ate 5 itens
// consecutivos, deixando o ritmo de consumo a merce do consumidor.

//e) Experimente com make(chan string, 2) (buffer de tamanho 2). O que acontece quando o buffer enche?

// O comportamento muda novamente, visto que o produtor pode enviar apenas 2 itens
// consecutivos antes de ter que esperar o consumo de um deles pelo consumidor.

