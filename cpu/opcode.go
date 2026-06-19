package cpu

// ExecuteFunc e' la funzione esecutiva associata a una voce del decoder.
type ExecuteFunc func(c *CPU8008, mem Memory, io IO, inst Instruction) error

// MachineCycle identifica il tipo di ciclo esterno dichiarato dall'8008.
type MachineCycle string

const (
	CyclePCI MachineCycle = "PCI" // fetch primo byte istruzione
	CyclePCR MachineCycle = "PCR" // lettura memoria o byte aggiuntivo
	CyclePCW MachineCycle = "PCW" // scrittura memoria
	CyclePCC MachineCycle = "PCC" // comando I/O
)

// Opcode descrive una voce della tabella decoder 8008.
type Opcode struct {
	Code       byte
	Mnemonic   string
	Length     byte
	MinStates  byte
	States     byte
	CycleCount byte
	Cycles     [3]MachineCycle
	Execute    ExecuteFunc
}

// MachineCycles restituisce una copia dei cicli usati dall'opcode.
func (o Opcode) MachineCycles() []MachineCycle {
	cycles := make([]MachineCycle, o.CycleCount)
	copy(cycles, o.Cycles[:o.CycleCount])
	return cycles
}

// InstructionTiming descrive il costo effettivo dell'ultima istruzione.
type InstructionTiming struct {
	States      byte
	WaitStates  uint64
	CycleCount  byte
	Cycles      [3]MachineCycle
	Conditional bool
	Taken       bool
}

// MachineCycles restituisce una copia dei cicli effettivi dell'istruzione.
func (t InstructionTiming) MachineCycles() []MachineCycle {
	cycles := make([]MachineCycle, t.CycleCount)
	copy(cycles, t.Cycles[:t.CycleCount])
	return cycles
}

// Instruction contiene l'opcode fetchato e gli eventuali byte operando.
type Instruction struct {
	PC           uint16
	Opcode       Opcode
	Operands     [2]byte
	OperandCount byte
}

func unimplementedExecute(_ *CPU8008, _ Memory, _ IO, inst Instruction) error {
	return &UnimplementedOpcodeError{
		PC:       inst.PC,
		Opcode:   inst.Opcode.Code,
		Mnemonic: inst.Opcode.Mnemonic,
		Length:   inst.Opcode.Length,
	}
}
