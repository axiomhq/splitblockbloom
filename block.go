package splitblockebloom

import (
	"encoding/binary"
	"io"
)

var internalHashSeeds = []uint64{0x44974d91, 0x47b6137b, 0xa2b7289d, 0x8824ad5b, 0x2df1424b, 0x705495c7, 0x5c6bfb31, 0x9efc4947}

func makeMask(hash uint64) [wordsPerBlock]uint32 {
	// TODO: This is a constant we can avoid a loop
	var result [wordsPerBlock]uint32
	for i := range result {
		result[i] = 1 << ((uint32(hash) * uint32(internalHashSeeds[i])) >> (32 - 5))
	}
	return result
}

type Block [wordsPerBlock]uint32

func (b *Block) Add(val []byte) {
	h := hash(val, blockSeed)
	for i, m := range makeMask(h) {
		b[i] |= m
	}
}

func (b *Block) Contains(val []byte) bool {
	h := hash(val, blockSeed)
	for i, m := range makeMask(h) {
		if b[i]&m == 0 {
			return false
		}
	}
	return true
}

func (b *Block) WriteTo(w io.Writer) (int, error) {
	byts := make([]byte, blockSizeInBytes)
	for i, v := range b {
		binary.LittleEndian.PutUint32(byts[i*4:], v)
	}
	return w.Write(byts)
}

func (b *Block) ReadFrom(r io.Reader) (int, error) {
	byts := make([]byte, blockSizeInBytes)
	n, err := r.Read(byts)
	if err != nil {
		return n, err
	}
	for i := range b {
		b[i] = binary.LittleEndian.Uint32(byts[i*4:])
	}
	return n, nil
}
