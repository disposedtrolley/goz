package machine

import (
	"fmt"
)

type Game []byte

type Machine struct {
	version int
	mem     []byte
}

func NewMachine(game Game) *Machine {
	return &Machine{
		mem: game,
	}
}

func (m *Machine) WithVersion(version int) {
	m.version = version
}

func (m *Machine) Start() error {
	if m.version == 0 {
		// Inspect the game file for the version.
		m.version = int(m.loadb(hVersion))
	}

	fmt.Printf("z%d gamefile weighing in at %d bytes\n", m.version, len(m.mem))

	fmt.Printf("beginning of static memory: %x\n", m.loadw(hStaticMemoryBegin))

	return nil
}

func (m *Machine) loadb(addr address) byte {
	return m.mem[addr]
}

func (m *Machine) loadw(addr address) uint16 {
	return uint16(m.mem[addr])<<8 | uint16(m.mem[addr+1])
}
