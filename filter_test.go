package splitblockbloom

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

// ByteSliceReadSeeker implements the io.ReadSeeker interface for a byte slice.
type ByteSliceReadSeeker struct {
	slice  []byte
	offset int64
}

// NewByteSliceReadSeeker creates a new ByteSliceReadSeeker.
func NewByteSliceReadSeeker(slice []byte) *ByteSliceReadSeeker {
	return &ByteSliceReadSeeker{slice: slice, offset: 0}
}

// Read implements the Read method of the io.ReadSeeker interface.
func (r *ByteSliceReadSeeker) Read(p []byte) (int, error) {
	if r.offset >= int64(len(r.slice)) {
		return 0, io.EOF // end of slice
	}

	n := copy(p, r.slice[r.offset:])
	r.offset += int64(n)
	return n, nil
}

// Seek implements the Seek method of the io.ReadSeeker interface.
func (r *ByteSliceReadSeeker) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = r.offset + offset
	case io.SeekEnd:
		newOffset = int64(len(r.slice)) + offset
	default:
		return 0, errors.New("invalid whence")
	}

	if newOffset < 0 {
		return 0, errors.New("negative position")
	}

	if newOffset > int64(len(r.slice)) {
		return 0, errors.New("position out of bounds")
	}

	r.offset = newOffset
	return newOffset, nil
}

func TestFilter(t *testing.T) {
	bb := NewFilter(1e6, 0.004)
	for i := 0; i < 1e6; i++ {
		bb.Add([]byte(fmt.Sprintf("val%d", i)))
	}

	for i := 0; i < 1e6; i++ {
		if !bb.Contains([]byte(fmt.Sprintf("val%d", i))) {
			t.Errorf("val%d should be in the filter", i)
		}
	}

	errs := 0
	for i := int(1e6); i < 10e6; i++ {
		if bb.Contains([]byte(fmt.Sprintf("val%d", i))) {
			errs++
			//t.Errorf("val%d should not be in the filter", i)
		}
	}

	t.Log("errs:", float64(errs)/1e6)
	t.Log("size:", bb.SizeInBytes())

	buf := bytes.NewBuffer(nil)
	n, err := bb.WriteTo(buf)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.Bytes()

	t.Log("wrote:", n, "bytes", "len:", len(b))

	for i := 0; i < 1e6; i++ {
		ok, err := ContainsFromStream(NewByteSliceReadSeeker(b), len(bb), []byte(fmt.Sprintf("val%d", i)))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("val%d should be in the filter", i)
		}
	}

	errs = 0
	for i := int(1e6); i < 10e6; i++ {
		ok, err := ContainsFromStream(NewByteSliceReadSeeker(b), len(bb), []byte(fmt.Sprintf("val%d", i)))
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			errs++
			//t.Errorf("val%d should not be in the filter", i)
		}
	}

	t.Log("errs:", float64(errs)/1e6)
}
