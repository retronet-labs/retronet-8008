package cpu

import "testing"

func TestADRegister(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x14
	c.B = 0x22
	mem.Write(0x0000, AD(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x36 {
		t.Fatalf("A = 0x%02X, want 0x36", c.A)
	}
	assertFlags(t, c, false, false, false, true)
}

func TestADCarryOutZeroParity(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0xFF
	c.B = 0x01
	mem.Write(0x0000, AD(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x00 {
		t.Fatalf("A = 0x%02X, want 0x00", c.A)
	}
	assertFlags(t, c, true, true, false, true)
}

func TestACUsesCarryIn(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x01
	c.B = 0x01
	c.Carry = true
	mem.Write(0x0000, AC(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x03 {
		t.Fatalf("A = 0x%02X, want 0x03", c.A)
	}
	assertFlags(t, c, false, false, false, true)
}

func TestSUWithoutBorrow(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x05
	c.B = 0x03
	mem.Write(0x0000, SU(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x02 {
		t.Fatalf("A = 0x%02X, want 0x02", c.A)
	}
	assertFlags(t, c, false, false, false, false)
}

func TestSUWithBorrow(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x03
	c.B = 0x05
	mem.Write(0x0000, SU(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0xFE {
		t.Fatalf("A = 0x%02X, want 0xFE", c.A)
	}
	assertFlags(t, c, true, false, true, false)
}

func TestSBUsesBorrowIn(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x05
	c.B = 0x03
	c.Carry = true
	mem.Write(0x0000, SB(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x01 {
		t.Fatalf("A = 0x%02X, want 0x01", c.A)
	}
	assertFlags(t, c, false, false, false, false)
}

func TestNDClearsCarryAndUpdatesFlags(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0xF0
	c.B = 0x0F
	c.Carry = true
	mem.Write(0x0000, ND(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x00 {
		t.Fatalf("A = 0x%02X, want 0x00", c.A)
	}
	assertFlags(t, c, false, true, false, true)
}

func TestXR(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0xFF
	c.B = 0x0F
	c.Carry = true
	mem.Write(0x0000, XR(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0xF0 {
		t.Fatalf("A = 0x%02X, want 0xF0", c.A)
	}
	assertFlags(t, c, false, false, true, true)
}

func TestOR(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x80
	c.B = 0x01
	c.Carry = true
	mem.Write(0x0000, OR(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x81 {
		t.Fatalf("A = 0x%02X, want 0x81", c.A)
	}
	assertFlags(t, c, false, false, true, true)
}

func TestCPDoesNotModifyAccumulator(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x10
	c.B = 0x10
	mem.Write(0x0000, CP(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x10 {
		t.Fatalf("A = 0x%02X, want unchanged 0x10", c.A)
	}
	assertFlags(t, c, false, true, false, true)
}

func TestCPBorrowDoesNotModifyAccumulator(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x00
	c.B = 0x01
	mem.Write(0x0000, CP(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x00 {
		t.Fatalf("A = 0x%02X, want unchanged 0x00", c.A)
	}
	assertFlags(t, c, true, false, true, true)
}

func TestADIImmediate(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x02
	mem.Write(0x0000, ADI())
	mem.Write(0x0001, 0x03)

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x05 {
		t.Fatalf("A = 0x%02X, want 0x05", c.A)
	}
	if c.PC != 0x0002 {
		t.Fatalf("PC = 0x%04X, want 0x0002", c.PC)
	}
	assertFlags(t, c, false, false, false, true)
}

func TestSBIImmediateWithBorrowIn(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x04
	c.Carry = true
	mem.Write(0x0000, SBI())
	mem.Write(0x0001, 0x02)

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x01 {
		t.Fatalf("A = 0x%02X, want 0x01", c.A)
	}
	assertFlags(t, c, false, false, false, false)
}

func TestALUReadsMAtHL(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.A = 0x02
	c.H = 0x40
	c.L = 0x20
	mem.Write(0x0020, 0x03)
	mem.Write(0x0000, AD(RegM))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x05 {
		t.Fatalf("A = 0x%02X, want 0x05", c.A)
	}
}

func TestINRDoesNotModifyCarry(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.B = 0xFF
	c.Carry = true
	mem.Write(0x0000, INR(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.B != 0x00 {
		t.Fatalf("B = 0x%02X, want 0x00", c.B)
	}
	assertFlags(t, c, true, true, false, true)
}

func TestDCRDoesNotModifyCarry(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.B = 0x00
	c.Carry = false
	mem.Write(0x0000, DCR(RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.B != 0xFF {
		t.Fatalf("B = 0x%02X, want 0xFF", c.B)
	}
	assertFlags(t, c, false, false, true, true)
}

func assertFlags(t *testing.T, c *CPU8008, carry, zero, sign, parity bool) {
	t.Helper()
	if c.Carry != carry || c.Zero != zero || c.Sign != sign || c.Parity != parity {
		t.Fatalf("flags C=%v Z=%v S=%v P=%v, want C=%v Z=%v S=%v P=%v", c.Carry, c.Zero, c.Sign, c.Parity, carry, zero, sign, parity)
	}
}
