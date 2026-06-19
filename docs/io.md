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
- `ErrInvalidInputPort` e `ErrInvalidOutputPort`

Le funzioni dell'interfaccia restano minimali per il futuro core istruzioni; i
validatori espliciti permettono invece a test, CLI e decoder di segnalare porte
non valide quando serve.

---

## Implementato ora

- Porte input `0..7`.
- Porte output `8..31`.
- Validazione esplicita dei range.
- Lettura/scrittura su implementazione `Ports`.
- Test sui limiti validi e invalidi.

---

## Da implementare

- Istruzioni `INP` e `OUT`.
- Callback o periferiche virtuali.
- Trace I/O.
- Interazione con una futura CLI runner.
