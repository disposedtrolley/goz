package zstring

import (
	"fmt"
	"strings"
)

type ZChar uint8 // it's really 5 bits, but we can only go as low as 8 natively.

type Alphabets []string

type Alphabet int

const (
	A0   Alphabet = 0
	A1   Alphabet = 1
	A2   Alphabet = 2
	A2v1 Alphabet = 3
)

// Alphabets begin at index 6.
// The final DefaultAlphabets is the A2 variation used by V1 of the Z-machine.
var DefaultAlphabets = Alphabets{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	" \n0123456789.,!?_#'\"/\\-:()",
	" 0123456789.,!?_#'\"/\\<-:()",
}

// Ztoa takes an array of Z-characters and returns its string representation.
func Ztoa(chars []ZChar, alphabets Alphabets, version int) string {
	lockAll := false
	lockNext := false
	currentAlphabet := A0

	var output strings.Builder

	for i, char := range chars {
		fmt.Println(i)
		if !lockAll && !lockNext {
			currentAlphabet = A0
		}

		// Space
		if char == 0 {
			output.WriteString(" ")
		}

		// Alphabet changes
		if char >= 2 && char <= 5 {
			currentAlphabet, lockNext, lockAll = transition(currentAlphabet, char, version)
			continue
		}

		// Alphabet reads
		if char >= 6 && char <= 31 {
			output.WriteByte(alphabets[currentAlphabet][char-6])
		}

		if lockNext {
			lockNext = false
		}
	}

	return output.String()
}

func transition(currAlphabet Alphabet, char ZChar, version int) (newAlphabet Alphabet, lockNext, lockAll bool) {
	// 				 from A0  from A1  from A2
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

	newAlphabet = transitions[currAlphabet][char]

	if version < 3 && char > 3 {
		lockAll = true
	}

	lockNext = true

	return newAlphabet, lockNext, lockAll
}
