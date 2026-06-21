# Memoria

L'Intel 8008 indirizza direttamente 16 KB di memoria, usando indirizzi a 14 bit.
Lo spazio memoria e' separato dallo spazio I/O.

---

## Cosa rappresenta

La memoria contiene istruzioni e dati visibili alla CPU. Il progetto offre sia
una memoria piatta per il core e i test, sia un bus macchina che distingue RAM,
ROM e intervalli non mappati.

---

## Come funziona nell'8008

Il program counter e gli indirizzi diretti sono larghi 14 bit. Ogni indirizzo
fuori range viene ricondotto ai 14 bit bassi:

- `0x4000` diventa `0x0000`
- `0xFFFF` diventa `0x3FFF`

Il pseudo-registro `M` usa la memoria puntata da `HL`, con `H` limitato ai
suoi 6 bit bassi.

---

## Come e' modellata nel progetto

Il package `cpu` espone:

- `Memory`, interfaccia con `Read(addr uint16) byte` e `Write(addr uint16, value byte)`
- `FlatMemory`, implementazione semplice da 16 KB
- `NewFlatMemory()`, costruttore con memoria inizializzata a zero
- `AddressSpaceSize`, dimensione dello spazio diretto

`FlatMemory` maschera sempre gli indirizzi con `AddressMask`.

Il package `machine` aggiunge:

- `MemoryBus`, bus a regioni che implementa `cpu.Memory`
- `NewMemoryBus(regions)`, con validazione di limiti e sovrapposizioni
- `Profile.NewMemory()`, che costruisce il bus previsto da un profilo
- `MemoryBus.LoadROM()`, caricamento privilegiato che protegge i byte caricati
- `MemoryBus.LoadBytes()`, caricamento raw che rispetta la protezione ROM

Una lettura da memoria non mappata restituisce `0xFF`. Una scrittura CPU verso
ROM o memoria non mappata viene ignorata, perche' l'interfaccia `cpu.Memory` non
prevede un valore di errore.

Le regioni `mixed` dei profili storici partono scrivibili. Quando viene caricata
una ROM locale, solo l'intervallo realmente occupato diventa ROM. Questa scelta
permette una macchina utilizzabile senza spacciare per storica una ripartizione
ROM/RAM ancora da verificare.

---

## Implementato ora

- Lettura e scrittura byte.
- Memoria piatta da 16 KB.
- Mascheramento degli indirizzi a 14 bit.
- Accesso al pseudo-registro `M` per load/move tramite `HL`.
- Uso di `M` come operando sorgente per la famiglia ALU.
- Caricamento di binari raw dalla CLI.
- Caricamento di ROM locali tramite slot di profilo.
- Bus `MemoryBus` che applica regioni RAM, ROM e mixed.
- Protezione in scrittura degli intervalli ROM caricati.
- Rifiuto di mappe sovrapposte e loader raw sopra ROM.
- Open bus convenzionale a `0xFF` per indirizzi non mappati.
- Test su zero init, lettura, scrittura e wrap degli indirizzi.

---

## Limiti e sviluppi

- Mappe ROM/RAM storiche verificate per SCELBI e Intellec.
- Eventuale bank switching a livello macchina.
