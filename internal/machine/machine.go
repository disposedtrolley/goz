package machine

import (
	"fmt"
)

type Game []byte

type Machine struct {
	version int
	mem     *memory
}

func NewMachine(game Game) *Machine {
	return &Machine{
		mem: newMemory(game),
	}
}

func (m *Machine) WithVersion(version int) {
	m.version = version
}

func (m *Machine) Start() error {
	if m.version == 0 {
		// Inspect the game file for the version.
		m.version = int(m.mem.readByte(hVersion))
	}

	fmt.Printf("z%d gamefile weighing in at %d bytes\n", m.version, m.mem.size())
	fmt.Printf("beginning of static memory: %x\n", m.mem.readWord(hStaticMemoryBegin))

	return nil
}
