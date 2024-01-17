package splitblockbloom

import (
	"encoding/binary"
	"io"
)

var internalHashSeeds = [...]uint64{
	0x44974d91,
	0x47b6137b,
	0xa2b7289d,
	0x8824ad5b,
	0x2df1424b,
	0x705495c7,
	0x5c6bfb31,
	0x9efc4947,
}

var (
	_ io.ReaderFrom = (*block)(nil)
	_ io.WriterTo   = (*block)(nil)
)

type block [wordsPerBlock]uint32

func (blk *block) AddHash(hash uint64) {
	for i, m := range internalHashSeeds {
		blk[i] |= 1 << ((uint32(hash) * uint32(m)) >> (32 - 5))
	}
}

func (blk *block) Contains(hash uint64) bool {
	for i, m := range internalHashSeeds {
		if blk[i]&(1<<((uint32(hash)*uint32(m))>>(32-5))) == 0 {
			return false
		}
	}
	return true
}

func (blk *block) WriteTo(w io.Writer) (int64, error) {
	b := make([]byte, blockSizeInBytes)
	for i, v := range blk {
		binary.LittleEndian.PutUint32(b[i*4:], v)
	}
	n, err := w.Write(b)
	return int64(n), err
}

func (blk *block) ReadFrom(r io.Reader) (int64, error) {
	b := make([]byte, blockSizeInBytes)
	n, err := r.Read(b)
	if err != nil {
		return int64(n), err
	}
	for i := range blk {
		blk[i] = binary.LittleEndian.Uint32(b[i*4:])
	}
	return int64(n), nil
}
