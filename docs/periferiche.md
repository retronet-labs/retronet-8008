# Periferiche configurabili

`PeripheralBus` assegna porte di `CallbackIO` a periferiche emulative senza
dipendere da mappe SCELBI o Intellec. Le callback di trace e debugger restano
osservatori e non possiedono le porte.

---

## Binding e ownership

Una `PeripheralBinding` contiene nome, callback input e callback output.
`Attach` valida prima l'intero binding:

- nome non vuoto
- porte nei range 8008 (`0..7` input, `8..31` output)
- callback non nil
- nessun duplicato interno
- nessun conflitto con periferiche o callback esterne

Solo dopo una validazione completa modifica il bus. Un errore non lascia quindi
binding parziali. `ErrPortInUse` include direzione, porta e proprietario.

`Detach(name)` libera tutte le porte della periferica ma non azzera i latch.
`Bindings()` restituisce una fotografia ordinata utile per UI e diagnostica.

Le periferiche vanno collegate prima di avviare `FrontPanel.Run`; il manager non
promette hot-plug concorrente durante una istruzione.

---

## Terminale configurabile

`Terminal.AttachPeripheral` usa `TerminalConfig`:

```go
terminal.AttachPeripheral(bus, "console", machine.TerminalConfig{
	InputPort:  2,
	OutputPort: 10,
})
```

`Terminal.Attach` conserva la convenzione predefinita input `0`, output `8`.
La CLI permette di scegliere le porte:

```powershell
go run ./cmd/retronet-8008 -bin programma.bin -terminal-input Q -terminal-in-port 2 -terminal-out-port 10
```

Le porte scelte sono configurazione emulativa, non una dichiarazione storica.

---

## Registro loopback

`RegisterPeripheral` e' un registro a 8 bit: una scrittura output aggiorna il
valore e una lettura input lo restituisce. Serve per test, handshake sintetici e
prototipi di periferiche.

La CLI accetta binding ripetibili `input=output`:

```powershell
go run ./cmd/retronet-8008 -bin programma.bin -loopback 1=9
```

Piu' loopback possono convivere se non condividono porte. Un conflitto con il
terminale o con un altro device termina con errore prima dell'esecuzione.

---

## Confine storico

Questa milestone non implementa cassette, UART, PROM programmer o controller
specifici. Tali periferiche dipendono da schemi, protocolli e mappe I/O
verificate; aggiungerle ora con porte arbitrarie le renderebbe fuorvianti.
