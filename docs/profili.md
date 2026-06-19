# Profili macchina

I profili macchina descrivono configurazioni 8008 ad alto livello sopra il core
`cpu`: indirizzo di caricamento predefinito, PC iniziale, limite di esecuzione e
slot ROM caricabili.

Il repository non include ROM storiche. I profili storici sono scheletri
documentali e accettano solo file locali forniti esplicitamente dall'utente.

---

## Profili disponibili

- `generic`: macchina piatta generica da 16 KB, senza slot ROM predefiniti.
- `intellec-8`: scheletro per sistemi Intel Intellec 8/MOD 8, con slot
  opzionale `monitor` a `0x0000`.
- `scelbi-8b`: scheletro per profilo SCELBI 8B compatibile con riferimenti
  SIMH, con slot opzionale `monitor` a `0x0000`.
- `scelbi-8h`: scheletro per sistemi SCELBI 8H, con slot opzionale `monitor` a
  `0x0000`.

Gli slot attuali usano come limite massimo lo spazio indirizzabile completo
dell'8008, cioe' 16 KB. Range piu' stretti verranno aggiunti solo quando il
progetto modellera' mappe memoria storiche piu' precise.

Tra i riferimenti software da provare in milestone successive ci sono SCELBI
Monitor, Editor, Assembler, SCELBAL, `forth-scelbi.bin` e le ROM di test usate
da altri simulatori 8008. Questi nomi non sono ancora slot distinti perche'
servono prima mappe memoria e convenzioni I/O verificate.

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

`-rom` e' ripetibile e usa il formato `nome=percorso`. E' possibile combinare
ROM di profilo e un binario raw:

```bash
go run ./cmd/retronet-8008 -profile intellec-8 -rom monitor=monitor.bin -bin programma.bin -addr 0x0100 -pc 0x0100
```

In questo caso la ROM viene caricata prima e il binario raw dopo. Se le regioni
si sovrappongono, l'ultimo caricamento vince.

---

## API

Il package `machine` espone:

- `Profiles()`: elenco ordinato dei profili disponibili.
- `Lookup(name)`: ricerca di un profilo per nome.
- `Profile.LoadROM(mem, name, data)`: caricamento di una ROM nello slot del
  profilo.
- `LoadBytes(mem, addr, data)`: caricamento raw con validazione 14 bit.
- `ValidateRange(addr, size)`: controllo di range senza wrap silenzioso.

I profili restituiti sono copie, quindi modifiche locali a `ROMSlots` non
alterano il catalogo globale.
