package cpu

import (
	"math/bits"

	"github.com/retronet-labs/retronet-hardware/bridge/i8008"
)

func isALUImmediateOpcode(code byte) bool {
	return code&0xC7 == 0x04
}

func isALURegisterOpcode(code byte) bool {
	return code&0xC0 == 0x80
}

func isIncrementOpcode(code byte) bool {
	if code&0xC7 != 0x00 {
		return false
	}
	r := Register((code >> 3) & 0x07)
	return isIncrementRegister(r)
}

func isDecrementOpcode(code byte) bool {
	if code&0xC7 != 0x01 {
		return false
	}
	r := Register((code >> 3) & 0x07)
	return isIncrementRegister(r)
}

func isIncrementRegister(r Register) bool {
	switch Register(regBits(r)) {
	case RegB, RegC, RegD, RegE, RegH, RegL:
		return true
	default:
		return false
	}
}

func executeALUImmediate(c *CPU8008, _ Memory, _ IO, inst Instruction) error {
	group := (inst.Opcode.Code >> 3) & 0x07
	c.executeALU(group, inst.Operands[0])
	return nil
}

func executeALURegister(c *CPU8008, mem Memory, _ IO, inst Instruction) error {
	group := (inst.Opcode.Code >> 3) & 0x07
	src := Register(inst.Opcode.Code & 0x07)
	value, err := c.readRegister(src, mem)
	if err != nil {
		return err
	}
	c.executeALU(group, value)
	return nil
}

func executeIncrement(c *CPU8008, mem Memory, _ IO, inst Instruction) error {
	r := Register((inst.Opcode.Code >> 3) & 0x07)
	value, err := c.readRegister(r, mem)
	if err != nil {
		return err
	}
	value++
	if err := c.writeRegister(r, value, mem); err != nil {
		return err
	}
	c.updateZeroSignParity(value)
	return nil
}

func executeDecrement(c *CPU8008, mem Memory, _ IO, inst Instruction) error {
	r := Register((inst.Opcode.Code >> 3) & 0x07)
	value, err := c.readRegister(r, mem)
	if err != nil {
		return err
	}
	value--
	if err := c.writeRegister(r, value, mem); err != nil {
		return err
	}
	c.updateZeroSignParity(value)
	return nil
}

// executeALU delega l'intero gruppo ALU dell'8008 alla ALU costruita a porte di
// RetroNet Logic, tramite l'adattatore i8008. L'adattatore restituisce risultato
// e flag già nella convenzione 8008 (incluso borrow = NOT carry per le
// sottrazioni). Per CMP il risultato non viene memorizzato.
func (c *CPU8008) executeALU(group byte, value byte) {
	result, flags := i8008.ALU(group, c.A, value, c.Carry)
	c.Carry = flags.Carry
	c.Zero = flags.Zero
	c.Sign = flags.Sign
	c.Parity = flags.Parity
	if group&0x07 != i8008.GroupCMP {
		c.A = result
	}
}

func (c *CPU8008) updateZeroSignParity(value byte) {
	c.Zero = value == 0
	c.Sign = value&0x80 != 0
	c.Parity = evenParity(value)
}

func evenParity(value byte) bool {
	return bits.OnesCount8(value)%2 == 0
}
