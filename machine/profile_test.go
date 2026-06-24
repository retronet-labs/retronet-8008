package machine

import (
	"errors"
	"strings"
	"testing"

	"github.com/retronet-labs/retronet-8008/cpu"
)

func TestProfilesAreLookupable(t *testing.T) {
	for _, profile := range Profiles() {
		got, ok := Lookup(profile.Name)
		if !ok {
			t.Fatalf("Lookup(%q) = false", profile.Name)
		}
		if got.Name != profile.Name {
			t.Fatalf("Lookup(%q).Name = %q", profile.Name, got.Name)
		}
	}
}

func TestProfilesReturnsCopies(t *testing.T) {
	profile, ok := Lookup("intellec-8")
	if !ok {
		t.Fatal("Lookup(intellec-8) = false")
	}
	profile.ROMSlots[0].Name = "mutated"
	profile.MemoryRegions[0].Name = "mutated"
	profile.IOPorts[0].Name = "mutated"
	profile.ROMHints[0].Name = "mutated"

	again, ok := Lookup("intellec-8")
	if !ok {
		t.Fatal("Lookup(intellec-8) second call = false")
	}
	if again.ROMSlots[0].Name != "monitor" {
		t.Fatalf("ROMSlots[0].Name = %q, want monitor", again.ROMSlots[0].Name)
	}
	if again.MemoryRegions[0].Name != "intellec-direct-memory" {
		t.Fatalf("MemoryRegions[0].Name = %q, want intellec-direct-memory", again.MemoryRegions[0].Name)
	}
	if again.IOPorts[0].Name != "callback-input-0" {
		t.Fatalf("IOPorts[0].Name = %q, want callback-input-0", again.IOPorts[0].Name)
	}
	if again.ROMHints[0].Name != "monitor" {
		t.Fatalf("ROMHints[0].Name = %q, want monitor", again.ROMHints[0].Name)
	}
}

func TestLoadBytes(t *testing.T) {
	mem := cpu.NewFlatMemory()

	if err := LoadBytes(mem, 0x0010, []byte{0xAA, 0xBB}); err != nil {
		t.Fatalf("LoadBytes = %v, want nil", err)
	}
	if got := mem.Read(0x0010); got != 0xAA {
		t.Fatalf("mem[0x0010] = 0x%02X, want 0xAA", got)
	}
	if got := mem.Read(0x0011); got != 0xBB {
		t.Fatalf("mem[0x0011] = 0x%02X, want 0xBB", got)
	}
}

func TestLoadBytesRejectsNilMemory(t *testing.T) {
	err := LoadBytes(nil, 0, []byte{0x00})
	if !errors.Is(err, cpu.ErrNilMemory) {
		t.Fatalf("LoadBytes(nil) = %v, want ErrNilMemory", err)
	}
}

func TestValidateRangeRejectsOverflow(t *testing.T) {
	err := ValidateRange(0x3FFF, 2)
	if err == nil {
		t.Fatal("ValidateRange = nil, want overflow error")
	}
}

func TestProfileLoadROM(t *testing.T) {
	profile, ok := Lookup("intellec-8")
	if !ok {
		t.Fatal("Lookup(intellec-8) = false")
	}
	mem := cpu.NewFlatMemory()

	if err := profile.LoadROM(mem, "monitor", []byte{0x06, 0x2A}); err != nil {
		t.Fatalf("LoadROM = %v, want nil", err)
	}
	if got := mem.Read(0x0000); got != 0x06 {
		t.Fatalf("mem[0x0000] = 0x%02X, want 0x06", got)
	}
}

func TestProfileLoadTestROM(t *testing.T) {
	profile, ok := Lookup("scelbi-8b")
	if !ok {
		t.Fatal("Lookup(scelbi-8b) = false")
	}
	mem := cpu.NewFlatMemory()

	if err := profile.LoadROM(mem, "test", []byte{cpu.INP(0), cpu.OUT(8), cpu.HLT()}); err != nil {
		t.Fatalf("LoadROM(test) = %v, want nil", err)
	}
	if got := mem.Read(0x0001); got != cpu.OUT(8) {
		t.Fatalf("mem[0x0001] = 0x%02X, want OUT 8", got)
	}
}

func TestProfileLoadROMRejectsUnknownSlot(t *testing.T) {
	profile, ok := Lookup("generic")
	if !ok {
		t.Fatal("Lookup(generic) = false")
	}
	err := profile.LoadROM(cpu.NewFlatMemory(), "monitor", []byte{0x00})
	if err == nil || !strings.Contains(err.Error(), "slot ROM") {
		t.Fatalf("LoadROM unknown slot = %v, want slot ROM error", err)
	}
}
