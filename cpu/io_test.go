package cpu

import (
	"errors"
	"testing"
)

func TestPortsInput(t *testing.T) {
	ports := NewPorts()

	if err := ports.SetInput(7, 0xA5); err != nil {
		t.Fatal(err)
	}

	if got := ports.Input(7); got != 0xA5 {
		t.Fatalf("Input(7) = 0x%02X, want 0xA5", got)
	}
}

func TestPortsOutput(t *testing.T) {
	ports := NewPorts()

	ports.Output(8, 0x11)
	ports.Output(31, 0x22)

	if got, err := ports.OutputValue(8); err != nil || got != 0x11 {
		t.Fatalf("OutputValue(8) = 0x%02X, %v; want 0x11, nil", got, err)
	}
	if got, err := ports.OutputValue(31); err != nil || got != 0x22 {
		t.Fatalf("OutputValue(31) = 0x%02X, %v; want 0x22, nil", got, err)
	}
}

func TestPortsImplementIO(t *testing.T) {
	var bus IO = NewPorts()
	bus.Output(8, 0x77)

	ports := bus.(*Ports)
	got, err := ports.OutputValue(8)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0x77 {
		t.Fatalf("OutputValue(8) = 0x%02X, want 0x77", got)
	}
}

func TestInputPortValidation(t *testing.T) {
	if err := ValidateInputPort(0); err != nil {
		t.Fatalf("ValidateInputPort(0) = %v, want nil", err)
	}
	if err := ValidateInputPort(7); err != nil {
		t.Fatalf("ValidateInputPort(7) = %v, want nil", err)
	}
	if err := ValidateInputPort(8); !errors.Is(err, ErrInvalidInputPort) {
		t.Fatalf("ValidateInputPort(8) = %v, want ErrInvalidInputPort", err)
	}
}

func TestOutputPortValidation(t *testing.T) {
	if err := ValidateOutputPort(8); err != nil {
		t.Fatalf("ValidateOutputPort(8) = %v, want nil", err)
	}
	if err := ValidateOutputPort(31); err != nil {
		t.Fatalf("ValidateOutputPort(31) = %v, want nil", err)
	}
	if err := ValidateOutputPort(7); !errors.Is(err, ErrInvalidOutputPort) {
		t.Fatalf("ValidateOutputPort(7) = %v, want ErrInvalidOutputPort", err)
	}
	if err := ValidateOutputPort(32); !errors.Is(err, ErrInvalidOutputPort) {
		t.Fatalf("ValidateOutputPort(32) = %v, want ErrInvalidOutputPort", err)
	}
}

func TestInvalidPortsAreNonMutating(t *testing.T) {
	ports := NewPorts()

	if err := ports.SetInput(8, 0xFF); !errors.Is(err, ErrInvalidInputPort) {
		t.Fatalf("SetInput(8) = %v, want ErrInvalidInputPort", err)
	}
	if got := ports.Input(8); got != 0 {
		t.Fatalf("Input(8) = 0x%02X, want 0", got)
	}

	ports.Output(7, 0xAA)
	if got, err := ports.OutputValue(8); err != nil || got != 0 {
		t.Fatalf("OutputValue(8) = 0x%02X, %v; want 0, nil", got, err)
	}
}
