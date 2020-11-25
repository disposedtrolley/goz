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
// at the provided offset. Abbreviations are resolved and included in the
// returned chars.
func (m *Machine) decodeZstring(offset memory.Address) (chars []zstring.ZChar) {
	chars = m.zStringToChars(offset)

	if m.version < 3 {
		// Only v3+ have abbreviations.
		return chars
	}

	for i := 0; i < len(chars); i++ {
		currChar := chars[i]
		if currChar >= 1 && currChar <= 3 && i < len(chars)-1 {
			// Abbreviation. Replace char[i] and char[i+1] with chars extracted
			// from the Z-string at the abbreviations table.
			nextChar := chars[i+1]
			abbreviationsTableOffset := 32*(currChar-1) + nextChar
			abbrevChars := m.zStringToChars(memory.HAbbreviationsTable + memory.Address(abbreviationsTableOffset))

			chars = append(chars, abbrevChars...)       // make room
			copy(chars[i+len(abbrevChars):], chars[i:]) // shift existing chars
			for j := 0; j < len(abbrevChars); j++ {
				// insert new chars
				chars[i+j] = abbrevChars[j]
			}
		}
	}

	return chars
}

// zStringToChars reads a Z-string beginning at the provided offset, returning
// an array of raw Z-characters found.
func (m *Machine) zStringToChars(offset memory.Address) (chars []zstring.ZChar) {
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

	return chars
}
