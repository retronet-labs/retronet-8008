# Contesto operativo per agenti

Questo file conserva le decisioni di progetto che servono per continuare lo
sviluppo di `retronet-8008` senza ricostruire ogni volta il contesto storico e
tecnico.

## Obiettivo

Implementare in Go un emulatore Intel 8008 didattico, testato e importabile,
con struttura coerente con `go-4004`. Il core deve restare utilizzabile senza
dipendere da una macchina storica specifica.

## Stato

Sono completate le milestone 0-15:

- core CPU, decoder e famiglie istruzionali 8008
- memoria e I/O separati
- control flow, stack interno, halt, stopped e jam instruction
- disassembler, trace e CLI
- profili `generic`, `intellec-8`, `scelbi-8h` e `scelbi-8b`
- callback I/O e smoke ROM locale
- bus memoria mappato con ROM protetta

La roadmap dettagliata e' in `docs/roadmap.md`.

## Architettura

- `cpu/`: core indipendente; espone `cpu.Memory` e `cpu.IO` come interfacce.
- `machine/`: profili, `MemoryBus`, loader ROM/raw e callback I/O.
- `cmd/retronet-8008/`: runner CLI e integrazione dei componenti.
- `docs/`: documentazione in italiano.
- `docs_do_not_commit/`: riferimenti locali ignorati da git; non committare.

Il package `cpu` non deve importare `machine`. I test Go vivono accanto al
package testato, senza una directory `tests/` separata.

## Decisioni da preservare

- Tutti gli indirizzi CPU sono mascherati a 14 bit (`0x0000-0x3FFF`).
- Il reset storico lascia la CPU in stato `Stopped/Halted`.
- L'opcode `0xFF` e' un alias di `HLT`, non un move `M,M`.
- `MemoryBus` restituisce `0xFF` per memoria non mappata.
- Le scritture CPU in ROM o memoria non mappata vengono ignorate.
- `MemoryKindMixed` resta scrivibile; `LoadROM` rende read-only soltanto i byte
  effettivamente caricati. Questo evita di inventare mappe storiche.
- Un binario raw non puo' sovrascrivere una ROM gia' caricata.
- Le porte callback input `0` e output `8` sono convenzioni dell'emulatore, non
  mappe SCELBI/Intellec storicamente verificate.
- Non aggiungere ROM storiche senza provenienza e licenza documentate.

## Verifica

Comando normale:

```powershell
go test -count=1 ./...
```

Se la cache Go globale non e' scrivibile nell'ambiente Codex:

```powershell
$env:GOCACHE='C:\work\source\retronet-8008\.gocache'
go test -count=1 ./...
```

Prima di un commit eseguire anche `gofmt` sui file Go modificati e
`git diff --check`. Mantenere commit piccoli e tematici.

## Prossimi passi

Ordine consigliato:

1. terminale virtuale collegato a `CallbackIO`
2. front panel con step, run, stop e jam/interrupt
3. mappe memoria e I/O storiche verificate
4. cassette e altre periferiche SCELBI/Intellec
5. timing e T-state
