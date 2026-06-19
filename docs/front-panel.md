# Front panel

`machine.FrontPanel` modella i controlli esterni necessari per guidare una
macchina 8008: selettori dati e indirizzo, lettura/deposito memoria, esecuzione,
arresto richiesto e jam instruction.

Non e' ancora la riproduzione estetica o elettrica di uno specifico pannello
SCELBI o Intellec. E' un coordinatore testabile che mantiene queste funzioni
fuori dal core CPU.

---

## Stato e controlli

Il pannello espone:

- otto switch dati tramite `SetSwitches` e `Switches`
- quattordici switch indirizzo tramite `SetAddress` e `Address`
- `Examine` per leggere il byte selezionato
- `Deposit` e `DepositSwitches` per scrivere memoria
- `Step` per una singola istruzione
- `Run` per eseguire fino a stop CPU, richiesta esterna o limite
- `Stop` per richiedere l'arresto del loop in modo concorrente
- `Reset` per applicare il reset storico
- `Jam` per forzare una istruzione esterna
- `InterruptRST(0..7)` per una jam vettorizzata
- `Snapshot` per fotografare CPU, switch, indirizzo e data bus

`Deposit` usa il normale bus: una scrittura verso ROM resta bloccata da
`MemoryBus`. Gli indirizzi del pannello vengono mascherati a 14 bit.

---

## Run e stop

`Run` restituisce un `PanelRunResult` con il numero di istruzioni e uno dei
motivi:

- `cpu-stopped`: la CPU ha eseguito `HLT` o era in stato stopped
- `requested`: `Stop()` ha richiesto l'arresto esterno
- `limit`: e' stato raggiunto il massimo di istruzioni

Una richiesta del pannello non modifica artificialmente `CPU.Halted` o
`CPU.Stopped`. Questi campi continuano a rappresentare soltanto lo stato del
processore. `PanelStepObserver` riceve una copia dello stato prima di ogni
istruzione ed e' usato dalla CLI per il trace.

---

## Jam e interrupt

L'8008 storico viene portato fuori dallo stato fermo mediante logica esterna che
forza una istruzione. `FrontPanel.Jam` delega a `CPU8008.Jam` mantenendo questo
modello.

`InterruptRST(n)` forza `RST n`: incrementa lo stack pointer interno, conserva
il PC corrente come ritorno e salta a `n*8`. Sono accettati solo i vettori
`0..7`.

---

## Switch come input

`AttachSwitches(ioBus, port)` collega dinamicamente gli switch dati a una porta
input. La CLI offre una forma deterministica:

```powershell
go run ./cmd/retronet-8008 -bin programma.bin -panel-switches 0x4B -panel-address 0x0000 -steps 100
```

`-panel-switches` alimenta il latch della porta convenzionale input `0` e
abilita la stampa del pannello. Se e' attivo anche il terminale, i byte accodati
hanno precedenza; a coda vuota il terminale torna al valore degli switch.

`-panel-address` seleziona il byte mostrato da `Snapshot`; il default e' il PC
iniziale. `-panel` stampa lo stato senza cambiare switch o indirizzo.

---

## Limiti

- Non esiste ancora una UI interattiva o una riproduzione grafica del pannello.
- Gli switch CLI sulla porta `0` sono una convenzione emulativa.
- `Stop` opera tra due istruzioni, non a livello di T-state.
- READY, interrupt elettrico e cicli bus saranno aggiunti con il timing.
