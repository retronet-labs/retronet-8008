package cpu

import "math/bits"

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

func (c *CPU8008) executeALU(group byte, value byte) {
	switch group & 0x07 {
	case 0: // ADD
		c.add(value, false)
	case 1: // ADC
		c.add(value, c.Carry)
	case 2: // SUB
		c.sub(value, false, true)
	case 3: // SBB
		c.sub(value, c.Carry, true)
	case 4: // AND
		c.A &= value
		c.Carry = false
		c.updateZeroSignParity(c.A)
	case 5: // XOR
		c.A ^= value
		c.Carry = false
		c.updateZeroSignParity(c.A)
	case 6: // OR
		c.A |= value
		c.Carry = false
		c.updateZeroSignParity(c.A)
	default: // CMP
		c.sub(value, false, false)
	}
}

func (c *CPU8008) add(value byte, carryIn bool) {
	result := uint16(c.A) + uint16(value)
	if carryIn {
		result++
	}
	c.A = byte(result)
	c.Carry = result > 0xFF
	c.updateZeroSignParity(c.A)
}

func (c *CPU8008) sub(value byte, borrowIn bool, store bool) {
	result := int(c.A) - int(value)
	if borrowIn {
		result--
	}
	out := byte(result)
	c.Carry = result < 0
	if store {
		c.A = out
	}
	c.updateZeroSignParity(out)
}

func (c *CPU8008) updateZeroSignParity(value byte) {
	c.Zero = value == 0
	c.Sign = value&0x80 != 0
	c.Parity = evenParity(value)
}

func evenParity(value byte) bool {
	return bits.OnesCount8(value)%2 == 0
}
