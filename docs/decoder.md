# Decoder e Step

Il decoder trasforma un byte opcode in metadata eseguibili dal core. La tabella
e il ciclo fetch-decode-execute sono completi: 250 encoding definiti agganciano
una funzione esecutiva, mentre i sei slot non definiti dall'ISA terminano con
un errore esplicito.

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
- `Disassemble(mem Memory, pc uint16) (Disassembly, error)`
- `Step(mem Memory, io IO) error`
- `Jam(mem Memory, io IO, code byte, operands ...byte) error`

`Step` legge l'opcode da `PC`, incrementa `PC`, legge gli eventuali operandi
secondo `Opcode.Length`, poi chiama la funzione esecutiva associata. Gli
encoding non definiti `22`, `2A`, `32`, `38`, `39` e `3A` restituiscono
`ErrUnimplementedOpcode`.

Se `Halted` o `Stopped` sono veri, `Step` non accede alla memoria e restituisce
`ErrCPUStopped`. `Jam` modella l'istruzione forzata da un interrupt esterno:
valida il numero di operandi, porta la CPU in stato running ed esegue
l'istruzione senza fetch da memoria.

`Disassemble` usa gli stessi metadata del decoder, ma legge solo i byte da
memoria e restituisce una rappresentazione testuale senza modificare stato CPU.

---

## Implementato ora

- Tabella decoder completa da 256 opcode.
- Metadata di lunghezza per istruzioni da 1, 2 e 3 byte.
- Mnemonici di base per le famiglie note.
- Disassembler con contesto memoria, bytes e `NextPC`.
- Trace CLI basato su disassembly del PC corrente prima di ogni `Step`.
- `Step` con fetch opcode e operandi.
- Esecuzione reale delle istruzioni load/move, ALU, rotate, control flow, HLT e I/O.
- `Jam` per eseguire una istruzione esterna da stato stopped.
- Wrap del `PC` a 14 bit.
- Errore `ErrCPUStopped` quando `Step` viene chiamato a CPU ferma.
- Errore `ErrNilIO` quando una istruzione I/O non riceve un bus I/O.
- Errore tipizzato `UnimplementedOpcodeError`.
- Test su tabella, lunghezze, mnemonici, fetch, PC e mapping porte.
- Matrice di esecuzione e metadata per tutti i 256 byte opcode.
- Fuzz test su decode, disassembly, fetch e limiti architetturali.

---

## Limiti

- La tabella descrive cicli macchina e costi in stati, non le transizioni dei
  singoli pin/T-state.
- Il confronto corrente usa oracle interni test-only; manca ancora un confronto
  differenziale con una seconda implementazione indipendente.
