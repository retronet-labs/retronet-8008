package machine

import "retronet-8008/cpu"

// InputCallback viene chiamata quando la CPU legge una porta di input.
//
// Il secondo argomento e' il valore latched sulla porta. La callback puo'
// restituirlo invariato o calcolare un valore dinamico.
type InputCallback func(port byte, latched byte) byte

// OutputCallback viene chiamata quando la CPU scrive una porta di output.
type OutputCallback func(port byte, value byte)

// CallbackIO implementa cpu.IO usando latch semplici e callback opzionali.
//
// Serve come ponte tra il core CPU, che conosce solo cpu.IO, e i profili
// macchina, che possono collegare terminali, front panel o periferiche simulate.
type CallbackIO struct {
	inputs          [8]byte
	outputs         [24]byte
	inputCallbacks  [8]InputCallback
	outputCallbacks [24]OutputCallback
}

// NewCallbackIO crea un bus I/O callback inizializzato a zero.
func NewCallbackIO() *CallbackIO {
	return &CallbackIO{}
}

// NewIO crea il bus I/O associato al profilo.
//
// In questa milestone tutti i profili usano CallbackIO. Il metodo mantiene il
// punto di estensione nel posto giusto per future periferiche SCELBI/Intellec.
func (Profile) NewIO() *CallbackIO {
	return NewCallbackIO()
}

// SetInput imposta il valore latched su una porta di input valida.
func (io *CallbackIO) SetInput(port byte, value byte) error {
	if err := cpu.ValidateInputPort(port); err != nil {
		return err
	}
	io.inputs[port] = value
	return nil
}

// OnInput registra o rimuove una callback per una porta input.
func (io *CallbackIO) OnInput(port byte, callback InputCallback) error {
	if err := cpu.ValidateInputPort(port); err != nil {
		return err
	}
	io.inputCallbacks[port] = callback
	return nil
}

// OnOutput registra o rimuove una callback per una porta output.
func (io *CallbackIO) OnOutput(port byte, callback OutputCallback) error {
	if err := cpu.ValidateOutputPort(port); err != nil {
		return err
	}
	io.outputCallbacks[port-8] = callback
	return nil
}

// Input legge una porta input. Porte non valide restituiscono zero.
func (io *CallbackIO) Input(port byte) byte {
	if !cpu.IsInputPort(port) {
		return 0
	}
	latched := io.inputs[port]
	if callback := io.inputCallbacks[port]; callback != nil {
		return callback(port, latched)
	}
	return latched
}

// Output scrive una porta output e invoca l'eventuale callback associata.
func (io *CallbackIO) Output(port byte, value byte) {
	if !cpu.IsOutputPort(port) {
		return
	}
	index := port - 8
	io.outputs[index] = value
	if callback := io.outputCallbacks[index]; callback != nil {
		callback(port, value)
	}
}

// OutputValue restituisce il valore latched su una porta output valida.
func (io *CallbackIO) OutputValue(port byte) (byte, error) {
	if err := cpu.ValidateOutputPort(port); err != nil {
		return 0, err
	}
	return io.outputs[port-8], nil
}
