package machine

import (
	"errors"
	"fmt"
	"sync/atomic"

	"retronet-8008/cpu"
)

var (
	ErrNilCPU               = errors.New("cpu 8008 non inizializzata")
	ErrFrontPanelRunning    = errors.New("front panel in esecuzione")
	ErrInvalidRestartVector = errors.New("vettore RST non valido")
)

// PanelStopReason descrive perche' un run del front panel e' terminato.
type PanelStopReason string

const (
	PanelStoppedByCPU     PanelStopReason = "cpu-stopped"
	PanelStoppedByRequest PanelStopReason = "requested"
	PanelStoppedByLimit   PanelStopReason = "limit"
)

// PanelRunResult riassume un'esecuzione sincrona del front panel.
type PanelRunResult struct {
	Steps  uint64
	Reason PanelStopReason
}

// PanelStepObserver viene chiamato prima di ogni istruzione con una copia
// dello stato CPU, utile per trace e debugger senza esporre mutazioni.
type PanelStepObserver func(step uint64, state cpu.CPU8008) error

// FrontPanelState e' una fotografia delle luci e dei selettori modellati.
type FrontPanelState struct {
	CPU           cpu.CPU8008
	Switches      byte
	Address       uint16
	Data          byte
	Running       bool
	StopRequested bool
}

// FrontPanel coordina il core, la memoria e l'I/O come dispositivo esterno.
// Solo Stop e i selettori atomici sono pensati per uso concorrente con Run.
type FrontPanel struct {
	cpu    *cpu.CPU8008
	memory cpu.Memory
	io     cpu.IO

	switches      atomic.Uint32
	address       atomic.Uint32
	running       atomic.Bool
	stopRequested atomic.Bool
}

// NewFrontPanel crea un pannello sopra componenti gia' configurati.
func NewFrontPanel(c *cpu.CPU8008, memory cpu.Memory, ioBus cpu.IO) (*FrontPanel, error) {
	if c == nil {
		return nil, ErrNilCPU
	}
	if memory == nil {
		return nil, cpu.ErrNilMemory
	}
	return &FrontPanel{cpu: c, memory: memory, io: ioBus}, nil
}

// SetSwitches imposta gli otto switch dati.
func (p *FrontPanel) SetSwitches(value byte) {
	p.switches.Store(uint32(value))
}

// Switches legge gli otto switch dati.
func (p *FrontPanel) Switches() byte {
	return byte(p.switches.Load())
}

// SetAddress imposta i quattordici switch indirizzo.
func (p *FrontPanel) SetAddress(addr uint16) {
	p.address.Store(uint32(addr & cpu.AddressMask))
}

// Address legge l'indirizzo selezionato.
func (p *FrontPanel) Address() uint16 {
	return uint16(p.address.Load())
}

// Examine legge la memoria all'indirizzo selezionato.
func (p *FrontPanel) Examine() byte {
	return p.memory.Read(p.Address())
}

// Deposit scrive un byte all'indirizzo selezionato. Il bus mantiene la propria
// policy: per esempio una MemoryBus ignora una scrittura in ROM.
func (p *FrontPanel) Deposit(value byte) error {
	if p.running.Load() {
		return ErrFrontPanelRunning
	}
	p.memory.Write(p.Address(), value)
	return nil
}

// DepositSwitches deposita il valore degli switch dati.
func (p *FrontPanel) DepositSwitches() error {
	return p.Deposit(p.Switches())
}

// AttachSwitches collega gli switch dati a una porta input callback.
func (p *FrontPanel) AttachSwitches(ioBus *CallbackIO, port byte) error {
	if ioBus == nil {
		return cpu.ErrNilIO
	}
	return ioBus.OnInput(port, func(_ byte, _ byte) byte {
		return p.Switches()
	})
}

// Reset applica il reset storico della CPU senza modificare i selettori.
func (p *FrontPanel) Reset() error {
	if p.running.Load() {
		return ErrFrontPanelRunning
	}
	p.stopRequested.Store(false)
	p.cpu.Reset()
	return nil
}

// Jam forza una istruzione esterna, come il circuito di interrupt dell'8008.
func (p *FrontPanel) Jam(code byte, operands ...byte) error {
	if p.running.Load() {
		return ErrFrontPanelRunning
	}
	return p.cpu.Jam(p.memory, p.io, code, operands...)
}

// InterruptRST forza un restart vettorizzato da 0 a 7.
func (p *FrontPanel) InterruptRST(vector byte) error {
	if vector > 7 {
		return fmt.Errorf("%w: %d", ErrInvalidRestartVector, vector)
	}
	return p.Jam(cpu.RST(vector))
}

// Step esegue una singola istruzione se la CPU e' in stato running.
func (p *FrontPanel) Step() error {
	return p.cpu.Step(p.memory, p.io)
}

// Stop richiede in modo concorrente l'arresto del prossimo ciclo Run.
func (p *FrontPanel) Stop() {
	p.stopRequested.Store(true)
}

// Run esegue fino a HLT/stopped, richiesta esterna o limite istruzioni.
func (p *FrontPanel) Run(limit uint64, observer PanelStepObserver) (PanelRunResult, error) {
	if !p.running.CompareAndSwap(false, true) {
		return PanelRunResult{}, ErrFrontPanelRunning
	}
	defer p.running.Store(false)

	var result PanelRunResult
	for result.Steps < limit {
		if p.stopRequested.Swap(false) {
			result.Reason = PanelStoppedByRequest
			return result, nil
		}
		if p.cpu.Halted || p.cpu.Stopped {
			result.Reason = PanelStoppedByCPU
			return result, nil
		}
		if observer != nil {
			if err := observer(result.Steps, *p.cpu); err != nil {
				return result, err
			}
		}
		if err := p.Step(); err != nil {
			if errors.Is(err, cpu.ErrCPUStopped) {
				result.Reason = PanelStoppedByCPU
				return result, nil
			}
			return result, err
		}
		result.Steps++
	}

	if p.cpu.Halted || p.cpu.Stopped {
		result.Reason = PanelStoppedByCPU
	} else if p.stopRequested.Swap(false) {
		result.Reason = PanelStoppedByRequest
	} else {
		result.Reason = PanelStoppedByLimit
	}
	return result, nil
}

// Snapshot fotografa CPU, selettori e byte memoria selezionato.
func (p *FrontPanel) Snapshot() FrontPanelState {
	address := p.Address()
	return FrontPanelState{
		CPU:           *p.cpu,
		Switches:      p.Switches(),
		Address:       address,
		Data:          p.memory.Read(address),
		Running:       p.running.Load(),
		StopRequested: p.stopRequested.Load(),
	}
}
