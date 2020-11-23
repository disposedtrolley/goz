package machine

import (
	"fmt"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
	"git.sr.ht/~disposedtrolley/go-zmachine/internal/zstring"
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

// decodeZstring returns an array of Z-characters found in memory beginning
// at the provided offset. Memory is read word-by-word, and the MSB of each
// word is checked to determine if the end of the Z-string has been reached.
func (m *Machine) decodeZstring(offset memory.Address) (chars []zstring.ZChar) {
	done := false
	for !done {
		word := m.mem.ReadWord(offset)
		//   --first byte-------   --second byte---
		//   7    6 5 4 3 2  1 0   7 6 5  4 3 2 1 0
		//   bit  --first--  --second---  --third--
		chars = append(chars, zstring.ZChar(word>>10&0x1F))
		chars = append(chars, zstring.ZChar(word>>5&0x1F))
		chars = append(chars, zstring.ZChar(word>>0&0x1F))

		done = word&(1<<15) != 0
		offset += 2
	}

	// TODO this needs to do some preprocessing to grab Z-strings from
	//      the abbreviations table if necessary. Ztoa will not have access
	//      to memory.
	return chars
}
