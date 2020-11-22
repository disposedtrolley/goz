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
		Name                    string
		File                    test.Gamefile
		ExpectedVersion         uint8
		ExpectedFlags           uint8
		ExpectedHighMemoryStart uint16
	}{
		{
			Name:                    "headers for Zork",
			File:                    test.ZorkZ3,
			ExpectedVersion:         3,
			ExpectedFlags:           0b00000000,
			ExpectedHighMemoryStart: 0x4e37,
		},
		{
			Name:                    "headers for Jigsaw",
			File:                    test.JigsawZ8,
			ExpectedVersion:         8,
			ExpectedFlags:           0b00000010,
			ExpectedHighMemoryStart: 0xbc60,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Load the game into memory.
			b, err := test.ReadGamfile(tc.File)
			require.Nil(t, err, "should not error when loading gamefile")
			mem := memory.NewMemory(b)

			require.NotPanics(t, func() {
				// Version
				version := mem.ReadByte(memory.HVersion)
				assert.Equal(t, tc.ExpectedVersion, version)

				// Flags bitfield
				flags := mem.ReadByte(memory.HFlags)
				assert.Equal(t, tc.ExpectedFlags, flags)
			})
		})
	}
}
