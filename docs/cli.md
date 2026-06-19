# CLI

La CLI `retronet-8008` e' un runner minimale per programmi binari raw 8008.
Carica i byte in `FlatMemory`, avvia la CPU tramite una jam di `JMP` al PC
iniziale, esegue un numero massimo di istruzioni e stampa un dump registri.

---

## Uso

```bash
go run ./cmd/retronet-8008 -bin programma.bin -steps 1000
```

Opzioni:

- `-bin`: percorso del binario raw da caricare. Obbligatorio.
- `-addr`: indirizzo di caricamento, decimale o `0xHEX`. Default `0x0000`.
- `-pc`: program counter iniziale, decimale o `0xHEX`. Default uguale ad
  `-addr`.
- `-steps`: numero massimo di istruzioni da eseguire. Default `1000`.

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
loaded=3 addr=0x0000 pc_start=0x0000 steps=2 limit_reached=false
A=0x2A B=0x00 C=0x00 D=0x00 E=0x00 H=0x00 L=0x00
PC=0x0003 SP=0 Halted=true Stopped=true
Flags C=false Z=false S=false P=false
Stack=[0x0003 0x0000 0x0000 0x0000 0x0000 0x0000 0x0000 0x0000]
```

---

## Limiti

- Il formato supportato e' solo binario raw.
- Non c'e' ancora disassembler nel dump.
- Le porte I/O usano `Ports`, l'implementazione semplice interna.
- ROM storiche e profili macchina arriveranno in milestone successive.
