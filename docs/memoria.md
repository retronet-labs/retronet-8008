# Memoria

L'Intel 8008 indirizza direttamente 16 KB di memoria, usando indirizzi a 14 bit.
Lo spazio memoria e' separato dallo spazio I/O.

---

## Cosa rappresenta

La memoria contiene istruzioni e dati visibili alla CPU. In questa fase il
progetto modella una memoria piatta da `0x0000` a `0x3FFF`, senza distinguere
ancora ROM e RAM.

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

---

## Implementato ora

- Lettura e scrittura byte.
- Memoria piatta da 16 KB.
- Mascheramento degli indirizzi a 14 bit.
- Accesso al pseudo-registro `M` per load/move tramite `HL`.
- Uso di `M` come operando sorgente per la famiglia ALU.
- Caricamento di binari raw dalla CLI.
- Caricamento di ROM locali tramite slot di profilo.
- Metadata `MemoryRegion` per descrivere le mappe macchina.
- Test su zero init, lettura, scrittura e wrap degli indirizzi.

---

## Da implementare

- Bus memoria che faccia rispettare regioni ROM in sola lettura.
- Mappe ROM/RAM storiche verificate per SCELBI e Intellec.
- Eventuale bank switching a livello macchina.
