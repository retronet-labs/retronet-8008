package machine

import (
	"errors"
	"testing"

	"retronet-8008/cpu"
)

func TestPeripheralBusRoutesAndListsBindings(t *testing.T) {
	ioBus := NewCallbackIO()
	bus, err := NewPeripheralBus(ioBus)
	if err != nil {
		t.Fatal(err)
	}
	var output byte
	err = bus.Attach(PeripheralBinding{
		Name: "device",
		Inputs: []PeripheralInput{{Port: 2, Handler: func(byte, byte) byte {
			return 0x5A
		}}},
		Outputs: []PeripheralOutput{{Port: 10, Handler: func(_ byte, value byte) {
			output = value
		}}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := ioBus.Input(2); got != 0x5A {
		t.Fatalf("Input(2) = 0x%02X, want 0x5A", got)
	}
	ioBus.Output(10, 0xA5)
	if output != 0xA5 {
		t.Fatalf("output = 0x%02X, want 0xA5", output)
	}
	bindings := bus.Bindings()
	if len(bindings) != 2 || bindings[0].Port != 2 || bindings[1].Port != 10 {
		t.Fatalf("Bindings = %+v", bindings)
	}
}

func TestPeripheralBusRejectsConflictWithoutPartialAttach(t *testing.T) {
	ioBus := NewCallbackIO()
	bus, err := NewPeripheralBus(ioBus)
	if err != nil {
		t.Fatal(err)
	}
	if err := bus.Attach(PeripheralBinding{
		Name:   "first",
		Inputs: []PeripheralInput{{Port: 1, Handler: func(byte, byte) byte { return 1 }}},
	}); err != nil {
		t.Fatal(err)
	}

	err = bus.Attach(PeripheralBinding{
		Name: "second",
		Inputs: []PeripheralInput{
			{Port: 2, Handler: func(byte, byte) byte { return 2 }},
			{Port: 1, Handler: func(byte, byte) byte { return 3 }},
		},
	})
	if !errors.Is(err, ErrPortInUse) {
		t.Fatalf("Attach conflict = %v, want ErrPortInUse", err)
	}
	if got := ioBus.Input(2); got != 0 {
		t.Fatalf("Input(2) = %d, partial binding leaked", got)
	}
}

func TestPeripheralBusDetachReleasesPorts(t *testing.T) {
	ioBus := NewCallbackIO()
	bus, err := NewPeripheralBus(ioBus)
	if err != nil {
		t.Fatal(err)
	}
	if err := bus.Attach(PeripheralBinding{
		Name:   "device",
		Inputs: []PeripheralInput{{Port: 0, Handler: func(byte, byte) byte { return 0xFF }}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := bus.Detach("device"); err != nil {
		t.Fatal(err)
	}
	if got := ioBus.Input(0); got != 0 {
		t.Fatalf("Input after detach = 0x%02X, want latch 0", got)
	}
	if err := bus.Detach("device"); !errors.Is(err, ErrPeripheralNotFound) {
		t.Fatalf("second Detach = %v, want ErrPeripheralNotFound", err)
	}
}

func TestRegisterPeripheralLoopback(t *testing.T) {
	ioBus := NewCallbackIO()
	bus, err := NewPeripheralBus(ioBus)
	if err != nil {
		t.Fatal(err)
	}
	register := NewRegisterPeripheral(0x11)
	if err := register.Attach(bus, "loopback", 1, 9); err != nil {
		t.Fatal(err)
	}

	if got := ioBus.Input(1); got != 0x11 {
		t.Fatalf("initial Input = 0x%02X, want 0x11", got)
	}
	ioBus.Output(9, 0x77)
	if got := ioBus.Input(1); got != 0x77 || register.Value() != 0x77 {
		t.Fatalf("loopback = 0x%02X register=0x%02X", got, register.Value())
	}
}

func TestPeripheralBusValidatesPorts(t *testing.T) {
	bus, err := NewPeripheralBus(NewCallbackIO())
	if err != nil {
		t.Fatal(err)
	}
	err = bus.Attach(PeripheralBinding{
		Name:   "bad",
		Inputs: []PeripheralInput{{Port: 8, Handler: func(byte, byte) byte { return 0 }}},
	})
	if !errors.Is(err, cpu.ErrInvalidInputPort) {
		t.Fatalf("Attach invalid input = %v", err)
	}
}
