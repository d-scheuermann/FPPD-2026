// Exercício 3 — Pipeline de processamento
//
// Complete as funções gerador, multiplicador e impressora
// para formar um pipeline de 3 estágios conectados por canais.

package main

import "fmt"

// gerador envia os números de 'valores' no canal 'out' e fecha o canal.
func gerador(valores []int, out chan<- int) {
	for _, v := range valores {
		out <- v
	}
	close(out)
}

// multiplicador lê do canal 'in', multiplica por 2, e envia no canal 'out'.
// Fecha 'out' ao terminar.
func multiplicador(in <-chan int, out chan<- int) {
	for v := range in {
		out <- v * 2
	}
	close(out)
}

// impressora lê do canal 'in' e imprime cada valor.
// Envia true no canal 'done' ao terminar.
func impressora(in <-chan int, done chan<- bool) {
	for v := range in {
		fmt.Println(v)
	}
	done <- true
}

func main3() {
	valores := []int{1, 2, 3, 4, 5}

	c1 := make(chan int)
	c2 := make(chan int)
	done := make(chan bool)

	go gerador(valores, c1)
	go multiplicador(c1, c2)
	go impressora(c2, done)

	<-done
	fmt.Println("Pipeline concluído!")
}

// PERGUNTA) os canais do pipeline devem ser com ou sem buffer? Justifique.

// Os canais do pipeline devem ser sem buffer, pois se houver buffer nesses dois canais,
// os valores que estao sendo transmitidos de canal a canal podem ficar assincronos,
// levando a resultados nao esperados.