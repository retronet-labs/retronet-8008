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

## Milestone successive

Le istruzioni saranno implementate per famiglie:

- ALU e flag
- rotate
- jump, call, return e restart
- halt e stato stopped
- input/output
- CLI con caricamento binario e dump registri
