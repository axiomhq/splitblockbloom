package splitblockbloom

import (
	"bytes"
	"math"
	"testing"

	"github.com/axiomhq/splitblockbloom/internal/byteseeker"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/stretchr/testify/require"
)

func TestNewFilter(t *testing.T) {
	var (
		ndv = uint64(1000)
		fpr = 0.01
		bvp = uint64(math.Ceil(RecommendedBitsPerValue(ndv, fpr)))
	)

	filter := NewFilter(ndv, bvp)
	require.Greater(t, len(filter), 0)
}

func TestAddAndContains(t *testing.T) {
	var (
		ndv = uint64(1000)
		fpr = 0.01
		bvp = uint64(math.Ceil(RecommendedBitsPerValue(ndv, fpr)))
	)

	filter := NewFilter(1000, bvp)

	hash := uint64(123456)
	filter.AddHash(hash)

	require.True(t, filter.Contains(hash))
	require.False(t, filter.Contains(1))
}

func TestWriteToAndReadFrom(t *testing.T) {
	var (
		ndv = uint64(1000)
		fpr = 0.01
		bvp = uint64(math.Ceil(RecommendedBitsPerValue(ndv, fpr)))
	)

	filter := NewFilter(1000, bvp)
	hash := uint64(123456)
	filter.AddHash(hash)

	buf := new(bytes.Buffer)
	_, err := filter.WriteTo(buf)
	require.NoError(t, err)

	newFilter := &Filter{}
	_, err = newFilter.ReadFrom(buf)
	require.NoError(t, err)

	require.True(t, newFilter.Contains(hash))
}

func TestFilterAll(t *testing.T) {
	fprs := []float64{0.004, 0.01, 0.1}
	bvps := []uint64{
		uint64(math.Ceil(RecommendedBitsPerValue(1e6, fprs[0]))),
		uint64(math.Ceil(RecommendedBitsPerValue(1e6, fprs[1]))),
		uint64(math.Ceil(RecommendedBitsPerValue(1e6, fprs[2]))),
	}
	count := int(1e6)
	for i, bvp := range bvps {
		bb := NewFilter(uint64(count), bvp)
		fpr := fprs[i]
		for i := 0; i < count; i++ {
			bb.AddHash(fnv1a.HashUint64(uint64(i)))
		}

		for i := 0; i < count; i++ {
			require.True(t, bb.Contains(fnv1a.HashUint64(uint64(i))))
		}

		errs := 0
		for i := count; i < count; i++ {
			if bb.Contains(fnv1a.HashUint64(uint64(i))) {
				errs++
			}
		}

		ratio := float64(errs) / float64(count)
		t.Log("errs ratio", ratio)
		t.Log("size:", bb.SizeInBytes())
		require.LessOrEqual(t, ratio, fpr)

		buf := bytes.NewBuffer(nil)
		n, err := bb.WriteTo(buf)
		require.NoError(t, err)

		b := buf.Bytes()

		t.Log("wrote:", n, "bytes", "len:", len(b))

		for i := 0; i < 1e6; i++ {
			ok, err := ContainsFromStream(&byteseeker.Buffer{B: b}, len(bb), fnv1a.HashUint64(uint64(i)))
			require.NoError(t, err)
			require.True(t, ok)
		}

		errs = 0
		for i := count; i < count; i++ {
			ok, err := ContainsFromStream(&byteseeker.Buffer{B: b}, len(bb), fnv1a.HashUint64(uint64(i)))
			require.NoError(t, err)
			if ok {
				errs++
			}
		}

		ratio = float64(errs) / float64(count)
		t.Log("errs ratio", ratio)
		t.Log("size:", bb.SizeInBytes())
		require.LessOrEqual(t, ratio, fpr)
	}
}
