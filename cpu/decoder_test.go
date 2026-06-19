package cpu

import "testing"

func TestOpcodeTableHas256Entries(t *testing.T) {
	table := OpcodeTable()

	if len(table) != 256 {
		t.Fatalf("len(OpcodeTable()) = %d, want 256", len(table))
	}

	for i, op := range table {
		if op.Code != byte(i) {
			t.Fatalf("OpcodeTable()[%d].Code = 0x%02X, want 0x%02X", i, op.Code, byte(i))
		}
		if op.Length < 1 || op.Length > 3 {
			t.Fatalf("OpcodeTable()[0x%02X].Length = %d, want 1..3", i, op.Length)
		}
		if op.Execute == nil {
			t.Fatalf("OpcodeTable()[0x%02X].Execute = nil", i)
		}
	}
}

func TestDecodeKnownInstructionLengths(t *testing.T) {
	tests := []struct {
		opcode byte
		want   byte
	}{
		{0xC1, 1}, // LAB
		{0x06, 2}, // LAI
		{0x3E, 2}, // LMI
		{0x04, 2}, // ADI
		{0x44, 3}, // JMP
		{0x46, 3}, // CAL
		{0x40, 3}, // JFC
		{0x42, 3}, // CFC
		{0x02, 1}, // RLC
		{0x05, 1}, // RST 0
	}

	for _, tt := range tests {
		if got := Decode(tt.opcode).Length; got != tt.want {
			t.Errorf("Decode(0x%02X).Length = %d, want %d", tt.opcode, got, tt.want)
		}
	}
}

func TestDecodeKnownMnemonics(t *testing.T) {
	tests := []struct {
		opcode byte
		want   string
	}{
		{0x00, "HLT"},
		{0x02, "RLC"},
		{0x04, "ADI"},
		{0x06, "LAI"},
		{0x3E, "LMI"},
		{0x44, "JMP"},
		{0x46, "CAL"},
		{0x41, "INP 0"},
		{0x4F, "INP 7"},
		{0x51, "OUT 8"},
		{0x61, "OUT 16"},
		{0x6F, "OUT 23"},
		{0x7F, "OUT 31"},
		{0x81, "ADB"},
		{0x87, "ADM"},
		{0xC1, "LAB"},
		{0xC0, "NOP"},
	}

	for _, tt := range tests {
		if got := Decode(tt.opcode).Mnemonic; got != tt.want {
			t.Errorf("Decode(0x%02X).Mnemonic = %q, want %q", tt.opcode, got, tt.want)
		}
	}
}
