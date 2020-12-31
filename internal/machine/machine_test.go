package machine

import (
	"bytes"
	"testing"

	"github.com/disposedtrolley/goz/internal/memory"
	"github.com/disposedtrolley/goz/internal/zstring"
	"github.com/disposedtrolley/goz/test"
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

func TestOutput(t *testing.T) {
	game, err := test.ReadGamfile(test.JigsawZ8)
	require.Nil(t, err, "should not error when reading gamefile")
	m := NewMachine(game)
	m.WithVersion(8)

	var buf bytes.Buffer
	m.SetOutput(&buf)
	m.Start()

	assert.Equal(t,
		"z8 gamefile weighing in at 304640 bytes\n\nbeginning of:\n  static memory: 6fc6\n  high memory: bc60\n",
		buf.String(),
		"should write to the supplied buffer")
}

func TestInput(t *testing.T) {
	tests := []struct {
		Name                 string
		Gamefile             test.Gamefile
		Version              int
		Input                string
		ExpectedEncodedInput []zstring.ZChar
	}{
		{
			Name:                 "when the example in §3.7 is encoded (v4+)",
			Gamefile:             test.JigsawZ8,
			Version:              8,
			Input:                "i",
			ExpectedEncodedInput: []zstring.ZChar{14, 5, 5, 5, 5, 5, 5, 5, 5},
		},
		{
			Name:                 "when the example in §3.7 is encoded (<v4)",
			Gamefile:             test.ZorkZ3,
			Version:              3,
			Input:                "i",
			ExpectedEncodedInput: []zstring.ZChar{14, 5, 5, 5, 5, 5},
		},
		{
			Name:                 "when a longer alpha string is encoded (v4+",
			Gamefile:             test.JigsawZ8,
			Version:              8,
			Input:                "examine",
			ExpectedEncodedInput: []zstring.ZChar{10, 29, 6, 18, 14, 19, 10, 5, 5},
		},
		{
			Name:                 "when a longer alpha string is encoded (<v4",
			Gamefile:             test.ZorkZ3,
			Version:              3,
			Input:                "examine",
			ExpectedEncodedInput: []zstring.ZChar{10, 29, 6, 18, 14, 19},
		},
		{
			Name:                 "when a string using other alphabets is encoded (<v4)",
			Gamefile:             test.ZorkZ3,
			Version:              3,
			Input:                "hello?",
			ExpectedEncodedInput: []zstring.ZChar{13, 10, 17, 17, 20, 5, 21},
		},
		{
			Name:                 "when a string using other alphabets is encoded (v4+)",
			Gamefile:             test.JigsawZ8,
			Version:              8,
			Input:                "hello?",
			ExpectedEncodedInput: []zstring.ZChar{13, 10, 17, 17, 20, 5, 21, 5, 5},
		},
		{
			Name:                 "when a ZSCII string is encoded (<v4)",
			Gamefile:             test.ZorkZ3,
			Version:              3,
			Input:                "áéíóú",
			ExpectedEncodedInput: []zstring.ZChar{5, 6, 6, 3, 5, 6, 5, 1},
		},
		{
			Name:                 "when a ZSCII string is encoded (v4+)",
			Gamefile:             test.JigsawZ8,
			Version:              8,
			Input:                "áéíóú",
			ExpectedEncodedInput: []zstring.ZChar{5, 6, 6, 3, 5, 6, 5, 1, 5, 6, 6, 3},
		},
		{
			Name:                 "when a mixed string is encoded (<v4)",
			Gamefile:             test.ZorkZ3,
			Version:              3,
			Input:                "á?l",
			ExpectedEncodedInput: []zstring.ZChar{5, 6, 6, 3, 5, 6, 5, 1},
		},
		{
			Name:                 "when a mixed string is encoded (v4+)",
			Gamefile:             test.JigsawZ8,
			Version:              8,
			Input:                "á?l",
			ExpectedEncodedInput: []zstring.ZChar{5, 6, 6, 3, 5, 6, 5, 1, 5, 21},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			game, err := test.ReadGamfile(tc.Gamefile)
			require.Nil(t, err, "should not error when reading gamefile")
			m := NewMachine(game)
			m.WithVersion(tc.Version)

			encodedInput := m.encodeInput(tc.Input)

			assert.Equal(t, tc.ExpectedEncodedInput, encodedInput)
		})
	}
}
