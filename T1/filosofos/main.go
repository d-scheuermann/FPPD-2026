package main

import (
	"flag"
	"fmt"
)

func main() {
	n := flag.Int("n", 5, "numero de filosofos")
	r := flag.Int("r", 10, "numero de refeicoes por filosofo")
	estrategia := flag.String("estrategia", "base",
		"estrategia: base, quebra ou semaforo")

	flag.Parse()

	fmt.Println("Filosofos:", *n)
	fmt.Println("Iteracoes:", *r)
	fmt.Println("Estrategia:", *estrategia)

	switch *estrategia {
	case "base":
		executarBase(*n, *r)

	case "quebra":
		executarQuebraSimetria(*n, *r)

	case "semaforo":
		executarSemaforo(*n, *r)

	default:
		fmt.Println("Estrategia invalida.")
	}
}