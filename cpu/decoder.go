package cpu

import "fmt"

var decoder = buildDecoder()

// Decode restituisce i metadata associati a un opcode 8008.
func Decode(code byte) Opcode {
	return decoder[code]
}

// OpcodeTable restituisce una copia della tabella decoder completa.
func OpcodeTable() [256]Opcode {
	return decoder
}

func buildDecoder() [256]Opcode {
	var table [256]Opcode
	for i := range table {
		code := byte(i)
		execute := ExecuteFunc(unimplementedExecute)
		if isInputOpcode(code) {
			execute = executeInput
		}
		if isOutputOpcode(code) {
			execute = executeOutput
		}
		if isControlFlowOpcode(code) {
			execute = executeControlFlow
		}
		if isRotateOpcode(code) {
			execute = executeRotate
		}
		if isALUImmediateOpcode(code) {
			execute = executeALUImmediate
		}
		if isALURegisterOpcode(code) {
			execute = executeALURegister
		}
		if isIncrementOpcode(code) {
			execute = executeIncrement
		}
		if isDecrementOpcode(code) {
			execute = executeDecrement
		}
		if isLoadImmediateOpcode(code) {
			execute = executeLoadImmediate
		}
		if isLoadMoveOpcode(code) {
			execute = executeLoadMove
		}
		// HLT va valutato per ultimo: 0xFF e' un alias documentato di HLT e
		// ricade anche nel pattern load/move (L M,M); l'halt deve prevalere.
		if isHaltOpcode(code) {
			execute = executeHalt
		}
		minStates, states := stateRangeFor(code)
		cycles, cycleCount := cyclesFor(code)
		table[i] = Opcode{
			Code:       code,
			Mnemonic:   mnemonicFor(code),
			Length:     lengthFor(code),
			MinStates:  minStates,
			States:     states,
			CycleCount: cycleCount,
			Cycles:     cycles,
			Execute:    execute,
		}
	}
	return table
}

func lengthFor(code byte) byte {
	switch {
	case code&0xC7 == 0x04: // ALU immediati: 00 AAA 100
		return 2
	case code&0xC7 == 0x06: // load immediato: 00 DDD 110
		return 2
	case code&0xC7 == 0x40: // jump condizionali: 010/011 CC 000
		return 3
	case code&0xC7 == 0x42: // call condizionali: 010/011 CC 010
		return 3
	case code&0xC7 == 0x44: // JMP alias: 01 XXX 100
		return 3
	case code&0xC7 == 0x46: // CAL alias: 01 XXX 110
		return 3
	default:
		return 1
	}
}

func stateRangeFor(code byte) (byte, byte) {
	if code&0xC7 == 0x40 || code&0xC7 == 0x42 {
		return 9, 11
	}
	if code&0xC7 == 0x03 {
		return 3, 5
	}
	switch lengthFor(code) {
	case 2:
		if code == 0x3E { // LMI
			return 9, 9
		}
		return 8, 8
	case 3:
		return 11, 11
	default:
		if code == 0x00 || code == 0x01 || code == 0xFF {
			return 4, 4
		}
		if code&0xC0 == 0xC0 && (code>>3)&0x07 == 0x07 {
			return 7, 7
		}
		if code&0xC0 == 0xC0 && code&0x07 == 0x07 {
			return 8, 8
		}
		if code&0xC0 == 0x80 && code&0x07 == 0x07 {
			return 8, 8
		}
		if isInputOpcode(code) {
			return 8, 8
		}
		if isOutputOpcode(code) {
			return 6, 6
		}
		return 5, 5
	}
}

func cyclesFor(code byte) ([3]MachineCycle, byte) {
	cycles := [3]MachineCycle{CyclePCI}
	switch {
	case code == 0x3E:
		cycles[1], cycles[2] = CyclePCR, CyclePCW
		return cycles, 3
	case isLoadImmediateOpcode(code), isALUImmediateOpcode(code):
		cycles[1] = CyclePCR
		return cycles, 2
	case code&0xC0 == 0xC0 && code != 0xFF && (code>>3)&0x07 == 0x07:
		cycles[1] = CyclePCW
		return cycles, 2
	case code&0xC0 == 0xC0 && code&0x07 == 0x07:
		cycles[1] = CyclePCR
		return cycles, 2
	case code&0xC0 == 0x80 && code&0x07 == 0x07:
		cycles[1] = CyclePCR
		return cycles, 2
	case lengthFor(code) == 3:
		cycles[1], cycles[2] = CyclePCR, CyclePCR
		return cycles, 3
	case isInputOpcode(code), isOutputOpcode(code):
		cycles[1] = CyclePCC
		return cycles, 2
	default:
		return cycles, 1
	}
}

func mnemonicFor(code byte) string {
	switch {
	case code == 0x00 || code == 0x01 || code == 0xFF:
		return "HLT"
	case code == 0x02:
		return "RLC"
	case code == 0x0A:
		return "RRC"
	case code == 0x12:
		return "RAL"
	case code == 0x1A:
		return "RAR"
	case code&0xC7 == 0x04:
		return immediateALUMnemonic((code >> 3) & 0x07)
	case code&0xC7 == 0x06:
		dst := (code >> 3) & 0x07
		if dst == byte(RegM) {
			return "LMI"
		}
		return fmt.Sprintf("L%sI", registerName(dst))
	case code&0xC7 == 0x00:
		return fmt.Sprintf("IN%s", registerName((code>>3)&0x07))
	case code&0xC7 == 0x01:
		return fmt.Sprintf("DC%s", registerName((code>>3)&0x07))
	case code&0xC7 == 0x03:
		return fmt.Sprintf("R%s%s", falseTruePrefix(code), conditionName((code>>3)&0x03))
	case code&0xC7 == 0x05:
		return fmt.Sprintf("RST %d", (code>>3)&0x07)
	case code&0xC7 == 0x07:
		return "RET"
	case code&0xC7 == 0x40:
		return fmt.Sprintf("J%s%s", falseTruePrefix(code), conditionName((code>>3)&0x03))
	case code&0xC7 == 0x42:
		return fmt.Sprintf("C%s%s", falseTruePrefix(code), conditionName((code>>3)&0x03))
	case code&0xC7 == 0x44:
		return "JMP"
	case code&0xC7 == 0x46:
		return "CAL"
	case isInputOpcode(code):
		return fmt.Sprintf("INP %d", inputPort(code))
	case isOutputOpcode(code):
		return fmt.Sprintf("OUT %d", outputPort(code))
	case code&0xC0 == 0x80:
		return registerALUMnemonic((code>>3)&0x07, code&0x07)
	case code&0xC0 == 0xC0:
		dst := (code >> 3) & 0x07
		src := code & 0x07
		if dst == src {
			return "NOP"
		}
		return fmt.Sprintf("L%s%s", registerName(dst), registerName(src))
	default:
		return fmt.Sprintf("??? 0x%02X", code)
	}
}

func immediateALUMnemonic(group byte) string {
	switch group {
	case 0:
		return "ADI"
	case 1:
		return "ACI"
	case 2:
		return "SUI"
	case 3:
		return "SBI"
	case 4:
		return "NDI"
	case 5:
		return "XRI"
	case 6:
		return "ORI"
	default:
		return "CPI"
	}
}

func registerALUMnemonic(group byte, src byte) string {
	prefix := [...]string{"AD", "AC", "SU", "SB", "ND", "XR", "OR", "CP"}[group]
	if src == byte(RegM) {
		return prefix + "M"
	}
	return prefix + registerName(src)
}

func registerName(code byte) string {
	switch code & 0x07 {
	case byte(RegA):
		return "A"
	case byte(RegB):
		return "B"
	case byte(RegC):
		return "C"
	case byte(RegD):
		return "D"
	case byte(RegE):
		return "E"
	case byte(RegH):
		return "H"
	case byte(RegL):
		return "L"
	default:
		return "M"
	}
}

func conditionName(code byte) string {
	switch Condition(code & 0x03) {
	case CondCarry:
		return "C"
	case CondZero:
		return "Z"
	case CondSign:
		return "S"
	default:
		return "P"
	}
}

func falseTruePrefix(code byte) string {
	if code&0x20 != 0 {
		return "T"
	}
	return "F"
}
