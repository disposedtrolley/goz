package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryRead(t *testing.T) {
	mem := newMemory([]byte{0xfe, 0xa2, 0x0d, 0x19, 0x00})

	assert.Equal(t, uint8(0xfe), mem.readByte(0), "should read a byte")
	assert.Equal(t, uint8(0x00), mem.readByte(4), "should read a byte at the end of memory")
	assert.Equal(t, uint16(0xfea2), mem.readWord(0), "should read a word")
	assert.Equal(t, uint16(0x1900), mem.readWord(3), "should read a word at the end of memory")
}

func TestMemoryRead_OutOfBounds(t *testing.T) {
	mem := newMemory([]byte{0xfe, 0xa2, 0x0d, 0x19, 0x00})

	assert.Panics(t, func() {
		mem.readWord(4)
	}, "should panic when reading a word from the last address")

	assert.Panics(t, func() {
		mem.readByte(5)
	}, "should panic when reading a byte past the end of memory")
}
