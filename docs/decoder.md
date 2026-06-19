# Decoder e Step

Il decoder trasforma un byte opcode in metadata eseguibili dal core. La tabella
e il ciclo fetch-decode-execute sono completi; le famiglie gia' implementate
agganciano una funzione esecutiva reale, mentre gli altri opcode terminano con
un errore esplicito di istruzione non implementata.

---

## Cosa rappresenta

Il decoder e' la mappa completa dei 256 opcode possibili dell'Intel 8008. Ogni
voce contiene:

- byte opcode
- mnemonico leggibile
- lunghezza istruzione, 1, 2 o 3 byte
- metadata di stati minimo
- funzione esecutiva

---

## Come funziona nell'8008

Il primo byte di ogni istruzione e' sempre l'opcode. Alcune istruzioni leggono
uno o due byte successivi:

- 1 byte: register move, ALU register, rotate, return, reset, I/O
- 2 byte: load immediati e ALU immediata
- 3 byte: jump e call

Il program counter resta a 14 bit e quindi avanza con wrap su `0x3FFF`.

---

## Come e' modellato nel progetto

Il package `cpu` espone:

- `Opcode`, metadata di una voce decoder
- `Instruction`, opcode fetchato con operandi
- `Decode(op byte) Opcode`
- `OpcodeTable() [256]Opcode`
- `Step(mem Memory, io IO) error`

`Step` legge l'opcode da `PC`, incrementa `PC`, legge gli eventuali operandi
secondo `Opcode.Length`, poi chiama la funzione esecutiva associata. Gli opcode
non ancora implementati restituiscono `ErrUnimplementedOpcode`.

---

## Implementato ora

- Tabella decoder completa da 256 opcode.
- Metadata di lunghezza per istruzioni da 1, 2 e 3 byte.
- Mnemonici di base per le famiglie note.
- `Step` con fetch opcode e operandi.
- Esecuzione reale delle istruzioni load/move e ALU.
- Wrap del `PC` a 14 bit.
- Errore tipizzato `UnimplementedOpcodeError`.
- Test su tabella, lunghezze, mnemonici, fetch e PC.

---

## Da implementare

- Funzioni esecutive reali per rotate, control flow, HLT e I/O.
- Disassembler con contesto memoria.
- Trace istruzione per istruzione.
- Metadata temporali piu' accurati.
