package machine

import (
	"encoding/json"
	"testing"

	"github.com/retronet-labs/retronet-8008/cpu"
)

func TestDebuggerProducesStructuredInstructionTrace(t *testing.T) {
	debugger, panel, _ := newDebugRig(t, []byte{cpu.NOP(), cpu.HLT()})
	var events []TraceEvent
	debugger.SetTraceSink(func(event TraceEvent) { events = append(events, event) })

	result, err := debugger.Run(8)
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 2 || result.Reason != DebugStoppedCPU {
		t.Fatalf("Run = %+v, want 2 cpu-stopped", result)
	}
	if len(events) != 2 || events[0].PC != 0 || events[1].PC != 1 {
		t.Fatalf("events = %+v", events)
	}
	if events[0].Timing.States != 5 || events[1].Timing.States != 4 {
		t.Fatalf("event timings = %d, %d", events[0].Timing.States, events[1].Timing.States)
	}
	if events[1].After.Halted != true || panel.Snapshot().CPU.PC != 2 {
		t.Fatalf("final state = %+v", panel.Snapshot())
	}
	if _, err := json.Marshal(events[0]); err != nil {
		t.Fatalf("json.Marshal = %v", err)
	}
}

func TestDebuggerStopsBeforePCBreakpoint(t *testing.T) {
	debugger, panel, _ := newDebugRig(t, []byte{cpu.NOP()})
	debugger.AddPCBreakpoint(0)

	result, err := debugger.Run(8)
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 0 || result.Reason != DebugStoppedBreakpoint {
		t.Fatalf("Run = %+v", result)
	}
	if panel.Snapshot().CPU.PC != 0 {
		t.Fatalf("PC = 0x%04X, want 0", panel.Snapshot().CPU.PC)
	}
}

func TestDebuggerStopsBeforeOpcodeBreakpoint(t *testing.T) {
	debugger, panel, _ := newDebugRig(t, []byte{cpu.NOP()})
	debugger.AddOpcodeBreakpoint(cpu.NOP())

	result, err := debugger.Run(8)
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 0 || result.Reason != DebugStoppedBreakpoint || panel.Snapshot().CPU.PC != 0 {
		t.Fatalf("Run = %+v PC=0x%04X", result, panel.Snapshot().CPU.PC)
	}
}

func TestDebuggerStopsAfterMemoryWatchpoint(t *testing.T) {
	debugger, panel, _ := newDebugRig(t, []byte{cpu.LI(cpu.RegM), 0xAA, cpu.HLT()})
	panel.cpu.H = 0x01
	panel.cpu.L = 0x00
	debugger.AddMemoryWatchpoint(0x0100)
	var event TraceEvent
	debugger.SetTraceSink(func(got TraceEvent) { event = got })

	result, err := debugger.Run(8)
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 1 || result.Reason != DebugStoppedWatchpoint {
		t.Fatalf("Run = %+v", result)
	}
	if len(event.MemoryWrites) != 1 || event.MemoryWrites[0].Address != 0x0100 || event.MemoryWrites[0].After != 0xAA {
		t.Fatalf("writes = %+v", event.MemoryWrites)
	}
}

func TestDebuggerStopsAfterOutputBreakpoint(t *testing.T) {
	debugger, _, _ := newDebugRig(t, []byte{cpu.LI(cpu.RegA), 0x5A, cpu.OUT(8), cpu.HLT()})
	if err := debugger.AddOutputBreakpoint(8); err != nil {
		t.Fatal(err)
	}
	var last TraceEvent
	debugger.SetTraceSink(func(event TraceEvent) { last = event })

	result, err := debugger.Run(8)
	if err != nil {
		t.Fatal(err)
	}
	if result.Steps != 2 || result.Reason != DebugStoppedIO {
		t.Fatalf("Run = %+v", result)
	}
	if len(last.IO) != 1 || last.IO[0].Port != 8 || last.IO[0].Value != 0x5A {
		t.Fatalf("IO = %+v", last.IO)
	}
}

func TestDebuggerEmitsWaitEvent(t *testing.T) {
	debugger, panel, _ := newDebugRig(t, []byte{cpu.NOP()})
	panel.SetReady(false)

	event, reason, err := debugger.Step()
	if err != nil {
		t.Fatal(err)
	}
	if reason != DebugStoppedWaiting || event.Kind != TraceWait || event.WaitCycle == nil || event.WaitCycle.Cycle != cpu.CyclePCI {
		t.Fatalf("event = %+v reason=%s", event, reason)
	}
}

func TestDebuggerTracesPendingInterrupt(t *testing.T) {
	debugger, panel, _ := newDebugRig(t, []byte{cpu.NOP()})
	if err := panel.RequestInterrupt(cpu.RST(2)); err != nil {
		t.Fatal(err)
	}

	event, _, err := debugger.Step()
	if err != nil {
		t.Fatal(err)
	}
	if !event.Interrupt || event.Opcode != cpu.RST(2) || event.After.PC != 0x0010 {
		t.Fatalf("interrupt event = %+v", event)
	}
}

func newDebugRig(t *testing.T, program []byte) (*Debugger, *FrontPanel, *ObservableMemory) {
	t.Helper()
	base := cpu.NewFlatMemory()
	memory, err := NewObservableMemory(base)
	if err != nil {
		t.Fatal(err)
	}
	if err := memory.LoadBytes(0, program); err != nil {
		t.Fatal(err)
	}
	ioBus := NewCallbackIO()
	panel, err := NewFrontPanel(cpu.NewCPU8008(), memory, ioBus)
	if err != nil {
		t.Fatal(err)
	}
	if err := panel.Jam(cpu.JMP(), 0, 0); err != nil {
		t.Fatal(err)
	}
	debugger, err := NewDebugger(panel, memory, ioBus)
	if err != nil {
		t.Fatal(err)
	}
	return debugger, panel, memory
}
