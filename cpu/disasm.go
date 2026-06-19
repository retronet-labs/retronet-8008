package cpu

import (
	"fmt"
	"strings"
)

// Disassembly descrive una istruzione letta dalla memoria con i suoi byte.
type Disassembly struct {
	PC       uint16
	Opcode   Opcode
	Bytes    [3]byte
	Length   byte
	Operand  string
	NextPC   uint16
	Operands [2]byte
}

// Disassemble legge da memoria l'istruzione all'indirizzo pc e ne restituisce
// una rappresentazione testuale minima. Gli operandi vengono letti dal contesto
// memoria e il next PC resta mascherato a 14 bit.
func Disassemble(mem Memory, pc uint16) (Disassembly, error) {
	if mem == nil {
		return Disassembly{}, ErrNilMemory
	}

	pc = addr14(pc)
	code := mem.Read(pc)
	op := Decode(code)
	d := Disassembly{
		PC:     pc,
		Opcode: op,
		Length: op.Length,
		NextPC: addr14(pc + uint16(op.Length)),
	}
	d.Bytes[0] = code

	for i := byte(1); i < op.Length; i++ {
		value := mem.Read(pc + uint16(i))
		d.Bytes[i] = value
		d.Operands[i-1] = value
	}
	d.Operand = disassemblyOperand(op.Code, d.Operands)
	return d, nil
}

// String formatta l'istruzione in forma compatta:
//
//	0000: 06 2A    LAI #0x2A
func (d Disassembly) String() string {
	bytes := make([]string, 0, d.Length)
	for i := byte(0); i < d.Length; i++ {
		bytes = append(bytes, fmt.Sprintf("%02X", d.Bytes[i]))
	}
	text := d.Opcode.Mnemonic
	if d.Operand != "" {
		text += " " + d.Operand
	}
	return fmt.Sprintf("%04X: %-8s %s", d.PC, strings.Join(bytes, " "), text)
}

func disassemblyOperand(code byte, operands [2]byte) string {
	switch {
	case code&0xC7 == 0x04:
		return fmt.Sprintf("#0x%02X", operands[0])
	case code&0xC7 == 0x06:
		return fmt.Sprintf("#0x%02X", operands[0])
	case code&0xC7 == 0x40:
		return fmt.Sprintf("0x%04X", disassemblyAddress(operands))
	case code&0xC7 == 0x42:
		return fmt.Sprintf("0x%04X", disassemblyAddress(operands))
	case code&0xC7 == 0x44:
		return fmt.Sprintf("0x%04X", disassemblyAddress(operands))
	case code&0xC7 == 0x46:
		return fmt.Sprintf("0x%04X", disassemblyAddress(operands))
	default:
		return ""
	}
}

func disassemblyAddress(operands [2]byte) uint16 {
	return (uint16(operands[1]&0x3F) << 8) | uint16(operands[0])
}
