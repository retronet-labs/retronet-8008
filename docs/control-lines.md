# READY e interrupt

Il front panel modella le due linee di controllo principali dell'8008 senza
dipendere da una macchina SCELBI o Intellec. Il modello opera al confine dei
cicli macchina dichiarati dal decoder; gli effetti dell'istruzione restano
atomici al termine della sequenza.

---

## READY e WAIT

`FrontPanel.SetReady(false)` impedisce al ciclo corrente di raggiungere `T3`.
`Step` restituisce `ErrCPUWaiting`, conserva PC e stato funzionale e registra
uno stato WAIT. Ogni nuovo tentativo con READY ancora basso aggiunge un altro
WAIT; quando READY torna alto, l'istruzione viene completata.

`SetReadyCallback` permette una decisione per ciclo tramite `CycleContext`:

- PC e opcode dell'istruzione
- indice del ciclo
- tipo `PCI`, `PCR`, `PCW` o `PCC`
- indicazione di interrupt acknowledge

Il pannello ricorda il ciclo sul quale si e' fermato. Un callback puo' quindi
lasciare passare `PCI` e bloccare, per esempio, il `PCC` di `OUT`.

`CPU8008.WaitStateCount` conta tutti i WAIT. `StateCount` include sia gli stati
base sia i WAIT; `LastTiming.WaitStates` associa l'attesa all'istruzione che la
completa.

---

## Interrupt

`RequestInterrupt(code, operands...)` accoda una jam instruction. Il pannello la
riconosce al prossimo confine `PCI`, usa un ciclo `T11` equivalente al primo
stato PCI e non avanza il program counter prima dell'istruzione forzata.

Questo permette di interrompere anche una CPU in `Stopped/Halted`. Per esempio
`RequestInterrupt(cpu.RST(3))` salva il PC corrente nello stack interno e salta
a `0x0018`. Una seconda richiesta viene rifiutata con `ErrInterruptPending`
finche' la prima non e' stata servita.

`InterruptRST` resta il comando immediato del front panel; `RequestInterrupt` e'
la forma sincronizzata da usare per simulare la linea hardware.

---

## Run e CLI

`FrontPanel.Run` restituisce il motivo `waiting` quando incontra READY basso,
senza consumare una istruzione. La CLI espone:

```powershell
go run ./cmd/retronet-8008 -bin programma.bin -ready=false
go run ./cmd/retronet-8008 -bin programma.bin -interrupt-rst 3
```

`-ready=false` e' utile per testare lo stato WAIT; non esiste ancora un input
interattivo che rialzi la linea durante lo stesso processo CLI.

`-interrupt-rst N` accoda `RST N` dopo la jam iniziale di avvio e prima del
primo fetch del programma. Il dump include `stop_reason`, conteggio WAIT e stato
del pannello.

---

## Limite dichiarato

Il modello distingue i cicli macchina e il punto di attesa, ma non emette ancora
ogni transizione di pin `S0/S1/S2/SYNC` o i due clock di ciascun T-state. Letture,
scritture e side-effect diventano visibili insieme al completamento
dell'istruzione.

La fonte primaria e' il manuale Intel 8008, sezioni *Processor Timing* e
*Processor Control Signals*.
