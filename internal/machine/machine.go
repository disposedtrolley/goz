package machine

import (
	"fmt"
	"io"
	"strings"

	"github.com/disposedtrolley/goz/internal/memory"
	"github.com/disposedtrolley/goz/internal/zstring"
)

type Game []byte

type Machine struct {
	version int
	mem     *memory.Memory
	output  io.Writer
	input   io.Reader
}

func NewMachine(game Game) *Machine {
	return &Machine{
		mem: memory.NewMemory(game),
	}
}

func (m *Machine) WithVersion(version int) {
	m.version = version
}

func (m *Machine) SetOutput(w io.Writer) {
	m.output = w
}

func (m *Machine) SetInput(r io.Reader) {
	m.input = r
}

func (m *Machine) Start() error {
	if m.version == 0 {
		// Inspect the game file for the version.
		m.version = int(m.mem.ReadByte(memory.HVersion))
	}

	m.output.Write([]byte(fmt.Sprintf("z%d gamefile weighing in at %d bytes\n", m.version, m.mem.Size())))

	m.output.Write([]byte(fmt.Sprintf(`
beginning of:
  static memory: %x
  high memory: %x
`, m.mem.ReadWord(memory.HStaticMemoryBegin), m.mem.ReadWord(memory.HHighMemoryBegin))))

	return nil
}

// decodeZstring returns a string representing the decoded Z-characters found
// at the provided memory offset.
func (m *Machine) decodeZstring(offset memory.Address) string {
	var output strings.Builder

	var decodeZstringHelper func(offset memory.Address, isAbbreviation bool)
	decodeZstringHelper = func(offset memory.Address, isAbbreviation bool) {
		chars := m.zStringToChars(offset)
		lock := false
		currentAlphabet := zstring.A0
		for i := 0; i < len(chars); i++ {
			char := chars[i]

			// ZSCII
			if currentAlphabet == zstring.A2 && char == 6 {
				zsciiCode := (uint16(chars[i+1]) << 5) | uint16(chars[i+2])

				switch c := zstring.ZSCIIChar(zsciiCode); {
				// Tab (v6 only)
				case c == zstring.ZSCIITab:
					output.WriteString("\t")
				// Sentence space (v6 only)
				case c == zstring.ZSCIISentenceSpace:
					output.WriteString(" ")
				// Newline
				case c == zstring.ZSCIINewline:
					output.WriteString("\n")
				// Standard ASCII
				case c >= 32 && c <= 126:
					output.WriteByte(byte(c))
				// Unicode
				case c >= 155 && c <= 251:
					output.WriteString(fmt.Sprintf("%c", zstring.UnicodeChars[c-155]))
				}

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
				if isAbbreviation {
					// Recursive abbreviation.
					panic(fmt.Errorf("recursive abbreviation at offset %v. Decoded: %s", offset, output.String()))
				}

				nextChar := chars[i+1]
				abbreviationsTableOffset := m.mem.WordAddress(memory.Address(32*(char-1) + nextChar))
				abbreviationAddress := m.mem.ByteAddress(memory.Address(m.mem.ReadWord(memory.HAbbreviationsTable))) + abbreviationsTableOffset
				stringAddress := m.mem.ReadWord(abbreviationAddress)
				// Addresses in the abbreviations table are all word addresses, see s1.2.2
				decodeZstringHelper(m.mem.WordAddress(memory.Address(stringAddress)), true)

				i++ // jump past the abbreviation
				currentAlphabet = zstring.A0
			}
		}
	}

	decodeZstringHelper(offset, false)

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

// encodeInput encodes the user input string s into an array of ZChars, comprised of one or
// more ZStrings.
func (m *Machine) encodeInput(s string) (zchars []zstring.ZChar) {
	// §3.7
	s = strings.ToLower(s)

	// Max length of the input word differs based on machine version.
	inputWordLength := 6
	if m.version >= 4 {
		inputWordLength = 9
	}

	// Characters above the max length are truncated.
	// The Z-machine's dictionary doesn't actually hold the word
	// `examine`, but `examin`.
	for i := 0; i < inputWordLength; i++ {
		if i >= len(s) {
			// We've run out of characters in the input word, so pad the rest.
			zchars = append(zchars, zstring.PAD)
			continue
		}

		// Process the current char from the input word.
		currChar := s[i]

		// Alphabet or ZSCII?
		idx := strings.Index(zstring.DefaultAlphabets[zstring.A0], string(currChar))
		if idx >= 0 {
			// Alphabet A0.
			zchars = append(zchars, zstring.ZChar(6+idx))
		} else {
			// ZSCII???
		}
	}

	return zchars
}
