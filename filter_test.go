package splitblockbloom

import (
	"bytes"
	"testing"

	"github.com/axiomhq/splitblockbloom/internal/byteseeker"
	"github.com/segmentio/fasthash/fnv1a"
)

func TestFilter(t *testing.T) {
	fpp := 0.004
	bb := NewFilter(1e6, fpp)
	for i := 0; i < 1e6; i++ {
		bb.AddHash(fnv1a.HashUint64(uint64(i)))
	}

	for i := 0; i < 1e6; i++ {
		if !bb.Contains(fnv1a.HashUint64(uint64(i))) {
			t.Errorf("val%d should be in the filter", i)
		}
	}

	errs := 0
	for i := int(1e6); i < 10e6; i++ {
		if bb.Contains(fnv1a.HashUint64(uint64(i))) {
			errs++
		}
	}

	ratio := float64(errs) / 1e6
	t.Log("errs ratio", ratio)
	t.Log("size:", bb.SizeInBytes())
	if ratio > fpp {
		t.Errorf("error ratio should be less than %f, got %f", fpp, ratio)
	}

	buf := bytes.NewBuffer(nil)
	n, err := bb.WriteTo(buf)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.Bytes()

	t.Log("wrote:", n, "bytes", "len:", len(b))

	for i := 0; i < 1e6; i++ {
		ok, err := ContainsFromStream(&byteseeker.Buffer{B: b}, len(bb), fnv1a.HashUint64(uint64(i)))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("val%d should be in the filter", i)
		}
	}

	errs = 0
	for i := int(1e6); i < 10e6; i++ {
		ok, err := ContainsFromStream(&byteseeker.Buffer{B: b}, len(bb), fnv1a.HashUint64(uint64(i)))
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			errs++
		}
	}

	ratio = float64(errs) / 1e6
	t.Log("errs ratio", ratio)
	t.Log("size:", bb.SizeInBytes())
	if ratio > fpp {
		t.Errorf("error ratio should be less than %f, got %f", fpp, ratio)
	}
}
