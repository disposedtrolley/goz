package zstring

type ZChar uint8 // it's really 5 bits, but we can only go as low as 8 natively.

type Alphabets []string

type Alphabet int

const (
	A0   Alphabet = 0
	A1   Alphabet = 1
	A2   Alphabet = 2
	A2v1 Alphabet = 3
)

type ZSCIIChar uint16

const (
	ZSCIITab           ZSCIIChar = 9
	ZSCIISentenceSpace ZSCIIChar = 11
	ZSCIINewline       ZSCIIChar = 13
)

// Alphabets begin at index 6.
// The final DefaultAlphabets is the A2 variation used by V1 of the Z-machine.
var DefaultAlphabets = Alphabets{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	" \n0123456789.,!?_#'\"/\\-:()",
	" 0123456789.,!?_#'\"/\\<-:()",
}

var UnicodeChars = []ZSCIIChar{
	0x00e4, 0x00f6, 0x00fc, 0x00c4, 0x00d6, 0x00dc, 0x00df, 0x00bb, 0x00ab,
	0x00eb, 0x00ef, 0x00ff, 0x00cb, 0x00cf, 0x00e1, 0x00e9, 0x00ed, 0x00f3,
	0x00fa, 0x00fd, 0x00c1, 0x00c9, 0x00cd, 0x00d3, 0x00da, 0x00dd, 0x00e0,
	0x00e8, 0x00ec, 0x00f2, 0x00f9, 0x00c0, 0x00c8, 0x00cc, 0x00d2, 0x00d9,
	0x00e2, 0x00ea, 0x00ee, 0x00f4, 0x00fb, 0x00c2, 0x00ca, 0x00ce, 0x00d4,
	0x00db, 0x00e5, 0x00c5, 0x00f8, 0x00d8, 0x00e3, 0x00f1, 0x00f5, 0x00c3,
	0x00d1, 0x00d5, 0x00e6, 0x00c6, 0x00e7, 0x00c7, 0x00fe, 0x00f0, 0x00de,
	0x00d0, 0x00a3, 0x0153, 0x0152, 0x00a1, 0x00bf,
}

func transitionsTable() map[Alphabet]map[ZChar]Alphabet {
	// 		 from A0  from A1  from A2
	// Z-char 2      A1       A2       A0  // next char only
	// Z-char 3      A2       A0       A1  // next char only
	// Z-char 4      A1       A2       A0  // permanent (<v3) next char only (v3+)
	// Z-char 5      A2       A0       A1  // permanent (<v3) next char only (v3+)
	var transitions = map[Alphabet]map[ZChar]Alphabet{}
	transitions[A0] = make(map[ZChar]Alphabet)
	transitions[A1] = make(map[ZChar]Alphabet)
	transitions[A2] = make(map[ZChar]Alphabet)
	transitions[A0][2] = A1
	transitions[A0][3] = A2
	transitions[A0][4] = A1
	transitions[A0][5] = A2
	transitions[A1][2] = A2
	transitions[A1][3] = A0
	transitions[A1][4] = A2
	transitions[A1][5] = A0
	transitions[A2][2] = A0
	transitions[A2][3] = A1
	transitions[A2][4] = A0
	transitions[A2][5] = A1

	return transitions
}

var transitions = transitionsTable()

func Transition(currAlphabet Alphabet, char ZChar, version int) (newAlphabet Alphabet, lock bool) {
	newAlphabet = transitions[currAlphabet][char]

	return newAlphabet, version < 3 && char > 3
}
