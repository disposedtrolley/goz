package memory

import (
	"fmt"
)

type Address uint

type Memory struct {
	content []byte
}

// NewMemory returns new memory initialised with contentToLoad.
func NewMemory(contentToLoad []byte) *Memory {
	return &Memory{
		content: contentToLoad,
	}
}

// ReadByte returns the byte at address a, or panics if a is
// out of bounds.
func (m *Memory) ReadByte(a Address) uint8 {
	m.boundsCheck(a)

	return m.content[a]
}

// ReadWord returns the byte at address a concatenated with the
// byte following it, or panics if a or a + 1 is out of bounds.
func (m *Memory) ReadWord(a Address) uint16 {
	m.boundsCheck(a)
	m.boundsCheck(a + 1)

	return uint16(m.content[a])<<8 | uint16(m.content[a+1])
}

func (m *Memory) WriteByte(a Address, b byte) {
	m.boundsCheck(a)
	m.readOnlyCheck(a)

	m.content[a] = b
}

func (m *Memory) WriteWord(a Address, w uint16) {
	m.boundsCheck(a)
	m.boundsCheck(a + 1)
	m.readOnlyCheck(a)

	m.content[a] = byte(w >> 8)
	m.content[a+1] = byte(w)
}

// boundsCheck ensures that a is within the region of m's
// content.
func (m *Memory) boundsCheck(a Address) {
	if int(a) >= len(m.content) {
		panic(fmt.Errorf("attempted to access address %x which is outside of initialised memory", a))
	}
}

func (m *Memory) readOnlyCheck(a Address) {
	if a >= HStaticMemoryBegin {
		panic(fmt.Errorf("attempted to write to address %x which is in static memory", a))
	}
}

// Size returns the length of the initialised memory.
func (m *Memory) Size() int {
	return len(m.content)
}
