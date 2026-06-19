# Terminale virtuale

`machine.Terminal` collega un flusso ASCII al bus I/O callback. E' pensato per
monitor, programmi didattici e test end-to-end senza introdurre dipendenze nel
core CPU.

---

## Porte

Il terminale usa le convenzioni correnti dei profili:

- input `0`: un byte viene consumato dalla coda quando la CPU esegue `INP 0`
- output `8`: il byte scritto con `OUT 8` viene inoltrato a un `io.Writer`

Queste porte non sono dichiarate come mappa storica SCELBI o Intellec. Potranno
essere rese configurabili quando saranno disponibili schemi verificati.

Quando la coda input e' vuota, il terminale restituisce il valore latched della
porta. Questo rende il comportamento deterministico e permette di combinare il
terminale con `-input 0=valore`.

---

## API Go

```go
terminal := machine.NewTerminal(output)
terminal.QueueInputString("HELLO\r")
if err := terminal.Attach(ioBus); err != nil {
	return err
}
```

API principali:

- `NewTerminal(output)`: crea la periferica; `nil` scarta l'output
- `Attach(ioBus)`: collega input `0` e output `8`
- `QueueInput(data)` e `QueueInputString(value)`: accodano byte
- `PendingInput()`: restituisce i byte non ancora consumati
- `Err()`: conserva il primo errore del writer

La coda e lo stato di errore sono protetti per consentire in futuro una sorgente
input interattiva. L'esecuzione CLI attuale resta buffered e deterministica.

---

## CLI

`-terminal` abilita il terminale in output. `-terminal-input` accoda una stringa
e abilita automaticamente il terminale:

```powershell
[IO.File]::WriteAllBytes("$env:TEMP\terminal-echo.bin", [byte[]](0x41, 0x51, 0x00))
go run ./cmd/retronet-8008 -profile scelbi-8b -rom "test=$env:TEMP\terminal-echo.bin" -terminal-input Z -steps 8
```

La ROM esegue `INP 0`, `OUT 8`, `HLT`, quindi stampa `Z` prima del dump CPU.
`-io-trace` puo' essere usato insieme al terminale: gli osservatori di trace non
sostituiscono le callback della periferica. Poiche' entrambi scrivono su stdout,
l'output ASCII e le righe di trace possono apparire adiacenti.

---

## Limiti

- Non legge ancora in modo interattivo dalla console durante l'esecuzione.
- Non modella baud rate, READY, handshake o temporizzazione seriale.
- Non applica conversioni di newline o encoding: i byte sono inoltrati tali e
  quali.
