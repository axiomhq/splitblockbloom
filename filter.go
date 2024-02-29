package splitblockbloom

import (
	"encoding/binary"
	"io"
)

var _ io.ReaderFrom = (*Filter)(nil)
var _ io.WriterTo = (*Filter)(nil)

type Filter []Block

// NewFilter creates a new blocked bloom filter.
func NewFilter(ndv uint64, fpp float64) Filter {
	numBytes := bytesNeeded(float64(ndv), fpp) + blockSizeInBytes
	return make([]Block, numBytes/blockSizeInBytes)
}

func (f Filter) SizeInBytes() int          { return len(f) * blockSizeInBytes }
func (f Filter) AddHash(hash uint64)       { f[hash%uint64(len(f))].AddHash(hash) }
func (f Filter) Contains(hash uint64) bool { return f[hash%uint64(len(f))].Contains(hash) }

func (f Filter) Merge(other Filter) {
	for i := range f {
		f[i].Merge(&other[i])
	}
}

func (f Filter) WriteTo(w io.Writer) (int64, error) {
	// write block count of filter
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(len(f)))
	wrote, err := w.Write(b)
	if err != nil {
		return int64(wrote), err
	}
	totalN := int64(wrote)
	// write each block
	for _, b := range f {
		n, err := b.WriteTo(w)
		totalN += n
		if err != nil {
			return totalN, err
		}
	}
	return totalN, nil
}

func (f *Filter) ReadFrom(r io.Reader) (int64, error) {
	// read length of filter
	b := make([]byte, 8)
	read, err := r.Read(b)
	if err != nil {
		return int64(read), err
	}
	n := binary.LittleEndian.Uint64(b)

	totalN := int64(read)
	// read each block
	*f = make([]Block, n)
	for i := range *f {
		n, err := (*f)[i].ReadFrom(r)
		totalN += n
		if err != nil {
			return totalN, err
		}
	}
	return totalN, nil
}
