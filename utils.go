package splitblockbloom

import (
	"math"
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

func calcFalsePositiveRatio(ndv, bytes float64) float64 {
	if ndv == 0 || bytes <= 0 || ndv/(bytes*8) > 3 {
		return 1.0
	}

	lam := wordsPerBlock * bitsPerWord / ((bytes * 8) / ndv)
	loglam := math.Log(lam)
	log1collide := -hashBits * math.Log(2.0)
	maxJ := 100

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

func RecommendedBitsPerValue(ndv uint64, desiredFalsePositiveRatio float64) float64 {
	// start with the optimal number of bits per value according ti standard bloom filter formula
	n := float64(ndv)
	m := -n * math.Log(desiredFalsePositiveRatio) / (math.Log(2) * math.Log(2))
	result := math.Ceil(m / blockSizeInBytes)

	lo, hi := result, result*2
	for lo < hi-1 {
		mid := lo + (hi-lo)/2
		if calcFalsePositiveRatio(n, mid) < desiredFalsePositiveRatio {
			hi = mid
		} else {
			lo = mid
		}
	}

	return float64(m / n)
}
