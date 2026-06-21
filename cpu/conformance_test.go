package cpu

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

var undefined8008Opcodes = map[byte]struct{}{
	0x22: {}, 0x2A: {}, 0x32: {}, 0x38: {}, 0x39: {}, 0x3A: {},
}

func TestIntel8008OpcodeMatrix(t *testing.T) {
	undefinedCount := 0
	for raw := 0; raw <= 0xFF; raw++ {
		code := byte(raw)
		op := Decode(code)
		wantLength, wantMinStates, wantStates, wantCycles := referenceOpcodeMetadata(code)

		if op.Length != wantLength || op.MinStates != wantMinStates || op.States != wantStates {
			t.Errorf("opcode 0x%02X metadata length/states = %d/%d..%d, want %d/%d..%d", code, op.Length, op.MinStates, op.States, wantLength, wantMinStates, wantStates)
		}
		gotCycles := op.MachineCycles()
		if !slices.Equal(gotCycles, wantCycles) {
			t.Errorf("opcode 0x%02X cycles = %v, want %v", code, gotCycles, wantCycles)
		}

		c := &CPU8008{SP: 1, Stack: [8]uint16{0x0123}}
		c.setPC(0)
		mem := NewFlatMemory()
		mem.Write(0, code)
		mem.Write(1, 0x34)
		mem.Write(2, 0x12)
		err := c.Step(mem, NewPorts())

		_, undefined := undefined8008Opcodes[code]
		if undefined {
			undefinedCount++
			if !errors.Is(err, ErrUnimplementedOpcode) {
				t.Errorf("undefined opcode 0x%02X Step error = %v, want ErrUnimplementedOpcode", code, err)
			}
			if !strings.HasPrefix(op.Mnemonic, "???") {
				t.Errorf("undefined opcode 0x%02X mnemonic = %q, want ???", code, op.Mnemonic)
			}
			continue
		}

		if err != nil {
			t.Errorf("defined opcode 0x%02X (%s) Step = %v", code, op.Mnemonic, err)
		}
		if strings.HasPrefix(op.Mnemonic, "???") {
			t.Errorf("defined opcode 0x%02X has placeholder mnemonic %q", code, op.Mnemonic)
		}
		if c.InstructionCount != 1 {
			t.Errorf("defined opcode 0x%02X instruction count = %d, want 1", code, c.InstructionCount)
		}
	}
	if undefinedCount != 6 {
		t.Fatalf("undefined opcode count = %d, want 6", undefinedCount)
	}
}

func TestALUDifferentialExhaustive(t *testing.T) {
	for group := byte(0); group < 8; group++ {
		for a := 0; a <= 0xFF; a++ {
			for operand := 0; operand <= 0xFF; operand++ {
				for _, carryIn := range []bool{false, true} {
					c := CPU8008{A: byte(a), Carry: carryIn}
					c.executeALU(group, byte(operand))
					want := referenceALU(group, byte(a), byte(operand), carryIn)
					if c.A != want.a || c.Carry != want.carry || c.Zero != want.zero || c.Sign != want.sign || c.Parity != want.parity {
						t.Fatalf("group=%d A=0x%02X operand=0x%02X carryIn=%v: got A=0x%02X C=%v Z=%v S=%v P=%v, want A=0x%02X C=%v Z=%v S=%v P=%v", group, a, operand, carryIn, c.A, c.Carry, c.Zero, c.Sign, c.Parity, want.a, want.carry, want.zero, want.sign, want.parity)
					}
				}
			}
		}
	}
}

func TestALUOpcodeMatrixAgainstReference(t *testing.T) {
	for group := byte(0); group < 8; group++ {
		for src := Register(0); src <= RegM; src++ {
			c := &CPU8008{A: 0x81, B: 0x11, C: 0x22, D: 0x33, E: 0x44, H: 0x12, L: 0x34, Carry: true}
			mem := NewFlatMemory()
			mem.Write(0, 0x80|(group<<3)|byte(src))
			mem.Write(0x1234, 0x66)
			operand := referenceRegisterValue(*c, mem.Read(0x1234), src)
			want := referenceALU(group, c.A, operand, c.Carry)

			if err := c.Step(mem, nil); err != nil {
				t.Fatalf("group=%d src=%d Step = %v", group, src, err)
			}
			assertALUReference(t, *c, want, fmt.Sprintf("group=%d src=%d", group, src))
		}

		c := &CPU8008{A: 0x81, Carry: true}
		mem := NewFlatMemory()
		mem.Write(0, 0x04|(group<<3))
		mem.Write(1, 0x66)
		want := referenceALU(group, c.A, 0x66, c.Carry)
		if err := c.Step(mem, nil); err != nil {
			t.Fatalf("immediate group=%d Step = %v", group, err)
		}
		assertALUReference(t, *c, want, fmt.Sprintf("immediate group=%d", group))
	}
}

func TestRotateDifferentialExhaustive(t *testing.T) {
	opcodes := []byte{0x02, 0x0A, 0x12, 0x1A}
	for _, code := range opcodes {
		for value := 0; value <= 0xFF; value++ {
			for _, carryIn := range []bool{false, true} {
				c := CPU8008{A: byte(value), Carry: carryIn, Zero: true, Sign: false, Parity: true}
				wantA, wantCarry := referenceRotate(code, byte(value), carryIn)
				if err := executeRotate(&c, nil, nil, Instruction{Opcode: Decode(code)}); err != nil {
					t.Fatal(err)
				}
				if c.A != wantA || c.Carry != wantCarry || !c.Zero || c.Sign || !c.Parity {
					t.Fatalf("opcode=0x%02X A=0x%02X carryIn=%v: got A=0x%02X C=%v Z=%v S=%v P=%v, want A=0x%02X C=%v with ZSP unchanged", code, value, carryIn, c.A, c.Carry, c.Zero, c.Sign, c.Parity, wantA, wantCarry)
				}
			}
		}
	}
}

func TestIncrementDecrementDifferentialExhaustive(t *testing.T) {
	registers := []Register{RegB, RegC, RegD, RegE, RegH, RegL}
	for _, reg := range registers {
		for value := 0; value <= 0xFF; value++ {
			for _, carryIn := range []bool{false, true} {
				for _, increment := range []bool{false, true} {
					c := CPU8008{Carry: carryIn}
					setReferenceRegister(&c, reg, byte(value))
					code := DCR(reg)
					wantValue := byte(value - 1)
					if increment {
						code = INR(reg)
						wantValue = byte(value + 1)
					}
					inst := Instruction{Opcode: Decode(code)}
					var err error
					if increment {
						err = executeIncrement(&c, nil, nil, inst)
					} else {
						err = executeDecrement(&c, nil, nil, inst)
					}
					if err != nil {
						t.Fatal(err)
					}
					zero, sign, parity := referenceZSP(wantValue)
					if got := referenceRegisterValue(c, 0, reg); got != wantValue || c.Carry != carryIn || c.Zero != zero || c.Sign != sign || c.Parity != parity {
						t.Fatalf("reg=%d value=0x%02X increment=%v carryIn=%v: got value=0x%02X C=%v Z=%v S=%v P=%v", reg, value, increment, carryIn, got, c.Carry, c.Zero, c.Sign, c.Parity)
					}
				}
			}
		}
	}
}

type referenceALUResult struct {
	a                         byte
	carry, zero, sign, parity bool
}

func referenceALU(group, a, operand byte, carryIn bool) referenceALUResult {
	result := referenceALUResult{a: a}
	flagValue := a
	switch group {
	case 0, 1:
		sum := uint16(a) + uint16(operand)
		if group == 1 && carryIn {
			sum++
		}
		result.a = byte(sum)
		flagValue = result.a
		result.carry = sum > 0xFF
	case 2, 3, 7:
		difference := int(a) - int(operand)
		if group == 3 && carryIn {
			difference--
		}
		flagValue = byte(difference)
		result.carry = difference < 0
		if group != 7 {
			result.a = flagValue
		}
	case 4:
		result.a = a & operand
		flagValue = result.a
	case 5:
		result.a = a ^ operand
		flagValue = result.a
	case 6:
		result.a = a | operand
		flagValue = result.a
	}
	result.zero, result.sign, result.parity = referenceZSP(flagValue)
	return result
}

func referenceZSP(value byte) (zero, sign, parity bool) {
	zero = value == 0
	sign = value&0x80 != 0
	parity = true
	for bit := byte(1); bit != 0; bit <<= 1 {
		if value&bit != 0 {
			parity = !parity
		}
	}
	return zero, sign, parity
}

func referenceRotate(code, value byte, carryIn bool) (byte, bool) {
	switch code {
	case 0x02:
		return value<<1 | value>>7, value&0x80 != 0
	case 0x0A:
		return value>>1 | value<<7, value&0x01 != 0
	case 0x12:
		carry := byte(0)
		if carryIn {
			carry = 1
		}
		return value<<1 | carry, value&0x80 != 0
	default:
		carry := byte(0)
		if carryIn {
			carry = 0x80
		}
		return value>>1 | carry, value&0x01 != 0
	}
}

func referenceOpcodeMetadata(code byte) (length, minStates, states byte, cycles []MachineCycle) {
	length = 1
	switch code & 0xC7 {
	case 0x04, 0x06:
		length = 2
	case 0x40, 0x42, 0x44, 0x46:
		length = 3
	}

	minStates, states = 5, 5
	cycles = []MachineCycle{CyclePCI}
	if code&0xC7 == 0x40 || code&0xC7 == 0x42 {
		return length, 9, 11, []MachineCycle{CyclePCI, CyclePCR, CyclePCR}
	}
	if code&0xC7 == 0x03 {
		return length, 3, 5, cycles
	}
	if code == 0x00 || code == 0x01 || code == 0xFF {
		return length, 4, 4, cycles
	}
	if code == 0x3E {
		return length, 9, 9, []MachineCycle{CyclePCI, CyclePCR, CyclePCW}
	}
	if code&0xC0 == 0xC0 && code != 0xFF && (code>>3)&0x07 == 0x07 {
		return length, 7, 7, []MachineCycle{CyclePCI, CyclePCW}
	}
	if code&0xC0 == 0xC0 && code&0x07 == 0x07 {
		return length, 8, 8, []MachineCycle{CyclePCI, CyclePCR}
	}
	if code&0xC0 == 0x80 && code&0x07 == 0x07 {
		return length, 8, 8, []MachineCycle{CyclePCI, CyclePCR}
	}
	if code&0xF0 == 0x40 && code&1 == 1 {
		return length, 8, 8, []MachineCycle{CyclePCI, CyclePCC}
	}
	if code&0xC0 == 0x40 && code&0x30 != 0 && code&1 == 1 {
		return length, 6, 6, []MachineCycle{CyclePCI, CyclePCC}
	}
	if length == 2 {
		return length, 8, 8, []MachineCycle{CyclePCI, CyclePCR}
	}
	if length == 3 {
		return length, 11, 11, []MachineCycle{CyclePCI, CyclePCR, CyclePCR}
	}
	return length, minStates, states, cycles
}

func assertALUReference(t *testing.T, c CPU8008, want referenceALUResult, context string) {
	t.Helper()
	if c.A != want.a || c.Carry != want.carry || c.Zero != want.zero || c.Sign != want.sign || c.Parity != want.parity {
		t.Fatalf("%s: got A=0x%02X C=%v Z=%v S=%v P=%v, want A=0x%02X C=%v Z=%v S=%v P=%v", context, c.A, c.Carry, c.Zero, c.Sign, c.Parity, want.a, want.carry, want.zero, want.sign, want.parity)
	}
}

func referenceRegisterValue(c CPU8008, memoryValue byte, reg Register) byte {
	switch reg {
	case RegA:
		return c.A
	case RegB:
		return c.B
	case RegC:
		return c.C
	case RegD:
		return c.D
	case RegE:
		return c.E
	case RegH:
		return c.H
	case RegL:
		return c.L
	default:
		return memoryValue
	}
}

func setReferenceRegister(c *CPU8008, reg Register, value byte) {
	switch reg {
	case RegB:
		c.B = value
	case RegC:
		c.C = value
	case RegD:
		c.D = value
	case RegE:
		c.E = value
	case RegH:
		c.H = value
	case RegL:
		c.L = value
	}
}
