# Architettura Intel 8008

L'Intel 8008 e' una CPU a 8 bit con program counter a 14 bit, sette registri
visibili al programmatore, quattro flag principali e uno stack interno di
indirizzi. Il core di questo progetto modella prima il comportamento a livello
di istruzione; timing, bus multiplexato e T-state saranno aggiunti solo dopo una
base funzionale solida.

---

## Cosa rappresenta

Questo documento descrive il modello generale del core CPU:

- registri `A`, `B`, `C`, `D`, `E`, `H`, `L`
- flag `Carry`, `Zero`, `Sign`, `Parity`
- `PC` a 14 bit
- stack interno da 8 voci a 14 bit
- stati `Halted` e `Stopped`

---

## Come funziona nell'8008

L'8008 indirizza direttamente 16 KB di memoria, da `0x0000` a `0x3FFF`.
Le istruzioni sono lunghe 1, 2 o 3 byte. Lo stack non e' memoria normale: e'
interno al processore e conserva indirizzi di ritorno per `CALL`, `RET` e `RST`.

All'accensione storica il processore entra in stato fermo e viene avviato da un
protocollo esterno basato su interrupt e jam instruction.

---

## Come e' modellato nel progetto

Il package `cpu` espone una struct `CPU8008` con campi pubblici e leggibili, in
stile `go-4004`. I vincoli hardware sono concentrati in helper piccoli:

- `addr14` maschera gli indirizzi a 14 bit
- `hlAddress` calcola l'indirizzo del pseudo-registro `M`
- `stackIndex` maschera lo stack pointer a 3 bit

---

## Implementato ora

- Stato CPU base.
- Reset deterministico.
- Costanti per registri e condizioni.
- Helper per indirizzi a 14 bit, HL e stack pointer.
- Test automatici sullo stato iniziale e sui mascheramenti.

---

## Da implementare

- Memoria e I/O.
- Fetch, decoder e `Step`.
- Istruzioni 8008.
- Semantica completa di `HLT`, `STOPPED`, interrupt e jam instruction.
- Timing e T-state.
