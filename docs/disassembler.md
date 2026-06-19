# Disassembler

Il package `cpu` espone un disassembler minimale basato sulla stessa tabella
decoder usata da `Step`.

---

## API

```go
d, err := cpu.Disassemble(mem, 0x0000)
```

`Disassemble` legge opcode e operandi dalla memoria, applica il wrap a 14 bit e
restituisce una struct `Disassembly` con:

- `PC`
- `Opcode`
- `Bytes`
- `Length`
- `Operand`
- `NextPC`

`String()` produce una riga compatta:

```text
0000: 06 2A    LAI #0x2A
```

Gli immediati sono mostrati come `#0xNN`. I target a 14 bit sono mostrati come
`0xNNNN`, con il byte alto mascherato come nell'esecuzione reale.

---

## CLI

La CLI usa la stessa API sia con l'opzione `-disasm N`, sia con `-trace` durante
l'esecuzione:

```bash
go run ./cmd/retronet-8008 -bin programma.bin -disasm 8
go run ./cmd/retronet-8008 -bin programma.bin -steps 1000 -trace
```

`-disasm` carica il binario raw, parte da `-pc` o da `-addr`, stampa `N`
istruzioni e termina senza eseguire il programma.

`-trace` invece disassembla il PC corrente prima di ogni `Step`, quindi segue
il flusso effettivo del programma.

---

## Limiti

- Non esiste ancora un formato simbolico con label.
- Il disassembler non annota stati, flag o side effect.
- Il trace non include ancora snapshot registri per ogni istruzione.
