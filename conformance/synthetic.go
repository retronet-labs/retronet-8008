package conformance

import (
	"fmt"

	"retronet-8008/cpu"
	"retronet-8008/machine"
)

// SyntheticSuite restituisce casi indipendenti da ROM e profili storici.
func SyntheticSuite() []Case {
	return []Case{
		{
			Name:    "load-move",
			Program: []byte{cpu.LI(cpu.RegA), 0x2A, cpu.L(cpu.RegB, cpu.RegA), cpu.HLT()},
			Verify: func(ctx *Context, run machine.DebugRunResult) error {
				if ctx.CPU.A != 0x2A || ctx.CPU.B != 0x2A || run.Reason != machine.DebugStoppedCPU {
					return fmt.Errorf("A=0x%02X B=0x%02X stop=%s", ctx.CPU.A, ctx.CPU.B, run.Reason)
				}
				return nil
			},
		},
		{
			Name:    "alu-flags",
			Program: []byte{cpu.LI(cpu.RegA), 0xFF, cpu.ADI(), 0x01, cpu.HLT()},
			Verify: func(ctx *Context, _ machine.DebugRunResult) error {
				if ctx.CPU.A != 0 || !ctx.CPU.Carry || !ctx.CPU.Zero || ctx.CPU.Sign || !ctx.CPU.Parity {
					return fmt.Errorf("A=0x%02X C=%v Z=%v S=%v P=%v", ctx.CPU.A, ctx.CPU.Carry, ctx.CPU.Zero, ctx.CPU.Sign, ctx.CPU.Parity)
				}
				return nil
			},
		},
		{
			Name: "memory-indirect",
			Program: []byte{
				cpu.LI(cpu.RegH), 0x01,
				cpu.LI(cpu.RegL), 0x00,
				cpu.LI(cpu.RegM), 0xA5,
				cpu.L(cpu.RegA, cpu.RegM),
				cpu.HLT(),
			},
			Verify: func(ctx *Context, _ machine.DebugRunResult) error {
				if ctx.CPU.A != 0xA5 || ctx.Memory.Read(0x0100) != 0xA5 {
					return fmt.Errorf("A=0x%02X M=0x%02X", ctx.CPU.A, ctx.Memory.Read(0x0100))
				}
				return nil
			},
		},
		{
			Name: "call-return",
			Program: []byte{
				cpu.CAL(), 0x06, 0x00,
				cpu.HLT(), 0x00, 0x00,
				cpu.LI(cpu.RegA), 0x42,
				cpu.RET(),
			},
			Verify: func(ctx *Context, _ machine.DebugRunResult) error {
				if ctx.CPU.A != 0x42 || ctx.CPU.SP != 0 || ctx.CPU.PC != 0x0004 {
					return fmt.Errorf("A=0x%02X SP=%d PC=0x%04X", ctx.CPU.A, ctx.CPU.SP, ctx.CPU.PC)
				}
				return nil
			},
		},
		{
			Name:    "conditional-jump-taken-timing",
			Program: []byte{cpu.JF(cpu.CondCarry), 0x06, 0x00, cpu.HLT(), 0x00, 0x00, cpu.HLT()},
			Verify: func(ctx *Context, _ machine.DebugRunResult) error {
				if ctx.CPU.StateCount != 26 {
					return fmt.Errorf("StateCount=%d want=26", ctx.CPU.StateCount)
				}
				return nil
			},
		},
		{
			Name:    "conditional-jump-not-taken-timing",
			Program: []byte{cpu.JF(cpu.CondCarry), 0x06, 0x00, cpu.HLT(), 0x00, 0x00, cpu.HLT()},
			Setup: func(ctx *Context) error {
				ctx.CPU.Carry = true
				return nil
			},
			Verify: func(ctx *Context, _ machine.DebugRunResult) error {
				if ctx.CPU.StateCount != 24 {
					return fmt.Errorf("StateCount=%d want=24", ctx.CPU.StateCount)
				}
				return nil
			},
		},
		{
			Name:    "rotate-carry",
			Program: []byte{cpu.LI(cpu.RegA), 0x81, cpu.RLC(), cpu.HLT()},
			Verify: func(ctx *Context, _ machine.DebugRunResult) error {
				if ctx.CPU.A != 0x03 || !ctx.CPU.Carry {
					return fmt.Errorf("A=0x%02X C=%v", ctx.CPU.A, ctx.CPU.Carry)
				}
				return nil
			},
		},
		{
			Name:    "io-echo",
			Program: []byte{cpu.INP(0), cpu.OUT(8), cpu.HLT()},
			Setup: func(ctx *Context) error {
				return ctx.IO.SetInput(0, 0x5A)
			},
			Verify: func(ctx *Context, _ machine.DebugRunResult) error {
				output, err := ctx.IO.OutputValue(8)
				if err != nil {
					return err
				}
				if ctx.CPU.A != 0x5A || output != 0x5A {
					return fmt.Errorf("A=0x%02X OUT8=0x%02X", ctx.CPU.A, output)
				}
				return nil
			},
		},
		{
			Name:      "stack-wrap",
			Program:   restartRing(),
			StepLimit: 8,
			Verify: func(ctx *Context, run machine.DebugRunResult) error {
				if run.Reason != machine.DebugStoppedLimit || ctx.CPU.SP != 0 {
					return fmt.Errorf("stop=%s SP=%d", run.Reason, ctx.CPU.SP)
				}
				return nil
			},
		},
		{
			Name:    "interrupt-rst",
			Program: interruptProgram(),
			Setup: func(ctx *Context) error {
				return ctx.Panel.RequestInterrupt(cpu.RST(1))
			},
			Verify: func(ctx *Context, run machine.DebugRunResult) error {
				if run.Steps != 2 || ctx.CPU.SP != 1 || ctx.CPU.Stack[0] != 0 {
					return fmt.Errorf("steps=%d SP=%d return=0x%04X", run.Steps, ctx.CPU.SP, ctx.CPU.Stack[0])
				}
				return nil
			},
		},
		{
			Name:    "ready-wait",
			Program: []byte{cpu.NOP()},
			Setup: func(ctx *Context) error {
				ctx.Panel.SetReady(false)
				return nil
			},
			Verify: func(ctx *Context, run machine.DebugRunResult) error {
				if run.Reason != machine.DebugStoppedWaiting || run.Steps != 0 || ctx.CPU.WaitStateCount != 1 {
					return fmt.Errorf("stop=%s steps=%d waits=%d", run.Reason, run.Steps, ctx.CPU.WaitStateCount)
				}
				return nil
			},
		},
	}
}

func restartRing() []byte {
	program := make([]byte, 64)
	for vector := byte(0); vector < 8; vector++ {
		program[int(vector)*8] = cpu.RST((vector + 1) & 0x07)
	}
	return program
}

func interruptProgram() []byte {
	program := make([]byte, 9)
	program[0] = cpu.NOP()
	program[8] = cpu.HLT()
	return program
}
