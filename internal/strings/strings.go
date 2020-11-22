package strings

import (
	"strings"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/helpers"
)

type alphabet []string

const (
	zsciiCharLength = 5
)

var (
	//              0   1    2    3    4    5    6    7    8    9    10   11   12   13   14   15   16   17   18   19   20   21   22   23   24   25   26   27   28   29   30   31
	a0   = alphabet{"", "^", "^", "^", "^", "^", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	a1   = alphabet{"", "^", "^", "^", "^", "^", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	a2   = alphabet{"", "^", "^", "^", "^", "^", "^", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".", ",", "!", "?", "_", "#", "'", `"`, "/", "\\", "<", "=", ":", "(", ")"}
	a2v2 = alphabet{"", "^", "^", "^", "^", "^", "^", "^", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".", ",", "!", "?", "_", "#", "'", `"`, "/", "\\", "=", ":", "(", ")"}
)

type ZStringProcessor struct {
	currentAlphabet alphabet
}

func NewZStringProcessor(version int) *ZStringProcessor {
	return &ZStringProcessor{
		currentAlphabet: a0,
	}
}

func (zp *ZStringProcessor) Ztoa(z []uint16) string {
	var output strings.Builder

	for _, word := range z {

		for start := 1; start <= zsciiCharLength*3; start += zsciiCharLength {
			char := helpers.Bits(uint(word), uint(start), uint(zsciiCharLength))

			output.WriteString(zp.currentAlphabet[char])
		}

		leftoverBit := (word & (1 << 0))
		if leftoverBit != 0 {
			// 1st bit is set -- end of ZSCII sequence reached.
			break
		}
	}

	return output.String()
}
