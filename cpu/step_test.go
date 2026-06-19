package cpu

import (
	"errors"
	"testing"
)

func TestStepRequiresMemory(t *testing.T) {
	c := NewCPU8008()
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
	c := NewCPU8008()
	mem := NewFlatMemory()
	mem.Write(0x0000, 0xC1) // LAB, 1 byte

	err := c.Step(mem, nil)

	if !errors.Is(err, ErrUnimplementedOpcode) {
		t.Fatalf("Step = %v, want ErrUnimplementedOpcode", err)
	}
	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}

	var opErr *UnimplementedOpcodeError
	if !errors.As(err, &opErr) {
		t.Fatalf("Step error type = %T, want *UnimplementedOpcodeError", err)
	}
	if opErr.PC != 0x0000 || opErr.Opcode != 0xC1 || opErr.Mnemonic != "LAB" || opErr.Length != 1 {
		t.Fatalf("unexpected opcode error: %+v", opErr)
	}
}

func TestStepConsumesOperandsBeforeUnimplementedError(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	mem.Write(0x0100, 0x44) // JMP, 3 byte
	mem.Write(0x0101, 0x34)
	mem.Write(0x0102, 0x12)
	c.setPC(0x0100)

	err := c.Step(mem, nil)

	if !errors.Is(err, ErrUnimplementedOpcode) {
		t.Fatalf("Step = %v, want ErrUnimplementedOpcode", err)
	}
	if c.PC != 0x0103 {
		t.Fatalf("PC = 0x%04X, want 0x0103", c.PC)
	}
}

func TestStepPCWrapsAt14Bits(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	mem.Write(0x3FFF, 0xC0) // NOP alias, 1 byte
	c.setPC(0x3FFF)

	err := c.Step(mem, nil)

	if !errors.Is(err, ErrUnimplementedOpcode) {
		t.Fatalf("Step = %v, want ErrUnimplementedOpcode", err)
	}
	if c.PC != 0x0000 {
		t.Fatalf("PC = 0x%04X, want wrap to 0x0000", c.PC)
	}
}

func TestStepOperandFetchWrapsAt14Bits(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	mem.Write(0x3FFE, 0x44) // JMP, 3 byte
	mem.Write(0x3FFF, 0xAA)
	mem.Write(0x0000, 0xBB)
	c.setPC(0x3FFE)

	err := c.Step(mem, nil)

	if !errors.Is(err, ErrUnimplementedOpcode) {
		t.Fatalf("Step = %v, want ErrUnimplementedOpcode", err)
	}
	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}
}
