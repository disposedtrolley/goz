package machine

import (
	"fmt"
)

type Game []byte

type Machine struct {
	game    Game
	version *int
}

func NewMachine(game Game) *Machine {
	return &Machine{
		game: game,
	}
}

func (m *Machine) WithVersion(version int) {
	m.version = &version
}

func (m *Machine) Start() error {
	if m.version == nil {
		// Inspect the game file for the version.
	}

	fmt.Printf("z%d gamefile weighing in at %d bytes", m.version, len(m.game))
}
