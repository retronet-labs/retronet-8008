package cpu

func isInputOpcode(code byte) bool {
	return code&0xF0 == 0x40 && code&0x01 == 0x01
}

func isOutputOpcode(code byte) bool {
	return code&0xC0 == 0x40 && code&0x30 != 0 && code&0x01 == 0x01
}

func executeInput(c *CPU8008, _ Memory, io IO, inst Instruction) error {
	if io == nil {
		return ErrNilIO
	}
	port := inputPort(inst.Opcode.Code)
	if err := ValidateInputPort(port); err != nil {
		return err
	}
	c.A = io.Input(port)
	return nil
}

func executeOutput(c *CPU8008, _ Memory, io IO, inst Instruction) error {
	if io == nil {
		return ErrNilIO
	}
	port := outputPort(inst.Opcode.Code)
	if err := ValidateOutputPort(port); err != nil {
		return err
	}
	io.Output(port, c.A)
	return nil
}

func inputPort(code byte) byte {
	return (code >> 1) & 0x07
}

func outputPort(code byte) byte {
	return (code >> 1) & 0x1F
}
