package byteseeker

import (
	"errors"
	"io"
)

var _ io.ReadSeeker = (*Buffer)(nil)

// Buffer implements the io.ReadSeeker interface for a byte slice.
type Buffer struct {
	B      []byte
	offset int64
}

// Read implements the Read method of the io.ReadSeeker interface.
func (r *Buffer) Read(p []byte) (int, error) {
	if r.offset >= int64(len(r.B)) {
		return 0, io.EOF // end of slice
	}

	n := copy(p, r.B[r.offset:])
	r.offset += int64(n)
	return n, nil
}

// Seek implements the Seek method of the io.ReadSeeker interface.
func (r *Buffer) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = r.offset + offset
	case io.SeekEnd:
		newOffset = int64(len(r.B)) + offset
	default:
		return 0, errors.New("invalid whence")
	}

	if newOffset < 0 {
		return 0, errors.New("negative position")
	}

	if newOffset > int64(len(r.B)) {
		return 0, errors.New("position out of bounds")
	}

	r.offset = newOffset
	return newOffset, nil
}
