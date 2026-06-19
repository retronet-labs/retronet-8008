package cpu

func isLoadImmediateOpcode(code byte) bool {
	return code&0xC7 == 0x06
}

func isLoadMoveOpcode(code byte) bool {
	return code&0xC0 == 0xC0
}

func executeLoadImmediate(c *CPU8008, mem Memory, _ IO, inst Instruction) error {
	dst := Register((inst.Opcode.Code >> 3) & 0x07)
	return c.writeRegister(dst, inst.Operands[0], mem)
}

func executeLoadMove(c *CPU8008, mem Memory, _ IO, inst Instruction) error {
	dst := Register((inst.Opcode.Code >> 3) & 0x07)
	src := Register(inst.Opcode.Code & 0x07)
	if dst == src {
		return nil
	}

	value, err := c.readRegister(src, mem)
	if err != nil {
		return err
	}
	return c.writeRegister(dst, value, mem)
}

func (c *CPU8008) readRegister(r Register, mem Memory) (byte, error) {
	switch Register(regBits(r)) {
	case RegA:
		return c.A, nil
	case RegB:
		return c.B, nil
	case RegC:
		return c.C, nil
	case RegD:
		return c.D, nil
	case RegE:
		return c.E, nil
	case RegH:
		return c.H, nil
	case RegL:
		return c.L, nil
	default:
		if mem == nil {
			return 0, ErrNilMemory
		}
		return mem.Read(c.HL()), nil
	}
}

func (c *CPU8008) writeRegister(r Register, value byte, mem Memory) error {
	switch Register(regBits(r)) {
	case RegA:
		c.A = value
	case RegB:
		c.B = value
	case RegC:
		c.C = value
	case RegD:
		c.D = value
	case RegE:
		c.E = value
	case RegH:
		c.H = value
	case RegL:
		c.L = value
	default:
		if mem == nil {
			return ErrNilMemory
		}
		mem.Write(c.HL(), value)
	}
	return nil
}
