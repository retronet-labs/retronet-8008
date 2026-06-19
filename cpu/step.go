package cpu

// Step esegue il ciclo fetch-decode-execute di una singola istruzione.
//
// In questa milestone il decoder conosce lunghezza e metadata degli opcode, ma
// le funzioni esecutive sono ancora segnaposto: Step consuma opcode e operandi,
// aggiorna il PC a 14 bit e poi restituisce ErrUnimplementedOpcode.
func (c *CPU8008) Step(mem Memory, io IO) error {
	if mem == nil {
		return ErrNilMemory
	}

	pcBefore := c.PC
	code := c.fetch(mem)
	op := Decode(code)
	inst := Instruction{
		PC:     pcBefore,
		Opcode: op,
	}

	for i := byte(1); i < op.Length; i++ {
		inst.Operands[i-1] = c.fetch(mem)
		inst.OperandCount++
	}

	return op.Execute(c, mem, io, inst)
}

func (c *CPU8008) fetch(mem Memory) byte {
	value := mem.Read(c.PC)
	c.setPC(c.PC + 1)
	return value
}
