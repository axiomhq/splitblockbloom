// Borrowed from https://github.com/segmentio/fasthash

package fnv1a

const prime64 = uint64(1099511628211)

// AddBytes64 adds the hash of b to the precomputed hash value h.
func AddBytes64(h uint64, b []byte) uint64 {
	for len(b) >= 8 {
		h = (h ^ uint64(b[0])) * prime64
		h = (h ^ uint64(b[1])) * prime64
		h = (h ^ uint64(b[2])) * prime64
		h = (h ^ uint64(b[3])) * prime64
		h = (h ^ uint64(b[4])) * prime64
		h = (h ^ uint64(b[5])) * prime64
		h = (h ^ uint64(b[6])) * prime64
		h = (h ^ uint64(b[7])) * prime64
		b = b[8:]
	}

	if len(b) >= 4 {
		h = (h ^ uint64(b[0])) * prime64
		h = (h ^ uint64(b[1])) * prime64
		h = (h ^ uint64(b[2])) * prime64
		h = (h ^ uint64(b[3])) * prime64
		b = b[4:]
	}

	if len(b) >= 2 {
		h = (h ^ uint64(b[0])) * prime64
		h = (h ^ uint64(b[1])) * prime64
		b = b[2:]
	}

	if len(b) > 0 {
		h = (h ^ uint64(b[0])) * prime64
	}

	return h
}
