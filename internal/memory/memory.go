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

// ByteAddress returns addr as a byte address.
func (m *Memory) ByteAddress(addr Address) Address {
	return addr
}

// WordAddress returns addr as a word address.
func (m *Memory) WordAddress(addr Address) Address {
	return addr * 2
}

// PackedAddress returns addr as a packed address. Versions 6 and 7 have
// additional constants for routine and string offsets and are not supported by
// this function.
func (m *Memory) PackedAddress(addr Address, version int) Address {
	switch v := version; {
	case v <= 3:
		return addr * 2
	case v <= 5:
		return addr * 4
	case v == 8:
		return addr * 8
	default:
		panic(fmt.Errorf("method PackedAddress not supported for version %v", version))
	}
}

// PackedAddressRoutine returns addr as a packed routine address. Only supports
// version 6 and 7.
func (m *Memory) PackedAddressRoutine(addr Address, version int) Address {
	if version != 6 && version != 7 {
		panic(fmt.Errorf("method not supported for version %v", version))
	}

	routineOffset := m.ReadWord(HRoutinesOffset)

	return Address(uint16(addr)*4 + 8*routineOffset)
}

// PackedAddressString returns addr as a packed string address. Only supports
// version 6 and 7.
func (m *Memory) PackedAddressString(addr Address, version int) Address {
	if version != 6 && version != 7 {
		panic(fmt.Errorf("method not supported for version %v", version))
	}

	stringOffset := m.ReadWord(HStaticStringsOffset)

	return Address(uint16(addr)*4 + 8*stringOffset)
}
