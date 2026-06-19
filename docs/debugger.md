# Trace strutturato e debugger

Il package `machine` offre un debugger instruction-level sopra `FrontPanel`,
senza aggiungere dipendenze al core CPU. Gli eventi combinano disassembly,
stato prima/dopo, timing e side-effect osservati.

---

## Eventi

`TraceEvent` e' serializzabile in JSON e contiene:

- sequenza, tipo evento, PC, opcode e byte istruzione
- disassembly e indicazione di jam da interrupt
- copia completa della CPU prima e dopo
- `InstructionTiming` effettivo
- scritture memoria con valore precedente, richiesto ed effettivo
- trasferimenti I/O con direzione, porta e valore
- ciclo macchina sul quale READY ha prodotto WAIT
- descrizione dell'eventuale breakpoint

I tipi correnti sono `instruction`, `wait` e `breakpoint`. Una istruzione `HLT`
produce comunque il proprio evento prima dello stop `cpu-stopped`.

---

## ObservableMemory

`ObservableMemory` decora qualsiasi `cpu.Memory`. Le scritture vengono delegate
al bus originale e poi notificate con il valore effettivamente leggibile.
Questo permette di distinguere una scrittura RAM riuscita da un tentativo ROM
ignorato.

I loader `LoadBytes` e `LoadROM` restano privilegiati e non generano eventi
runtime. La CLI usa sempre questo wrapper, anche quando il debugger non e'
attivo.

---

## Breakpoint e watchpoint

`Debugger` supporta:

- breakpoint PC prima dell'istruzione
- breakpoint opcode prima dell'istruzione
- watchpoint dopo una scrittura a un indirizzo
- breakpoint dopo input o output su una porta
- stop distinti per WAIT, richiesta pannello, HLT e limite

I breakpoint PC/opcode non modificano lo stato CPU. Watchpoint memoria e I/O
completano invece l'istruzione responsabile e poi fermano il run.

---

## CLI

Le opzioni sono ripetibili:

```powershell
go run ./cmd/retronet-8008 -bin programma.bin -break 0x0100
go run ./cmd/retronet-8008 -bin programma.bin -break-opcode 0x00
go run ./cmd/retronet-8008 -bin programma.bin -watch 0x0200
go run ./cmd/retronet-8008 -bin programma.bin -break-input 0 -break-output 8
go run ./cmd/retronet-8008 -bin programma.bin -trace-json trace.jsonl
```

`-trace-json` crea o sostituisce un file JSON Lines, un oggetto per evento. Il
dump umano resta su stdout e non contamina il file strutturato. Il campo
`stop_reason` puo' assumere anche `breakpoint`, `watchpoint` e `io-breakpoint`.

Il vecchio `-trace` testuale resta disponibile e conserva il formato compatto
prima dell'esecuzione.

---

## Limiti

- I breakpoint sono instruction-level, non pin/T-state-level.
- Non esistono ancora simboli o sorgenti assembly associate agli indirizzi.
- Il trace include le copie CPU complete e privilegia chiarezza e stabilita'
  rispetto alla dimensione minima del file.
