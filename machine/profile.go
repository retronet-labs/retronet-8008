// Package machine definisce profili macchina sopra il core CPU 8008.
package machine

import (
	"fmt"
	"sort"

	"retronet-8008/cpu"
)

const DefaultStepLimit = uint64(1000)

// MemoryKind descrive il comportamento di una regione memoria.
type MemoryKind string

const (
	MemoryKindRAM   MemoryKind = "ram"
	MemoryKindROM   MemoryKind = "rom"
	MemoryKindMixed MemoryKind = "mixed"
)

// IODirection descrive la direzione di una porta I/O 8008.
type IODirection string

const (
	IODirectionInput  IODirection = "input"
	IODirectionOutput IODirection = "output"
)

// ROMSlot descrive una regione caricabile associata a un profilo macchina.
type ROMSlot struct {
	Name        string
	Address     uint16
	MaxSize     int
	Required    bool
	Description string
}

// MemoryRegion descrive una regione di memoria prevista da un profilo.
type MemoryRegion struct {
	Name        string
	Start       uint16
	End         uint16
	Kind        MemoryKind
	Description string
}

// IOPort descrive una porta usata o riservata da un profilo macchina.
type IOPort struct {
	Port        byte
	Direction   IODirection
	Name        string
	Historical  bool
	Description string
}

// ROMHint descrive una ROM locale utile per validare un profilo.
type ROMHint struct {
	Name        string
	Slot        string
	Included    bool
	Description string
}

// Profile descrive una macchina 8008 ad alto livello.
//
// I profili storici sono volutamente conservativi: non includono ROM reali e
// usano slot documentali per permettere caricamenti locali espliciti.
type Profile struct {
	Name               string
	Description        string
	HistoricalNote     string
	DefaultLoadAddress uint16
	DefaultStartPC     uint16
	DefaultStepLimit   uint64
	ROMSlots           []ROMSlot
	MemoryRegions      []MemoryRegion
	IOPorts            []IOPort
	ROMHints           []ROMHint
}

var profiles = []Profile{
	{
		Name:               "generic",
		Description:        "Macchina piatta generica: 16 KB, nessuna ROM predefinita.",
		HistoricalNote:     "Profilo didattico senza macchina storica associata.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
		MemoryRegions:      flatMemoryRegion("direct-memory", MemoryKindRAM, "Spazio diretto 8008 usato come RAM piatta per test e binari raw."),
	},
	{
		Name:               "intellec-8",
		Description:        "Profilo iniziale per sistemi Intel Intellec 8/MCS-8; ROM locali non incluse.",
		HistoricalNote:     "Intellec era la linea Intel di sistemi di sviluppo per microprocessori; Intellec 8 serviva a sviluppare e provare software/firmware 8008, non a essere un home computer generico.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
		ROMSlots:           monitorAndTestSlots(),
		MemoryRegions:      flatMemoryRegion("intellec-direct-memory", MemoryKindMixed, "Spazio diretto massimo dell'8008; resta scrivibile salvo gli intervalli occupati da ROM locali."),
		IOPorts:            callbackConsolePorts(),
		ROMHints:           monitorAndTestHints("Intel monitor o ROM diagnostica locale per Intellec 8."),
	},
	{
		Name:               "scelbi-8h",
		Description:        "Profilo iniziale per sistemi SCELBI 8H; ROM locali non incluse.",
		HistoricalNote:     "SCELBI 8H era un microcomputer/kit basato su Intel 8008, venduto anche assemblato; il sistema base era front-panel oriented, con periferiche opzionali come terminale o cassette.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
		ROMSlots:           monitorAndTestSlots(),
		MemoryRegions:      flatMemoryRegion("scelbi-direct-memory", MemoryKindMixed, "Spazio diretto fino a 16 KB; resta scrivibile salvo gli intervalli occupati da ROM locali."),
		IOPorts:            callbackConsolePorts(),
		ROMHints:           monitorAndTestHints("SCELBI Monitor, Editor, Assembler o SCELBAL convertiti in binario locale."),
	},
	{
		Name:               "scelbi-8b",
		Description:        "Profilo iniziale SCELBI 8B compatibile con riferimenti SIMH; ROM locali non incluse.",
		HistoricalNote:     "SCELBI 8B e' il profilo piu' comodo per confronti con simulatori esistenti e software come SCELBAL/Forth, ma qui resta ancora una mappa conservativa.",
		DefaultLoadAddress: 0x0000,
		DefaultStartPC:     0x0000,
		DefaultStepLimit:   DefaultStepLimit,
		ROMSlots:           monitorAndTestSlots(),
		MemoryRegions:      flatMemoryRegion("scelbi-direct-memory", MemoryKindMixed, "Spazio diretto fino a 16 KB; resta scrivibile salvo gli intervalli occupati da ROM locali."),
		IOPorts:            callbackConsolePorts(),
		ROMHints:           monitorAndTestHints("SCELBI Monitor, SCELBAL o forth-scelbi.bin caricati come file locali."),
	},
}

func flatMemoryRegion(name string, kind MemoryKind, description string) []MemoryRegion {
	return []MemoryRegion{
		{
			Name:        name,
			Start:       0x0000,
			End:         cpu.AddressMask,
			Kind:        kind,
			Description: description,
		},
	}
}

func monitorAndTestSlots() []ROMSlot {
	return []ROMSlot{
		{
			Name:        "monitor",
			Address:     0x0000,
			MaxSize:     cpu.AddressSpaceSize,
			Required:    false,
			Description: "Monitor/bootstrap locale caricato dall'utente.",
		},
		{
			Name:        "test",
			Address:     0x0000,
			MaxSize:     cpu.AddressSpaceSize,
			Required:    false,
			Description: "ROM locale di smoke test caricata a 0x0000.",
		},
	}
}

func callbackConsolePorts() []IOPort {
	return []IOPort{
		{
			Port:        0,
			Direction:   IODirectionInput,
			Name:        "callback-input-0",
			Historical:  false,
			Description: "Porta input convenzionale per test, terminale o front panel emulato via callback.",
		},
		{
			Port:        8,
			Direction:   IODirectionOutput,
			Name:        "callback-output-8",
			Historical:  false,
			Description: "Porta output convenzionale per test, terminale o front panel emulato via callback.",
		},
	}
}

func monitorAndTestHints(monitorDescription string) []ROMHint {
	return []ROMHint{
		{
			Name:        "monitor",
			Slot:        "monitor",
			Included:    false,
			Description: monitorDescription,
		},
		{
			Name:        "io-smoke",
			Slot:        "test",
			Included:    false,
			Description: "ROM locale minima: INP 0, OUT 8, HLT; valida loader e callback I/O senza dipendere da ROM storiche.",
		},
	}
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
	p.MemoryRegions = append([]MemoryRegion(nil), p.MemoryRegions...)
	p.IOPorts = append([]IOPort(nil), p.IOPorts...)
	p.ROMHints = append([]ROMHint(nil), p.ROMHints...)
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

// NewMemory crea il bus memoria mappato associato al profilo.
func (p Profile) NewMemory() (*MemoryBus, error) {
	return NewMemoryBus(p.MemoryRegions)
}

type byteLoader interface {
	LoadBytes(addr uint16, data []byte) error
}

type romLoader interface {
	LoadROM(addr uint16, data []byte) error
}

// LoadBytes carica data in memoria a partire da addr, senza wrap silenzioso.
func LoadBytes(mem cpu.Memory, addr uint16, data []byte) error {
	if mem == nil {
		return cpu.ErrNilMemory
	}
	if err := ValidateRange(addr, len(data)); err != nil {
		return err
	}
	if loader, ok := mem.(byteLoader); ok {
		return loader.LoadBytes(addr, data)
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
	if loader, ok := mem.(romLoader); ok {
		return loader.LoadROM(slot.Address, data)
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
