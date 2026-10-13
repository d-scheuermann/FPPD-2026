# T1 - Problemas Clássicos de Concorrência

FPPD - Fundamentos de Processamento Paralelo e Distribuído (98713-04), 2026/2

Alunos: Daniel Scheuermann, Davi Bandeira de Lima, João Gabriel Curtinaz Maron, Gabriel Ulbrich Bartzen

## Estrutura

```
T1/
├── go.mod
├── filosofos/                    Problema 1: Jantar dos Filósofos
│   ├── main.go                   leitura dos parâmetros e escolha da estratégia
│   ├── base.go                   versão base (demonstra o deadlock)
│   ├── estrategia-quebrasim.go   estratégia 1: quebra de simetria
│   ├── estrategia-semaforo.go    estratégia 2: no máximo N-1 filósofos à mesa
│   └── metricas.go               refeições e tempo médio de espera
└── produtor-consumidor/          Problema 2: Produtor/Consumidor com buffer limitado
    ├── channel/v-channel.go      versão com channel
    └── semaforos/v-semaforos.go  versão com semáforos
```

## Requisitos

- Go 1.22 ou superior
- Nenhuma biblioteca externa

Todos os comandos abaixo devem ser executados a partir da pasta `T1/`. O parâmetro `-race` ativa o detector de data race; ele pode ser omitido para execuções mais rápidas.

## Compilação

```
go build ./...
```

Para gerar executáveis:

```
go build -o bin/filosofos ./filosofos
go build -o bin/pc-channel ./produtor-consumidor/channel
go build -o bin/pc-semaforos ./produtor-consumidor/semaforos
```

## Problema 1: Jantar dos Filósofos

```
go run -race ./filosofos -estrategia base     -n 5 -r 10
go run -race ./filosofos -estrategia quebra   -n 5 -r 100
go run -race ./filosofos -estrategia semaforo -n 5 -r 100
```

| Parâmetro | Descrição | Padrão |
|---|---|---|
| `-n` | número de filósofos (mínimo 2) | 5 |
| `-r` | número de refeições por filósofo | 10 |
| `-estrategia` | `base`, `quebra` ou `semaforo` | `base` |

- `base`: todos pegam primeiro o garfo esquerdo. O programa força a espera circular, detecta o deadlock e encerra sem deixar goroutines bloqueadas.
- `quebra`: o último filósofo pega os garfos na ordem inversa (hierarquia de recursos).
- `semaforo`: um semáforo contador limita a N-1 os filósofos que disputam garfos ao mesmo tempo.

Ao final, o programa imprime o número de refeições e o tempo médio de espera de cada filósofo.

## Problema 2: Produtor/Consumidor

### Versão com channel

```
go run -race ./produtor-consumidor/channel
go run -race ./produtor-consumidor/channel -k 10 -p 4 -c 4 -n 500
go run -race ./produtor-consumidor/channel -prod 2ms -cons 1ms
go run ./produtor-consumidor/channel -csv > resultados.csv
```

| Parâmetro | Descrição | Padrão |
|---|---|---|
| `-p` | número de produtores | 4 |
| `-c` | número de consumidores | 4 |
| `-k` | capacidade do buffer (0 testa K = 1, 10 e 100) | 0 |
| `-n` | itens gerados por produtor | 500 |
| `-prod` | tempo médio de produção por item | 1ms |
| `-cons` | tempo médio de consumo por item | 1ms |
| `-timeout` | timeout do consumidor 0 (`select` com `time.After`) | 2ms |
| `-reps` | repetições por valor de K (reporta a média) | 3 |
| `-csv` | saída em formato CSV | false |

Saída: tempo, throughput, ocupação média do buffer, timeouts do consumidor 0 e itens por consumidor. Ao final, o programa verifica que o total produzido é igual ao consumido e que cada item foi consumido exatamente uma vez.

### Versão com semáforos

```
go run -race ./produtor-consumidor/semaforos -p 4 -c 4 -k 10 -s 3
```

| Parâmetro | Descrição | Padrão |
|---|---|---|
| `-p` | número de produtores | 100 |
| `-c` | número de consumidores | 140 |
| `-k` | capacidade do buffer | 1000 |
| `-s` | segundos de produção | 20 |

Os produtores trabalham durante `-s` segundos. Depois disso, os consumidores esvaziam o buffer e encerram (o término pode levar até 1 s a mais, intervalo do timeout dos consumidores). Ao final, o programa imprime o total produzido, o total consumido, o throughput e a ocupação média do buffer.
