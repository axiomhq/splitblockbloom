package splitblockbloom

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/axiomhq/splitblockbloom/internal/byteseeker"
)

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
		ok, err := ContainsFromStream(&byteseeker.Buffer{B: b}, len(bb), []byte(fmt.Sprintf("val%d", i)))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("val%d should be in the filter", i)
		}
	}

	errs = 0
	for i := int(1e6); i < 10e6; i++ {
		ok, err := ContainsFromStream(&byteseeker.Buffer{B: b}, len(bb), []byte(fmt.Sprintf("val%d", i)))
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
