// Package machine definisce profili macchina sopra il core CPU 8008.
package machine

import (
	"fmt"
	"sort"

	"retronet-8008/cpu"
)

const DefaultStepLimit = uint64(1000)

// ROMSlot descrive una regione caricabile associata a un profilo macchina.
type ROMSlot struct {
	Name        string
	Address     uint16
	MaxSize     int
	Required    bool
	Description string
}

// Profile descrive una macchina 8008 ad alto livello.
//
// I profili storici sono volutamente conservativi: non includono ROM reali e
// usano slot documentali per permettere caricamenti locali espliciti.
type Profile struct {
	Name               string
	Description        string
	DefaultLoadAddress uint16
	DefaultStartPC     uint16
	DefaultStepLimit   uint64
	ROMSlots           []ROMSlot
}

var profiles = []Profile{
	{
		Name:               "generic",
		Description:        "Macchina piatta generica: 16 KB, nessuna ROM predefinita.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
	},
	{
		Name:               "intellec-8",
		Description:        "Scheletro per sistemi Intel Intellec 8/MOD 8; ROM locali non incluse.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
		ROMSlots: []ROMSlot{
			{
				Name:        "monitor",
				Address:     0x0000,
				MaxSize:     cpu.AddressSpaceSize,
				Required:    false,
				Description: "Monitor/bootstrap locale caricato dall'utente.",
			},
		},
	},
	{
		Name:               "scelbi-8h",
		Description:        "Scheletro per sistemi SCELBI 8H; ROM locali non incluse.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
		ROMSlots: []ROMSlot{
			{
				Name:        "monitor",
				Address:     0x0000,
				MaxSize:     cpu.AddressSpaceSize,
				Required:    false,
				Description: "Monitor/bootstrap locale caricato dall'utente.",
			},
		},
	},
	{
		Name:               "scelbi-8b",
		Description:        "Scheletro per profilo SCELBI 8B compatibile con riferimenti SIMH; ROM locali non incluse.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
		ROMSlots: []ROMSlot{
			{
				Name:        "monitor",
				Address:     0x0000,
				MaxSize:     cpu.AddressSpaceSize,
				Required:    false,
				Description: "Monitor/bootstrap locale caricato dall'utente.",
			},
		},
	},
}

// Profiles restituisce una copia ordinata dei profili disponibili.
func Profiles() []Profile {
	out := make([]Profile, 0, len(profiles))
	for _, profile := range profiles {
		out = append(out, cloneProfile(profile))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

// Lookup restituisce un profilo per nome.
func Lookup(name string) (Profile, bool) {
	for _, p := range profiles {
		if p.Name == name {
			return cloneProfile(p), true
		}
	}
	return Profile{}, false
}

func cloneProfile(p Profile) Profile {
	p.ROMSlots = append([]ROMSlot(nil), p.ROMSlots...)
	return p
}

// ROMSlot restituisce uno slot ROM del profilo per nome.
func (p Profile) ROMSlot(name string) (ROMSlot, bool) {
	for _, slot := range p.ROMSlots {
		if slot.Name == name {
			return slot, true
		}
	}
	return ROMSlot{}, false
}

// LoadBytes carica data in memoria a partire da addr, senza wrap silenzioso.
func LoadBytes(mem cpu.Memory, addr uint16, data []byte) error {
	if mem == nil {
		return cpu.ErrNilMemory
	}
	if err := ValidateRange(addr, len(data)); err != nil {
		return err
	}
	for i, b := range data {
		mem.Write(addr+uint16(i), b)
	}
	return nil
}

// LoadROM carica data nello slot ROM indicato.
func (p Profile) LoadROM(mem cpu.Memory, name string, data []byte) error {
	slot, ok := p.ROMSlot(name)
	if !ok {
		return fmt.Errorf("slot ROM %q non presente nel profilo %q", name, p.Name)
	}
	if slot.MaxSize > 0 && len(data) > slot.MaxSize {
		return fmt.Errorf("ROM %q: %d byte superano il limite %d", name, len(data), slot.MaxSize)
	}
	return LoadBytes(mem, slot.Address, data)
}

// ValidateRange controlla che un caricamento non esca dallo spazio 14 bit.
func ValidateRange(addr uint16, size int) error {
	if size < 0 {
		return fmt.Errorf("dimensione negativa %d", size)
	}
	if size > cpu.AddressSpaceSize {
		return fmt.Errorf("%d byte superano lo spazio indirizzabile %d byte", size, cpu.AddressSpaceSize)
	}
	if int(addr)+size > cpu.AddressSpaceSize {
		return fmt.Errorf("%d byte a 0x%04X superano 0x%04X", size, addr, cpu.AddressMask)
	}
	return nil
}
