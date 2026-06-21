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

## Conformance esaustiva del core

La milestone 25 aggiunge test direttamente nel package `cpu`:

- tutti i 256 byte attraversano decoder, metadata ed esecuzione
- i 250 encoding definiti devono completare una istruzione
- `22`, `2A`, `32`, `38`, `39` e `3A` devono produrre
  `ErrUnimplementedOpcode`
- un oracle matematico separato verifica tutte le combinazioni ALU: otto gruppi,
  256 accumulatori, 256 operandi e due valori Carry
- rotate e incremento/decremento sono verificati su ogni byte possibile
- la matrice ALU controlla anche selezione registro, `M` e immediati

Gli oracle sono codice test-only e non chiamano le primitive ALU per calcolare
il risultato atteso. Questa e' una verifica differenziale interna: puo'
individuare divergenze tra core e modello di test, ma non e' una certificazione
esterna dell'ISA.

Esecuzione normale:

```powershell
go test -count=1 ./cpu
```

I fuzz target includono un corpus iniziale con tutti i 256 opcode. Nel normale
`go test` vengono eseguiti i seed; una campagna esplicita puo' essere avviata con:

```powershell
go test ./cpu -run=^$ -fuzz=FuzzDecodeDisassemble -fuzztime=30s
go test ./cpu -run=^$ -fuzz=FuzzStepMaintainsArchitecturalBounds -fuzztime=30s
```

Resta utile un futuro confronto con SIMH o un'altra implementazione indipendente,
quando ambiente, versione di riferimento e casi di scambio saranno fissati in
modo riproducibile.

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
