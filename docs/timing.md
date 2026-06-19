# Timing e stati macchina

Il core registra il costo delle istruzioni secondo la tabella *Internal
Processor Operation* del manuale Intel 8008. Il modello resta
instruction-level, ma espone abbastanza metadata per READY, trace e futuri
avanzamenti per ciclo.

---

## Stati e clock

Intel descrive gli stati interni `T1`, `T2`, `T3`, `T4` e `T5`. Ogni stato
richiede due periodi del clock bifase. In questo progetto `StateCount` conta gli
stati Intel, non i singoli impulsi di clock.

Un ciclo tipico usa cinque stati, ma l'8008 salta `T4` e `T5` quando non servono.
Per questo le istruzioni hanno costi diversi anche con la stessa lunghezza.

---

## Cicli macchina

`MachineCycle` distingue i quattro codici emessi nei bit di controllo:

- `PCI`: fetch del primo byte istruzione
- `PCR`: lettura memoria, operando o byte indirizzo
- `PCW`: scrittura memoria
- `PCC`: comando I/O

Ogni `Opcode` contiene `MinStates`, `States`, `CycleCount` e una sequenza fissa
di massimo tre cicli. `States` e' il costo massimo; `MinStates` differisce per
le istruzioni condizionali.

Esempi verificati sulla tabella Intel:

| Famiglia | Stati | Cicli |
|---|---:|---|
| registro-registro | 5 | `PCI` |
| registro da `M` | 8 | `PCI/PCR` |
| `M` da registro | 7 | `PCI/PCW` |
| immediato registro | 8 | `PCI/PCR` |
| `LMI` | 9 | `PCI/PCR/PCW` |
| jump/call condizionale | 9 o 11 | `PCI/PCR/PCR` |
| return condizionale | 3 o 5 | `PCI` |
| `INP` | 8 | `PCI/PCC` |
| `OUT` | 6 | `PCI/PCC` |
| `HLT` | 4 | `PCI` |

Per jump, call e return condizionali il timing viene deciso usando i flag prima
dell'esecuzione. `InstructionTiming.Conditional` e `Taken` conservano il
risultato effettivo.

---

## Contatori CPU

`CPU8008` espone:

- `InstructionCount`: istruzioni completate
- `StateCount`: stati Intel completati
- `LastTiming`: costo e cicli dell'ultima istruzione
- `WaitStateCount`: stati WAIT aggiunti dalla logica READY

Le istruzioni forzate con `Jam` sono incluse nei contatori. La jam iniziale di
`JMP` usata dalla CLI compare quindi nel totale, mentre il campo `steps` del
runner continua a contare solo le istruzioni eseguite dal programma.

Gli errori prima del completamento dell'istruzione non incrementano i contatori.
`Reset` azzera tutto il timing.

Gli stati WAIT vengono aggiunti subito a `StateCount` e associati alla successiva
istruzione completata tramite `LastTiming.WaitStates`.

---

## Fonte

La fonte primaria e' *Intel 8008 8 Bit Parallel Central Processor Unit Users
Manual*, revisione novembre 1973/maggio 1974, pagine interne 15-17. Il manuale
locale resta in `docs_do_not_commit/` e non viene distribuito dal repository.
