package cpu

import "fmt"

// IO e' il bus di input/output separato dalla memoria.
//
// L'8008 dispone di 8 porte di input e 24 porte di output. Le istruzioni INP e
// OUT usano i validatori di questo file prima di accedere al bus.
type IO interface {
	Input(port byte) byte
	Output(port byte, value byte)
}

// Ports e' una implementazione semplice dell'I/O 8008.
type Ports struct {
	InputPorts  [8]byte
	OutputPorts [24]byte
}

// NewPorts crea porte I/O inizializzate a zero.
func NewPorts() *Ports {
	return &Ports{}
}

// IsInputPort restituisce true per le porte di input valide, 0..7.
func IsInputPort(port byte) bool {
	return port <= 7
}

// IsOutputPort restituisce true per le porte di output valide, 8..31.
func IsOutputPort(port byte) bool {
	return port >= 8 && port <= 31
}

// ValidateInputPort controlla che port sia nel range di input 0..7.
func ValidateInputPort(port byte) error {
	if IsInputPort(port) {
		return nil
	}
	return fmt.Errorf("%w: %d", ErrInvalidInputPort, port)
}

// ValidateOutputPort controlla che port sia nel range di output 8..31.
func ValidateOutputPort(port byte) error {
	if IsOutputPort(port) {
		return nil
	}
	return fmt.Errorf("%w: %d", ErrInvalidOutputPort, port)
}

// SetInput imposta il valore letto da una porta di input valida.
func (p *Ports) SetInput(port byte, value byte) error {
	if err := ValidateInputPort(port); err != nil {
		return err
	}
	p.InputPorts[port] = value
	return nil
}

// Input legge una porta di input. Le porte non valide restituiscono zero.
func (p *Ports) Input(port byte) byte {
	if !IsInputPort(port) {
		return 0
	}
	return p.InputPorts[port]
}

// Output scrive una porta di output. Le porte non valide vengono ignorate.
func (p *Ports) Output(port byte, value byte) {
	if !IsOutputPort(port) {
		return
	}
	p.OutputPorts[port-8] = value
}

// OutputValue restituisce il valore latched su una porta di output valida.
func (p *Ports) OutputValue(port byte) (byte, error) {
	if err := ValidateOutputPort(port); err != nil {
		return 0, err
	}
	return p.OutputPorts[port-8], nil
}
