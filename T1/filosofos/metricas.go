package main

import (
	"fmt"
	"time"
)

type MetricaFilosofo struct {
	Refeicoes   int
	EsperaTotal time.Duration
}

func imprimirMetricas(metricas []MetricaFilosofo) {
	fmt.Println("\n--- METRICAS ---")

	for i, m := range metricas {
		var media time.Duration

		if m.Refeicoes > 0 {
			media = m.EsperaTotal / time.Duration(m.Refeicoes)
		}

		fmt.Printf(
			"Filosofo %d: refeicoes=%d | espera media=%v\n",
			i,
			m.Refeicoes,
			media,
		)
	}
}
