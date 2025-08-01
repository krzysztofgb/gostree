# gostree

A Go implementation of an order-statistic tree.

## What is an Order-Statistic Tree?

An [order-statistic tree](https://en.wikipedia.org/wiki/Order_statistic_tree) is a variant of a binary search tree (or B-tree) that maintains subtree sizes at each node.
This augmentation enables efficient O(log n) operations for:
- **Select(k)**: Find the k-th smallest element
- **Rank(x)**: Determine the rank (position) of element x

## Compatibility

This library aims to be a drop-in replacement for [google/btree](https://github.com/google/btree) with additional order-statistic operations.

## CI/CD Workflows

- **Testing**: Runs on every push and PR
- **Releases**: Tags matching `v*.*.*` trigger automatic GitHub releases with generated notes

## Installation

```bash
go get github.com/krzysztofgb/gostree
```

## Usage

```go
import "github.com/krzysztofgb/gostree"

// Usage examples coming soon
```

## License

SEE [LICENSE](LICENSE) file for details.