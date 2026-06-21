# I/O

L'Intel 8008 separa lo spazio I/O dalla memoria. Le porte non sono memory-mapped:
vengono raggiunte da istruzioni dedicate.

---

## Cosa rappresenta

Il bus I/O collega la CPU a periferiche esterne. L'8008 ha una mappa asimmetrica:

- 8 porte di input, numerate `0..7`
- 24 porte di output, numerate `8..31`

---

## Come funziona nell'8008

Le istruzioni di input leggono un byte da una porta di ingresso e lo portano in
`A`. Le istruzioni di output scrivono il contenuto di `A` su una porta di uscita.

Il fatto che input e output abbiano range diversi e' parte del modello storico:
non vanno trattati come 32 porte bidirezionali equivalenti.

---

## Come e' modellato nel progetto

Il package `cpu` espone:

- `IO`, interfaccia con `Input(port byte) byte` e `Output(port byte, value byte)`
- `Ports`, implementazione semplice con 8 input e 24 output
- `ValidateInputPort` e `ValidateOutputPort`
- `ErrInvalidInputPort`, `ErrInvalidOutputPort` ed `ErrNilIO`
- helper opcode `INP(port)` e `OUT(port)`

`INP` usa il formato `0100 MMM1`, quindi raggiunge le porte input `0..7`.
`OUT` usa il formato `01 RRMMM1` con `RR != 00`, quindi raggiunge le porte
output `8..31`. Il decoder conserva questa asimmetria invece di trattare lo
spazio come 32 porte bidirezionali.

Il package `machine` espone inoltre `CallbackIO`, una implementazione di `cpu.IO`
con latch e callback opzionali per ogni porta:

- `SetInput(port, value)` inizializza un input.
- `OnInput(port, callback)` calcola o osserva il valore letto.
- `OnOutput(port, callback)` osserva o inoltra il valore scritto.
- `ObserveInput(port, observer)` osserva il valore finale letto senza cambiarlo.
- `ObserveOutput(port, observer)` osserva una scrittura senza sostituire la
  callback della periferica.
- `OutputValue(port)` legge l'ultimo valore latched in uscita.
- `Profile.NewIO()` crea il bus associato al profilo macchina.

Le callback permettono di collegare terminali, front panel, cassette o semplici
trace senza introdurre dipendenze di macchina nel core CPU.

Gli osservatori sono separati dalle callback: `-io-trace` puo' quindi convivere
con una periferica che produce input o riceve output sulla stessa porta.

---

## Implementato ora

- Porte input `0..7`.
- Porte output `8..31`.
- Validazione esplicita dei range.
- Lettura/scrittura su implementazione `Ports`.
- Istruzione `INP`: legge una porta input in `A`.
- Istruzione `OUT`: scrive `A` su una porta output.
- Helper `INP` e `OUT`.
- Errore `ErrNilIO` per istruzioni I/O senza bus.
- Test sui limiti validi e invalidi.
- Test sul mapping completo delle 8 porte input e 24 porte output.
- Bus `machine.CallbackIO`.
- Callback input/output validate per porta.
- Opzioni CLI `-input porta=valore` e `-io-trace`.
- Test di integrazione con ROM locale `INP 0`, `OUT 8`, `HLT`.
- Terminale ASCII buffered sulle porte convenzionali `0` e `8`.
- Osservatori I/O componibili usati dal trace.
- Binding periferiche con ownership, conflitti e detach.
- Terminale configurabile e registro loopback generico.

---

## Limiti e sviluppi

- Periferiche storiche complete: cassette e interfacce verificate.
- Mappe I/O storiche verificate per SCELBI e Intellec.
- Eventuale trace strutturato oltre all'output testuale CLI.
