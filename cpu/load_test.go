package cpu

import "testing"

func TestLRegisterToRegister(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	c.B = 0x42
	c.Carry = true
	c.Zero = true
	mem.Write(0x0000, L(RegA, RegB))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x42 {
		t.Fatalf("A = 0x%02X, want 0x42", c.A)
	}
	if c.B != 0x42 {
		t.Fatalf("B = 0x%02X, want unchanged 0x42", c.B)
	}
	if !c.Carry || !c.Zero {
		t.Fatal("load register should not modify flags")
	}
	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}
}

func TestLRegisterSelfIsNOP(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	c.A = 0x9A
	c.Carry = true
	mem.Write(0x0000, NOP())

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.A != 0x9A {
		t.Fatalf("A = 0x%02X, want unchanged 0x9A", c.A)
	}
	if !c.Carry {
		t.Fatal("NOP/load self should not modify flags")
	}
	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}
}

func TestLIRegisterImmediate(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	mem.Write(0x0000, LI(RegD))
	mem.Write(0x0001, 0x8A)

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.D != 0x8A {
		t.Fatalf("D = 0x%02X, want 0x8A", c.D)
	}
	if c.PC != 0x0002 {
		t.Fatalf("PC = 0x%04X, want 0x0002", c.PC)
	}
}

func TestLFromMReadsMemoryAtHL(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	c.H = 0x80
	c.L = 0x20
	mem.Write(0x0020, 0xA7)
	mem.Write(0x0000, L(RegE, RegM))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if c.E != 0xA7 {
		t.Fatalf("E = 0x%02X, want 0xA7", c.E)
	}
}

func TestLToMWritesMemoryAtHL(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	c.A = 0xB6
	c.H = 0x3F
	c.L = 0x34
	mem.Write(0x0000, L(RegM, RegA))

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if got := mem.Read(0x3F34); got != 0xB6 {
		t.Fatalf("mem[0x3F34] = 0x%02X, want 0xB6", got)
	}
}

func TestLMIWritesImmediateToMemoryAtHL(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	c.H = 0xFF
	c.L = 0x34
	mem.Write(0x0000, LI(RegM))
	mem.Write(0x0001, 0x5C)

	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}

	if got := mem.Read(0x3F34); got != 0x5C {
		t.Fatalf("mem[0x3F34] = 0x%02X, want 0x5C", got)
	}
	if c.PC != 0x0002 {
		t.Fatalf("PC = 0x%04X, want 0x0002", c.PC)
	}
}
