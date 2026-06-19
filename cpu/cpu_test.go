package cpu

import "testing"

func TestNewCPU8008ResetState(t *testing.T) {
	c := NewCPU8008()

	assertResetState(t, c)
}

func TestResetClearsStateAndStopsCPU(t *testing.T) {
	c := NewCPU8008()
	c.A = 1
	c.B = 2
	c.C = 3
	c.D = 4
	c.E = 5
	c.H = 0xFF
	c.L = 0xAA
	c.Carry = true
	c.Zero = true
	c.Sign = true
	c.Parity = true
	c.PC = 0x3ABC
	c.SP = 7
	c.Halted = false
	c.Stopped = false
	for i := range c.Stack {
		c.Stack[i] = uint16(i + 1)
	}

	c.Reset()

	assertResetState(t, c)
}

func TestSetPCMasksTo14Bits(t *testing.T) {
	c := NewCPU8008()

	c.setPC(0x4000)
	if c.PC != 0x0000 {
		t.Fatalf("PC = 0x%04X, want 0x0000", c.PC)
	}

	c.setPC(0xFFFF)
	if c.PC != 0x3FFF {
		t.Fatalf("PC = 0x%04X, want 0x3FFF", c.PC)
	}
}

func TestHLUsesOnlyLowerSixBitsOfH(t *testing.T) {
	c := NewCPU8008()
	c.H = 0xFF
	c.L = 0x34

	if got := c.HL(); got != 0x3F34 {
		t.Fatalf("HL() = 0x%04X, want 0x3F34", got)
	}

	c.H = 0x40
	c.L = 0x12
	if got := c.HL(); got != 0x0012 {
		t.Fatalf("HL() = 0x%04X, want 0x0012", got)
	}
}

func TestSetSPMasksToThreeBits(t *testing.T) {
	c := NewCPU8008()
	c.setSP(0xFF)

	if c.SP != 7 {
		t.Fatalf("SP = %d, want 7", c.SP)
	}
}

func TestSetStackMasksSlotAndAddress(t *testing.T) {
	c := NewCPU8008()
	c.setStack(9, 0xFFFF)

	if c.Stack[1] != 0x3FFF {
		t.Fatalf("Stack[1] = 0x%04X, want 0x3FFF", c.Stack[1])
	}
}

func assertResetState(t *testing.T, c *CPU8008) {
	t.Helper()

	if c.A != 0 || c.B != 0 || c.C != 0 || c.D != 0 || c.E != 0 || c.H != 0 || c.L != 0 {
		t.Fatalf("registri non azzerati: A=%d B=%d C=%d D=%d E=%d H=%d L=%d", c.A, c.B, c.C, c.D, c.E, c.H, c.L)
	}
	if c.Carry || c.Zero || c.Sign || c.Parity {
		t.Fatalf("flag non azzerati: C=%v Z=%v S=%v P=%v", c.Carry, c.Zero, c.Sign, c.Parity)
	}
	if c.PC != 0 {
		t.Fatalf("PC = 0x%04X, want 0x0000", c.PC)
	}
	if c.SP != 0 {
		t.Fatalf("SP = %d, want 0", c.SP)
	}
	for i, v := range c.Stack {
		if v != 0 {
			t.Fatalf("Stack[%d] = 0x%04X, want 0", i, v)
		}
	}
	if !c.Halted {
		t.Fatal("Halted = false, want true after reset")
	}
	if !c.Stopped {
		t.Fatal("Stopped = false, want true after reset")
	}
}
