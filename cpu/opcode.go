package cpu

// ExecuteFunc e' la funzione esecutiva associata a una voce del decoder.
type ExecuteFunc func(c *CPU8008, mem Memory, io IO, inst Instruction) error

// Opcode descrive una voce della tabella decoder 8008.
type Opcode struct {
	Code     byte
	Mnemonic string
	Length   byte
	States   byte
	Execute  ExecuteFunc
}

// Instruction contiene l'opcode fetchato e gli eventuali byte operando.
type Instruction struct {
	PC           uint16
	Opcode       Opcode
	Operands     [2]byte
	OperandCount byte
}

func unimplementedExecute(_ *CPU8008, _ Memory, _ IO, inst Instruction) error {
	return &UnimplementedOpcodeError{
		PC:       inst.PC,
		Opcode:   inst.Opcode.Code,
		Mnemonic: inst.Opcode.Mnemonic,
		Length:   inst.Opcode.Length,
	}
}
