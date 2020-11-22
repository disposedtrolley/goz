package memory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.sr.ht/~disposedtrolley/go-zmachine/internal/memory"
)

func TestMemoryRead(t *testing.T) {
	mem := memory.NewMemory([]byte{0xfe, 0xa2, 0x0d, 0x19, 0x00})

	assert.Equal(t, uint8(0xfe), mem.ReadByte(0), "should read a byte")
	assert.Equal(t, uint8(0x00), mem.ReadByte(4), "should read a byte at the end of memory")
	assert.Equal(t, uint16(0xfea2), mem.ReadWord(0), "should read a word")
	assert.Equal(t, uint16(0x1900), mem.ReadWord(3), "should read a word at the end of memory")
}

func TestMemoryRead_OutOfBounds(t *testing.T) {
	mem := memory.NewMemory([]byte{0xfe, 0xa2, 0x0d, 0x19, 0x00})

	assert.Panics(t, func() {
		mem.ReadWord(4)
	}, "should panic when reading a word from the last address")

	assert.Panics(t, func() {
		mem.ReadByte(5)
	}, "should panic when reading a byte past the end of memory")
}

func TestMemoryWrite(t *testing.T) {
	mem := memory.NewMemory([]byte{0xfe, 0xa2, 0x0d, 0x19, 0x00})

	mem.WriteByte(0, 0xff)
	assert.Equal(t, byte(0xff), mem.ReadByte(0))

	mem.WriteWord(3, 0xaabb)
	assert.Equal(t, uint16(0xaabb), mem.ReadWord(3))
}

func TestMemoryWrite_OutOfBounds(t *testing.T) {
	mem := memory.NewMemory([]byte{0xfe, 0xa2, 0x0d, 0x19, 0x00})

	assert.Panics(t, func() {
		mem.WriteWord(4, 0xffdd)
	}, "should panic when writing a word to the last address")

	assert.Panics(t, func() {
		mem.WriteByte(5, 0x11)
	}, "should panic when writing a byte past the end of memory")
}
