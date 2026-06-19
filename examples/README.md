# Esempi

Questa cartella conterra' programmi dimostrativi scritti in Go che costruiscono
piccole ROM usando gli helper del package `cpu`, come avviene nel modulo
`go-4004`.

Per ora il runner CLI puo' eseguire piccoli binari raw. Un programma minimale
che carica `0x2A` in `A` e poi ferma la CPU e':

```text
06 2A 00
```

Su PowerShell:

```powershell
[IO.File]::WriteAllBytes("$env:TEMP\load-a.bin", [byte[]](0x06, 0x2A, 0x00))
go run ./cmd/retronet-8008 -bin "$env:TEMP\load-a.bin" -steps 8
```

Per vedere le istruzioni senza eseguirle:

```powershell
go run ./cmd/retronet-8008 -bin "$env:TEMP\load-a.bin" -disasm 2
```

Per eseguire e vedere il flusso reale:

```powershell
go run ./cmd/retronet-8008 -bin "$env:TEMP\load-a.bin" -steps 8 -trace
```

Per vedere i profili macchina disponibili:

```powershell
go run ./cmd/retronet-8008 -profiles
```

Per caricare una ROM locale in uno slot di profilo:

```powershell
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -steps 1000
```

Esempi Go veri e propri arriveranno con assembler e programmi dimostrativi
versionati.
