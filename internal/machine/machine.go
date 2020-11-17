package machine

import (
	"fmt"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
)

type Game []byte

type Machine struct {
	game    Game
	version int
	mem     *memory.Memory
}

func NewMachine(game Game) *Machine {
	return &Machine{
		game: game,
		mem:  memory.NewMemory(),
	}
}

func (m *Machine) WithVersion(version int) {
	m.version = version
}

func (m *Machine) Start() error {
	m.mem.Load(m.game)

	if m.version == 0 {
		// Inspect the game file for the version.
		m.version = int(m.mem.ReadByte(0))
	}

	fmt.Printf("z%d gamefile weighing in at %d bytes", m.version, m.mem.Size())

	return nil
}
