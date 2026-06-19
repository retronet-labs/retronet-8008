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

# esegue un binario raw con la CLI minimale
go run ./cmd/retronet-8008 -bin programma.bin -steps 1000

# disassembla 8 istruzioni senza eseguire
go run ./cmd/retronet-8008 -bin programma.bin -disasm 8

# esegue stampando ogni istruzione
go run ./cmd/retronet-8008 -bin programma.bin -steps 1000 -trace

# elenca profili macchina e slot ROM locali
go run ./cmd/retronet-8008 -profiles

# carica una ROM locale nello slot monitor del profilo Intellec
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -steps 1000
```

---

# Stato attuale

Il progetto ha completato le prime milestone fondamentali:

- struttura Go iniziale
- package `cpu` con stato base dell'Intel 8008
- memoria piatta da 16 KB
- I/O separato con 8 porte input e 24 porte output
- decoder tabellare da 256 opcode
- ciclo `Step` con fetch e incremento del PC
- istruzioni load/move
- istruzioni ALU e gestione flag
- istruzioni rotate dell'accumulatore
- istruzioni jump, call, return e restart
- `HLT`, stato `Stopped` e jam instruction esterna
- istruzioni I/O `INP` e `OUT`
- CLI runner per binari raw con dump registri
- disassembler minimale con contesto memoria
- trace istruzione per istruzione nella CLI
- profili macchina base e caricamento ROM locali
- documentazione italiana iniziale

Sono gia' modellati registri, flag, program counter a 14 bit, stack interno,
reset storico in stato fermo, memoria diretta, porte I/O, metadata del decoder e
le famiglie istruzionali iniziali eseguibili. I profili storici sono scheletri
senza ROM distribuite nel repository. Gli opcode non ancora implementati
restituiscono un errore esplicito.

---

# Struttura progetto

```text
retronet-8008/
|-- go.mod
|-- readme.md
|-- cmd/
|   `-- retronet-8008/
|       |-- main.go
|       `-- main_test.go
|-- machine/
|   |-- profile.go
|   `-- profile_test.go
|-- cpu/
|   |-- alu.go
|   |-- control.go
|   |-- cpu.go
|   |-- decoder.go
|   |-- disasm.go
|   |-- errors.go
|   |-- halt.go
|   |-- helpers.go
|   |-- io.go
|   |-- io_instructions.go
|   |-- jam.go
|   |-- load.go
|   |-- memory.go
|   |-- opcode.go
|   |-- opcodes.go
|   |-- rotate.go
|   |-- step.go
|   |-- alu_test.go
|   |-- control_test.go
|   |-- cpu_test.go
|   |-- decoder_test.go
|   |-- disasm_test.go
|   |-- halt_test.go
|   |-- helpers_test.go
|   |-- io_test.go
|   |-- io_instructions_test.go
|   |-- load_test.go
|   |-- memory_test.go
|   |-- rotate_test.go
|   |-- step_test.go
|   `-- test_helpers_test.go
|-- docs/
|   |-- architettura.md
|   |-- cli.md
|   |-- decoder.md
|   |-- disassembler.md
|   |-- flags.md
|   |-- io.md
|   |-- istruzioni.md
|   |-- memoria.md
|   |-- profili.md
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
5. Load e move. Completato.
6. ALU e flag. Completato.
7. Rotate. Completato.
8. Control flow e stack interno. Completato.
9. HLT, stopped e jam instruction. Completato.
10. I/O istruzionale. Completato.
11. CLI runner e tooling minimo. Completato.
12. Disassembler minimale. Completato.
13. Trace istruzione per istruzione. Completato.
14. Profili macchina e caricamento ROM locali. Completato.
15. ROM storiche reali, mappe memoria piu' precise, timing e T-state.

La roadmap dettagliata vive in `docs/roadmap.md`.

---

# Limiti noti

- La CLI carica binari raw e ROM locali via profilo, senza formati ROM
  strutturati.
- Le famiglie istruzionali principali sono implementate a livello funzionale.
- Il repository non include ROM storiche.
- Timing, T-state e dettagli elettrici dell'interrupt reale sono rimandati a milestone future.
