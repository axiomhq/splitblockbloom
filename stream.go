package splitblockbloom

import "io"

func BlockFromStream(stream io.ReadSeeker, blockCount, idx int) (*Block, error) {
	if idx >= blockCount {
		return nil, io.EOF
	}
	// Seek to the correct block
	if _, err := stream.Seek(int64(idx*blockSizeInBytes)+8, io.SeekStart); err != nil {
		return nil, err
	}
	// Read the block
	block := &Block{}
	if _, err := block.ReadFrom(stream); err != nil {
		return nil, err
	}
	return block, nil
}

func ContainsFromStream(stream io.ReadSeeker, blockCount int, val []byte) (bool, error) {
	blockIdx := hash(val, filterSeed) % uint64(blockCount)
	block, err := BlockFromStream(stream, blockCount, int(blockIdx))
	if err != nil {
		return false, err
	}
	return block.Contains(val), nil
}
