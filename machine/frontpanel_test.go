package machine

import (
	"errors"
	"testing"

	"retronet-8008/cpu"
)

func TestFrontPanelRejectsMissingComponents(t *testing.T) {
	if _, err := NewFrontPanel(nil, cpu.NewFlatMemory(), nil); !errors.Is(err, ErrNilCPU) {
		t.Fatalf("NewFrontPanel(nil CPU) = %v, want ErrNilCPU", err)
	}
	if _, err := NewFrontPanel(cpu.NewCPU8008(), nil, nil); !errors.Is(err, cpu.ErrNilMemory) {
		t.Fatalf("NewFrontPanel(nil memory) = %v, want ErrNilMemory", err)
	}
}

func TestFrontPanelSwitchesExamineDepositAndSnapshot(t *testing.T) {
	mem := cpu.NewFlatMemory()
	panel := newTestFrontPanel(t, mem, nil)
	panel.SetAddress(0x4010)
	panel.SetSwitches(0xA5)
	if err := panel.DepositSwitches(); err != nil {
		t.Fatal(err)
	}

	if got := panel.Address(); got != 0x0010 {
		t.Fatalf("Address = 0x%04X, want 0x0010", got)
	}
	if got := panel.Examine(); got != 0xA5 {
		t.Fatalf("Examine = 0x%02X, want 0xA5", got)
	}
	state := panel.Snapshot()
	if state.Switches != 0xA5 || state.Address != 0x0010 || state.Data != 0xA5 {
		t.Fatalf("Snapshot = %+v", state)
	}
}

func TestFrontPanelDepositRespectsROMProtection(t *testing.T) {
	mem, err := NewMemoryBus([]MemoryRegion{
		{Name: "memory", Start: 0x0000, End: cpu.AddressMask, Kind: MemoryKindMixed},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := mem.LoadROM(0x0020, []byte{0xAA}); err != nil {
		t.Fatal(err)
	}
	panel := newTestFrontPanel(t, mem, nil)
	panel.SetAddress(0x0020)

	if err := panel.Deposit(0x55); err != nil {
		t.Fatal(err)
	}
	if got := panel.Examine(); got != 0xAA {
		t.Fatalf("ROM after Deposit = 0x%02X, want 0xAA", got)
	}
}

func TestFrontPanelAttachesSwitchesToInput(t *testing.T) {
	ioBus := NewCallbackIO()
	panel := newTestFrontPanel(t, cpu.NewFlatMemory(), ioBus)
	panel.SetSwitches(0x5A)
	if err := panel.AttachSwitches(ioBus, 2); err != nil {
		t.Fatal(err)
	}

	if got := ioBus.Input(2); got != 0x5A {
		t.Fatalf("Input(2) = 0x%02X, want switches 0x5A", got)
	}
}

func TestFrontPanelRunStopsOnCPUHalt(t *testing.T) {
	mem := cpu.NewFlatMemory()
	mem.Write(0x0000, cpu.NOP())
	mem.Write(0x0001, cpu.HLT())
	panel := newTestFrontPanel(t, mem, nil)
	if err := panel.Jam(cpu.JMP(), 0x00, 0x00); err != nil {
		t.Fatal(err)
	}
	var observed []uint16

	result, err := panel.Run(8, func(_ uint64, state cpu.CPU8008) error {
		observed = append(observed, state.PC)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 2 || result.Reason != PanelStoppedByCPU {
		t.Fatalf("Run = %+v, want 2 steps cpu-stopped", result)
	}
	if len(observed) != 2 || observed[0] != 0 || observed[1] != 1 {
		t.Fatalf("observed PCs = %v, want [0 1]", observed)
	}
}

func TestFrontPanelRunStopsAtLimit(t *testing.T) {
	mem := cpu.NewFlatMemory()
	mem.Write(0x0000, cpu.NOP())
	mem.Write(0x0001, cpu.NOP())
	panel := newTestFrontPanel(t, mem, nil)
	if err := panel.Jam(cpu.JMP(), 0x00, 0x00); err != nil {
		t.Fatal(err)
	}

	result, err := panel.Run(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 1 || result.Reason != PanelStoppedByLimit {
		t.Fatalf("Run = %+v, want 1 step limit", result)
	}
}

func TestFrontPanelHonorsStopRequest(t *testing.T) {
	panel := newTestFrontPanel(t, cpu.NewFlatMemory(), nil)
	if err := panel.Jam(cpu.NOP()); err != nil {
		t.Fatal(err)
	}
	panel.Stop()

	result, err := panel.Run(10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 0 || result.Reason != PanelStoppedByRequest {
		t.Fatalf("Run = %+v, want requested before first step", result)
	}
}

func TestFrontPanelInterruptRST(t *testing.T) {
	panel := newTestFrontPanel(t, cpu.NewFlatMemory(), nil)
	if err := panel.Jam(cpu.JMP(), 0x34, 0x12); err != nil {
		t.Fatal(err)
	}
	if err := panel.InterruptRST(3); err != nil {
		t.Fatal(err)
	}
	state := panel.Snapshot().CPU
	if state.PC != 0x0018 || state.SP != 1 || state.Stack[0] != 0x1234 {
		t.Fatalf("CPU after RST = PC 0x%04X SP %d Stack[0] 0x%04X", state.PC, state.SP, state.Stack[0])
	}
	if err := panel.InterruptRST(8); !errors.Is(err, ErrInvalidRestartVector) {
		t.Fatalf("InterruptRST(8) = %v, want ErrInvalidRestartVector", err)
	}
}

func newTestFrontPanel(t *testing.T, mem cpu.Memory, ioBus cpu.IO) *FrontPanel {
	t.Helper()
	panel, err := NewFrontPanel(cpu.NewCPU8008(), mem, ioBus)
	if err != nil {
		t.Fatal(err)
	}
	return panel
}
