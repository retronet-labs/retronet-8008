# Conformance sintetica

Il package `conformance` verifica il core con piccoli programmi costruiti dagli
helper opcode del progetto. Non usa ROM storiche, mappe SCELBI/Intellec o file
esterni.

---

## Runner

Ogni `Case` contiene:

- nome e byte programma
- indirizzi di caricamento e avvio
- limite istruzioni
- setup opzionale dei componenti
- verifica finale

`RunCase` crea CPU, memoria, I/O, front panel e debugger nuovi. Un fallimento non
interrompe `RunSuite`: `SuiteResult` conserva tutti gli esiti con passi, motivo
di stop ed errore.

La CLI esegue la suite integrata senza richiedere `-bin` o `-rom`:

```powershell
go run ./cmd/retronet-8008 -conformance
```

---

## Casi integrati

`SyntheticSuite` copre:

1. load e move
2. ALU con carry, zero, sign e parity
3. memoria indiretta tramite `HL/M`
4. call e return con stack interno
5. salto condizionale preso con timing
6. salto condizionale non preso con timing
7. rotate e carry
8. echo I/O input `0`/output `8`
9. wrap dello stack interno dopo otto restart
10. interrupt jammed con `RST`
11. READY basso e stato WAIT

La suite non sostituisce un riferimento indipendente cycle-perfect, ma protegge
le invarianti del progetto senza dipendere dalle milestone storiche rinviate.

---

## ROM locali opzionali

`VerifyLocalROM` calcola dimensione e SHA-256 in streaming. `ROMExpectation`
puo' controllare entrambi o ignorarne uno (`ExpectedSize=-1`, hash vuoto).

La CLI espone lo stesso controllo:

```powershell
go run ./cmd/retronet-8008 -verify-rom monitor.bin -rom-size 2048 -rom-sha256 HASH
```

Il comando non carica e non esegue la ROM. Un match SHA-256 dimostra soltanto
che il file e' quello atteso: non dimostra provenienza, autenticita' storica o
diritto di redistribuzione. Licenza e fonte devono essere documentate
separatamente prima di aggiungere qualsiasi binario al repository.
