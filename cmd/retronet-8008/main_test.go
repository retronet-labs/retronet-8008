package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"retronet-8008/cpu"
)

func TestRunLoadsProgramAndPrintsDump(t *testing.T) {
	bin := writeTempProgram(t, []byte{
		cpu.LI(cpu.RegA), 0x2A,
		cpu.HLT(),
	})
	var stdout, stderr bytes.Buffer

	code := run([]string{"-bin", bin, "-steps", "8"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit = %d, stderr = %s", code, stderr.String())
	}
	out := stdout.String()
	wantParts := []string{
		"loaded=3 addr=0x0000 pc_start=0x0000 steps=2 limit_reached=false",
		"A=0x2A",
		"PC=0x0003",
		"Halted=true Stopped=true",
	}
	for _, part := range wantParts {
		if !strings.Contains(out, part) {
			t.Fatalf("output missing %q:\n%s", part, out)
		}
	}
}

func TestRunLoadsAtAddressAndStartsAtPC(t *testing.T) {
	bin := writeTempProgram(t, []byte{
		cpu.LI(cpu.RegB), 0x33,
		cpu.L(cpu.RegA, cpu.RegB),
		cpu.HLT(),
	})
	var stdout, stderr bytes.Buffer

	code := run([]string{"-bin", bin, "-addr", "0x0010", "-pc", "0x0010", "-steps", "8"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit = %d, stderr = %s", code, stderr.String())
	}
	out := stdout.String()
	wantParts := []string{
		"loaded=4 addr=0x0010 pc_start=0x0010 steps=3 limit_reached=false",
		"A=0x33 B=0x33",
		"PC=0x0014",
	}
	for _, part := range wantParts {
		if !strings.Contains(out, part) {
			t.Fatalf("output missing %q:\n%s", part, out)
		}
	}
}

func TestRunReportsLimitReached(t *testing.T) {
	bin := writeTempProgram(t, []byte{cpu.NOP(), cpu.NOP(), cpu.NOP()})
	var stdout, stderr bytes.Buffer

	code := run([]string{"-bin", bin, "-steps", "2"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit = %d, stderr = %s", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "steps=2 limit_reached=true") {
		t.Fatalf("output = %s, want limit reached", out)
	}
	if !strings.Contains(out, "PC=0x0002") {
		t.Fatalf("output = %s, want PC after two NOPs", out)
	}
}

func TestRunDisassemblesWithoutExecution(t *testing.T) {
	bin := writeTempProgram(t, []byte{
		cpu.LI(cpu.RegA), 0x2A,
		cpu.JMP(), 0x00, 0x10,
	})
	var stdout, stderr bytes.Buffer

	code := run([]string{"-bin", bin, "-disasm", "2"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit = %d, stderr = %s", code, stderr.String())
	}
	out := stdout.String()
	want := "0000: 06 2A    LAI #0x2A\n0002: 44 00 10 JMP 0x1000\n"
	if out != want {
		t.Fatalf("stdout = %q, want %q", out, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %s, want empty", stderr.String())
	}
}

func TestRunTracesExecution(t *testing.T) {
	bin := writeTempProgram(t, []byte{
		cpu.LI(cpu.RegA), 0x2A,
		cpu.HLT(),
	})
	var stdout, stderr bytes.Buffer

	code := run([]string{"-bin", bin, "-steps", "8", "-trace"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit = %d, stderr = %s", code, stderr.String())
	}
	out := stdout.String()
	wantParts := []string{
		"trace=0 0000: 06 2A    LAI #0x2A\n",
		"trace=1 0002: 00       HLT\n",
		"loaded=3 addr=0x0000 pc_start=0x0000 steps=2 limit_reached=false",
		"A=0x2A",
	}
	for _, part := range wantParts {
		if !strings.Contains(out, part) {
			t.Fatalf("output missing %q:\n%s", part, out)
		}
	}
}

func TestRunTraceHonorsStepLimit(t *testing.T) {
	bin := writeTempProgram(t, []byte{cpu.NOP(), cpu.NOP()})
	var stdout, stderr bytes.Buffer

	code := run([]string{"-bin", bin, "-steps", "1", "-trace"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit = %d, stderr = %s", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "trace=0 0000: C0       NOP\n") {
		t.Fatalf("output = %s, want trace for first NOP", out)
	}
	if strings.Contains(out, "trace=1") {
		t.Fatalf("output = %s, did not expect trace beyond step limit", out)
	}
	if !strings.Contains(out, "steps=1 limit_reached=true") {
		t.Fatalf("output = %s, want limit reached", out)
	}
}

func TestRunRequiresBinaryPath(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(nil, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("run exit = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "flag -bin obbligatorio") {
		t.Fatalf("stderr = %s, want missing -bin error", stderr.String())
	}
}

func TestRunHelpExitsCleanly(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"-h"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run -h exit = %d, want 0", code)
	}
	if !strings.Contains(stderr.String(), "Usage of retronet-8008") {
		t.Fatalf("stderr = %s, want usage", stderr.String())
	}
}

func TestParseAddressRejectsOutside14BitSpace(t *testing.T) {
	if _, err := parseAddress("0x4000"); err == nil {
		t.Fatal("parseAddress(0x4000) = nil, want error")
	}
}

func writeTempProgram(t *testing.T, program []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "program.bin")
	if err := os.WriteFile(path, program, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
