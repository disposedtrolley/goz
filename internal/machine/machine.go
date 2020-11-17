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
		m.version = int(m.loadb(0))
	}

	fmt.Printf("z%d gamefile weighing in at %d bytes", m.version, len(m.mem))

	return nil
}

func (m *Machine) loadb(addr int) byte {
	return m.mem[addr]
}
