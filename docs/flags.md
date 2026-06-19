# Flags

L'8008 espone quattro flag principali: Carry, Zero, Sign e Parity.

---

## Cosa rappresentano

- `Carry`: riporto o prestito nelle operazioni aritmetiche.
- `Zero`: risultato uguale a zero.
- `Sign`: bit piu' alto del risultato impostato.
- `Parity`: parita' pari del risultato.

---

## Come funzionano nell'8008

Le istruzioni aritmetiche e logiche aggiornano i flag in base al risultato.
Alcune famiglie hanno regole speciali:

- `INR` e `DCR` non modificano Carry.
- Le rotate modificano solo Carry.
- `CMP` aggiorna i flag come una sottrazione, ma non modifica `A`.

Le istruzioni condizionali selezionano un flag con due bit:

| Codice | Flag |
|--------|------|
| `00` | Carry |
| `01` | Zero |
| `10` | Sign |
| `11` | Parity |

---

## Come sono modellati nel progetto

`CPU8008` contiene i campi booleani `Carry`, `Zero`, `Sign` e `Parity`. Le
costanti `CondCarry`, `CondZero`, `CondSign` e `CondParity` fissano i codici
condizione per decoder e helper futuri.

---

## Implementato ora

- Campi flag nello stato CPU.
- Reset dei flag a `false`.
- Costanti condizione.

---

## Da implementare

- Aggiornamento flag nelle istruzioni ALU.
- Parity helper per calcolare la parita' pari.
- Semantica di Carry per sottrazioni e confronti.
- Salti, call e return condizionali.
