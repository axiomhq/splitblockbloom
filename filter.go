package splitblockebloom

import (
	"encoding/binary"
	"io"
)

type Filter []Block

// NewFilter creates a new blocked bloom filter.
func NewFilter(ndv uint64, fpp float64) Filter {
	return make([]Block, (blockBytesNeeded(float64(ndv), fpp)/blockSizeInBytes)*wordsPerBlock)
}

func (f Filter) SizeInBytes() int { return len(f) * blockSizeInBytes }
func (f Filter) Add(val []byte)   { f[hash(val, filterSeed)%uint64(len(f))].Add(val) }
func (f Filter) Contains(val []byte) bool {
	return f[hash(val, filterSeed)%uint64(len(f))].Contains(val)
}

func (f Filter) WriteTo(w io.Writer) (int64, error) {
	// write block count of filter
	byts := make([]byte, 8)
	binary.LittleEndian.PutUint64(byts, uint64(len(f)))
	totalN, err := w.Write(byts)
	if err != nil {
		return int64(totalN), err
	}
	// write each block
	for _, b := range f {
		n, err := b.WriteTo(w)
		totalN += n
		if err != nil {
			return int64(totalN), err
		}
	}
	return int64(totalN), nil
}

func (f Filter) ReaderFrom(r io.Reader) (int64, error) {
	// read length of filter
	byts := make([]byte, 8)
	totalN, err := r.Read(byts)
	if err != nil {
		return int64(totalN), err
	}
	n := binary.LittleEndian.Uint64(byts)

	// read each block
	f = make([]Block, n)
	for i := range f {
		n, err := f[i].ReadFrom(r)
		totalN += n
		if err != nil {
			return int64(totalN), err
		}
	}
	return int64(totalN), nil
}
