package machine

import (
	"fmt"
	"strings"

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

// decodeZstring returns a string representing the decoded Z-characters found
// at the provided memory offset.
func (m *Machine) decodeZstring(offset memory.Address) string {
	chars := m.zStringToChars(offset)
	var output strings.Builder
	lock := false
	currentAlphabet := zstring.A0
	for i := 0; i < len(chars); i++ {
		char := chars[i]

		// ZSCII
		if currentAlphabet == zstring.A2 && char == 6 {
			zsciiCode := (uint16(chars[i+1]) << 5) | uint16(chars[i+2])
			output.WriteByte(byte(zsciiCode))
			i += 2
			currentAlphabet = zstring.A0

			continue
		}

		// Alphabet reads
		if char >= 6 && char <= 31 {
			output.WriteByte(zstring.DefaultAlphabets[currentAlphabet][char-6])
		}

		// Space
		if char == 0 {
			output.WriteString(" ")
		}

		// Reset the alphabet as necessary
		if !lock {
			currentAlphabet = zstring.A0
		}

		// Alphabet changes
		if char >= 2 && char <= 5 {
			currentAlphabet, lock = zstring.Transition(currentAlphabet, char, m.version)
		}

		// Abbreviation
		if char >= 1 && char <= 3 && i < len(chars)-1 {
			nextChar := chars[i+1]
			abbreviationsTableOffset := uint32(32*(char-1) + nextChar)

			// TODO make constructors for byte, word, and packed addresses.
			abbreviationAddress := memory.Address(uint32(m.mem.ReadWord(memory.HAbbreviationsTable)) + abbreviationsTableOffset * 2)  // The offset is a word address, so multiply by 2.
			stringAddress := m.mem.ReadWord(abbreviationAddress)
			// Addresses in the abbreviations table are all word addresses, see s1.2.2
			output.WriteString(m.decodeZstring(memory.Address(stringAddress * 2)))

			i++  // jump past the abbreviation
			currentAlphabet = zstring.A0
		}
	}

	return output.String()
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
