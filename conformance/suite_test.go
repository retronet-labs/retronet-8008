package conformance

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"retronet-8008/cpu"
	"retronet-8008/machine"
)

func TestSyntheticSuitePasses(t *testing.T) {
	result := RunSuite(SyntheticSuite())
	if result.Failed != 0 || result.Passed != len(result.Cases) {
		for _, test := range result.Cases {
			if !test.Passed {
				t.Errorf("%s: %s", test.Name, test.Error)
			}
		}
		t.Fatalf("suite = passed %d failed %d", result.Passed, result.Failed)
	}
}

func TestRunSuiteKeepsFailures(t *testing.T) {
	cases := []Case{
		{Name: "pass", Program: []byte{cpu.HLT()}},
		{
			Name:    "fail",
			Program: []byte{cpu.HLT()},
			Verify: func(*Context, machine.DebugRunResult) error {
				return os.ErrInvalid
			},
		},
	}

	result := RunSuite(cases)
	if result.Passed != 1 || result.Failed != 1 || !strings.Contains(result.Cases[1].Error, os.ErrInvalid.Error()) {
		t.Fatalf("suite = %+v", result)
	}
}

func TestVerifyLocalROM(t *testing.T) {
	data := []byte{0x41, 0x51, 0x00}
	path := filepath.Join(t.TempDir(), "rom.bin")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(data)
	expectedHash := hex.EncodeToString(hash[:])

	result, err := VerifyLocalROM(path, ROMExpectation{
		Name:           "io-smoke",
		ExpectedSize:   int64(len(data)),
		ExpectedSHA256: expectedHash,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Matches || result.ActualSHA256 != expectedHash {
		t.Fatalf("verification = %+v", result)
	}
}

func TestVerifyLocalROMReportsMismatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rom.bin")
	if err := os.WriteFile(path, []byte{0x00}, 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := VerifyLocalROM(path, ROMExpectation{ExpectedSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if result.Matches {
		t.Fatalf("verification = %+v, want mismatch", result)
	}
}
