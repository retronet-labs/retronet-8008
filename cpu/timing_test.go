package cpu

import "testing"

func TestOpcodeTimingMetadata(t *testing.T) {
	tests := []struct {
		name      string
		code      byte
		minStates byte
		states    byte
		cycles    []MachineCycle
	}{
		{"register move", L(RegA, RegB), 5, 5, []MachineCycle{CyclePCI}},
		{"load from memory", L(RegA, RegM), 8, 8, []MachineCycle{CyclePCI, CyclePCR}},
		{"store to memory", L(RegM, RegA), 7, 7, []MachineCycle{CyclePCI, CyclePCW}},
		{"immediate", LI(RegA), 8, 8, []MachineCycle{CyclePCI, CyclePCR}},
		{"immediate store", LI(RegM), 9, 9, []MachineCycle{CyclePCI, CyclePCR, CyclePCW}},
		{"conditional jump", JF(CondCarry), 9, 11, []MachineCycle{CyclePCI, CyclePCR, CyclePCR}},
		{"conditional return", RF(CondCarry), 3, 5, []MachineCycle{CyclePCI}},
		{"input", INP(0), 8, 8, []MachineCycle{CyclePCI, CyclePCC}},
		{"output", OUT(8), 6, 6, []MachineCycle{CyclePCI, CyclePCC}},
		{"halt", HLT(), 4, 4, []MachineCycle{CyclePCI}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := Decode(tt.code)
			if op.MinStates != tt.minStates || op.States != tt.states {
				t.Fatalf("state range = %d..%d, want %d..%d", op.MinStates, op.States, tt.minStates, tt.states)
			}
			gotCycles := op.MachineCycles()
			if len(gotCycles) != len(tt.cycles) {
				t.Fatalf("cycles = %v, want %v", gotCycles, tt.cycles)
			}
			for i := range gotCycles {
				if gotCycles[i] != tt.cycles[i] {
					t.Fatalf("cycles = %v, want %v", gotCycles, tt.cycles)
				}
			}
		})
	}
}

func TestConditionalTimingUsesPreInstructionFlags(t *testing.T) {
	tests := []struct {
		name       string
		carry      bool
		wantTaken  bool
		wantStates byte
	}{
		{"false condition taken", false, true, 11},
		{"false condition not taken", true, false, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRunningCPU(t)
			c.Carry = tt.carry
			mem := NewFlatMemory()
			mem.Write(0, JF(CondCarry))
			mem.Write(1, 0x10)
			mem.Write(2, 0x00)

			if err := c.Step(mem, nil); err != nil {
				t.Fatal(err)
			}
			if c.LastTiming.Taken != tt.wantTaken || c.LastTiming.States != tt.wantStates {
				t.Fatalf("timing = %+v, want taken=%v states=%d", c.LastTiming, tt.wantTaken, tt.wantStates)
			}
		})
	}
}

func TestTimingCountersIncludeStepAndJam(t *testing.T) {
	c := NewCPU8008()
	mem := NewFlatMemory()
	mem.Write(0, HLT())

	if err := c.Jam(mem, nil, NOP()); err != nil {
		t.Fatal(err)
	}
	if err := c.Step(mem, nil); err != nil {
		t.Fatal(err)
	}
	if c.InstructionCount != 2 || c.StateCount != 9 {
		t.Fatalf("counts = instructions %d states %d, want 2 and 9", c.InstructionCount, c.StateCount)
	}
	if c.LastTiming.States != 4 {
		t.Fatalf("last states = %d, want 4", c.LastTiming.States)
	}

	c.Reset()
	if c.InstructionCount != 0 || c.StateCount != 0 || c.WaitStateCount != 0 || c.LastTiming.States != 0 {
		t.Fatalf("timing after Reset = count %d states %d last %+v", c.InstructionCount, c.StateCount, c.LastTiming)
	}
}

func TestRecordWaitStateIsAttachedToNextInstruction(t *testing.T) {
	c := NewCPU8008()
	c.RecordWaitState()
	if err := c.Jam(NewFlatMemory(), nil, NOP()); err != nil {
		t.Fatal(err)
	}

	if c.StateCount != 6 || c.WaitStateCount != 1 || c.LastTiming.WaitStates != 1 {
		t.Fatalf("timing = states %d waits %d last %+v", c.StateCount, c.WaitStateCount, c.LastTiming)
	}
}
