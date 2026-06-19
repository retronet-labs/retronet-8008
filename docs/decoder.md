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

- 1 byte: register move, ALU register, rotate, return, reset, halt, I/O
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
- `Jam(mem Memory, io IO, code byte, operands ...byte) error`

`Step` legge l'opcode da `PC`, incrementa `PC`, legge gli eventuali operandi
secondo `Opcode.Length`, poi chiama la funzione esecutiva associata. Gli opcode
non ancora implementati restituiscono `ErrUnimplementedOpcode`.

Se `Halted` o `Stopped` sono veri, `Step` non accede alla memoria e restituisce
`ErrCPUStopped`. `Jam` modella l'istruzione forzata da un interrupt esterno:
valida il numero di operandi, porta la CPU in stato running ed esegue
l'istruzione senza fetch da memoria.

---

## Implementato ora

- Tabella decoder completa da 256 opcode.
- Metadata di lunghezza per istruzioni da 1, 2 e 3 byte.
- Mnemonici di base per le famiglie note.
- `Step` con fetch opcode e operandi.
- Esecuzione reale delle istruzioni load/move, ALU, rotate, control flow, HLT e I/O.
- `Jam` per eseguire una istruzione esterna da stato stopped.
- Wrap del `PC` a 14 bit.
- Errore `ErrCPUStopped` quando `Step` viene chiamato a CPU ferma.
- Errore `ErrNilIO` quando una istruzione I/O non riceve un bus I/O.
- Errore tipizzato `UnimplementedOpcodeError`.
- Test su tabella, lunghezze, mnemonici, fetch, PC e mapping porte.

---

## Da implementare

- Disassembler con contesto memoria.
- Trace istruzione per istruzione.
- Metadata temporali piu' accurati.
