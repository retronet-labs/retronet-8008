package cpu

// AddressSpaceSize e' la memoria direttamente indirizzabile dall'Intel 8008:
// 16 KB, derivati dal program counter e dal bus indirizzi a 14 bit.
const AddressSpaceSize = int(AddressMask) + 1

// Memory e' il bus memoria visto dal core CPU.
//
// L'8008 separa memoria e I/O: queste funzioni riguardano solo lo spazio
// diretto 0x0000-0x3FFF.
type Memory interface {
	Read(addr uint16) byte
	Write(addr uint16, value byte)
}

// FlatMemory modella una RAM/ROM piatta da 16 KB.
//
// Per ora non distingue RAM e ROM: e' il supporto minimo per test, esempi e
// futuro fetch istruzioni. Gli indirizzi sono sempre mascherati a 14 bit.
type FlatMemory struct {
	Data [AddressSpaceSize]byte
}

// NewFlatMemory crea una memoria piatta inizializzata a zero.
func NewFlatMemory() *FlatMemory {
	return &FlatMemory{}
}

// Read legge un byte dalla memoria diretta, mascherando l'indirizzo a 14 bit.
func (m *FlatMemory) Read(addr uint16) byte {
	return m.Data[addr14(addr)]
}

// Write scrive un byte nella memoria diretta, mascherando l'indirizzo a 14 bit.
func (m *FlatMemory) Write(addr uint16, value byte) {
	m.Data[addr14(addr)] = value
}
