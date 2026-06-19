package cpu

func isRotateOpcode(code byte) bool {
	switch code {
	case 0x02, 0x0A, 0x12, 0x1A:
		return true
	default:
		return false
	}
}

func executeRotate(c *CPU8008, _ Memory, _ IO, inst Instruction) error {
	switch inst.Opcode.Code {
	case RLC():
		old := c.A
		c.Carry = old&0x80 != 0
		c.A = (old << 1) | (old >> 7)
	case RRC():
		old := c.A
		c.Carry = old&0x01 != 0
		c.A = (old >> 1) | (old << 7)
	case RAL():
		old := c.A
		carryIn := byte(0)
		if c.Carry {
			carryIn = 1
		}
		c.Carry = old&0x80 != 0
		c.A = (old << 1) | carryIn
	case RAR():
		old := c.A
		carryIn := byte(0)
		if c.Carry {
			carryIn = 0x80
		}
		c.Carry = old&0x01 != 0
		c.A = (old >> 1) | carryIn
	}
	return nil
}
