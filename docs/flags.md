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
condizione per decoder, helper e istruzioni condizionali.

La funzione interna che aggiorna `Zero`, `Sign` e `Parity` usa il byte risultato.
`Parity` vale true quando il numero di bit a 1 e' pari.

---

## Implementato ora

- Campi flag nello stato CPU.
- Reset dei flag a `false`.
- Costanti condizione.
- Aggiornamento flag nelle istruzioni ALU.
- Carry su addizione come riporto oltre 8 bit.
- Carry su sottrazione e compare come borrow.
- Carry azzerato dalle operazioni logiche `ND`, `XR` e `OR`.
- `CMP`/`CP` aggiorna i flag senza modificare `A`.
- `INR` e `DCR` aggiornano `Zero`, `Sign` e `Parity` senza modificare Carry.
- Rotate aggiorna solo Carry e lascia invariati `Zero`, `Sign` e `Parity`.
- Jump, call e return condizionali leggono i flag senza modificarli.

---

## Da implementare

- Eventuali dettagli temporali delle istruzioni condizionali.
