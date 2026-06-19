# CLI

La CLI `retronet-8008` e' un runner minimale per programmi 8008. Carica byte
raw e ROM locali nel bus memoria del profilo, avvia la CPU tramite una jam di
`JMP` al PC iniziale, esegue un numero massimo di istruzioni e stampa un dump
registri.

---

## Uso

```bash
go run ./cmd/retronet-8008 -bin programma.bin -steps 1000
```

Opzioni:

- `-bin`: percorso del binario raw da caricare.
- `-addr`: indirizzo di caricamento, decimale o `0xHEX`. Default `0x0000`.
- `-pc`: program counter iniziale, decimale o `0xHEX`. Default uguale ad
  `-addr`.
- `-profile`: profilo macchina da usare. Default `generic`.
- `-profiles`: elenca i profili macchina disponibili e termina.
- `-rom`: carica una ROM di profilo nel formato `nome=percorso`. Ripetibile.
- `-input`: inizializza una porta input nel formato `porta=valore`. Ripetibile.
- `-io-trace`: stampa letture e scritture I/O effettuate tramite callback.
- `-terminal`: collega un terminale ASCII buffered alle porte `0` e `8`.
- `-terminal-input`: accoda testo al terminale e abilita `-terminal`.
- `-panel`: stampa lo snapshot del front panel dopo il run.
- `-panel-switches`: imposta gli switch dati e il latch input `0`.
- `-panel-address`: seleziona l'indirizzo esaminato dal pannello.
- `-ready`: livello READY globale; con `false` il run termina in WAIT.
- `-interrupt-rst`: forza un vettore `RST 0..7` prima del primo fetch.
- `-trace-json`: scrive eventi strutturati in un file JSON Lines.
- `-break`: breakpoint PC; ripetibile.
- `-break-opcode`: breakpoint opcode; ripetibile.
- `-watch`: watchpoint scrittura memoria; ripetibile.
- `-break-input`, `-break-output`: breakpoint I/O; ripetibili.
- `-steps`: numero massimo di istruzioni da eseguire. Default `1000`.
- `-disasm`: disassembla N istruzioni dal PC iniziale e termina senza eseguire.
- `-trace`: stampa ogni istruzione prima dell'esecuzione.

Per eseguire serve almeno un `-bin` o un `-rom`, tranne quando si usa
`-profiles`.

Gli indirizzi sono limitati allo spazio 14 bit dell'8008, quindi
`0x0000-0x3FFF`. Il loader rifiuta binari che superano la fine dello spazio
indirizzabile invece di fare wrap silenzioso.

---

## Esempio

Questo programma carica `0x2A` in `A` e poi esegue `HLT`:

```powershell
[IO.File]::WriteAllBytes("$env:TEMP\load-a.bin", [byte[]](0x06, 0x2A, 0x00))
go run ./cmd/retronet-8008 -bin "$env:TEMP\load-a.bin" -steps 8
```

Output atteso, con formato compatto:

```text
profile=generic loaded=3 roms=0 addr=0x0000 pc_start=0x0000 steps=2 limit_reached=false
A=0x2A B=0x00 C=0x00 D=0x00 E=0x00 H=0x00 L=0x00
PC=0x0003 SP=0 Halted=true Stopped=true
Flags C=false Z=false S=false P=false
Stack=[0x0003 0x0000 0x0000 0x0000 0x0000 0x0000 0x0000 0x0000]
```

---

## Disassembly

Lo stesso binario puo' essere listato senza esecuzione:

```powershell
go run ./cmd/retronet-8008 -bin "$env:TEMP\load-a.bin" -disasm 2
```

Output:

```text
0000: 06 2A    LAI #0x2A
0002: 00       HLT
```

---

## Trace

Durante l'esecuzione si puo' stampare ogni istruzione realmente eseguita:

```powershell
go run ./cmd/retronet-8008 -bin "$env:TEMP\load-a.bin" -steps 8 -trace
```

Output iniziale:

```text
trace=0 0000: 06 2A    LAI #0x2A
trace=1 0002: 00       HLT
profile=generic loaded=3 roms=0 addr=0x0000 pc_start=0x0000 steps=2 limit_reached=false
```

Il trace usa il PC corrente prima dello `Step`, quindi segue salti, call, return
e restart invece di limitarsi alla sequenza lineare in memoria.

---

## Profili e ROM locali

I profili disponibili si consultano con:

```bash
go run ./cmd/retronet-8008 -profiles
```

Un profilo storico puo' ricevere ROM locali tramite slot nominati:

```bash
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -steps 1000
```

Le ROM vengono caricate prima del binario raw passato con `-bin`. I byte ROM
diventano read-only: un binario successivo puo' usare la RAM libera, ma una
sovrapposizione viene rifiutata con un errore esplicito.

---

## I/O callback e ROM di test

Questa ROM legge input `0`, lo copia su output `8` e si ferma:

```powershell
[IO.File]::WriteAllBytes("$env:TEMP\io-smoke.bin", [byte[]](0x41, 0x51, 0x00))
go run ./cmd/retronet-8008 -profile scelbi-8b -rom "test=$env:TEMP\io-smoke.bin" -input 0=0x5A -io-trace -steps 8
```

Output rilevante:

```text
io in port=0 value=0x5A
io out port=8 value=0x5A
profile=scelbi-8b loaded=0 roms=1 addr=0x0000 pc_start=0x0000 steps=3 limit_reached=false
```

`-input` usa solo le porte `0..7`. Il trace output usa le porte `8..31`, come
richiesto dall'I/O asimmetrico dell'8008.

---

## Terminale buffered

La stessa ROM puo' ricevere il carattere `Z` dalla coda del terminale e
ristamparlo su stdout:

```powershell
go run ./cmd/retronet-8008 -profile scelbi-8b -rom "test=$env:TEMP\io-smoke.bin" -terminal-input Z -steps 8
```

`-terminal-input` implica `-terminal`. Il terminale e `-io-trace` possono essere
attivi insieme; il testo raw puo' risultare adiacente alle righe di trace.

---

## Front panel

La CLI usa sempre `machine.FrontPanel` per jam iniziale e ciclo di esecuzione.
Lo snapshot e' opzionale:

```powershell
go run ./cmd/retronet-8008 -bin programma.bin -panel -panel-switches 0x4B -panel-address 0x0100
```

Gli switch alimentano la porta input convenzionale `0`; un eventuale
`-input 0=...` precedente viene sostituito. L'indirizzo del pannello serve solo
per la lettura mostrata e non modifica il PC.

Il campo `stop_reason` distingue `cpu-stopped`, `requested`, `waiting` e
`limit`. READY basso non viene riportato come errore e non incrementa `steps`.

Breakpoint e watchpoint attivano `machine.Debugger`. Il run puo' terminare anche
con `breakpoint`, `watchpoint` o `io-breakpoint`. Il file `-trace-json` contiene
solo JSON Lines; il dump finale resta su stdout.

---

## Limiti

- Il formato supportato e' binario raw; non ci sono ancora container ROM.
- Le mappe storiche ROM/RAM non sono ancora verificate; i profili proteggono
  gli intervalli delle immagini ROM effettivamente caricate.
- Le porte I/O usano `machine.CallbackIO`; non ci sono ancora periferiche
  storiche complete.
- Il terminale e' buffered: non legge ancora dalla console durante il run.
- Il repository non include ROM storiche: i file devono essere forniti
  localmente dall'utente.
