package cpu

func isControlFlowOpcode(code byte) bool {
	switch {
	case code&0xC7 == 0x40: // JFc/JTc
		return true
	case code&0xC7 == 0x42: // CFc/CTc
		return true
	case code&0xC7 == 0x44: // JMP aliases
		return true
	case code&0xC7 == 0x46: // CAL aliases
		return true
	case code&0xC7 == 0x03: // RFc/RTc
		return true
	case code&0xC7 == 0x07: // RET aliases
		return true
	case code&0xC7 == 0x05: // RST
		return true
	default:
		return false
	}
}

func executeControlFlow(c *CPU8008, _ Memory, _ IO, inst Instruction) error {
	code := inst.Opcode.Code
	switch {
	case code&0xC7 == 0x40:
		if c.conditionTaken(code) {
			c.setPC(inst.addressOperand())
		}
	case code&0xC7 == 0x42:
		if c.conditionTaken(code) {
			c.call(inst.addressOperand())
		}
	case code&0xC7 == 0x44:
		c.setPC(inst.addressOperand())
	case code&0xC7 == 0x46:
		c.call(inst.addressOperand())
	case code&0xC7 == 0x03:
		if c.conditionTaken(code) {
			c.ret()
		}
	case code&0xC7 == 0x07:
		c.ret()
	case code&0xC7 == 0x05:
		c.call(uint16((code>>3)&0x07) << 3)
	}
	return nil
}

func (inst Instruction) addressOperand() uint16 {
	low := uint16(inst.Operands[0])
	high := uint16(inst.Operands[1] & 0x3F)
	return (high << 8) | low
}

func (c *CPU8008) conditionTaken(code byte) bool {
	value := c.conditionValue(Condition((code >> 3) & 0x03))
	if code&0x20 != 0 {
		return value
	}
	return !value
}

func (c *CPU8008) conditionValue(cond Condition) bool {
	switch cond {
	case CondCarry:
		return c.Carry
	case CondZero:
		return c.Zero
	case CondSign:
		return c.Sign
	default:
		return c.Parity
	}
}

func (c *CPU8008) call(target uint16) {
	c.setSP(c.SP + 1)
	c.setPC(target)
}

func (c *CPU8008) ret() {
	c.setSP(c.SP - 1)
	c.setPC(c.Stack[c.SP])
}
