# Split Block Bloom Filters

This repository contains an implementation of the Split Block Bloom Filters as described in the paper ["Split Block Bloom Filters"](https://arxiv.org/pdf/2101.01719.pdf) by Jim Apple.

## Overview

The Split Block Bloom Filter is an advanced implementation of the classic Bloom Filter, a probabilistic data structure used for set membership testing. It is designed to offer high efficiency and speed for modern computing architectures. This implementation is based on the research presented in the paper "Split Block Bloom Filters" by Jim Apple.

## Features

- **Enhanced Cache Efficiency:** Reduces cache line accesses to just one, significantly improving cache performance.
- **Optimized Hash Functions:** Utilizes a split approach and eight hash functions for efficient bit setting and checking.
- **Balanced Speed and Accuracy:** Offers a practical trade-off between a slightly higher false positive rate and significantly faster operations. This trade-off is favorable in scenarios where speed is more critical than having the absolute lowest false positive rate, especially within the practical range of ε (0.40% to 19%).


## Usage
```go
package main

import (
    "fmt"
    "github.com/axiomhq/splitblockbloom" // Ensure this is the correct path to your package
)

func main() {
    // Creating a new Split Block Bloom Filter
    // Parameters: Number of distinct values (ndv), False positive probability (fpp)
    // Example: Creating a filter for 1000 distinct values with a 1% false positive probability
    filter := splitblockbloom.NewFilter(1000, 0.01)

    // Adding elements to the filter
    // The AddHash method takes a hash of the item you want to add
    // Example: Adding elements with hash values 123456 and 789012
    filter.AddHash(123456)
    filter.AddHash(789012)

    // Checking if elements are in the filter
    // The Contains method returns a boolean indicating whether the element
    // (by its hash) is possibly in the filter
    // Remember: Bloom Filters can have false positives
    fmt.Println("123456 in filter:", filter.Contains(123456)) // true (element is added)
    fmt.Println("111111 in filter:", filter.Contains(111111)) // false (probably, as it wasn't added)
}
```