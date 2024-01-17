package splitblockbloom

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/axiomhq/splitblockbloom/internal/byteseeker"
	"github.com/segmentio/fasthash/fnv1a"
)

func TestNewFilter(t *testing.T) {
	ndv := uint64(1000)
	fpp := 0.01

	filter := NewFilter(ndv, fpp)
	if len(filter) == 0 {
		t.Error("NewFilter created an empty filter")
	}
}

func TestAddAndContains(t *testing.T) {
	filter := NewFilter(1000, 0.01)

	hash := uint64(123456)
	filter.AddHash(hash)

	if !filter.Contains(hash) {
		t.Errorf("Filter should contain the hash %d", hash)
	}

	randomHash := uint64(rand.Int63()) // Generate a random hash
	if filter.Contains(randomHash) {
		t.Errorf("Filter falsely contains random hash %d", randomHash)
	}
}

func TestWriteToAndReadFrom(t *testing.T) {
	filter := NewFilter(1000, 0.01)
	hash := uint64(123456)
	filter.AddHash(hash)

	buf := new(bytes.Buffer)
	_, err := filter.WriteTo(buf)
	if err != nil {
		t.Errorf("WriteTo failed: %v", err)
	}

	newFilter := &Filter{}
	_, err = newFilter.ReadFrom(buf)
	if err != nil {
		t.Errorf("ReadFrom failed: %v", err)
	}

	if !newFilter.Contains(hash) {
		t.Error("ReadFrom filter does not contain hash that was written")
	}
}

func TestFilterAll(t *testing.T) {
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
