package machine

import (
	"bytes"
	"errors"
	"testing"

	"retronet-8008/cpu"
)

func TestTerminalConsumesQueuedInputAndFallsBackToLatch(t *testing.T) {
	ioBus := NewCallbackIO()
	if err := ioBus.SetInput(TerminalInputPort, 0x7F); err != nil {
		t.Fatal(err)
	}
	terminal := NewTerminal(nil)
	terminal.QueueInputString("AB")
	if err := terminal.Attach(ioBus); err != nil {
		t.Fatal(err)
	}

	for i, want := range []byte{'A', 'B', 0x7F} {
		if got := ioBus.Input(TerminalInputPort); got != want {
			t.Fatalf("Input call %d = 0x%02X, want 0x%02X", i, got, want)
		}
	}
	if got := terminal.PendingInput(); got != 0 {
		t.Fatalf("PendingInput = %d, want 0", got)
	}
}

func TestTerminalWritesOutput(t *testing.T) {
	var output bytes.Buffer
	ioBus := NewCallbackIO()
	terminal := NewTerminal(&output)
	if err := terminal.Attach(ioBus); err != nil {
		t.Fatal(err)
	}

	ioBus.Output(TerminalOutputPort, 'O')
	ioBus.Output(TerminalOutputPort, 'K')

	if got := output.String(); got != "OK" {
		t.Fatalf("terminal output = %q, want OK", got)
	}
	if err := terminal.Err(); err != nil {
		t.Fatalf("Terminal.Err = %v, want nil", err)
	}
}

func TestTerminalCapturesWriterError(t *testing.T) {
	want := errors.New("write failed")
	terminal := NewTerminal(errorWriter{err: want})
	ioBus := NewCallbackIO()
	if err := terminal.Attach(ioBus); err != nil {
		t.Fatal(err)
	}

	ioBus.Output(TerminalOutputPort, 'X')
	if !errors.Is(terminal.Err(), want) {
		t.Fatalf("Terminal.Err = %v, want %v", terminal.Err(), want)
	}
}

func TestTerminalRejectsNilBus(t *testing.T) {
	if err := NewTerminal(nil).Attach(nil); !errors.Is(err, cpu.ErrNilIO) {
		t.Fatalf("Attach(nil) = %v, want ErrNilIO", err)
	}
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}
