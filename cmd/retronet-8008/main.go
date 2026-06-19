// Comando retronet-8008: runner minimale dell'emulatore Intel 8008.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"retronet-8008/cpu"
	"retronet-8008/machine"
)

type runConfig struct {
	binPath      string
	profileName  string
	listProfiles bool
	loadAt       uint16
	startPC      uint16
	steps        uint64
	disasm       uint64
	trace        bool
	ioTrace      bool
	roms         romFlags
	inputs       inputFlags
}

type romSpec struct {
	name string
	path string
}

type romFlags []romSpec

func (r *romFlags) String() string {
	if r == nil || len(*r) == 0 {
		return ""
	}
	parts := make([]string, 0, len(*r))
	for _, spec := range *r {
		parts = append(parts, spec.name+"="+spec.path)
	}
	return strings.Join(parts, ",")
}

func (r *romFlags) Set(value string) error {
	name, path, ok := strings.Cut(value, "=")
	if !ok {
		return errors.New("usa nome=percorso")
	}
	name = strings.TrimSpace(name)
	path = strings.TrimSpace(path)
	if name == "" || path == "" {
		return errors.New("nome e percorso ROM sono obbligatori")
	}
	*r = append(*r, romSpec{name: name, path: path})
	return nil
}

type inputSpec struct {
	port  byte
	value byte
}

type inputFlags []inputSpec

func (i *inputFlags) String() string {
	if i == nil || len(*i) == 0 {
		return ""
	}
	parts := make([]string, 0, len(*i))
	for _, spec := range *i {
		parts = append(parts, fmt.Sprintf("%d=0x%02X", spec.port, spec.value))
	}
	return strings.Join(parts, ",")
}

func (i *inputFlags) Set(value string) error {
	portText, valueText, ok := strings.Cut(value, "=")
	if !ok {
		return errors.New("usa porta=valore")
	}
	port, err := parsePort(portText)
	if err != nil {
		return err
	}
	if err := cpu.ValidateInputPort(port); err != nil {
		return err
	}
	n, err := strconv.ParseUint(strings.TrimSpace(valueText), 0, 8)
	if err != nil {
		return err
	}
	*i = append(*i, inputSpec{port: port, value: byte(n)})
	return nil
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	cfg, err := parseFlags(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(stderr, "errore: %v\n", err)
		return 2
	}

	if cfg.listProfiles {
		printProfiles(stdout)
		return 0
	}

	profile, ok := machine.Lookup(cfg.profileName)
	if !ok {
		fmt.Fprintf(stderr, "errore profilo: profilo %q non disponibile\n", cfg.profileName)
		return 2
	}

	c := cpu.NewCPU8008()
	mem := cpu.NewFlatMemory()
	ports := profile.NewIO()

	for _, spec := range cfg.inputs {
		if err := ports.SetInput(spec.port, spec.value); err != nil {
			fmt.Fprintf(stderr, "errore input I/O: %v\n", err)
			return 2
		}
	}
	if cfg.ioTrace {
		if err := registerIOTrace(stdout, ports); err != nil {
			fmt.Fprintf(stderr, "errore trace I/O: %v\n", err)
			return 2
		}
	}

	for _, spec := range cfg.roms {
		data, err := os.ReadFile(spec.path)
		if err != nil {
			fmt.Fprintf(stderr, "errore caricamento ROM %s: %v\n", spec.name, err)
			return 1
		}
		if err := profile.LoadROM(mem, spec.name, data); err != nil {
			fmt.Fprintf(stderr, "errore caricamento ROM %s: %v\n", spec.name, err)
			return 1
		}
	}

	loaded := 0
	if cfg.binPath != "" {
		program, err := os.ReadFile(cfg.binPath)
		if err != nil {
			fmt.Fprintf(stderr, "errore caricamento binario: %v\n", err)
			return 1
		}
		if err := machine.LoadBytes(mem, cfg.loadAt, program); err != nil {
			fmt.Fprintf(stderr, "errore caricamento binario: %v\n", err)
			return 1
		}
		loaded = len(program)
	}

	if cfg.disasm > 0 {
		if err := printDisassembly(stdout, mem, cfg.startPC, cfg.disasm); err != nil {
			fmt.Fprintf(stderr, "errore disassembly: %v\n", err)
			return 1
		}
		return 0
	}

	if err := c.Jam(mem, ports, cpu.JMP(), byte(cfg.startPC), byte(cfg.startPC>>8)); err != nil {
		fmt.Fprintf(stderr, "errore avvio CPU: %v\n", err)
		return 1
	}

	var trace io.Writer
	if cfg.trace {
		trace = stdout
	}
	executed, limitReached, err := runSteps(c, mem, ports, cfg.steps, trace)
	printDump(stdout, c, cfg, loaded, len(cfg.roms), executed, limitReached)
	if err != nil {
		fmt.Fprintf(stderr, "errore esecuzione: %v\n", err)
		return 1
	}
	return 0
}

func parseFlags(args []string, stderr io.Writer) (runConfig, error) {
	fs := flag.NewFlagSet("retronet-8008", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg runConfig
	loadAt := fs.String("addr", "0x0000", "indirizzo di caricamento, decimale o 0xHEX")
	startPC := fs.String("pc", "", "program counter iniziale, default uguale ad -addr")
	fs.StringVar(&cfg.binPath, "bin", "", "percorso del binario da caricare")
	fs.StringVar(&cfg.profileName, "profile", "generic", "profilo macchina da usare")
	fs.BoolVar(&cfg.listProfiles, "profiles", false, "elenca i profili macchina disponibili")
	fs.Var(&cfg.roms, "rom", "carica una ROM di profilo nel formato nome=percorso; ripetibile")
	fs.Var(&cfg.inputs, "input", "inizializza una porta input nel formato porta=valore; ripetibile")
	fs.Uint64Var(&cfg.steps, "steps", machine.DefaultStepLimit, "numero massimo di istruzioni da eseguire")
	fs.Uint64Var(&cfg.disasm, "disasm", 0, "disassembla N istruzioni e termina senza eseguire")
	fs.BoolVar(&cfg.trace, "trace", false, "stampa ogni istruzione prima dell'esecuzione")
	fs.BoolVar(&cfg.ioTrace, "io-trace", false, "stampa letture e scritture I/O tramite callback")

	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	if cfg.listProfiles {
		return cfg, nil
	}
	if cfg.binPath == "" && len(cfg.roms) == 0 {
		fs.Usage()
		return cfg, errors.New("flag -bin o -rom obbligatorio")
	}

	addr, err := parseAddress(*loadAt)
	if err != nil {
		return cfg, fmt.Errorf("addr non valido: %w", err)
	}
	cfg.loadAt = addr

	if *startPC == "" {
		cfg.startPC = cfg.loadAt
		return cfg, nil
	}
	pc, err := parseAddress(*startPC)
	if err != nil {
		return cfg, fmt.Errorf("pc non valido: %w", err)
	}
	cfg.startPC = pc
	return cfg, nil
}

func parseAddress(value string) (uint16, error) {
	value = strings.TrimSpace(value)
	n, err := strconv.ParseUint(value, 0, 16)
	if err != nil {
		return 0, err
	}
	if uint16(n)&^cpu.AddressMask != 0 {
		return 0, fmt.Errorf("0x%04X fuori dallo spazio 14 bit", n)
	}
	return uint16(n), nil
}

func parsePort(value string) (byte, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(value), 0, 8)
	if err != nil {
		return 0, err
	}
	return byte(n), nil
}

func runSteps(c *cpu.CPU8008, mem cpu.Memory, ioBus cpu.IO, limit uint64, trace io.Writer) (uint64, bool, error) {
	var executed uint64
	for executed < limit {
		if c.Halted || c.Stopped {
			return executed, false, nil
		}
		if trace != nil {
			d, err := cpu.Disassemble(mem, c.PC)
			if err != nil {
				return executed, false, err
			}
			fmt.Fprintf(trace, "trace=%d %s\n", executed, d.String())
		}
		err := c.Step(mem, ioBus)
		if err != nil {
			if errors.Is(err, cpu.ErrCPUStopped) {
				return executed, false, nil
			}
			return executed, false, err
		}
		executed++
	}
	return executed, !(c.Halted || c.Stopped), nil
}

func printDisassembly(w io.Writer, mem cpu.Memory, pc uint16, count uint64) error {
	for i := uint64(0); i < count; i++ {
		d, err := cpu.Disassemble(mem, pc)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, d.String())
		pc = d.NextPC
	}
	return nil
}

func registerIOTrace(w io.Writer, ioBus *machine.CallbackIO) error {
	for p := byte(0); p <= 7; p++ {
		port := p
		if err := ioBus.OnInput(port, func(port byte, value byte) byte {
			fmt.Fprintf(w, "io in port=%d value=0x%02X\n", port, value)
			return value
		}); err != nil {
			return err
		}
	}
	for p := byte(8); p <= 31; p++ {
		port := p
		if err := ioBus.OnOutput(port, func(port byte, value byte) {
			fmt.Fprintf(w, "io out port=%d value=0x%02X\n", port, value)
		}); err != nil {
			return err
		}
	}
	return nil
}

func printProfiles(w io.Writer) {
	for _, profile := range machine.Profiles() {
		fmt.Fprintf(w, "%s: %s\n", profile.Name, profile.Description)
		if profile.HistoricalNote != "" {
			fmt.Fprintf(w, "  note %s\n", profile.HistoricalNote)
		}
		for _, region := range profile.MemoryRegions {
			fmt.Fprintf(w, "  mem %s 0x%04X-0x%04X %s - %s\n", region.Name, region.Start, region.End, region.Kind, region.Description)
		}
		for _, slot := range profile.ROMSlots {
			required := "optional"
			if slot.Required {
				required = "required"
			}
			fmt.Fprintf(w, "  rom %s @0x%04X max=%d %s - %s\n", slot.Name, slot.Address, slot.MaxSize, required, slot.Description)
		}
		for _, port := range profile.IOPorts {
			historical := "emu"
			if port.Historical {
				historical = "historical"
			}
			fmt.Fprintf(w, "  io %s %d %s %s - %s\n", port.Direction, port.Port, port.Name, historical, port.Description)
		}
		for _, hint := range profile.ROMHints {
			included := "external"
			if hint.Included {
				included = "included"
			}
			fmt.Fprintf(w, "  hint %s slot=%s %s - %s\n", hint.Name, hint.Slot, included, hint.Description)
		}
	}
}

func printDump(w io.Writer, c *cpu.CPU8008, cfg runConfig, loaded int, roms int, executed uint64, limitReached bool) {
	fmt.Fprintf(w, "profile=%s loaded=%d roms=%d addr=0x%04X pc_start=0x%04X steps=%d limit_reached=%v\n", cfg.profileName, loaded, roms, cfg.loadAt, cfg.startPC, executed, limitReached)
	fmt.Fprintf(w, "A=0x%02X B=0x%02X C=0x%02X D=0x%02X E=0x%02X H=0x%02X L=0x%02X\n", c.A, c.B, c.C, c.D, c.E, c.H, c.L)
	fmt.Fprintf(w, "PC=0x%04X SP=%d Halted=%v Stopped=%v\n", c.PC, c.SP, c.Halted, c.Stopped)
	fmt.Fprintf(w, "Flags C=%v Z=%v S=%v P=%v\n", c.Carry, c.Zero, c.Sign, c.Parity)
	fmt.Fprintf(w, "Stack=%s\n", formatStack(c.Stack))
}

func formatStack(stack [8]uint16) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, addr := range stack {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "0x%04X", addr)
	}
	b.WriteByte(']')
	return b.String()
}
