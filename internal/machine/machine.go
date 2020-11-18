package machine

import (
	"fmt"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
)

type Game []byte

type Machine struct {
	version int
	mem     *memory.Memory
}

func NewMachine(game Game) *Machine {
	return &Machine{
		mem: memory.NewMemory(game),
	}
}

func (m *Machine) WithVersion(version int) {
	m.version = version
}

func (m *Machine) Start() error {
	if m.version == 0 {
		// Inspect the game file for the version.
		m.version = int(m.mem.ReadByte(memory.HVersion))
	}

	fmt.Printf("z%d gamefile weighing in at %d bytes\n", m.version, m.mem.Size())
	fmt.Printf(`
beginning of:
  static memory: %x
  high memory: %x
`, m.mem.ReadWord(memory.HStaticMemoryBegin), m.mem.ReadWord(memory.HHighMemoryBegin))

	return nil
}
