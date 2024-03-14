package splitblockbloom

import (
	"encoding/binary"
	"io"
	"math/bits"
)

var salt = [...]uint64{
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
	_ io.ReaderFrom = (*Block)(nil)
	_ io.WriterTo   = (*Block)(nil)
)

type Block [wordsPerBlock]uint32

func (blk *Block) AddHash(hash uint64) {
	blk[0] |= 1 << ((uint32(hash) * uint32(salt[0])) >> (27))
	blk[1] |= 1 << ((uint32(hash) * uint32(salt[1])) >> (27))
	blk[2] |= 1 << ((uint32(hash) * uint32(salt[2])) >> (27))
	blk[3] |= 1 << ((uint32(hash) * uint32(salt[3])) >> (27))
	blk[4] |= 1 << ((uint32(hash) * uint32(salt[4])) >> (27))
	blk[5] |= 1 << ((uint32(hash) * uint32(salt[5])) >> (27))
	blk[6] |= 1 << ((uint32(hash) * uint32(salt[6])) >> (27))
	blk[7] |= 1 << ((uint32(hash) * uint32(salt[7])) >> (27))
}

func (blk *Block) Contains(hash uint64) bool {
	return blk[0]&(1<<((uint32(hash)*uint32(salt[0]))>>(27))) != 0 &&
		blk[1]&(1<<((uint32(hash)*uint32(salt[1]))>>(27))) != 0 &&
		blk[2]&(1<<((uint32(hash)*uint32(salt[2]))>>(27))) != 0 &&
		blk[3]&(1<<((uint32(hash)*uint32(salt[3]))>>(27))) != 0 &&
		blk[4]&(1<<((uint32(hash)*uint32(salt[4]))>>(27))) != 0 &&
		blk[5]&(1<<((uint32(hash)*uint32(salt[5]))>>(27))) != 0 &&
		blk[6]&(1<<((uint32(hash)*uint32(salt[6]))>>(27))) != 0 &&
		blk[7]&(1<<((uint32(hash)*uint32(salt[7]))>>(27))) != 0
}

func (blk *Block) AddHashIfNotContains(hash uint64) bool {
	changed := false
	changed = changed || blk[0]&(1<<((uint32(hash)*uint32(salt[0]))>>(27))) == 0
	changed = changed || blk[1]&(1<<((uint32(hash)*uint32(salt[1]))>>(27))) == 0
	changed = changed || blk[2]&(1<<((uint32(hash)*uint32(salt[2]))>>(27))) == 0
	changed = changed || blk[3]&(1<<((uint32(hash)*uint32(salt[3]))>>(27))) == 0
	changed = changed || blk[4]&(1<<((uint32(hash)*uint32(salt[4]))>>(27))) == 0
	changed = changed || blk[5]&(1<<((uint32(hash)*uint32(salt[5]))>>(27))) == 0
	changed = changed || blk[6]&(1<<((uint32(hash)*uint32(salt[6]))>>(27))) == 0
	changed = changed || blk[7]&(1<<((uint32(hash)*uint32(salt[7]))>>(27))) == 0
	blk[0] |= 1 << ((uint32(hash) * uint32(salt[0])) >> (27))
	blk[1] |= 1 << ((uint32(hash) * uint32(salt[1])) >> (27))
	blk[2] |= 1 << ((uint32(hash) * uint32(salt[2])) >> (27))
	blk[3] |= 1 << ((uint32(hash) * uint32(salt[3])) >> (27))
	blk[4] |= 1 << ((uint32(hash) * uint32(salt[4])) >> (27))
	blk[5] |= 1 << ((uint32(hash) * uint32(salt[5])) >> (27))
	blk[6] |= 1 << ((uint32(hash) * uint32(salt[6])) >> (27))
	blk[7] |= 1 << ((uint32(hash) * uint32(salt[7])) >> (27))
	return changed
}

func (blk *Block) Merge(other *Block) {
	blk[0] |= other[0]
	blk[1] |= other[1]
	blk[2] |= other[2]
	blk[3] |= other[3]
	blk[4] |= other[4]
	blk[5] |= other[5]
	blk[6] |= other[6]
	blk[7] |= other[7]
}

func (blk *Block) WriteTo(w io.Writer) (int64, error) {
	b := make([]byte, blockSizeInBytes)
	binary.LittleEndian.PutUint32(b[0*4:], blk[0])
	binary.LittleEndian.PutUint32(b[1*4:], blk[1])
	binary.LittleEndian.PutUint32(b[2*4:], blk[2])
	binary.LittleEndian.PutUint32(b[3*4:], blk[3])
	binary.LittleEndian.PutUint32(b[4*4:], blk[4])
	binary.LittleEndian.PutUint32(b[5*4:], blk[5])
	binary.LittleEndian.PutUint32(b[6*4:], blk[6])
	binary.LittleEndian.PutUint32(b[7*4:], blk[7])
	n, err := w.Write(b)
	if n != blockSizeInBytes {
		return int64(n), io.ErrShortWrite
	}
	return int64(n), err
}

func (blk *Block) ReadFrom(r io.Reader) (int64, error) {
	b := make([]byte, blockSizeInBytes)
	n, err := io.ReadFull(r, b)
	if err != nil {
		return int64(n), err
	}
	if n != blockSizeInBytes {
		return int64(n), io.ErrUnexpectedEOF
	}
	blk[0] = binary.LittleEndian.Uint32(b[0*4:])
	blk[1] = binary.LittleEndian.Uint32(b[1*4:])
	blk[2] = binary.LittleEndian.Uint32(b[2*4:])
	blk[3] = binary.LittleEndian.Uint32(b[3*4:])
	blk[4] = binary.LittleEndian.Uint32(b[4*4:])
	blk[5] = binary.LittleEndian.Uint32(b[5*4:])
	blk[6] = binary.LittleEndian.Uint32(b[6*4:])
	blk[7] = binary.LittleEndian.Uint32(b[7*4:])
	return int64(n), nil
}

func (blk *Block) EstimateCardinality() int {
	var count int
	count += bits.OnesCount32(blk[0])
	count += bits.OnesCount32(blk[1])
	count += bits.OnesCount32(blk[2])
	count += bits.OnesCount32(blk[3])
	count += bits.OnesCount32(blk[4])
	count += bits.OnesCount32(blk[5])
	count += bits.OnesCount32(blk[6])
	count += bits.OnesCount32(blk[7])
	return count / 8
}
