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
