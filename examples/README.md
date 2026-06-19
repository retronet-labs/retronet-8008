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

Esempi Go veri e propri arriveranno con assembler/disassembler e profili
macchina.
