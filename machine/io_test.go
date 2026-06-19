package machine

import (
	"errors"
	"testing"

	"retronet-8008/cpu"
)

func TestCallbackIOInputUsesLatchedValue(t *testing.T) {
	ioBus := NewCallbackIO()

	if err := ioBus.SetInput(0, 0x5A); err != nil {
		t.Fatalf("SetInput = %v, want nil", err)
	}

	if got := ioBus.Input(0); got != 0x5A {
		t.Fatalf("Input(0) = 0x%02X, want 0x5A", got)
	}
}

func TestCallbackIOInputUsesCallback(t *testing.T) {
	ioBus := NewCallbackIO()
	if err := ioBus.SetInput(1, 0x10); err != nil {
		t.Fatalf("SetInput = %v, want nil", err)
	}
	if err := ioBus.OnInput(1, func(port byte, latched byte) byte {
		if port != 1 {
			t.Fatalf("callback port = %d, want 1", port)
		}
		return latched + 1
	}); err != nil {
		t.Fatalf("OnInput = %v, want nil", err)
	}

	if got := ioBus.Input(1); got != 0x11 {
		t.Fatalf("Input(1) = 0x%02X, want 0x11", got)
	}
}

func TestCallbackIOOutputUsesCallbackAndLatch(t *testing.T) {
	ioBus := NewCallbackIO()
	var gotPort, gotValue byte
	if err := ioBus.OnOutput(8, func(port byte, value byte) {
		gotPort = port
		gotValue = value
	}); err != nil {
		t.Fatalf("OnOutput = %v, want nil", err)
	}

	ioBus.Output(8, 0xA5)

	if gotPort != 8 || gotValue != 0xA5 {
		t.Fatalf("callback = port %d value 0x%02X, want port 8 value 0xA5", gotPort, gotValue)
	}
	latched, err := ioBus.OutputValue(8)
	if err != nil {
		t.Fatalf("OutputValue = %v, want nil", err)
	}
	if latched != 0xA5 {
		t.Fatalf("OutputValue(8) = 0x%02X, want 0xA5", latched)
	}
}

func TestCallbackIORejectsInvalidPorts(t *testing.T) {
	ioBus := NewCallbackIO()

	if err := ioBus.SetInput(8, 0); !errors.Is(err, cpu.ErrInvalidInputPort) {
		t.Fatalf("SetInput invalid = %v, want ErrInvalidInputPort", err)
	}
	if err := ioBus.OnOutput(7, nil); !errors.Is(err, cpu.ErrInvalidOutputPort) {
		t.Fatalf("OnOutput invalid = %v, want ErrInvalidOutputPort", err)
	}
	if err := ioBus.ObserveInput(8, nil); !errors.Is(err, cpu.ErrInvalidInputPort) {
		t.Fatalf("ObserveInput invalid = %v, want ErrInvalidInputPort", err)
	}
}

func TestCallbackIOObserversSeeCallbackValues(t *testing.T) {
	ioBus := NewCallbackIO()
	if err := ioBus.OnInput(0, func(_ byte, value byte) byte { return value + 1 }); err != nil {
		t.Fatal(err)
	}
	var observedInput, observedOutput byte
	if err := ioBus.ObserveInput(0, func(_ byte, value byte) { observedInput = value }); err != nil {
		t.Fatal(err)
	}
	if err := ioBus.ObserveOutput(8, func(_ byte, value byte) { observedOutput = value }); err != nil {
		t.Fatal(err)
	}

	if got := ioBus.Input(0); got != 1 || observedInput != 1 {
		t.Fatalf("Input = %d observed = %d, want 1, 1", got, observedInput)
	}
	ioBus.Output(8, 0xA5)
	if observedOutput != 0xA5 {
		t.Fatalf("observed output = 0x%02X, want 0xA5", observedOutput)
	}
}

func TestProfileNewIOReturnsCallbackBus(t *testing.T) {
	profile, ok := Lookup("intellec-8")
	if !ok {
		t.Fatal("Lookup(intellec-8) = false")
	}

	var ioBus cpu.IO = profile.NewIO()
	if ioBus == nil {
		t.Fatal("Profile.NewIO() = nil")
	}
}
