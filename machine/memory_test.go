package machine

import (
	"errors"
	"testing"

	"retronet-8008/cpu"
)

func TestGenericProfileNewMemoryIsWritable(t *testing.T) {
	profile, ok := Lookup("generic")
	if !ok {
		t.Fatal("Lookup(generic) = false")
	}
	mem, err := profile.NewMemory()
	if err != nil {
		t.Fatalf("NewMemory = %v, want nil", err)
	}

	mem.Write(0x0010, 0xA5)
	if got := mem.Read(0x4010); got != 0xA5 {
		t.Fatalf("Read(0x4010) = 0x%02X, want 0xA5", got)
	}
	kind, mapped := mem.Kind(0x0010)
	if !mapped || kind != MemoryKindRAM {
		t.Fatalf("Kind(0x0010) = %q, %v, want ram, true", kind, mapped)
	}
}

func TestProfileLoadROMProtectsLoadedRange(t *testing.T) {
	profile, ok := Lookup("scelbi-8b")
	if !ok {
		t.Fatal("Lookup(scelbi-8b) = false")
	}
	mem, err := profile.NewMemory()
	if err != nil {
		t.Fatalf("NewMemory = %v, want nil", err)
	}
	if err := profile.LoadROM(mem, "test", []byte{0x06, 0x2A}); err != nil {
		t.Fatalf("LoadROM = %v, want nil", err)
	}

	mem.Write(0x0000, 0xFF)
	if got := mem.Read(0x0000); got != 0x06 {
		t.Fatalf("ROM after Write = 0x%02X, want 0x06", got)
	}
	if kind, _ := mem.Kind(0x0001); kind != MemoryKindROM {
		t.Fatalf("Kind(0x0001) = %q, want rom", kind)
	}

	mem.Write(0x0002, 0x55)
	if got := mem.Read(0x0002); got != 0x55 {
		t.Fatalf("mixed byte after ROM = 0x%02X, want 0x55", got)
	}
}

func TestLoadBytesRejectsROMWithoutPartialWrite(t *testing.T) {
	bus, err := NewMemoryBus([]MemoryRegion{
		{Name: "ram", Start: 0x0000, End: 0x0000, Kind: MemoryKindRAM},
		{Name: "rom", Start: 0x0001, End: 0x0001, Kind: MemoryKindROM},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = LoadBytes(bus, 0x0000, []byte{0x11, 0x22})
	if !errors.Is(err, ErrReadOnlyMemory) {
		t.Fatalf("LoadBytes = %v, want ErrReadOnlyMemory", err)
	}
	if got := bus.Read(0x0000); got != 0x00 {
		t.Fatalf("RAM changed after rejected load: 0x%02X", got)
	}
}

func TestMemoryBusUnmappedReadsOpenBus(t *testing.T) {
	bus, err := NewMemoryBus([]MemoryRegion{
		{Name: "ram", Start: 0x0100, End: 0x01FF, Kind: MemoryKindRAM},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := bus.Read(0x0000); got != 0xFF {
		t.Fatalf("unmapped Read = 0x%02X, want 0xFF", got)
	}
	bus.Write(0x0000, 0x12)
	if got := bus.Read(0x0000); got != 0xFF {
		t.Fatalf("unmapped Read after Write = 0x%02X, want 0xFF", got)
	}
}

func TestNewMemoryBusRejectsOverlappingRegions(t *testing.T) {
	_, err := NewMemoryBus([]MemoryRegion{
		{Name: "first", Start: 0x0000, End: 0x00FF, Kind: MemoryKindRAM},
		{Name: "second", Start: 0x0080, End: 0x01FF, Kind: MemoryKindROM},
	})
	if err == nil {
		t.Fatal("NewMemoryBus overlap = nil, want error")
	}
}

func TestNewMemoryBusRejectsInvalidRegion(t *testing.T) {
	_, err := NewMemoryBus([]MemoryRegion{
		{Name: "outside", Start: cpu.AddressMask + 1, End: cpu.AddressMask + 1, Kind: MemoryKindRAM},
	})
	if err == nil {
		t.Fatal("NewMemoryBus outside address space = nil, want error")
	}
}
