package cpu

import "errors"

var (
	// ErrInvalidInputPort segnala una porta di input fuori dal range 0..7.
	ErrInvalidInputPort = errors.New("porta input 8008 non valida")

	// ErrInvalidOutputPort segnala una porta di output fuori dal range 8..31.
	ErrInvalidOutputPort = errors.New("porta output 8008 non valida")
)
