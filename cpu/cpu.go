package cpu

// CPU8008 rappresenta lo stato interno essenziale del processore Intel 8008.
//
// Il core e' volutamente esplicito e didattico, nello stile di go-4004:
// i registri sono campi pubblici, mentre gli helper interni mantengono i
// vincoli architetturali come indirizzi a 14 bit e stack pointer a 3 bit.
type CPU8008 struct {
	A uint8 // Accumulatore
	B uint8 // Registro generale
	C uint8 // Registro generale
	D uint8 // Registro generale
	E uint8 // Registro generale
	H uint8 // Registro alto per HL; in indirizzamento usa solo 6 bit
	L uint8 // Registro basso per HL

	Carry  bool // Carry flag
	Zero   bool // Zero flag
	Sign   bool // Sign flag
	Parity bool // Parity flag

	PC uint16 // Program Counter, 14 bit

	// Stack interno indirizzi: 8 voci da 14 bit.
	// Una voce rappresenta il PC corrente, quindi CALL/RST hanno 7 livelli utili.
	Stack [8]uint16
	SP    uint8 // Stack pointer interno, 3 bit

	// Halted e Stopped modellano lo stato fermo storico del chip. Reset e HLT
	// fermano la CPU; Jam simula l'istruzione forzata da un interrupt esterno.
	Halted  bool
	Stopped bool

	// InstructionCount e StateCount includono anche istruzioni jammed.
	InstructionCount  uint64
	StateCount        uint64
	WaitStateCount    uint64
	LastTiming        InstructionTiming
	pendingWaitStates uint64
}

// NewCPU8008 crea una CPU nello stato di reset storico: registri azzerati e
// processore fermo, in attesa di una jam instruction esterna.
func NewCPU8008() *CPU8008 {
	c := &CPU8008{}
	c.Reset()
	return c
}

// Reset azzera registri, flag, PC, stack e SP.
//
// A differenza di molte CPU successive, il comportamento storico dell'8008 al
// power-on porta il processore in stato fermo. Per questo Reset imposta sia
// Halted sia Stopped a true; l'uscita dallo stop avviene tramite Jam.
func (c *CPU8008) Reset() {
	*c = CPU8008{
		Halted:  true,
		Stopped: true,
	}
}

// setPC aggiorna il program counter applicando il limite a 14 bit.
func (c *CPU8008) setPC(addr uint16) {
	c.PC = addr14(addr)
	c.Stack[c.SP] = c.PC
}

// setSP aggiorna lo stack pointer interno applicando il limite a 3 bit.
func (c *CPU8008) setSP(sp uint8) {
	c.SP = stackIndex(sp)
}

// setStack salva un indirizzo in una voce dello stack interno, mascherandolo a
// 14 bit.
func (c *CPU8008) setStack(slot uint8, addr uint16) {
	c.Stack[stackIndex(slot)] = addr14(addr)
}
