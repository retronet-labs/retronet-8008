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
