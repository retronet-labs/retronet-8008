package cpu

// addr14 conserva solo i 14 bit di indirizzo visibili all'Intel 8008.
func addr14(v uint16) uint16 {
	return v & AddressMask
}

// stackIndex conserva solo i 3 bit del puntatore stack interno.
func stackIndex(v uint8) uint8 {
	return v & 0x07
}

// hlAddress costruisce l'indirizzo puntato dalla coppia H/L.
//
// Sul 8008 H contribuisce solo con i 6 bit bassi; i bit 6 e 7 sono ignorati
// quando HL viene usato per accedere al pseudo-registro M.
func hlAddress(h, l uint8) uint16 {
	return (uint16(h&0x3F) << 8) | uint16(l)
}

// HL restituisce l'indirizzo di memoria puntato dai registri H e L.
func (c *CPU8008) HL() uint16 {
	return hlAddress(c.H, c.L)
}

// regBits restituisce i tre bit di codifica di un registro 8008.
func regBits(r Register) byte {
	return byte(r) & 0x07
}

// condBits restituisce i due bit di codifica di una condizione 8008.
func condBits(c Condition) byte {
	return byte(c) & 0x03
}

// NOP restituisce un opcode no-op della famiglia load.
//
// Sul decoder 8008 i load con sorgente e destinazione uguali non modificano lo
// stato. Usiamo LAA come no-op leggibile per test ed esempi.
func NOP() byte { return L(RegA, RegA) }

// L costruisce un opcode di trasferimento tra registri o pseudo-registro M.
//
// Esempi:
//
//	L(RegA, RegB) carica B in A.
//	L(RegM, RegA) scrive A nella memoria puntata da HL.
//	L(RegA, RegM) legge in A la memoria puntata da HL.
func L(dst, src Register) byte {
	return 0xC0 | (regBits(dst) << 3) | regBits(src)
}

// LI costruisce il primo byte di un load immediato.
//
// Il byte immediato va scritto subito dopo l'opcode in memoria. Con dst=RegM
// produce LMI, cioe' scrittura immediata nella memoria puntata da HL.
func LI(dst Register) byte {
	return 0x06 | (regBits(dst) << 3)
}

// INR costruisce l'opcode di incremento registro.
//
// L'8008 definisce questa famiglia per B, C, D, E, H e L. L'incremento aggiorna
// Zero, Sign e Parity, ma non modifica Carry.
func INR(r Register) byte {
	return regBits(r) << 3
}

// DCR costruisce l'opcode di decremento registro.
//
// L'8008 definisce questa famiglia per B, C, D, E, H e L. Il decremento aggiorna
// Zero, Sign e Parity, ma non modifica Carry.
func DCR(r Register) byte {
	return 0x01 | (regBits(r) << 3)
}

// AD somma all'accumulatore il registro o M indicato.
func AD(src Register) byte { return aluRegisterOpcode(0, src) }

// AC somma all'accumulatore il registro o M indicato, includendo Carry.
func AC(src Register) byte { return aluRegisterOpcode(1, src) }

// SU sottrae dall'accumulatore il registro o M indicato.
func SU(src Register) byte { return aluRegisterOpcode(2, src) }

// SB sottrae dall'accumulatore il registro o M indicato, includendo il borrow in Carry.
func SB(src Register) byte { return aluRegisterOpcode(3, src) }

// ND applica AND tra accumulatore e registro o M indicato.
func ND(src Register) byte { return aluRegisterOpcode(4, src) }

// XR applica XOR tra accumulatore e registro o M indicato.
func XR(src Register) byte { return aluRegisterOpcode(5, src) }

// OR applica OR tra accumulatore e registro o M indicato.
func OR(src Register) byte { return aluRegisterOpcode(6, src) }

// CP confronta accumulatore e registro o M indicato senza modificare A.
func CP(src Register) byte { return aluRegisterOpcode(7, src) }

// ADI costruisce il primo byte di add immediato.
func ADI() byte { return aluImmediateOpcode(0) }

// ACI costruisce il primo byte di add immediato con Carry.
func ACI() byte { return aluImmediateOpcode(1) }

// SUI costruisce il primo byte di subtract immediato.
func SUI() byte { return aluImmediateOpcode(2) }

// SBI costruisce il primo byte di subtract immediato con borrow.
func SBI() byte { return aluImmediateOpcode(3) }

// NDI costruisce il primo byte di AND immediato.
func NDI() byte { return aluImmediateOpcode(4) }

// XRI costruisce il primo byte di XOR immediato.
func XRI() byte { return aluImmediateOpcode(5) }

// ORI costruisce il primo byte di OR immediato.
func ORI() byte { return aluImmediateOpcode(6) }

// CPI costruisce il primo byte di compare immediato.
func CPI() byte { return aluImmediateOpcode(7) }

func aluRegisterOpcode(group byte, src Register) byte {
	return 0x80 | ((group & 0x07) << 3) | regBits(src)
}

func aluImmediateOpcode(group byte) byte {
	return 0x04 | ((group & 0x07) << 3)
}

// RLC ruota A a sinistra: bit 7 va in bit 0 e in Carry.
func RLC() byte { return 0x02 }

// RRC ruota A a destra: bit 0 va in bit 7 e in Carry.
func RRC() byte { return 0x0A }

// RAL ruota A a sinistra attraverso Carry.
func RAL() byte { return 0x12 }

// RAR ruota A a destra attraverso Carry.
func RAR() byte { return 0x1A }

// JMP costruisce un jump incondizionato a 3 byte.
//
// I due byte indirizzo vanno scritti dopo l'opcode, low byte prima e high byte
// poi. Del byte alto vengono usati solo i 6 bit bassi.
func JMP() byte { return 0x44 }

// JF costruisce un jump condizionato se il flag indicato e' false.
func JF(cond Condition) byte { return 0x40 | (condBits(cond) << 3) }

// JT costruisce un jump condizionato se il flag indicato e' true.
func JT(cond Condition) byte { return 0x60 | (condBits(cond) << 3) }

// CAL costruisce una call incondizionata a 3 byte.
func CAL() byte { return 0x46 }

// CF costruisce una call condizionata se il flag indicato e' false.
func CF(cond Condition) byte { return 0x42 | (condBits(cond) << 3) }

// CT costruisce una call condizionata se il flag indicato e' true.
func CT(cond Condition) byte { return 0x62 | (condBits(cond) << 3) }

// RET costruisce un return incondizionato.
func RET() byte { return 0x07 }

// RF costruisce un return condizionato se il flag indicato e' false.
func RF(cond Condition) byte { return 0x03 | (condBits(cond) << 3) }

// RT costruisce un return condizionato se il flag indicato e' true.
func RT(cond Condition) byte { return 0x23 | (condBits(cond) << 3) }

// RST costruisce una restart/call a vettore pagina zero n*8.
func RST(n byte) byte { return 0x05 | ((n & 0x07) << 3) }
