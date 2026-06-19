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

## Milestone successive

- disassembler con contesto memoria
- trace istruzione per istruzione
- profili macchina e ROM storiche
- timing e T-state
