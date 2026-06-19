package cpu

import "testing"

func TestFlatMemoryStartsZeroed(t *testing.T) {
	mem := NewFlatMemory()

	if got := mem.Read(0x0000); got != 0 {
		t.Fatalf("Read(0x0000) = 0x%02X, want 0", got)
	}
	if got := mem.Read(0x3FFF); got != 0 {
		t.Fatalf("Read(0x3FFF) = 0x%02X, want 0", got)
	}
}

func TestFlatMemoryReadWrite(t *testing.T) {
	mem := NewFlatMemory()

	mem.Write(0x1234, 0xAB)

	if got := mem.Read(0x1234); got != 0xAB {
		t.Fatalf("Read(0x1234) = 0x%02X, want 0xAB", got)
	}
}

func TestFlatMemoryMasksReadAddress(t *testing.T) {
	mem := NewFlatMemory()
	mem.Write(0x0000, 0x11)
	mem.Write(0x3FFF, 0x22)

	if got := mem.Read(0x4000); got != 0x11 {
		t.Fatalf("Read(0x4000) = 0x%02X, want 0x11", got)
	}
	if got := mem.Read(0xFFFF); got != 0x22 {
		t.Fatalf("Read(0xFFFF) = 0x%02X, want 0x22", got)
	}
}

func TestFlatMemoryMasksWriteAddress(t *testing.T) {
	mem := NewFlatMemory()

	mem.Write(0x4000, 0x44)
	if got := mem.Read(0x0000); got != 0x44 {
		t.Fatalf("Read(0x0000) = 0x%02X, want 0x44", got)
	}

	mem.Write(0xFFFF, 0x55)
	if got := mem.Read(0x3FFF); got != 0x55 {
		t.Fatalf("Read(0x3FFF) = 0x%02X, want 0x55", got)
	}
}

func TestFlatMemoryImplementsMemory(t *testing.T) {
	var mem Memory = NewFlatMemory()
	mem.Write(0x0020, 0x80)

	if got := mem.Read(0x0020); got != 0x80 {
		t.Fatalf("Read(0x0020) = 0x%02X, want 0x80", got)
	}
}
