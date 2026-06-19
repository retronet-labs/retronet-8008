# Roadmap RetroNet 8008

Questa roadmap guida lo sviluppo dell'emulatore Intel 8008 in milestone piccole,
testabili e documentate. Ogni milestone deve compilare e passare `go test ./...`.

---

## Milestone 0 - Bootstrap progetto

Stato: completata.

Contenuto:

- modulo Go
- struttura directory coerente con `go-4004`
- README iniziale
- documentazione roadmap
- CLI minimale

---

## Milestone 1 - Stato CPU base

Stato: completata.

Contenuto:

- registri `A`, `B`, `C`, `D`, `E`, `H`, `L`
- flag `Carry`, `Zero`, `Sign`, `Parity`
- program counter a 14 bit
- stack interno a 8 voci da 14 bit
- stati `Halted` e `Stopped`
- `Reset()`
- helper di mascheramento indirizzi e indirizzo `HL`
- test sul reset e sui vincoli a 14 bit

---

## Milestone 2 - Memoria e I/O

Stato: completata.

Contenuto:

- memoria piatta da 16 KB
- mascheramento indirizzi a `0x3FFF`
- I/O separato dalla memoria
- 8 porte input
- 24 porte output
- test memoria e I/O

---

## Milestone 3 - Fetch, decoder e Step

Stato: completata.

Contenuto:

- fetch opcode da `PC`
- incremento `PC` mascherato a 14 bit
- tabella decoder da 256 opcode
- errore esplicito per opcode non implementati
- `Step()`
- test su decoder, fetch, PC e opcode non implementati

---

## Milestone 4 - Load e move

Stato: completata.

Contenuto:

- trasferimenti tra registri
- load immediati
- pseudo-registro `M`
- lettura e scrittura tramite `HL`
- helper mini-assembler `L`, `LI` e `NOP`
- test su registri, immediati e memoria

---

## Milestone 5 - ALU e flags

Stato: completata.

Contenuto:

- ADD, ADC, SUB, SBB
- AND, XOR, OR, CMP
- operandi registro, `M` e immediati
- `INR` e `DCR`
- aggiornamento `Carry`, `Zero`, `Sign` e `Parity`
- `CMP` senza modifica di `A`
- test su edge case di carry, borrow, zero, sign e parity

---

## Milestone 6 - Rotate

Stato: completata.

Contenuto:

- `RLC`
- `RRC`
- `RAL`
- `RAR`
- modifica del solo flag Carry
- test con Carry iniziale 0 e 1

---

## Milestone 7 - Control flow e stack

Stato: completata.

Contenuto:

- `JMP`, `JF` e `JT`
- `CAL`, `CF` e `CT`
- `RET`, `RF` e `RT`
- `RST n`
- target a 14 bit con low byte e high byte mascherato
- stack interno con `SP` a 3 bit e PC corrente in `Stack[SP]`
- overflow ciclico senza fault
- helper mini-assembler per control flow
- test su salto, call, return, restart e profondita' stack

---

## Milestone 8 - HLT, stopped e jam instruction

Stato: completata.

Contenuto:

- `HLT` e alias `0x00`/`0x01`
- helper mini-assembler `HLT()`
- `Step` bloccato da `Halted` o `Stopped`
- errore `ErrCPUStopped`
- `Jam(mem, io, code, operands...)` per istruzioni forzate dall'esterno
- validazione operandi jammed con `ErrInvalidJamInstruction`
- ripartenza da reset o halt tramite jam di `NOP` o `RST`
- test su stop, halt, jam e return PC nello stack interno

---

## Milestone 9 - I/O istruzionale

Stato: completata.

Contenuto:

- `INP` da porte input `0..7` verso `A`
- `OUT` da `A` verso porte output `8..31`
- helper mini-assembler `INP(port)` e `OUT(port)`
- decoder corretto per pattern `0100 MMM1` e `01 RRMMM1`
- errore `ErrNilIO` per esecuzione senza bus I/O
- test su flags invariati, jam I/O e mapping completo delle porte

---

## Milestone 10 - CLI runner minimale

Stato: completata.

Contenuto:

- comando `cmd/retronet-8008`
- caricamento di binari raw in `FlatMemory`
- opzioni `-bin`, `-addr`, `-pc` e `-steps`
- avvio tramite jam di `JMP` al PC iniziale
- esecuzione con limite di istruzioni
- stop pulito su `HLT`/`Stopped`
- dump finale di registri, flag, PC, SP e stack interno
- test unitari della CLI senza avviare processi esterni

---

## Milestone 11 - Disassembler minimale

Stato: completata.

Contenuto:

- `cpu.Disassemble(mem, pc)`
- struct `Disassembly` con PC, opcode, bytes, operando e `NextPC`
- formattazione compatta `0000: 06 2A    LAI #0x2A`
- lettura operandi da memoria con wrap a 14 bit
- target 14 bit mascherati come in esecuzione
- opzione CLI `-disasm N`
- test su istruzioni 1/2/3 byte, wrap e CLI

---

## Milestone 12 - Trace istruzione per istruzione

Stato: completata.

Contenuto:

- opzione CLI `-trace`
- disassembly del PC corrente prima di ogni `Step`
- trace coerente con salti, call, return e halt
- integrazione con limite `-steps`
- dump finale invariato dopo il trace
- test su trace di esecuzione e limite step

---

## Milestone 13 - Profili macchina e ROM locali

Stato: completata.

Contenuto:

- package `machine`
- catalogo profili `generic`, `intellec-8`, `scelbi-8b` e `scelbi-8h`
- slot ROM `monitor` per profili storici iniziali
- caricamento ROM locali con validazione dello spazio 14 bit
- opzioni CLI `-profile`, `-profiles` e `-rom`
- nessuna ROM storica inclusa nel repository
- test su profili, caricamento ROM e integrazione CLI

---

## Milestone 14 - SCELBI/Intellec e I/O callback

Stato: completata.

Contenuto:

- spiegazione storica di SCELBI 8H/8B e Intel Intellec 8
- metadata `MemoryRegion`, `IOPort` e `ROMHint`
- slot ROM locale `test` separato dallo slot `monitor`
- bus `machine.CallbackIO` con latch e callback per porta
- opzioni CLI `-input porta=valore` e `-io-trace`
- ROM locale di smoke test `INP 0`, `OUT 8`, `HLT`
- porte callback `0` e `8` dichiarate come convenzioni non storiche
- test su callback, copie dei profili e integrazione ROM/I/O

---

## Milestone 15 - Bus memoria mappato

Stato: completata.

Contenuto:

- `machine.MemoryBus` compatibile con `cpu.Memory`
- regioni RAM, ROM e mixed con validazione delle sovrapposizioni
- `Profile.NewMemory()` usato dalla CLI
- caricamento privilegiato delle immagini ROM
- protezione delle ROM dalle scritture CPU e dai loader raw
- open bus convenzionale a `0xFF` per memoria non mappata
- profilo `generic` interamente RAM
- profili storici mixed senza mappe ROM/RAM inventate
- test unitari e integrazione CLI

---

## Milestone 16 - Terminale virtuale

Stato: completata.

Contenuto:

- `machine.Terminal` buffered e indipendente dal core CPU
- coda input ASCII sulla porta convenzionale `0`
- output ASCII verso `io.Writer` dalla porta convenzionale `8`
- fallback al latch input quando la coda e' vuota
- osservatori `CallbackIO` componibili con le callback delle periferiche
- opzioni CLI `-terminal` e `-terminal-input`
- compatibilita' con `-io-trace`
- test unitari e ROM echo end-to-end

---

## Milestone 17 - Front panel

Stato: completata.

Contenuto:

- `machine.FrontPanel` sopra CPU, memoria e I/O esistenti
- switch dati a 8 bit e indirizzo a 14 bit
- examine e deposit con rispetto della protezione ROM
- step e run con motivi di arresto distinti
- richiesta `Stop()` separata dallo stato HLT/stopped della CPU
- reset, jam instruction e interrupt vettorizzato `RST 0..7`
- snapshot del pannello e observer pre-istruzione
- CLI coordinata dal front panel
- opzioni `-panel`, `-panel-switches` e `-panel-address`
- test unitari e integrazione CLI

---

## Milestone successive

- ROM storiche reali, solo quando provenance e licenze saranno chiare
- mappe memoria storiche piu' precise per Intellec e SCELBI
- cassette e altre periferiche virtuali
- timing e T-state
