package splitblockbloom

import (
	"bytes"
	"math"
	"testing"
)

func TestBlockFromStreamAndContainsFromStream(t *testing.T) {
	// Create a filter and add a hash
	fpp := 0.01
	bvp := RecommendedBitsPerValue(1000, fpp)
	filter := NewFilter(1000, uint64(math.Ceil(bvp)))
	hash := uint64(123456)
	filter.AddHash(hash)

	// Serialize the filter to a buffer
	buf := new(bytes.Buffer)
	_, err := filter.WriteTo(buf)
	if err != nil {
		t.Fatalf("Failed to write filter to buffer: %v", err)
	}

	// Test blockFromStream
	blockCount := len(filter)
	block, err := blockFromStream(bytes.NewReader(buf.Bytes()), blockCount, 0)
	if err != nil {
		t.Errorf("blockFromStream failed: %v", err)
	}
	if block == nil {
		t.Error("blockFromStream returned a nil block")
	}

	// Test ContainsFromStream
	exists, err := ContainsFromStream(bytes.NewReader(buf.Bytes()), blockCount, hash)
	if err != nil {
		t.Errorf("ContainsFromStream failed: %v", err)
	}
	if !exists {
		t.Error("ContainsFromStream should have found the hash")
	}

	// Test with a hash that was not added
	randomHash := uint64(654321)
	exists, err = ContainsFromStream(bytes.NewReader(buf.Bytes()), blockCount, randomHash)
	if err != nil {
		t.Errorf("ContainsFromStream failed: %v", err)
	}
	if exists {
		t.Error("ContainsFromStream falsely found a hash that was not added")
	}
}
