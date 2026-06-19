# Profili macchina

I profili macchina descrivono configurazioni 8008 ad alto livello sopra il core
`cpu`: indirizzo di caricamento predefinito, PC iniziale, limite di esecuzione e
slot ROM caricabili, regioni memoria e porte I/O convenzionali.

Il repository non include ROM storiche. I profili storici sono scheletri
eseguibili ma conservativi: accettano file locali forniti esplicitamente
dall'utente e non pretendono ancora di riprodurre ogni scheda elettronica.

---

## Che cosa sono SCELBI e Intellec

### SCELBI

SCELBI Computer Consulting fu una delle prime aziende a vendere un
microcomputer basato su Intel 8008. Il modello SCELBI 8H, commercializzato nel
1974 come kit o macchina assemblata, era pensato per sperimentatori e utenti
tecnici. Il sistema base era fortemente legato al front panel; terminale,
teleprinter, cassette e altre periferiche erano espansioni.

Lo SCELBI 8B fu una revisione successiva. Oggi e' particolarmente interessante
per l'emulazione perche' esistono software e riferimenti compatibili, tra cui
SCELBI Monitor, Editor, Assembler, SCELBAL e immagini usate da simulatori come
SIMH.

### Intel Intellec

Intellec era la famiglia Intel di sistemi di sviluppo per i suoi primi
microprocessori. Intellec 8 era rivolto all'Intel 8008/MCS-8: serviva agli
ingegneri per sviluppare, provare e trasferire firmware verso PROM/EPROM e
sistemi embedded.

Non va interpretato come un normale personal computer domestico. Era soprattutto
un banco di sviluppo con front panel, memoria e interfacce opzionali, utile per
costruire altre macchine basate sull'8008.

---

## Profili disponibili

- `generic`: macchina piatta didattica da 16 KB.
- `intellec-8`: profilo iniziale per Intel Intellec 8/MCS-8.
- `scelbi-8b`: profilo iniziale per software e confronti SCELBI/SIMH.
- `scelbi-8h`: profilo iniziale per la macchina SCELBI 8H.

I profili storici espongono due slot alternativi a `0x0000`:

- `monitor`: monitor/bootstrap locale.
- `test`: ROM locale di smoke test.

Gli slot accettano al massimo i 16 KB indirizzabili dall'8008. Sono alternativi:
caricarli entrambi sovrappone i byte a partire da `0x0000`.

---

## Memoria e I/O del profilo

La regione `0x0000-0x3FFF` e' descritta come `mixed`: puo' contenere ROM e RAM,
ma la separazione non viene ancora fatta rispettare dal bus memoria.

Per SCELBI e Intellec sono documentate due porte convenzionali:

- input `0`: `callback-input-0`
- output `8`: `callback-output-8`

Queste porte sono una convenzione dell'emulatore per test e terminali virtuali,
non una dichiarazione di fedelta' storica. Il campo `Historical=false` nel
profilo rende esplicita questa distinzione. Le mappe storiche definitive
verranno aggiunte solo dopo verifica di schemi, monitor e periferiche.

---

## CLI

Per elencare i profili:

```bash
go run ./cmd/retronet-8008 -profiles
```

Per caricare una ROM locale nello slot `monitor` del profilo Intellec:

```bash
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -steps 1000
```

Per vedere anche le operazioni I/O:

```bash
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -io-trace
```

`-rom` e' ripetibile e usa il formato `nome=percorso`. E' possibile combinare
ROM di profilo e un binario raw:

```bash
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -bin programma.bin -addr 0x0100 -pc 0x0100
```

In questo caso la ROM viene caricata prima e il binario raw dopo. Se le regioni
si sovrappongono, l'ultimo caricamento vince.

---

## ROM locale di smoke test

La ROM minima usata dai test di integrazione contiene:

```text
41 51 00
```

Le istruzioni sono:

```text
INP 0
OUT 8
HLT
```

Su PowerShell si puo' creare un vero file binario locale ed eseguirlo:

```powershell
[IO.File]::WriteAllBytes("$env:TEMP\io-smoke.bin", [byte[]](0x41, 0x51, 0x00))
go run ./cmd/retronet-8008 -profile scelbi-8b -rom "test=$env:TEMP\io-smoke.bin" -input 0=0x5A -io-trace -steps 8
```

Output I/O atteso:

```text
io in port=0 value=0x5A
io out port=8 value=0x5A
```

Questa e' una ROM di test reale come file locale, ma non e' una ROM storica.
Serve a verificare insieme loader, profilo, istruzioni I/O e callback.

---

## API

Il package `machine` espone:

- `Profiles()`: elenco ordinato dei profili disponibili.
- `Lookup(name)`: ricerca di un profilo per nome.
- `MemoryRegion`, `IOPort` e `ROMHint`: metadata di macchina leggibili.
- `Profile.LoadROM(mem, name, data)`: caricamento di una ROM nello slot del
  profilo.
- `LoadBytes(mem, addr, data)`: caricamento raw con validazione 14 bit.
- `ValidateRange(addr, size)`: controllo di range senza wrap silenzioso.
- `CallbackIO`: bus I/O con latch e callback per porta.
- `Profile.NewIO()`: crea il bus callback associato al profilo.
- `OnInput` e `OnOutput`: collegamento di periferiche o trace esterni.

I profili restituiti sono copie profonde delle slice, quindi modifiche locali a
slot, regioni, porte o suggerimenti ROM non alterano il catalogo globale.

---

## Limiti dichiarati

- Nessuna ROM storica e' inclusa nel repository.
- `FlatMemory` non protegge ancora le regioni ROM dalle scritture.
- Le porte callback `0` e `8` sono convenzioni di test, non mappe storiche.
- Front panel, terminale, cassette, PROM programmer e bank switching non sono
  ancora periferiche emulate.
- Timing, READY, INTERRUPT e bus cycle reali restano fuori da questa milestone.

---

## Riferimenti

- Intel, *MCS-8 / 8008 User Manual*, conservato anche nei riferimenti locali di
  sviluppo.
- [SCELBI Computer Museum](https://www.scelbi.com/).
- [Intel Intellec](https://en.wikipedia.org/wiki/Intellec), panoramica storica
  della famiglia di sistemi di sviluppo.
