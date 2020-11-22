package zstring

import (
	"strings"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
)

type alphabet []string

var (
	//              0   1    2    3    4    5    6    7    8    9    10   11   12   13   14   15   16   17   18   19   20   21   22   23   24   25   26   27   28   29   30   31
	a0   = alphabet{"", "^", "^", "^", "^", "^", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	a1   = alphabet{"", "^", "^", "^", "^", "^", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	a2v1 = alphabet{"", "^", "^", "^", "^", "^", "^", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".", ",", "!", "?", "_", "#", "'", `"`, "/", "\\", "<", "=", ":", "(", ")"}
	a2   = alphabet{"", "^", "^", "^", "^", "^", "^", "^", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".", ",", "!", "?", "_", "#", "'", `"`, "/", "\\", "=", ":", "(", ")"}
)

// TODO Ztoa should take an array of ZSCII characters and output an ASCII string,
//      rather than decoding from a memory offset.
//      The flow should be:
//      	1) CPU encounters a Z-string which needs to be decoded at address x.
//      	2) The CPU walks the memory from address x, collecting ZSCII characters
//      	   into an array.
//      	3) The ZSCII characters stay in the array until they need to be converted
//      	   into ASCII for display purposes, at which time Ztoa should be called.
func Ztoa(mem *memory.Memory, addr memory.Address, version int) string {
	currentAlphabet := a0
	var output strings.Builder

	done := false
	for !done {
		word := mem.ReadWord(addr)
		//   --first byte-------   --second byte---
		//   7    6 5 4 3 2  1 0   7 6 5  4 3 2 1 0
		//   bit  --first--  --second---  --third--
		output.WriteString(currentAlphabet[word>>10&0x1F])
		output.WriteString(currentAlphabet[word>>5&0x1F])
		output.WriteString(currentAlphabet[word>>0&0x1F])

		leftoverBit := word & (1 << 15)
		done = leftoverBit != 0

		addr += 2
	}

	return output.String()
}
