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
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actualAscii := zstring.Ztoa(tc.InputZChars, tc.Alphabets, tc.ZVersion)
			assert.Equal(t, tc.ExpectedAscii, actualAscii)
		})
	}

}
