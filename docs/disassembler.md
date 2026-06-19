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

La CLI usa la stessa API con l'opzione `-disasm N`:

```bash
go run ./cmd/retronet-8008 -bin programma.bin -disasm 8
```

`-disasm` carica il binario raw, parte da `-pc` o da `-addr`, stampa `N`
istruzioni e termina senza eseguire il programma.

---

## Limiti

- Non esiste ancora un formato simbolico con label.
- Il disassembler non annota stati, flag o side effect.
- Non c'e' ancora trace durante l'esecuzione; questa sara' una milestone
  separata.
