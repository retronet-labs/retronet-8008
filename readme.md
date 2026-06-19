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

Il progetto e' nella fase di stato CPU base:

- struttura Go iniziale
- package `cpu` con stato base dell'Intel 8008
- documentazione italiana iniziale
- nessun decoder e nessuna istruzione implementata

Sono gia' modellati registri, flag, program counter a 14 bit, stack interno e
reset storico in stato fermo.

---

# Struttura progetto

```text
retronet-8008/
├── go.mod
├── readme.md
├── cmd/
│   └── retronet-8008/
│       └── main.go
├── cpu/
├── docs/
│   ├── architettura.md
│   ├── flags.md
│   ├── registri.md
│   ├── roadmap.md
│   └── stack.md
├── examples/
│   └── README.md
├── testdata/
│   └── README.md
└── cpu/
    ├── cpu.go
    ├── helpers.go
    ├── opcodes.go
    ├── cpu_test.go
    └── helpers_test.go
```

Il layout segue volutamente quello di `go-4004`: il package `cpu` resta alla
radice ed e' importabile da CLI, esempi e test.

---

# Roadmap breve

1. Bootstrap del progetto. Completato.
2. Stato CPU base. Completato.
3. Memoria diretta a 16 KB e I/O separato.
4. Fetch, decoder e `Step`.
5. Famiglie istruzionali, una alla volta, con test e documentazione.

La roadmap dettagliata vive in `docs/roadmap.md`.

---

# Limiti noti

- La CLI non carica ancora programmi.
- Non esistono ancora memoria, I/O, decoder o istruzioni.
- Timing, T-state e interrupt/jam instruction sono rimandati a milestone future.
