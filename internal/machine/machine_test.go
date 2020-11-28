package machine

import (
	"testing"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
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
			Name:          "when a string with abbreviations is decoded - v3",
			Gamefile:      test.ZorkZ3,
			Version:       3,
			MemoryOffset:  0x6EE4,
			ExpectedASCII: "ZORK I: The Great Underground Empire\nCopyright (c) 1981, 1982, 1983 Infocom, Inc. ",
		},
		{
			Name:          "when a long Z-string is decoded - v3",
			Gamefile:      test.ZorkZ3,
			Version:       3,
			MemoryOffset:  0x10EEE,
			ExpectedASCII: "\"WELCOME TO ZORK!\n\nZORK is a game of adventure, danger, and low cunning. In it you will explore some of the most amazing territory ever seen by mortals. No computer should be without one!\"\n",
		},
		{
			Name:          "when a string with ZSCII characters is decoded - v3",
			Gamefile:      test.ZorkZ3,
			Version:       3,
			MemoryOffset:  0x5908,
			ExpectedASCII: ">",
		},
		{
			Name:          "when a string with abbreviations is decoded - v8",
			Gamefile:      test.JigsawZ8,
			Version:       8,
			MemoryOffset:  0x38631,
			ExpectedASCII: "               Welcome to JIGSAW\n",
		},
		{
			Name:          "when a long Z-string is decoded - v8",
			Gamefile:      test.JigsawZ8,
			Version:       8,
			MemoryOffset:  0x314BC,
			ExpectedASCII: "\nNew Year's Eve, 1999, a quarter to midnight and where else to be but Century Park! Fireworks cascade across the sky, your stomach rumbles uneasily, music and lasers howl across the parkland... Not exactly your ideal party (especially as that rather attractive stranger in black has slipped back into the crowds) - but cheer up, you won't live to see the next.\n\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			game, err := test.ReadGamfile(tc.Gamefile)
			require.Nil(t, err, "should not error when reading gamefile")
			m := NewMachine(game)
			m.WithVersion(tc.Version)
			str := m.decodeZstring(tc.MemoryOffset)
			assert.Equal(t, tc.ExpectedASCII, str)
		})
	}
}
