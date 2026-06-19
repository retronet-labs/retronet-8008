package cpu

// AddressMask limita gli indirizzi allo spazio diretto del 8008:
// 14 bit, quindi 0x0000-0x3FFF.
const AddressMask uint16 = 0x3FFF

// Register identifica i codici registro usati dall'ISA Intel 8008.
//
// Il valore 111 non e' un registro fisico: indica il pseudo-registro M,
// cioe' il byte di memoria puntato da HL.
type Register byte

const (
	RegA Register = 0b000
	RegB Register = 0b001
	RegC Register = 0b010
	RegD Register = 0b011
	RegE Register = 0b100
	RegH Register = 0b101
	RegL Register = 0b110
	RegM Register = 0b111
)

// Condition identifica i quattro flag selezionabili dalle istruzioni
// condizionali dell'8008.
type Condition byte

const (
	CondCarry  Condition = 0b00
	CondZero   Condition = 0b01
	CondSign   Condition = 0b10
	CondParity Condition = 0b11
)
