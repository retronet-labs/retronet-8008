package cpu

import "testing"

func TestAddr14(t *testing.T) {
	tests := []struct {
		in   uint16
		want uint16
	}{
		{0x0000, 0x0000},
		{0x1234, 0x1234},
		{0x3FFF, 0x3FFF},
		{0x4000, 0x0000},
		{0xFFFF, 0x3FFF},
	}

	for _, tt := range tests {
		if got := addr14(tt.in); got != tt.want {
			t.Errorf("addr14(0x%04X) = 0x%04X, want 0x%04X", tt.in, got, tt.want)
		}
	}
}

func TestHLAddressMasksH(t *testing.T) {
	if got := hlAddress(0xBF, 0x80); got != 0x3F80 {
		t.Fatalf("hlAddress(0xBF, 0x80) = 0x%04X, want 0x3F80", got)
	}
}

func TestRegisterCodes(t *testing.T) {
	tests := []struct {
		name string
		reg  Register
		want byte
	}{
		{"A", RegA, 0b000},
		{"B", RegB, 0b001},
		{"C", RegC, 0b010},
		{"D", RegD, 0b011},
		{"E", RegE, 0b100},
		{"H", RegH, 0b101},
		{"L", RegL, 0b110},
		{"M", RegM, 0b111},
	}

	for _, tt := range tests {
		if got := regBits(tt.reg); got != tt.want {
			t.Errorf("regBits(%s) = %03b, want %03b", tt.name, got, tt.want)
		}
	}
}

func TestConditionCodes(t *testing.T) {
	tests := []struct {
		name string
		cond Condition
		want byte
	}{
		{"Carry", CondCarry, 0b00},
		{"Zero", CondZero, 0b01},
		{"Sign", CondSign, 0b10},
		{"Parity", CondParity, 0b11},
	}

	for _, tt := range tests {
		if got := condBits(tt.cond); got != tt.want {
			t.Errorf("condBits(%s) = %02b, want %02b", tt.name, got, tt.want)
		}
	}
}

func TestLoadOpcodeHelpers(t *testing.T) {
	if got := L(RegA, RegB); got != 0xC1 {
		t.Fatalf("L(RegA, RegB) = 0x%02X, want 0xC1", got)
	}
	if got := L(RegM, RegA); got != 0xF8 {
		t.Fatalf("L(RegM, RegA) = 0x%02X, want 0xF8", got)
	}
	if got := LI(RegD); got != 0x1E {
		t.Fatalf("LI(RegD) = 0x%02X, want 0x1E", got)
	}
	if got := LI(RegM); got != 0x3E {
		t.Fatalf("LI(RegM) = 0x%02X, want 0x3E", got)
	}
	if got := NOP(); got != 0xC0 {
		t.Fatalf("NOP() = 0x%02X, want 0xC0", got)
	}
	if got := HLT(); got != 0x00 {
		t.Fatalf("HLT() = 0x%02X, want 0x00", got)
	}
}

func TestALUOpcodeHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  byte
		want byte
	}{
		{"AD B", AD(RegB), 0x81},
		{"AC M", AC(RegM), 0x8F},
		{"SU C", SU(RegC), 0x92},
		{"SB D", SB(RegD), 0x9B},
		{"ND E", ND(RegE), 0xA4},
		{"XR H", XR(RegH), 0xAD},
		{"OR L", OR(RegL), 0xB6},
		{"CP A", CP(RegA), 0xB8},
		{"ADI", ADI(), 0x04},
		{"ACI", ACI(), 0x0C},
		{"SUI", SUI(), 0x14},
		{"SBI", SBI(), 0x1C},
		{"NDI", NDI(), 0x24},
		{"XRI", XRI(), 0x2C},
		{"ORI", ORI(), 0x34},
		{"CPI", CPI(), 0x3C},
		{"INR B", INR(RegB), 0x08},
		{"DCR L", DCR(RegL), 0x31},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = 0x%02X, want 0x%02X", tt.name, tt.got, tt.want)
		}
	}
}

func TestRotateOpcodeHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  byte
		want byte
	}{
		{"RLC", RLC(), 0x02},
		{"RRC", RRC(), 0x0A},
		{"RAL", RAL(), 0x12},
		{"RAR", RAR(), 0x1A},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = 0x%02X, want 0x%02X", tt.name, tt.got, tt.want)
		}
	}
}

func TestControlFlowOpcodeHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  byte
		want byte
	}{
		{"JMP", JMP(), 0x44},
		{"JFC", JF(CondCarry), 0x40},
		{"JTP", JT(CondParity), 0x78},
		{"CAL", CAL(), 0x46},
		{"CFZ", CF(CondZero), 0x4A},
		{"CTS", CT(CondSign), 0x72},
		{"RET", RET(), 0x07},
		{"RFC", RF(CondCarry), 0x03},
		{"RTP", RT(CondParity), 0x3B},
		{"RST 3", RST(3), 0x1D},
		{"RST masks vector", RST(9), 0x0D},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = 0x%02X, want 0x%02X", tt.name, tt.got, tt.want)
		}
	}
}
