package machine

import (
	"io"
	"sync"

	"retronet-8008/cpu"
)

const (
	TerminalInputPort  byte = 0
	TerminalOutputPort byte = 8
)

// Terminal e' una periferica ASCII buffered collegabile a CallbackIO.
// Quando la coda input e' vuota restituisce il valore latched della porta.
type Terminal struct {
	mu     sync.Mutex
	input  []byte
	output io.Writer
	err    error
}

// NewTerminal crea un terminale che scrive su output. Un output nil scarta i
// byte, mantenendo utile la sola parte input.
func NewTerminal(output io.Writer) *Terminal {
	if output == nil {
		output = io.Discard
	}
	return &Terminal{output: output}
}

// Attach collega il terminale alle porte convenzionali input 0 e output 8.
func (t *Terminal) Attach(ioBus *CallbackIO) error {
	if ioBus == nil {
		return cpu.ErrNilIO
	}
	if err := ioBus.OnInput(TerminalInputPort, t.readInput); err != nil {
		return err
	}
	return ioBus.OnOutput(TerminalOutputPort, t.writeOutput)
}

// QueueInput accoda una copia dei byte che verranno consumati da INP 0.
func (t *Terminal) QueueInput(data []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.input = append(t.input, data...)
}

// QueueInputString accoda testo senza applicare conversioni di encoding.
func (t *Terminal) QueueInputString(value string) {
	t.QueueInput([]byte(value))
}

// PendingInput restituisce il numero di byte ancora disponibili.
func (t *Terminal) PendingInput() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.input)
}

// Err restituisce il primo errore prodotto dal writer di output.
func (t *Terminal) Err() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.err
}

func (t *Terminal) readInput(_ byte, latched byte) byte {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.input) == 0 {
		return latched
	}
	value := t.input[0]
	t.input = t.input[1:]
	return value
}

func (t *Terminal) writeOutput(_ byte, value byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.err != nil {
		return
	}
	_, t.err = t.output.Write([]byte{value})
}
