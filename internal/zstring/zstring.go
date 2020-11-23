package zstring

import (
	"strings"
)

type ZChar uint8 // it's really 5 bits, but we can only go as low as 8 natively.

// Alphabets begin at index 6.
// The final alphabet is the A2 variation used by V1 of the Z-machine.
var alphabet = []string{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	" \n0123456789.,!?_#'\"/\\-:()",
	" 0123456789.,!?_#'\"/\\<-:()",
}

// Ztoa takes an array of Z-characters and returns its string representation.
func Ztoa(chars []ZChar, version int) string {
	currentAlphabet := alphabet[0]
	var output strings.Builder

	for _, char := range chars {
		if char >= 6 {
			output.WriteByte(currentAlphabet[char-6])
		}
	}

	return output.String()
}
