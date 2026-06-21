package cpu

// Step esegue il ciclo fetch-decode-execute di una singola istruzione.
//
// Step consuma opcode e operandi, aggiorna il PC a 14 bit e invoca la funzione
// esecutiva registrata nel decoder. I sei encoding non definiti dall'ISA
// restituiscono ErrUnimplementedOpcode. Se la CPU e' ferma, Step non accede al
// bus e restituisce ErrCPUStopped.
func (c *CPU8008) Step(mem Memory, io IO) error {
	if c.Halted || c.Stopped {
		return ErrCPUStopped
	}
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

	timing := c.instructionTiming(op)
	if err := op.Execute(c, mem, io, inst); err != nil {
		return err
	}
	c.recordTiming(timing)
	return nil
}

func (c *CPU8008) fetch(mem Memory) byte {
	value := mem.Read(c.PC)
	c.setPC(c.PC + 1)
	return value
}

func (c *CPU8008) instructionTiming(op Opcode) InstructionTiming {
	timing := InstructionTiming{
		States:     op.States,
		CycleCount: op.CycleCount,
		Cycles:     op.Cycles,
		Taken:      true,
	}
	code := op.Code
	if code&0xC7 == 0x40 || code&0xC7 == 0x42 || code&0xC7 == 0x03 {
		timing.Conditional = true
		timing.Taken = c.conditionTaken(code)
		if !timing.Taken {
			timing.States = op.MinStates
		}
	}
	return timing
}

func (c *CPU8008) recordTiming(timing InstructionTiming) {
	timing.WaitStates = c.pendingWaitStates
	c.pendingWaitStates = 0
	c.InstructionCount++
	c.StateCount += uint64(timing.States)
	c.LastTiming = timing
}

// RecordWaitState registra uno stato WAIT richiesto dalla logica macchina.
func (c *CPU8008) RecordWaitState() {
	c.StateCount++
	c.WaitStateCount++
	c.pendingWaitStates++
}
