# retronet-8008 - Emulatore Intel 8008

Un emulatore dell'Intel 8008 scritto in Go, sviluppato passo passo con approccio
didattico e coerente con il modulo RetroNet dedicato al 4004.

L'obiettivo e' costruire un core 8008 istruzione-accurato, ben testato e
documentato in italiano. Il progetto e' separato dagli emulatori 4004: puo'
seguirne stile, organizzazione e disciplina, ma non importa codice da altri
repository.

---

# Quick Start

```bash
# lancia tutti i test
go test ./...

# esegue la CLI minimale
go run ./cmd/retronet-8008
```

---

# Stato attuale

Il progetto ha completato le prime quattro basi:

- struttura Go iniziale
- package `cpu` con stato base dell'Intel 8008
- memoria piatta da 16 KB
- I/O separato con 8 porte input e 24 porte output
- decoder tabellare da 256 opcode
- ciclo `Step` con fetch e incremento del PC
- documentazione italiana iniziale

Non esistono ancora istruzioni eseguibili: `Step` consuma opcode e operandi, poi
restituisce un errore esplicito di opcode non implementato. Sono gia' modellati
registri, flag, program counter a 14 bit, stack interno, reset storico in stato
fermo, memoria diretta, porte I/O e metadata del decoder.

---

# Struttura progetto

```text
retronet-8008/
|-- go.mod
|-- readme.md
|-- cmd/
|   `-- retronet-8008/
|       `-- main.go
|-- cpu/
|   |-- cpu.go
|   |-- decoder.go
|   |-- errors.go
|   |-- helpers.go
|   |-- io.go
|   |-- memory.go
|   |-- opcode.go
|   |-- opcodes.go
|   |-- step.go
|   |-- cpu_test.go
|   |-- decoder_test.go
|   |-- helpers_test.go
|   |-- io_test.go
|   |-- memory_test.go
|   `-- step_test.go
|-- docs/
|   |-- architettura.md
|   |-- decoder.md
|   |-- flags.md
|   |-- io.md
|   |-- memoria.md
|   |-- registri.md
|   |-- roadmap.md
|   `-- stack.md
|-- examples/
|   `-- README.md
`-- testdata/
    `-- README.md
```

Il layout segue volutamente quello di `go-4004`: il package `cpu` resta alla
radice ed e' importabile da CLI, esempi e test.

---

# Roadmap breve

1. Bootstrap del progetto. Completato.
2. Stato CPU base. Completato.
3. Memoria diretta a 16 KB e I/O separato. Completato.
4. Fetch, decoder e `Step`. Completato.
5. Load e move.
6. Famiglie istruzionali successive, una alla volta, con test e documentazione.

La roadmap dettagliata vive in `docs/roadmap.md`.

---

# Limiti noti

- La CLI non carica ancora programmi.
- Non esistono ancora istruzioni implementate.
- Timing, T-state e interrupt/jam instruction sono rimandati a milestone future.
