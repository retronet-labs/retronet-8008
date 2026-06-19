package cpu

import (
	"errors"
	"testing"
)

func TestDisassembleOneByteInstruction(t *testing.T) {
	mem := NewFlatMemory()
	mem.Write(0x0000, L(RegA, RegB))

	d, err := Disassemble(mem, 0x0000)
	if err != nil {
		t.Fatalf("Disassemble = %v, want nil", err)
	}
	if d.PC != 0x0000 || d.NextPC != 0x0001 {
		t.Fatalf("PC/NextPC = 0x%04X/0x%04X, want 0x0000/0x0001", d.PC, d.NextPC)
	}
	if d.Length != 1 || d.Opcode.Mnemonic != "LAB" || d.Operand != "" {
		t.Fatalf("unexpected disassembly: %+v", d)
	}
	if got := d.String(); got != "0000: C1       LAB" {
		t.Fatalf("String() = %q", got)
	}
}

func TestDisassembleImmediateInstruction(t *testing.T) {
	mem := NewFlatMemory()
	mem.Write(0x0010, LI(RegA))
	mem.Write(0x0011, 0x2A)

	d, err := Disassemble(mem, 0x0010)
	if err != nil {
		t.Fatalf("Disassemble = %v, want nil", err)
	}
	if d.NextPC != 0x0012 {
		t.Fatalf("NextPC = 0x%04X, want 0x0012", d.NextPC)
	}
	if d.Operand != "#0x2A" {
		t.Fatalf("Operand = %q, want #0x2A", d.Operand)
	}
	if got := d.String(); got != "0010: 06 2A    LAI #0x2A" {
		t.Fatalf("String() = %q", got)
	}
}

func TestDisassembleAddressInstructionMasksHighByte(t *testing.T) {
	mem := NewFlatMemory()
	mem.Write(0x0020, JMP())
	mem.Write(0x0021, 0x34)
	mem.Write(0x0022, 0xFF)

	d, err := Disassemble(mem, 0x0020)
	if err != nil {
		t.Fatalf("Disassemble = %v, want nil", err)
	}
	if d.Operand != "0x3F34" {
		t.Fatalf("Operand = %q, want 0x3F34", d.Operand)
	}
	if got := d.String(); got != "0020: 44 34 FF JMP 0x3F34" {
		t.Fatalf("String() = %q", got)
	}
}

func TestDisassembleWrapsOperandReadsAndNextPC(t *testing.T) {
	mem := NewFlatMemory()
	mem.Write(0x3FFF, LI(RegB))
	mem.Write(0x0000, 0x77)

	d, err := Disassemble(mem, 0x3FFF)
	if err != nil {
		t.Fatalf("Disassemble = %v, want nil", err)
	}
	if d.NextPC != 0x0001 {
		t.Fatalf("NextPC = 0x%04X, want 0x0001", d.NextPC)
	}
	if d.Bytes[1] != 0x77 || d.Operand != "#0x77" {
		t.Fatalf("operand byte/operand = 0x%02X/%q, want 0x77/#0x77", d.Bytes[1], d.Operand)
	}
}

func TestDisassembleRequiresMemory(t *testing.T) {
	_, err := Disassemble(nil, 0x0000)
	if !errors.Is(err, ErrNilMemory) {
		t.Fatalf("Disassemble(nil) = %v, want ErrNilMemory", err)
	}
}
