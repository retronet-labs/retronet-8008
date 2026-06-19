# Registri

L'8008 ha sette registri a 8 bit visibili al programmatore: `A`, `B`, `C`, `D`,
`E`, `H`, `L`.

---

## Cosa rappresentano

- `A` e' l'accumulatore usato dalle istruzioni aritmetiche e logiche.
- `B`, `C`, `D`, `E`, `H`, `L` sono registri generali.
- `M` non e' un registro fisico: indica la memoria puntata da `HL`.

---

## Come funzionano nell'8008

Molte istruzioni codificano registri con tre bit:

| Codice | Registro |
|--------|----------|
| `000` | `A` |
| `001` | `B` |
| `010` | `C` |
| `011` | `D` |
| `100` | `E` |
| `101` | `H` |
| `110` | `L` |
| `111` | `M` |

Quando viene usato `M`, l'indirizzo e' formato da `H:L`. Solo i 6 bit bassi di
`H` partecipano all'indirizzo; i bit 6 e 7 sono ignorati.

---

## Come sono modellati nel progetto

`CPU8008` contiene campi pubblici `A`, `B`, `C`, `D`, `E`, `H`, `L`. Il metodo
`HL()` restituisce l'indirizzo a 14 bit puntato da `H` e `L`.

Le costanti `RegA` ... `RegM` rappresentano i codici registro usati dal decoder
e dagli helper opcode.

---

## Implementato ora

- Registri fisici a 8 bit.
- Costanti registro.
- Calcolo `HL()` con mascheramento di `H & 0x3F`.
- Load e move tra registri.
- Load immediati.
- Lettura e scrittura del pseudo-registro `M`.
- ALU e flag sui registri.

---

## Da implementare

- Nessuna semantica registro di base nota da aggiungere in questa fase.
