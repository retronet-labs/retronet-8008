# CLI

La CLI `retronet-8008` e' un runner minimale per programmi 8008. Carica byte
raw e ROM locali in `FlatMemory`, avvia la CPU tramite una jam di `JMP` al PC
iniziale, esegue un numero massimo di istruzioni e stampa un dump registri.

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

Le ROM vengono caricate prima del binario raw passato con `-bin`. Se le regioni
si sovrappongono, il caricamento successivo sovrascrive i byte precedenti.

---

## Limiti

- Il formato supportato e' binario raw; non ci sono ancora container ROM.
- Le porte I/O usano `Ports`, l'implementazione semplice interna.
- Il repository non include ROM storiche: i file devono essere forniti
  localmente dall'utente.
