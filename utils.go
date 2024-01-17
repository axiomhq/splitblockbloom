package splitblockbloom

import (
	"math"

	"github.com/segmentio/fasthash/fnv1a"
)

const (
	bitsPerWord      = 32
	wordsPerBlock    = 8
	hashBits         = 32
	blockSizeInBits  = bitsPerWord * wordsPerBlock
	blockSizeInBytes = blockSizeInBits / 8
	filterSeed       = 0x9747b28c
	blockSeed        = 0x5c6bfb31
)

func hash(val []byte, seed uint64) uint64 {
	return fnv1a.AddBytes64(seed, val)
}

func calcFalsePositiveRatio(ndv, bytes float64) float64 {
	if ndv == 0 || bytes <= 0 || ndv/(bytes*8) > 3 {
		return 1.0
	}

	lam := wordsPerBlock * bitsPerWord / ((bytes * 8) / ndv)
	loglam := math.Log(lam)
	log1collide := -hashBits * math.Log(2.0)
	maxJ := 10000

	var result float64
	for j := 0; j < maxJ; j++ {
		i := float64(maxJ - 1 - j)
		lGamma, _ := math.Lgamma(i + 1)
		logp := i*loglam - lam - lGamma
		logfinner := wordsPerBlock * math.Log(1.0-math.Pow(1.0-1.0/bitsPerWord, i))
		logcollide := math.Log(i) + log1collide
		result += math.Exp(logp+logfinner) + math.Exp(logp+logcollide)
	}

	return math.Min(result, 1.0)
}

func blockBytesNeeded(ndv, desiredFalsePositiveRatio float64) uint64 {
	result := 1.0
	for calcFalsePositiveRatio(ndv, result) > desiredFalsePositiveRatio {
		result *= 2
	}
	if result <= blockSizeInBytes {
		return blockSizeInBytes
	}

	lo, hi := 0.0, result
	for lo < hi-1 {
		mid := lo + (hi-lo)/2
		if calcFalsePositiveRatio(ndv, mid) < desiredFalsePositiveRatio {
			hi = mid
		} else {
			lo = mid
		}
	}
	return uint64(math.Ceil(hi/blockSizeInBytes) * blockSizeInBytes)
}
