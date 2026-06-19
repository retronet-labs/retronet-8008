package cpu

import "testing"

func TestJMP(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	writeAddressedOpcode(mem, 0x0000, JMP(), 0x1234)

	if err := c.Step(mem, nil); err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.PC != 0x1234 {
		t.Fatalf("PC = 0x%04X, want 0x1234", c.PC)
	}
	if c.Stack[0] != 0x1234 {
		t.Fatalf("Stack[0] = 0x%04X, want current PC 0x1234", c.Stack[0])
	}
}

func TestJMPMasksTargetAddress(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	mem.Write(0x0000, JMP())
	mem.Write(0x0001, 0x34)
	mem.Write(0x0002, 0xFF)

	if err := c.Step(mem, nil); err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.PC != 0x3F34 {
		t.Fatalf("PC = 0x%04X, want 0x3F34", c.PC)
	}
}

func TestConditionalJumpTakenAndNotTaken(t *testing.T) {
	tests := []struct {
		name   string
		opcode byte
		carry  bool
		wantPC uint16
	}{
		{"JF takes when false", JF(CondCarry), false, 0x0200},
		{"JF skips when true", JF(CondCarry), true, 0x0003},
		{"JT takes when true", JT(CondCarry), true, 0x0200},
		{"JT skips when false", JT(CondCarry), false, 0x0003},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRunningCPU(t)
			c.Carry = tt.carry
			mem := NewFlatMemory()
			writeAddressedOpcode(mem, 0x0000, tt.opcode, 0x0200)

			if err := c.Step(mem, nil); err != nil {
				t.Fatalf("Step = %v, want nil", err)
			}
			if c.PC != tt.wantPC {
				t.Fatalf("PC = 0x%04X, want 0x%04X", c.PC, tt.wantPC)
			}
			if c.SP != 0 {
				t.Fatalf("SP = %d, want 0", c.SP)
			}
		})
	}
}

func TestCALAndRET(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	writeAddressedOpcode(mem, 0x0000, CAL(), 0x0010)
	mem.Write(0x0010, RET())

	if err := c.Step(mem, nil); err != nil {
		t.Fatalf("CAL Step = %v, want nil", err)
	}
	if c.PC != 0x0010 {
		t.Fatalf("PC after CAL = 0x%04X, want 0x0010", c.PC)
	}
	if c.SP != 1 {
		t.Fatalf("SP after CAL = %d, want 1", c.SP)
	}
	if c.Stack[0] != 0x0003 {
		t.Fatalf("Stack[0] = 0x%04X, want return PC 0x0003", c.Stack[0])
	}
	if c.Stack[1] != 0x0010 {
		t.Fatalf("Stack[1] = 0x%04X, want current PC 0x0010", c.Stack[1])
	}

	if err := c.Step(mem, nil); err != nil {
		t.Fatalf("RET Step = %v, want nil", err)
	}
	if c.PC != 0x0003 {
		t.Fatalf("PC after RET = 0x%04X, want 0x0003", c.PC)
	}
	if c.SP != 0 {
		t.Fatalf("SP after RET = %d, want 0", c.SP)
	}
}

func TestConditionalCallTakenAndNotTaken(t *testing.T) {
	tests := []struct {
		name   string
		opcode byte
		zero   bool
		wantPC uint16
		wantSP uint8
	}{
		{"CF takes when false", CF(CondZero), false, 0x0120, 1},
		{"CF skips when true", CF(CondZero), true, 0x0003, 0},
		{"CT takes when true", CT(CondZero), true, 0x0120, 1},
		{"CT skips when false", CT(CondZero), false, 0x0003, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRunningCPU(t)
			c.Zero = tt.zero
			mem := NewFlatMemory()
			writeAddressedOpcode(mem, 0x0000, tt.opcode, 0x0120)

			if err := c.Step(mem, nil); err != nil {
				t.Fatalf("Step = %v, want nil", err)
			}
			if c.PC != tt.wantPC {
				t.Fatalf("PC = 0x%04X, want 0x%04X", c.PC, tt.wantPC)
			}
			if c.SP != tt.wantSP {
				t.Fatalf("SP = %d, want %d", c.SP, tt.wantSP)
			}
			if tt.wantSP == 1 && c.Stack[0] != 0x0003 {
				t.Fatalf("Stack[0] = 0x%04X, want return PC 0x0003", c.Stack[0])
			}
		})
	}
}

func TestConditionalReturnTakenAndNotTaken(t *testing.T) {
	tests := []struct {
		name    string
		opcode  byte
		parity  bool
		wantPC  uint16
		wantSP  uint8
		wantTop uint16
	}{
		{"RF takes when false", RF(CondParity), false, 0x0222, 0, 0x0222},
		{"RF skips when true", RF(CondParity), true, 0x0101, 1, 0x0101},
		{"RT takes when true", RT(CondParity), true, 0x0222, 0, 0x0222},
		{"RT skips when false", RT(CondParity), false, 0x0101, 1, 0x0101},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRunningCPU(t)
			c.Parity = tt.parity
			c.setStack(0, 0x0222)
			c.setSP(1)
			c.setPC(0x0100)
			mem := NewFlatMemory()
			mem.Write(0x0100, tt.opcode)

			if err := c.Step(mem, nil); err != nil {
				t.Fatalf("Step = %v, want nil", err)
			}
			if c.PC != tt.wantPC {
				t.Fatalf("PC = 0x%04X, want 0x%04X", c.PC, tt.wantPC)
			}
			if c.SP != tt.wantSP {
				t.Fatalf("SP = %d, want %d", c.SP, tt.wantSP)
			}
			if c.Stack[c.SP] != tt.wantTop {
				t.Fatalf("Stack[SP] = 0x%04X, want 0x%04X", c.Stack[c.SP], tt.wantTop)
			}
		})
	}
}

func TestRST(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	mem.Write(0x0000, RST(3))

	if err := c.Step(mem, nil); err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.PC != 0x0018 {
		t.Fatalf("PC = 0x%04X, want vector 0x0018", c.PC)
	}
	if c.SP != 1 {
		t.Fatalf("SP = %d, want 1", c.SP)
	}
	if c.Stack[0] != 0x0001 {
		t.Fatalf("Stack[0] = 0x%04X, want return PC 0x0001", c.Stack[0])
	}
	if c.Stack[1] != 0x0018 {
		t.Fatalf("Stack[1] = 0x%04X, want current PC 0x0018", c.Stack[1])
	}
}

func TestStackDepthSevenUsefulReturns(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()

	for i := uint16(0); i < 7; i++ {
		writeAddressedOpcode(mem, i*0x10, CAL(), (i+1)*0x10)
	}

	for i := 0; i < 7; i++ {
		if err := c.Step(mem, nil); err != nil {
			t.Fatalf("call %d Step = %v, want nil", i, err)
		}
	}

	if c.SP != 7 {
		t.Fatalf("SP = %d, want 7", c.SP)
	}
	if c.PC != 0x0070 {
		t.Fatalf("PC = 0x%04X, want 0x0070", c.PC)
	}
	if c.Stack[0] != 0x0003 {
		t.Fatalf("Stack[0] = 0x%04X, want first return 0x0003", c.Stack[0])
	}
	if c.Stack[6] != 0x0063 {
		t.Fatalf("Stack[6] = 0x%04X, want seventh return 0x0063", c.Stack[6])
	}
	if c.Stack[7] != 0x0070 {
		t.Fatalf("Stack[7] = 0x%04X, want current PC 0x0070", c.Stack[7])
	}
}

func TestStackOverflowWrapsSilently(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()

	for i := uint16(0); i < 8; i++ {
		writeAddressedOpcode(mem, i*0x10, CAL(), (i+1)*0x10)
	}

	for i := 0; i < 8; i++ {
		if err := c.Step(mem, nil); err != nil {
			t.Fatalf("call %d Step = %v, want nil", i, err)
		}
	}

	if c.SP != 0 {
		t.Fatalf("SP = %d, want wrapped 0", c.SP)
	}
	if c.PC != 0x0080 {
		t.Fatalf("PC = 0x%04X, want 0x0080", c.PC)
	}
	if c.Stack[0] != 0x0080 {
		t.Fatalf("Stack[0] = 0x%04X, want wrapped current PC 0x0080", c.Stack[0])
	}
	if c.Stack[7] != 0x0073 {
		t.Fatalf("Stack[7] = 0x%04X, want eighth return PC 0x0073", c.Stack[7])
	}
}

func writeAddressedOpcode(mem *FlatMemory, at uint16, opcode byte, target uint16) {
	mem.Write(at, opcode)
	mem.Write(at+1, byte(target))
	mem.Write(at+2, byte(target>>8))
}
