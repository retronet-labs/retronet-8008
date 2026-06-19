package cpu

import (
	"errors"
	"fmt"
)

var (
	// ErrNilMemory segnala che Step e' stato chiamato senza bus memoria.
	ErrNilMemory = errors.New("memoria 8008 non inizializzata")

	// ErrUnimplementedOpcode segnala un opcode valido nel decoder ma non ancora eseguito.
	ErrUnimplementedOpcode = errors.New("opcode 8008 non implementato")

	// ErrInvalidInputPort segnala una porta di input fuori dal range 0..7.
	ErrInvalidInputPort = errors.New("porta input 8008 non valida")

	// ErrInvalidOutputPort segnala una porta di output fuori dal range 8..31.
	ErrInvalidOutputPort = errors.New("porta output 8008 non valida")
)

// UnimplementedOpcodeError conserva il contesto dell'opcode letto da Step.
type UnimplementedOpcodeError struct {
	PC       uint16
	Opcode   byte
	Mnemonic string
	Length   byte
}

func (e *UnimplementedOpcodeError) Error() string {
	return fmt.Sprintf("%v: PC=0x%04X opcode=0x%02X mnemonic=%s length=%d", ErrUnimplementedOpcode, e.PC, e.Opcode, e.Mnemonic, e.Length)
}

func (e *UnimplementedOpcodeError) Unwrap() error {
	return ErrUnimplementedOpcode
}
