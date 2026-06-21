package cpu

import (
	"errors"
	"testing"
)

func FuzzDecodeDisassemble(f *testing.F) {
	for code := 0; code <= 0xFF; code++ {
		f.Add([]byte{byte(code), 0x34, 0x12, byte(code), ^byte(code)})
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 5 {
			return
		}
		pc := (uint16(data[3])<<8 | uint16(data[4])) & AddressMask
		mem := NewFlatMemory()
		for i := 0; i < 3; i++ {
			mem.Write(pc+uint16(i), data[i])
		}

		d, err := Disassemble(mem, pc)
		if err != nil {
			t.Fatal(err)
		}
		op := Decode(data[0])
		if d.PC != pc || d.Opcode.Code != data[0] || d.Length != op.Length {
			t.Fatalf("disassembly header = PC 0x%04X opcode 0x%02X length %d, want PC 0x%04X opcode 0x%02X length %d", d.PC, d.Opcode.Code, d.Length, pc, data[0], op.Length)
		}
		if d.NextPC != (pc+uint16(op.Length))&AddressMask {
			t.Fatalf("NextPC = 0x%04X, want 0x%04X", d.NextPC, (pc+uint16(op.Length))&AddressMask)
		}
		for i := byte(0); i < op.Length; i++ {
			if d.Bytes[i] != data[i] {
				t.Fatalf("byte %d = 0x%02X, want 0x%02X", i, d.Bytes[i], data[i])
			}
		}
		if d.String() == "" {
			t.Fatal("empty disassembly string")
		}
	})
}

func FuzzStepMaintainsArchitecturalBounds(f *testing.F) {
	for code := 0; code <= 0xFF; code++ {
		f.Add([]byte{byte(code), 0x34, 0x12, byte(code), ^byte(code), byte(code * 17)})
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 6 {
			return
		}
		pc := (uint16(data[3])<<8 | uint16(data[4])) & AddressMask
		c := &CPU8008{
			A: data[5], B: data[4], C: data[3], D: data[2], E: data[1],
			H: data[0], L: data[5], Carry: data[5]&1 != 0,
			Zero: data[5]&2 != 0, Sign: data[5]&4 != 0, Parity: data[5]&8 != 0,
		}
		c.setSP(data[5])
		for i := range c.Stack {
			c.setStack(uint8(i), uint16(data[(i+1)%len(data)])<<8|uint16(data[i%len(data)]))
		}
		c.setPC(pc)

		mem := NewFlatMemory()
		for i := 0; i < 3; i++ {
			mem.Write(pc+uint16(i), data[i])
		}
		err := c.Step(mem, NewPorts())
		_, undefined := undefined8008Opcodes[data[0]]
		if undefined {
			if !errors.Is(err, ErrUnimplementedOpcode) {
				t.Fatalf("undefined opcode 0x%02X Step = %v", data[0], err)
			}
		} else if err != nil {
			t.Fatalf("defined opcode 0x%02X Step = %v", data[0], err)
		}

		if c.PC > AddressMask || c.SP > 7 {
			t.Fatalf("architectural bounds violated: PC=0x%04X SP=%d", c.PC, c.SP)
		}
		for i, address := range c.Stack {
			if address > AddressMask {
				t.Fatalf("Stack[%d] = 0x%04X outside 14-bit space", i, address)
			}
		}
	})
}
