# Istruzioni

Questo documento raccoglie lo stato di implementazione delle famiglie
istruzionali dell'Intel 8008.

Gli esempi che chiamano `Step` assumono che la CPU sia gia' in stato running,
ad esempio dopo una jam iniziale di `NOP`.

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

## Rotate

Stato: implementato.

Le rotate lavorano solo sull'accumulatore `A` e sul flag `Carry`. Non modificano
`Zero`, `Sign` o `Parity`.

---

## Operazioni implementate

| Istruzione | Helper | Effetto |
|------------|--------|---------|
| RLC | `RLC()` | ruota `A` a sinistra; bit 7 va in bit 0 e in Carry |
| RRC | `RRC()` | ruota `A` a destra; bit 0 va in bit 7 e in Carry |
| RAL | `RAL()` | ruota `A` a sinistra attraverso Carry |
| RAR | `RAR()` | ruota `A` a destra attraverso Carry |

---

## Esempio

```go
c.A = 0b1000_0000
c.Carry = true
mem.Write(0x0000, cpu.RAL())
```

Dopo `Step`, `A = 0b0000_0001` e `Carry = true`: il vecchio bit 7 esce in
Carry, mentre il vecchio Carry entra nel bit 0.

---

## Test coperti

- `RLC` e `RRC`
- `RAL` e `RAR` con Carry iniziale 0 e 1
- Carry in uscita
- `Zero`, `Sign` e `Parity` invariati
- avanzamento del PC
- helper opcode

---

## Control Flow e Stack

Stato: implementato.

Le istruzioni di controllo di flusso modificano `PC` e, per call/restart/return,
usano lo stack interno a 8 voci dell'8008.

---

## Operazioni implementate

| Famiglia | Helper | Effetto |
|----------|--------|---------|
| Jump incondizionato | `JMP()` | `PC = target` |
| Jump se flag false | `JF(cond)` | salta se il flag selezionato e' `false` |
| Jump se flag true | `JT(cond)` | salta se il flag selezionato e' `true` |
| Call incondizionata | `CAL()` | salva ritorno e salta a `target` |
| Call se flag false | `CF(cond)` | call condizionata su flag `false` |
| Call se flag true | `CT(cond)` | call condizionata su flag `true` |
| Return incondizionato | `RET()` | ripristina il PC dallo stack interno |
| Return se flag false | `RF(cond)` | return condizionato su flag `false` |
| Return se flag true | `RT(cond)` | return condizionato su flag `true` |
| Restart | `RST(n)` | call al vettore `n * 8` |

I target a 3 byte sono codificati little-endian: byte basso, poi byte alto. Del
byte alto vengono usati solo i 6 bit bassi, quindi l'indirizzo resta sempre nel
range `0x0000-0x3FFF`.

---

## Esempio

```go
mem.Write(0x0000, cpu.CAL())
mem.Write(0x0001, 0x10)
mem.Write(0x0002, 0x00)
mem.Write(0x0010, cpu.RET())
```

Dopo la `CAL`, `PC = 0x0010`, `SP = 1` e `Stack[0] = 0x0003`. Dopo `RET`,
`PC = 0x0003` e `SP = 0`.

---

## Test coperti

- jump incondizionato e indirizzi a 14 bit
- jump condizionali presi e non presi
- call e return con indirizzo di ritorno
- call e return condizionali presi e non presi
- restart verso `n * 8`
- profondita' utile 7 dello stack interno
- overflow ciclico senza errore
- helper opcode

---

## HLT, Stopped e Jam

Stato: implementato.

`HLT` ferma la CPU impostando sia `Halted` sia `Stopped`. Dopo il reset storico
la CPU parte gia' in questo stato: una chiamata diretta a `Step` non effettua
fetch da memoria e restituisce `ErrCPUStopped`.

---

## Come riparte la CPU

L'8008 viene riavviato da una istruzione forzata dall'esterno. Nel progetto
questo e' modellato da `Jam(mem, io, code, operands...)`, che:

- valida il numero di operandi atteso dal decoder
- cancella `Halted` e `Stopped`
- esegue l'istruzione fornita senza leggere l'opcode dalla memoria

Per avviare un programma in memoria durante test o esempi si puo' usare una jam
di `NOP`, che lascia invariato il `PC`:

```go
c := cpu.NewCPU8008()
_ = c.Jam(nil, nil, cpu.NOP())
```

Per modellare un interrupt reale e' piu' tipico usare `RST(n)`, cosi' il PC
corrente resta nello stack interno e l'esecuzione riparte dal vettore `n * 8`.

---

## Test coperti

- `HLT` sugli encoding `0x00`, `0x01` e `0xFF`
- `Step` bloccato dopo reset o halt
- jam di `NOP` per entrare in stato running
- validazione del numero di operandi jammed
- jam di `RST` da stato stopped
- conservazione del return PC dopo `HLT` e `RST`

---

## Input/Output

Stato: implementato.

Le istruzioni I/O usano il bus separato dalla memoria:

| Istruzione | Helper | Effetto |
|------------|--------|---------|
| INP | `INP(port)` | legge la porta input `0..7` in `A` |
| OUT | `OUT(port)` | scrive `A` sulla porta output `8..31` |

`INP` usa il pattern opcode `0100 MMM1`, dove `MMM` seleziona la porta input.
`OUT` usa `01 RRMMM1` con `RR != 00`: i cinque bit `RRMMM` selezionano una
porta output da `8` a `31`.

Entrambe le istruzioni lasciano invariati i flag. Se vengono eseguite senza bus
I/O, restituiscono `ErrNilIO`.

---

## Esempio

```go
ports := cpu.NewPorts()
_ = ports.SetInput(3, 0xA5)
mem.Write(0x0000, cpu.INP(3))
mem.Write(0x0001, cpu.OUT(16))
```

Dopo due `Step`, `A = 0xA5` e la porta output `16` contiene `0xA5`.

---

## Test coperti

- lettura `INP` in accumulatore
- scrittura `OUT` da accumulatore
- flags invariati
- errore `ErrNilIO`
- esecuzione via `Jam`
- mapping completo delle 8 porte input e 24 porte output
- helper opcode

---

## Verifica trasversale

- La matrice della milestone 25 esegue tutti i 250 encoding definiti e controlla
  che i sei slot non definiti producano l'errore previsto.
- Gli oracle test-only coprono esaustivamente ALU, flag, rotate e
  incremento/decremento.
- Il trace JSON include stato CPU prima/dopo, timing, memoria e I/O.
