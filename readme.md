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
go run ./cmd/retronet-8008 -bin programma.rom -steps 1000

# disassembla 8 istruzioni senza eseguire
go run ./cmd/retronet-8008 -bin programma.rom -disasm 8

# esegue stampando ogni istruzione
go run ./cmd/retronet-8008 -bin programma.rom -steps 1000 -trace

# elenca profili macchina e slot ROM locali
go run ./cmd/retronet-8008 -profiles

# carica una ROM locale nello slot monitor del profilo Intellec
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -steps 1000

# esegue una ROM locale con input e trace I/O a callback
go run ./cmd/retronet-8008 -profile scelbi-8b -rom test=io-smoke.bin -input 0=0x5A -io-trace -steps 8

# accoda input ASCII e collega l'output del terminale
go run ./cmd/retronet-8008 -profile scelbi-8b -rom test=io-smoke.bin -terminal-input Z -steps 8

# mostra front panel, switch dati e byte memoria selezionato
go run ./cmd/retronet-8008 -bin programma.rom -panel-switches 0x4B -panel-address 0x0100
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
- profili SCELBI/Intellec documentati con metadata memoria e I/O
- bus I/O a callback con input configurabile e trace CLI
- bus memoria mappato con protezione ROM
- terminale ASCII buffered componibile con il trace I/O
- front panel con jam, step/run/stop, switch ed examine/deposit
- timing Intel con cicli macchina e contatori cumulativi
- READY/WAIT per ciclo macchina e interrupt sincronizzato
- trace JSON e debugger con breakpoint/watchpoint
- conformance sintetica e verifica SHA-256 di ROM locali
- periferiche configurabili con ownership e loopback generico
- conformance esaustiva ALU e matrice completa dei 256 opcode
- fuzz test per decoder, disassembler e vincoli architetturali
- documentazione italiana iniziale

Sono modellati registri, flag, program counter a 14 bit, stack interno, reset
storico in stato fermo, memoria diretta, porte I/O, metadata del decoder e tutte
le famiglie istruzionali documentate. Dei 256 byte possibili, 250 sono encoding
definiti; `22`, `2A`, `32`, `38`, `39` e `3A` sono slot non definiti e
restituiscono `ErrUnimplementedOpcode`. I profili storici restano scheletri
senza ROM distribuite nel repository.

---

# Struttura progetto

```text
retronet-8008/
|-- go.mod
|-- AGENTS.md
|-- readme.md
|-- cmd/
|   `-- retronet-8008/
|       |-- main.go
|       `-- main_test.go
|-- conformance/
|   |-- rom.go
|   |-- suite.go
|   |-- suite_test.go
|   `-- synthetic.go
|-- machine/
|   |-- debugger.go
|   |-- debugger_test.go
|   |-- frontpanel.go
|   |-- frontpanel_test.go
|   |-- io.go
|   |-- io_test.go
|   |-- memory.go
|   |-- memory_test.go
|   |-- observable_memory.go
|   |-- observable_memory_test.go
|   |-- peripheral.go
|   |-- peripheral_test.go
|   |-- profile.go
|   |-- profile_test.go
|   |-- terminal.go
|   `-- terminal_test.go
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
|   |-- control-lines.md
|   |-- conformance.md
|   |-- debugger.md
|   |-- decoder.md
|   |-- disassembler.md
|   |-- flags.md
|   |-- front-panel.md
|   |-- io.md
|   |-- istruzioni.md
|   |-- memoria.md
|   |-- periferiche.md
|   |-- profili.md
|   |-- registri.md
|   |-- roadmap.md
|   |-- stack.md
|   |-- terminale.md
|   `-- timing.md
|-- examples/
|   `-- README.md
`-- testdata/
    `-- README.md
```

Il layout segue volutamente quello di `go-4004`: il package `cpu` resta alla
radice ed e' importabile da CLI, esempi e test.

---

# Roadmap breve

- Milestone 0-17: bootstrap, core, tooling, profili e front panel. Completate.
- Milestone 18: mappe storiche verificate. Rinviata in attesa delle fonti.
- Milestone 19: ROM storiche. Rinviata in attesa di provenienza e licenze.
- Milestone 20: timing Intel e cicli macchina. Completata.
- Milestone 21: READY e interrupt al confine PCI. Completata.
- Milestone 22: trace strutturato e debugger. Completata.
- Milestone 23: conformance sintetica e verifica ROM locale. Completata.
- Milestone 24: periferiche generiche configurabili. Completata.
- Milestone 25: matrice opcode, oracle esaustivi e fuzz test. Completata.

La roadmap dettagliata vive in `docs/roadmap.md`.

La checklist che definisce il perimetro e i gate della prima release pubblica
e' in [`docs/release-v0.1.0.md`](docs/release-v0.1.0.md).

---

# Limiti noti

- La CLI carica binari raw e ROM locali via profilo, senza formati ROM
  strutturati.
- Le famiglie istruzionali documentate sono implementate a livello funzionale,
  ma non ancora validate contro un secondo emulatore indipendente.
- Il repository non include ROM storiche.
- Le porte callback `0` e `8` sono convenzioni di test, non mappe storiche
  definitive.
- I profili storici proteggono le ROM caricate, ma non dichiarano ancora una
  ripartizione ROM/RAM storicamente verificata.
- I costi in stati e i cicli macchina sono modellati; le singole transizioni di
  pin/T-state e i dettagli elettrici dell'interrupt reale non lo sono ancora.
