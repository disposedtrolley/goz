package zstring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/zstring"
)

func TestZtoa(t *testing.T) {
	tests := []struct {
		Name          string
		ZVersion      int
		InputZChars   []zstring.ZChar
		Alphabets     zstring.Alphabets
		ExpectedAscii string
	}{
		{
			Name: "when 'hello' is decoded",
			// [0 01101 01010 10001] [1 10001 10100 00101]
			// [0011010101010001] [1100011010000000]
			ZVersion:      3,
			InputZChars:   []zstring.ZChar{0b01101, 0b01010, 0b10001, 0b10001, 0b10100, 0b00101},
			Alphabets:     zstring.DefaultAlphabets,
			ExpectedAscii: "hello",
		},
		{
			Name:          "when something else is decoded",
			ZVersion:      3,
			InputZChars:   []zstring.ZChar{0b100, 0b11111, 0b100, 0b10100, 0b100, 0b10111, 0b100, 0b10000, 0b0, 0b100, 0b1110, 0b101, 0b11101, 0b0, 0b1, 0b1, 0b100, 0b1100, 0b10111, 0b1010, 0b110, 0b11001, 0b0, 0b100, 0b11010, 0b10011, 0b1001, 0b1010, 0b10111, 0b11, 0b1011, 0b100, 0b1010, 0b10010, 0b10101, 0b1110, 0b10111, 0b1010, 0b101, 0b111, 0b100, 0b1000, 0b10100, 0b10101, 0b11110, 0b10111, 0b1110, 0b1100, 0b1101, 0b11001, 0b0, 0b101, 0b11110, 0b1000, 0b101, 0b11111, 0b0, 0b101, 0b1001, 0b101, 0b10001, 0b101, 0b10000, 0b101, 0b1001, 0b1, 0b11, 0b101, 0b1001, 0b101, 0b10001, 0b101, 0b10000, 0b101, 0b1010, 0b1, 0b11, 0b101, 0b1001, 0b101, 0b10001, 0b101, 0b10000, 0b101, 0b1011, 0b0, 0b100, 0b1110, 0b10011, 0b1011, 0b10100, 0b1000, 0b10100, 0b10010, 0b1, 0b11, 0b100, 0b1110, 0b10011, 0b1000, 0b1, 0b1100},
			ExpectedAscii: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actualAscii := zstring.Ztoa(tc.InputZChars, tc.Alphabets, tc.ZVersion)
			assert.Equal(t, tc.ExpectedAscii, actualAscii)
		})
	}

}
