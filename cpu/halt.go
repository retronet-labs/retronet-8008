package cpu

func isHaltOpcode(code byte) bool {
	// L'8008 ferma la CPU su 0x00 e 0x01 e, come alias documentato, su 0xFF
	// (lo slot L M,M, privo di senso come trasferimento registro-registro).
	return code == 0x00 || code == 0x01 || code == 0xFF
}

func executeHalt(c *CPU8008, _ Memory, _ IO, _ Instruction) error {
	c.Halted = true
	c.Stopped = true
	return nil
}
