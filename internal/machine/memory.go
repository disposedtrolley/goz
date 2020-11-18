package machine

import (
	"fmt"
)

type address uint

// constants for some addresses of interest.
const (
	hVersion           address = 0
	hStaticMemoryBegin address = 0x0E
	hHighMemoryBegin   address = 0x04
)

type memory struct {
	content []byte
}

// newMemory returns new memory initialised with contentToLoad.
func newMemory(contentToLoad []byte) *memory {
	return &memory{
		content: contentToLoad,
	}
}

// readByte returns the byte at address a, or panics if a is
// out of bounds.
func (m *memory) readByte(a address) uint8 {
	m.boundsCheck(a)

	return m.content[a]
}

// readWord returns the byte at address a concatenated with the
// byte following it, or panics if a or a + 1 is out of bounds.
func (m *memory) readWord(a address) uint16 {
	m.boundsCheck(a)
	m.boundsCheck(a + 1)

	return uint16(m.content[a])<<8 | uint16(m.content[a+1])
}

func (m *memory) writeByte(a address, b byte) {
	m.boundsCheck(a)
	m.readOnlyCheck(a)

	m.content[a] = b
}

func (m *memory) writeWord(a address, w uint16) {
	m.boundsCheck(a)
	m.boundsCheck(a + 1)
	m.readOnlyCheck(a)

	m.content[a] = byte(w >> 8)
	m.content[a+1] = byte(w)
}

// boundsCheck ensures that a is within the region of m's
// content.
func (m *memory) boundsCheck(a address) {
	if int(a) >= len(m.content) {
		panic(fmt.Errorf("attempted to access address %x which is outside of initialised memory", a))
	}
}

func (m *memory) readOnlyCheck(a address) {
	if a >= hStaticMemoryBegin {
		panic(fmt.Errorf("attempted to write to address %x which is in static memory", a))
	}
}

// size returns the length of the initialised memory.
func (m *memory) size() int {
	return len(m.content)
}
