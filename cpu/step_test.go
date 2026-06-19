package cpu

import (
	"errors"
	"testing"
)

func TestStepRequiresMemory(t *testing.T) {
	c := newRunningCPU(t)
	c.setPC(0x1234)

	err := c.Step(nil, nil)

	if !errors.Is(err, ErrNilMemory) {
		t.Fatalf("Step(nil, nil) = %v, want ErrNilMemory", err)
	}
	if c.PC != 0x1234 {
		t.Fatalf("PC = 0x%04X, want unchanged 0x1234", c.PC)
	}
}

func TestStepFetchesOpcodeAndAdvancesPC(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	mem.Write(0x0000, NOP())

	err := c.Step(mem, nil)

	if err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}
	if c.Stack[0] != 0x0001 {
		t.Fatalf("Stack[0] = 0x%04X, want current PC 0x0001", c.Stack[0])
	}
}

func TestStepConsumesOperandsBeforeExecution(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	mem.Write(0x0100, 0x44) // JMP, 3 byte
	mem.Write(0x0101, 0x34)
	mem.Write(0x0102, 0x12)
	c.setPC(0x0100)

	err := c.Step(mem, nil)

	if err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.PC != 0x1234 {
		t.Fatalf("PC = 0x%04X, want 0x1234", c.PC)
	}
}

func TestStepPCWrapsAt14Bits(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	mem.Write(0x3FFF, NOP())
	c.setPC(0x3FFF)

	err := c.Step(mem, nil)

	if err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.PC != 0x0000 {
		t.Fatalf("PC = 0x%04X, want wrap to 0x0000", c.PC)
	}
}

func TestStepOperandFetchWrapsAt14Bits(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	mem.Write(0x3FFE, 0x44) // JMP, 3 byte
	mem.Write(0x3FFF, 0xAA)
	mem.Write(0x0000, 0xBB)
	c.setPC(0x3FFE)

	err := c.Step(mem, nil)

	if err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.PC != 0x3BAA {
		t.Fatalf("PC = 0x%04X, want 0x3BAA", c.PC)
	}
}
