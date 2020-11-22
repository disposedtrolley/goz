package zstring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
	"git.sr.ht/~disposedtrolley/go-zmachine/internal/zstring"
)

func TestZtoa(t *testing.T) {
	tests := []struct {
		Name          string
		ZVersion      int
		InputZscii    []byte
		ExpectedAscii string
	}{
		{
			Name: "when 'hello' is decoded",
			// [0 01101 01010 10001] [1 10001 10100 00000]
			// [0011010101010001] [1100011010000000]
			ZVersion:      3,
			InputZscii:    []byte{0b00110101, 0b01010001, 0b11000110, 0b10000000},
			ExpectedAscii: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			mem := memory.NewMemory(tc.InputZscii)

			actualAscii := zstring.Ztoa(mem, 0, tc.ZVersion)
			assert.Equal(t, tc.ExpectedAscii, actualAscii)
		})
	}

}
