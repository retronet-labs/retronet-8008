package cpu

func isHaltOpcode(code byte) bool {
	return code == 0x00 || code == 0x01
}

func executeHalt(c *CPU8008, _ Memory, _ IO, _ Instruction) error {
	c.Halted = true
	c.Stopped = true
	return nil
}
