package machine

import (
	"testing"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
	"git.sr.ht/~disposedtrolley/go-zmachine/internal/zstring"
	"git.sr.ht/~disposedtrolley/go-zmachine/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeZString(t *testing.T) {
	tests := []struct {
		Name          string
		Gamefile      test.Gamefile
		Version       int
		MemoryOffset  memory.Address
		ExpectedASCII string
	}{
		{
			Name:          "when a string with abbreviations is decoded",
			Gamefile:      test.ZorkZ3,
			Version:       3,
			MemoryOffset:  0x6EE4,
			ExpectedASCII: "foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			game, err := test.ReadGamfile(tc.Gamefile)
			require.Nil(t, err, "should not error when reading gamefile")
			m := NewMachine(game)
			chars := m.decodeZstring(tc.MemoryOffset)
			assert.Equal(t, tc.ExpectedASCII, zstring.Ztoa(chars, zstring.DefaultAlphabets, tc.Version))
		})
	}
}
