package cpu

import "testing"

func TestRLC(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0b1000_0001
	c.Zero = true
	c.Sign = true
	c.Parity = false
	mem.Write(0x0000, RLC())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0b0000_0011 {
		t.Fatalf("A = %08b, want 00000011", c.A)
	}
	if !c.Carry {
		t.Fatal("Carry = false, want true")
	}
	assertNonCarryFlags(t, c, true, true, false)
}

func TestRRC(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0b1000_0001
	c.Zero = true
	c.Sign = false
	c.Parity = true
	mem.Write(0x0000, RRC())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0b1100_0000 {
		t.Fatalf("A = %08b, want 11000000", c.A)
	}
	if !c.Carry {
		t.Fatal("Carry = false, want true")
	}
	assertNonCarryFlags(t, c, true, false, true)
}

func TestRALWithCarryClear(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0b0100_0001
	c.Carry = false
	c.Zero = false
	c.Sign = true
	c.Parity = true
	mem.Write(0x0000, RAL())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0b1000_0010 {
		t.Fatalf("A = %08b, want 10000010", c.A)
	}
	if c.Carry {
		t.Fatal("Carry = true, want false")
	}
	assertNonCarryFlags(t, c, false, true, true)
}

func TestRALWithCarrySetAndCarryOut(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0b1000_0000
	c.Carry = true
	c.Zero = true
	c.Sign = false
	c.Parity = false
	mem.Write(0x0000, RAL())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0b0000_0001 {
		t.Fatalf("A = %08b, want 00000001", c.A)
	}
	if !c.Carry {
		t.Fatal("Carry = false, want true")
	}
	assertNonCarryFlags(t, c, true, false, false)
}

func TestRARWithCarryClear(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0b0000_0010
	c.Carry = false
	c.Zero = true
	c.Sign = true
	c.Parity = false
	mem.Write(0x0000, RAR())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0b0000_0001 {
		t.Fatalf("A = %08b, want 00000001", c.A)
	}
	if c.Carry {
		t.Fatal("Carry = true, want false")
	}
	assertNonCarryFlags(t, c, true, true, false)
}

func TestRARWithCarrySetAndCarryOut(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0b0000_0001
	c.Carry = true
	c.Zero = false
	c.Sign = false
	c.Parity = true
	mem.Write(0x0000, RAR())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0b1000_0000 {
		t.Fatalf("A = %08b, want 10000000", c.A)
	}
	if !c.Carry {
		t.Fatal("Carry = false, want true")
	}
	assertNonCarryFlags(t, c, false, false, true)
}

func TestRotateAdvancesPC(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	mem.Write(0x0000, RLC())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}
}

func assertNonCarryFlags(t *testing.T, c *CPU8008, zero, sign, parity bool) {
	t.Helper()
	if c.Zero != zero || c.Sign != sign || c.Parity != parity {
		t.Fatalf("flags Z=%v S=%v P=%v, want Z=%v S=%v P=%v", c.Zero, c.Sign, c.Parity, zero, sign, parity)
	}
}
