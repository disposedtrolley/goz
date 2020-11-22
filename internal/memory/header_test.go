package memory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
	"git.sr.ht/~disposedtrolley/go-zmachine/test"
)

func TestHeaderCommon(t *testing.T) {
	tests := []struct {
		Name            string
		File            test.Gamefile
		ExpectedVersion uint8
	}{
		{
			Name:            "common headers for a z3 game",
			File:            test.ZorkZ3,
			ExpectedVersion: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Load the game into memory.
			b, err := test.ReadGamfile(tc.File)
			require.Nil(t, err, "should not error when loading gamefile")
			mem := memory.NewMemory(b)

			require.NotPanics(t, func() {
				version := mem.ReadByte(memory.HVersion)
				assert.Equal(t, tc.ExpectedVersion, version)
			})
		})
	}
}
