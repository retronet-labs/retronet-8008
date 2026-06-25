package cpu

import "github.com/retronet-labs/retronet-hardware/bridge/i8008"

func isRotateOpcode(code byte) bool {
	switch code {
	case 0x02, 0x0A, 0x12, 0x1A:
		return true
	default:
		return false
	}
}

func executeRotate(c *CPU8008, _ Memory, _ IO, inst Instruction) error {
	// Le rotazioni sono delegate allo shifter a gate di RetroNet Logic, tramite
	// l'adattatore i8008. Toccano solo il flag Carry.
	switch inst.Opcode.Code {
	case RLC():
		c.A, c.Carry = i8008.RotateLeftCircular(c.A)
	case RRC():
		c.A, c.Carry = i8008.RotateRightCircular(c.A)
	case RAL():
		c.A, c.Carry = i8008.RotateLeftThroughCarry(c.A, c.Carry)
	case RAR():
		c.A, c.Carry = i8008.RotateRightThroughCarry(c.A, c.Carry)
	}
	return nil
}
