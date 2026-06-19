package cpu

import "testing"

func newRunningCPU(t *testing.T) *CPU8008 {
	t.Helper()

	c := NewCPU8008()
	jamNOP(t, c)
	return c
}

func jamNOP(t *testing.T, c *CPU8008) {
	t.Helper()

	if err := c.Jam(nil, nil, NOP()); err != nil {
		t.Fatalf("Jam(NOP) = %v, want nil", err)
	}
}
