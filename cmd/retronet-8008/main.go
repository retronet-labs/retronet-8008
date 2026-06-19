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
)

const defaultStepLimit = uint64(1000)

type runConfig struct {
	binPath string
	loadAt  uint16
	startPC uint16
	steps   uint64
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

	program, err := os.ReadFile(cfg.binPath)
	if err != nil {
		fmt.Fprintf(stderr, "errore caricamento binario: %v\n", err)
		return 1
	}
	if err := validateLoadRange(cfg.loadAt, len(program)); err != nil {
		fmt.Fprintf(stderr, "errore caricamento binario: %v\n", err)
		return 1
	}

	c := cpu.NewCPU8008()
	mem := cpu.NewFlatMemory()
	ports := cpu.NewPorts()
	for i, b := range program {
		mem.Write(cfg.loadAt+uint16(i), b)
	}

	if err := c.Jam(mem, ports, cpu.JMP(), byte(cfg.startPC), byte(cfg.startPC>>8)); err != nil {
		fmt.Fprintf(stderr, "errore avvio CPU: %v\n", err)
		return 1
	}

	executed, limitReached, err := runSteps(c, mem, ports, cfg.steps)
	printDump(stdout, c, cfg, len(program), executed, limitReached)
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
	fs.Uint64Var(&cfg.steps, "steps", defaultStepLimit, "numero massimo di istruzioni da eseguire")

	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	if cfg.binPath == "" {
		fs.Usage()
		return cfg, errors.New("flag -bin obbligatorio")
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

func validateLoadRange(addr uint16, size int) error {
	if size > cpu.AddressSpaceSize {
		return fmt.Errorf("%d byte superano lo spazio indirizzabile %d byte", size, cpu.AddressSpaceSize)
	}
	if int(addr)+size > cpu.AddressSpaceSize {
		return fmt.Errorf("%d byte a 0x%04X superano 0x%04X", size, addr, cpu.AddressMask)
	}
	return nil
}

func runSteps(c *cpu.CPU8008, mem cpu.Memory, ioBus cpu.IO, limit uint64) (uint64, bool, error) {
	var executed uint64
	for executed < limit {
		if c.Halted || c.Stopped {
			return executed, false, nil
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

func printDump(w io.Writer, c *cpu.CPU8008, cfg runConfig, loaded int, executed uint64, limitReached bool) {
	fmt.Fprintf(w, "loaded=%d addr=0x%04X pc_start=0x%04X steps=%d limit_reached=%v\n", loaded, cfg.loadAt, cfg.startPC, executed, limitReached)
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
