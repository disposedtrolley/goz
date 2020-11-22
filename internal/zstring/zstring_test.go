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
		InputZscii    []uint16
		ExpectedAscii string
	}{
		{
			Name: "when 'hello' is decoded",
			// [0 01101 01010 10001] [1 10001 10100 00000]
			// [0011010101010001] [1100011010000000]
			ZVersion:      3,
			InputZscii:    []uint16{0b0011010101010001, 0b1100011010000000},
			ExpectedAscii: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			zp := zstring.NewZStringProcessor(tc.ZVersion)

			actualAscii := zp.Ztoa(tc.InputZscii)
			assert.Equal(t, tc.ExpectedAscii, actualAscii)
		})
	}

}
