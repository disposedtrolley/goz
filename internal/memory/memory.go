package memory

type Memory struct {
	bits []byte
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Load(bits []byte) {
	m.bits = bits
}

func (m *Memory) Size() int {
	return len(m.bits)
}

func (m *Memory) ReadByte(addr int) byte {
	return m.bits[addr]
}

func (m *Memory) ReadBit(addr int, bit int) byte {
	b := m.bits[addr]
	return (b & (1<<bit - 1))
}
