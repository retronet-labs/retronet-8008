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
- memoria diretta da 16 KB
- bus I/O separato
- decoder opcode e ciclo `Step`
- disassembler con contesto memoria
- prime istruzioni load/move
- ALU instruction-level e gestione flag
- rotate dell'accumulatore
- control flow e stack interno
- halt, stopped e jam instruction esterna
- istruzioni I/O separate dalla memoria

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
- `FlatMemory` modella lo spazio diretto da 16 KB
- `Ports` modella le porte input/output separate
- `Decode` e `Step` gestiscono fetch e dispatch istruzione
- `Disassemble` legge opcode e operandi da memoria senza modificare stato CPU
- `Jam` esegue un'istruzione fornita dall'esterno senza fetch da memoria
- `L`, `LI` e `NOP` costruiscono opcode load/move per test ed esempi
- `HLT` costruisce l'opcode halt canonico
- gli helper ALU costruiscono opcode aritmetici e logici leggibili
- `RLC`, `RRC`, `RAL` e `RAR` costruiscono opcode rotate leggibili
- `JMP`, `JF`, `JT`, `CAL`, `CF`, `CT`, `RET`, `RF`, `RT` e `RST` costruiscono
  opcode di control flow leggibili
- `INP` e `OUT` costruiscono opcode I/O leggibili

---

## Implementato ora

- Stato CPU base.
- Reset deterministico.
- Costanti per registri e condizioni.
- Helper per indirizzi a 14 bit, HL e stack pointer.
- Memoria piatta da 16 KB.
- I/O separato con 8 porte input e 24 porte output.
- Decoder tabellare da 256 opcode.
- Disassembler minimale con bytes, operando e `NextPC`.
- `Step` con fetch opcode, fetch operandi e incremento `PC`.
- `Step` bloccato da `Halted` o `Stopped` con `ErrCPUStopped`.
- `Jam` come modello didattico dell'istruzione forzata da interrupt esterno.
- Istruzioni load/move tra registri, immediati e `M`.
- Istruzioni ALU con registri, `M` e immediati.
- Istruzioni rotate dell'accumulatore.
- Istruzioni jump, call, return e restart con stack interno.
- `HLT` e alias `0x00`/`0x01`.
- Istruzioni `INP` e `OUT` su bus I/O separato.
- Test automatici sullo stato iniziale e sui mascheramenti.

---

## Da implementare

- Dettagli elettrici e temporali di interrupt e jam instruction reali.
- Timing e T-state.
