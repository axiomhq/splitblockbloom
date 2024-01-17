package splitblockbloom

import "io"

func blockFromStream(stream io.ReadSeeker, blockCount, idx int) (*block, error) {
	if idx >= blockCount {
		return nil, io.EOF
	}
	// Seek to the correct block
	if _, err := stream.Seek(int64(idx*blockSizeInBytes)+8, io.SeekStart); err != nil {
		return nil, err
	}
	// Read the block
	block := &block{}
	if _, err := block.ReadFrom(stream); err != nil {
		return nil, err
	}
	return block, nil
}

func ContainsFromStream(stream io.ReadSeeker, blockCount int, hash uint64) (bool, error) {
	blockIdx := hash % uint64(blockCount)
	block, err := blockFromStream(stream, blockCount, int(blockIdx))
	if err != nil {
		return false, err
	}
	return block.Contains(hash), nil
}
