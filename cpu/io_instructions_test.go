package cpu

import (
	"errors"
	"fmt"
	"testing"
)

func TestINPReadsInputPortIntoAccumulator(t *testing.T) {
	c := newRunningCPU(t)
	c.Carry = true
	c.Sign = true
	c.Parity = true
	mem := NewFlatMemory()
	ports := NewPorts()
	if err := ports.SetInput(7, 0xA5); err != nil {
		t.Fatal(err)
	}
	mem.Write(0x0000, INP(7))

	if err := c.Step(mem, ports); err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	if c.A != 0xA5 {
		t.Fatalf("A = 0x%02X, want 0xA5", c.A)
	}
	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}
	if !c.Carry || !c.Sign || !c.Parity || c.Zero {
		t.Fatalf("flags changed unexpectedly: C=%v Z=%v S=%v P=%v", c.Carry, c.Zero, c.Sign, c.Parity)
	}
}

func TestOUTWritesAccumulatorToOutputPort(t *testing.T) {
	c := newRunningCPU(t)
	c.A = 0x3C
	c.Carry = true
	c.Zero = true
	mem := NewFlatMemory()
	ports := NewPorts()
	mem.Write(0x0000, OUT(31))

	if err := c.Step(mem, ports); err != nil {
		t.Fatalf("Step = %v, want nil", err)
	}
	got, err := ports.OutputValue(31)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0x3C {
		t.Fatalf("OutputValue(31) = 0x%02X, want 0x3C", got)
	}
	if c.A != 0x3C {
		t.Fatalf("A = 0x%02X, want unchanged 0x3C", c.A)
	}
	if c.PC != 0x0001 {
		t.Fatalf("PC = 0x%04X, want 0x0001", c.PC)
	}
	if !c.Carry || !c.Zero || c.Sign || c.Parity {
		t.Fatalf("flags changed unexpectedly: C=%v Z=%v S=%v P=%v", c.Carry, c.Zero, c.Sign, c.Parity)
	}
}

func TestIOInstructionsRequireBus(t *testing.T) {
	tests := []struct {
		name   string
		opcode byte
	}{
		{"INP", INP(0)},
		{"OUT", OUT(8)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRunningCPU(t)
			mem := NewFlatMemory()
			mem.Write(0x0000, tt.opcode)

			err := c.Step(mem, nil)

			if !errors.Is(err, ErrNilIO) {
				t.Fatalf("Step = %v, want ErrNilIO", err)
			}
			if c.PC != 0x0001 {
				t.Fatalf("PC = 0x%04X, want 0x0001 after fetch", c.PC)
			}
		})
	}
}

func TestJamCanExecuteIOInstructions(t *testing.T) {
	c := NewCPU8008()
	ports := NewPorts()
	if err := ports.SetInput(3, 0xC7); err != nil {
		t.Fatal(err)
	}

	if err := c.Jam(nil, ports, INP(3)); err != nil {
		t.Fatalf("Jam(INP 3) = %v, want nil", err)
	}
	if c.A != 0xC7 {
		t.Fatalf("A = 0x%02X, want 0xC7", c.A)
	}
	if c.Halted || c.Stopped {
		t.Fatalf("Halted=%v Stopped=%v, want running", c.Halted, c.Stopped)
	}

	c.A = 0x5A
	if err := c.Jam(nil, ports, OUT(16)); err != nil {
		t.Fatalf("Jam(OUT 16) = %v, want nil", err)
	}
	got, err := ports.OutputValue(16)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0x5A {
		t.Fatalf("OutputValue(16) = 0x%02X, want 0x5A", got)
	}
}

func TestInputOpcodesCoverAllInputPorts(t *testing.T) {
	for port := byte(0); port <= 7; port++ {
		opcode := INP(port)
		op := Decode(opcode)

		if !isInputOpcode(opcode) {
			t.Fatalf("INP(%d) opcode 0x%02X not recognized as input", port, opcode)
		}
		if inputPort(opcode) != port {
			t.Fatalf("inputPort(0x%02X) = %d, want %d", opcode, inputPort(opcode), port)
		}
		if op.Mnemonic != fmt.Sprintf("INP %d", port) {
			t.Fatalf("Decode(0x%02X).Mnemonic = %q", opcode, op.Mnemonic)
		}
		if op.Length != 1 || op.States != 8 {
			t.Fatalf("Decode(0x%02X) length/states = %d/%d, want 1/8", opcode, op.Length, op.States)
		}
	}
}

func TestOutputOpcodesCoverAllOutputPorts(t *testing.T) {
	for port := byte(8); port <= 31; port++ {
		opcode := OUT(port)
		op := Decode(opcode)

		if !isOutputOpcode(opcode) {
			t.Fatalf("OUT(%d) opcode 0x%02X not recognized as output", port, opcode)
		}
		if outputPort(opcode) != port {
			t.Fatalf("outputPort(0x%02X) = %d, want %d", opcode, outputPort(opcode), port)
		}
		if op.Mnemonic != fmt.Sprintf("OUT %d", port) {
			t.Fatalf("Decode(0x%02X).Mnemonic = %q", opcode, op.Mnemonic)
		}
		if op.Length != 1 || op.States != 6 {
			t.Fatalf("Decode(0x%02X) length/states = %d/%d, want 1/6", opcode, op.Length, op.States)
		}
	}
}
