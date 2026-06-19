package cpu

import "fmt"

// Jam simula una jam instruction fornita dall'esterno, tipicamente in risposta
// a un interrupt. L'istruzione viene eseguita senza fetch da memoria e riporta
// la CPU in stato running prima del dispatch.
func (c *CPU8008) Jam(mem Memory, io IO, code byte, operands ...byte) error {
	op := Decode(code)
	wantOperands := int(op.Length - 1)
	if len(operands) != wantOperands {
		return fmt.Errorf("%w: opcode=0x%02X operands=%d want=%d", ErrInvalidJamInstruction, code, len(operands), wantOperands)
	}

	inst := Instruction{
		PC:           c.PC,
		Opcode:       op,
		OperandCount: byte(wantOperands),
	}
	copy(inst.Operands[:], operands)

	c.Halted = false
	c.Stopped = false
	return op.Execute(c, mem, io, inst)
}
