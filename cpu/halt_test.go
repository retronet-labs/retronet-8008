package cpu

import (
	"errors"
	"testing"
)

func TestStepStoppedAfterResetDoesNotFetch(t *testing.T) {
	c := NewCPU8008()

	err := c.Step(nil, nil)

	if !errors.Is(err, ErrCPUStopped) {
		t.Fatalf("Step(nil, nil) = %v, want ErrCPUStopped", err)
	}
	if c.PC != 0 {
		t.Fatalf("PC = 0x%04X, want unchanged 0x0000", c.PC)
	}
}

func TestHLTSetsHaltedAndStopped(t *testing.T) {
	tests := []byte{0x00, 0x01, 0xFF}

	for _, opcode := range tests {
		t.Run(Decode(opcode).Mnemonic, func(t *testing.T) {
			c := newRunningCPU(t)
			mem := NewFlatMemory()
			mem.Write(0x0000, opcode)
			mem.Write(0x0001, LI(RegA))
			mem.Write(0x0002, 0x55)

			if err := c.Step(mem, nil); err != nil {
				t.Fatalf("HLT Step = %v, want nil", err)
			}
			if !c.Halted || !c.Stopped {
				t.Fatalf("Halted=%v Stopped=%v, want both true", c.Halted, c.Stopped)
			}
			if c.PC != 0x0001 {
				t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
			}

			err := c.Step(mem, nil)
			if !errors.Is(err, ErrCPUStopped) {
				t.Fatalf("Step while stopped = %v, want ErrCPUStopped", err)
			}
			if c.A != 0 {
				t.Fatalf("A = 0x%02X, want unchanged 0x00", c.A)
			}
			if c.PC != 0x0001 {
				t.Fatalf("PC after stopped Step = 0x%04X, want unchanged 0x0001", c.PC)
			}
		})
	}
}

func TestHLTAlias0xFFDecodesAsHalt(t *testing.T) {
	op := Decode(0xFF)
	if op.Mnemonic != "HLT" {
		t.Fatalf("Decode(0xFF).Mnemonic = %q, want HLT", op.Mnemonic)
	}
	if op.Length != 1 {
		t.Fatalf("Decode(0xFF).Length = %d, want 1", op.Length)
	}

	mem := NewFlatMemory()
	mem.Write(0x0000, 0xFF)
	d, err := Disassemble(mem, 0x0000)
	if err != nil {
		t.Fatalf("Disassemble = %v, want nil", err)
	}
	if d.Opcode.Mnemonic != "HLT" {
		t.Fatalf("disasm 0xFF = %q, want HLT (non un move L M,M)", d.Opcode.Mnemonic)
	}
}

func TestJamClearsStoppedStateAndDoesNotFetch(t *testing.T) {
	c := NewCPU8008()
	c.setPC(0x0123)

	if err := c.Jam(nil, nil, NOP()); err != nil {
		t.Fatalf("Jam(NOP) = %v, want nil", err)
	}
	if c.Halted || c.Stopped {
		t.Fatalf("Halted=%v Stopped=%v, want both false", c.Halted, c.Stopped)
	}
	if c.PC != 0x0123 {
		t.Fatalf("PC = 0x%04X, want unchanged 0x0123", c.PC)
	}
}

func TestJamValidatesOperandCount(t *testing.T) {
	c := NewCPU8008()

	err := c.Jam(nil, nil, JMP(), 0x34)

	if !errors.Is(err, ErrInvalidJamInstruction) {
		t.Fatalf("Jam(JMP, one operand) = %v, want ErrInvalidJamInstruction", err)
	}
	if !c.Halted || !c.Stopped {
		t.Fatalf("Halted=%v Stopped=%v, want still stopped", c.Halted, c.Stopped)
	}
}

func TestJamCanExecuteRSTFromStoppedCPU(t *testing.T) {
	c := NewCPU8008()

	if err := c.Jam(nil, nil, RST(2)); err != nil {
		t.Fatalf("Jam(RST 2) = %v, want nil", err)
	}
	if c.Halted || c.Stopped {
		t.Fatalf("Halted=%v Stopped=%v, want running", c.Halted, c.Stopped)
	}
	if c.PC != 0x0010 {
		t.Fatalf("PC = 0x%04X, want vector 0x0010", c.PC)
	}
	if c.SP != 1 {
		t.Fatalf("SP = %d, want 1", c.SP)
	}
	if c.Stack[0] != 0x0000 {
		t.Fatalf("Stack[0] = 0x%04X, want return PC 0x0000", c.Stack[0])
	}
	if c.Stack[1] != 0x0010 {
		t.Fatalf("Stack[1] = 0x%04X, want current PC 0x0010", c.Stack[1])
	}
}

func TestJamResumesAfterHLTWithReturnAddress(t *testing.T) {
	c := newRunningCPU(t)
	mem := NewFlatMemory()
	c.setPC(0x0100)
	mem.Write(0x0100, HLT())

	if err := c.Step(mem, nil); err != nil {
		t.Fatalf("HLT Step = %v, want nil", err)
	}
	if c.PC != 0x0101 || c.Stack[0] != 0x0101 {
		t.Fatalf("after HLT PC=0x%04X Stack[0]=0x%04X, want 0x0101", c.PC, c.Stack[0])
	}

	if err := c.Jam(nil, nil, RST(1)); err != nil {
		t.Fatalf("Jam(RST 1) = %v, want nil", err)
	}
	if c.Halted || c.Stopped {
		t.Fatalf("Halted=%v Stopped=%v, want running", c.Halted, c.Stopped)
	}
	if c.PC != 0x0008 {
		t.Fatalf("PC = 0x%04X, want vector 0x0008", c.PC)
	}
	if c.SP != 1 {
		t.Fatalf("SP = %d, want 1", c.SP)
	}
	if c.Stack[0] != 0x0101 {
		t.Fatalf("Stack[0] = 0x%04X, want return PC 0x0101", c.Stack[0])
	}
	if c.Stack[1] != 0x0008 {
		t.Fatalf("Stack[1] = 0x%04X, want current PC 0x0008", c.Stack[1])
	}
}
