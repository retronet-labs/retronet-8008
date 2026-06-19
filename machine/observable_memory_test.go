package machine

import (
	"testing"

	"retronet-8008/cpu"
)

func TestObservableMemoryReportsEffectiveWrite(t *testing.T) {
	base, err := NewMemoryBus([]MemoryRegion{
		{Name: "mixed", Start: 0, End: cpu.AddressMask, Kind: MemoryKindMixed},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := base.LoadROM(0x0010, []byte{0xAA}); err != nil {
		t.Fatal(err)
	}
	memory, err := NewObservableMemory(base)
	if err != nil {
		t.Fatal(err)
	}
	var writes []MemoryWrite
	memory.ObserveWrites(func(write MemoryWrite) { writes = append(writes, write) })

	memory.Write(0x4010, 0x55)

	if len(writes) != 1 {
		t.Fatalf("writes = %v, want one", writes)
	}
	want := MemoryWrite{Address: 0x0010, Before: 0xAA, Requested: 0x55, After: 0xAA}
	if writes[0] != want {
		t.Fatalf("write = %+v, want %+v", writes[0], want)
	}
}

func TestObservableMemoryLoadDoesNotEmitRuntimeWrite(t *testing.T) {
	base := cpu.NewFlatMemory()
	memory, err := NewObservableMemory(base)
	if err != nil {
		t.Fatal(err)
	}
	called := false
	memory.ObserveWrites(func(MemoryWrite) { called = true })

	if err := memory.LoadBytes(0x0020, []byte{0x12}); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("loader emitted runtime write")
	}
	if got := memory.Read(0x0020); got != 0x12 {
		t.Fatalf("Read = 0x%02X, want 0x12", got)
	}
}
