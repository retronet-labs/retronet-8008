# Istruzioni

Questo documento raccoglie lo stato di implementazione delle famiglie
istruzionali dell'Intel 8008.

---

## Load e Move

Stato: implementato.

Le istruzioni load e move trasferiscono byte tra registri, immediati e
pseudo-registro `M`.

---

## Cosa rappresentano

- `L(dst, src)`: copia il contenuto di `src` in `dst`.
- `LI(dst)`: carica nel registro `dst` il byte immediato successivo.
- `M`: non e' un registro fisico; indica la memoria puntata da `HL`.

I load non modificano i flag.

---

## Come funzionano nell'8008

I trasferimenti registro-registro usano il formato `11 DDD SSS`, dove `DDD` e'
il registro destinazione e `SSS` e' il registro sorgente. Il codice `111`
seleziona `M`.

I load immediati usano il formato `00 DDD 110` e consumano un secondo byte.
Quando `DDD = 111`, l'istruzione e' `LMI`: il byte immediato viene scritto nella
memoria puntata da `HL`.

Se sorgente e destinazione coincidono, l'istruzione non modifica lo stato ed e'
trattata come no-op.

---

## Come sono modellate nel progetto

Il package `cpu` espone helper mini-assembler:

```go
cpu.L(cpu.RegA, cpu.RegB) // LAB: A = B
cpu.L(cpu.RegM, cpu.RegA) // LMA: mem[HL] = A
cpu.L(cpu.RegA, cpu.RegM) // LAM: A = mem[HL]
cpu.LI(cpu.RegD)          // LDI: byte successivo -> D
cpu.LI(cpu.RegM)          // LMI: byte successivo -> mem[HL]
cpu.NOP()                 // LAA, no-op leggibile
```

`Step` esegue queste istruzioni direttamente; gli altri opcode restano
collegati allo stub `ErrUnimplementedOpcode`.

---

## Test coperti

- trasferimenti registro-registro
- no-op da load self
- load immediato su registro
- lettura `M` tramite `HL`
- scrittura `M` tramite `HL`
- `LMI`
- helper opcode

---

## ALU e Flags

Stato: implementato.

Le istruzioni ALU lavorano sull'accumulatore `A`, su un operando registro o `M`,
oppure su un byte immediato. Aggiornano `Carry`, `Zero`, `Sign` e `Parity`
secondo il risultato.

---

## Operazioni implementate

| Famiglia | Helper registro/M | Helper immediato | Effetto |
|----------|-------------------|------------------|---------|
| Add | `AD(src)` | `ADI()` | `A = A + src` |
| Add con Carry | `AC(src)` | `ACI()` | `A = A + src + Carry` |
| Subtract | `SU(src)` | `SUI()` | `A = A - src` |
| Subtract con borrow | `SB(src)` | `SBI()` | `A = A - src - Carry` |
| AND | `ND(src)` | `NDI()` | `A = A & src` |
| XOR | `XR(src)` | `XRI()` | `A = A ^ src` |
| OR | `OR(src)` | `ORI()` | `A = A | src` |
| Compare | `CP(src)` | `CPI()` | flags da `A - src`, `A` invariato |

Sono implementati anche:

- `INR(r)`: incrementa `B`, `C`, `D`, `E`, `H` o `L`
- `DCR(r)`: decrementa `B`, `C`, `D`, `E`, `H` o `L`

`INR` e `DCR` aggiornano `Zero`, `Sign` e `Parity`, ma non modificano `Carry`.

---

## Esempio

```go
mem.Write(0x0000, cpu.LI(cpu.RegA))
mem.Write(0x0001, 0x02)
mem.Write(0x0002, cpu.ADI())
mem.Write(0x0003, 0x03)
```

Dopo due `Step`, `A = 0x05`, `Carry = false`, `Zero = false`, `Sign = false`,
`Parity = true`.

---

## Test coperti

- addizione con e senza carry out
- addizione con carry in
- sottrazione con e senza borrow
- subtract con borrow in
- operazioni logiche e azzeramento Carry
- compare senza modifica dell'accumulatore
- immediati
- operando `M` tramite `HL`
- `INR`/`DCR` con Carry invariato
- helper opcode

---

## Da implementare

- Rotate.
- Control flow.
- HLT/STOPPED.
- I/O instructions.
