# gostree

A generic [order-statistic tree](https://en.wikipedia.org/wiki/Order_statistic_tree) tree implementation in Go.

## Installation

```bash
go get github.com/krzysztofgb/gostree
```

## Usage

```go
import "github.com/krzysztofgb/gostree"

// Create a new tree with initial capacity
tree := gostree.New[int](16)

// Insert elements
tree.Insert(5)
tree.Insert(3)
tree.Insert(7)
tree.Insert(1)
tree.Insert(9)

// Get the k-th smallest element (0-indexed)
val, found := tree.Select(2)  // Returns 5, true (3rd smallest)

// Get the rank (number of elements less than x)
rank := tree.Rank(6)  // Returns 3 (elements 1, 3, 5 are less than 6)

// Delete an element
deleted := tree.Delete(3)  // Returns true

// Get the size
size := tree.Size()  // Returns 4
```

## API Reference

### Core Methods

```go
// New creates a new order-statistic tree with the specified initial capacity
func New[T cmp.Ordered](capacity int) *Tree[T]

// Insert adds an item to the tree
func (t *Tree[T]) Insert(item T)

// Delete removes an item from the tree, returns true if found and deleted
func (t *Tree[T]) Delete(item T) bool

// Select returns the k-th smallest element (0-indexed)
func (t *Tree[T]) Select(k int) (T, bool)

// Rank returns the number of elements less than the given item
func (t *Tree[T]) Rank(item T) int

// Size returns the number of elements in the tree
func (t *Tree[T]) Size() int
```

## Concurrency Safety

**Write operations are NOT concurrent safe.**
The following methods modify the tree structure and require external synchronization when used concurrently:
- `Insert()`
- `Delete()`

**Read operations ARE concurrent safe.**
Multiple goroutines can safely call these methods simultaneously without external synchronization:
- `Select()`
- `Rank()`
- `Size()`

If you need to use this tree in a concurrent environment with both readers and writers, you must implement your own synchronization (e.g., using `sync.RWMutex`).

## Performance

### Benchmark Results

Benchmarks were run against:
* [ajwerner/orderstat](https://github.com/ajwerner/orderstat)
* [google/btree](https://github.com/google/btree)

Note, benchmarks for the "Rank" and "Select" methods were omitted for `google/btree` (not supported).

<details>

<summary>Insert</summary>

```text
                                                      │   sec/op    │
Insert/krzysztofgb/gostree/100_elements-12              3.230µ ± 1%
Insert/ajwerner/orderstat/100_elements-12               10.45µ ± 2%
Insert/google/btree/100_elements-12                     13.56µ ± 1%
Insert/krzysztofgb/gostree/1000_elements-12             52.21µ ± 1%
Insert/ajwerner/orderstat/1000_elements-12              186.5µ ± 1%
Insert/google/btree/1000_elements-12                    213.6µ ± 0%
Insert/krzysztofgb/gostree/10000_elements-12            889.2µ ± 1%
Insert/ajwerner/orderstat/10000_elements-12             2.660m ± 0%
Insert/google/btree/10000_elements-12                   2.917m ± 1%

                                                      │      B/op      │
Insert/krzysztofgb/gostree/100_elements-12              4.734Ki ± 0%
Insert/ajwerner/orderstat/100_elements-12               8.023Ki ± 0%
Insert/google/btree/100_elements-12                     8.914Ki ± 0%
Insert/krzysztofgb/gostree/1000_elements-12             46.92Ki ± 0%
Insert/ajwerner/orderstat/1000_elements-12              71.13Ki ± 0%
Insert/google/btree/1000_elements-12                    94.31Ki ± 0%
Insert/krzysztofgb/gostree/10000_elements-12            468.8Ki ± 0%
Insert/ajwerner/orderstat/10000_elements-12             1.076Mi ± 0%
Insert/google/btree/10000_elements-12                   920.8Ki ± 0%

                                                      │   allocs/op   │
Insert/krzysztofgb/gostree/100_elements-12               101.0 ± 0%
Insert/ajwerner/orderstat/100_elements-12                71.00 ± 0%
Insert/google/btree/100_elements-12                      282.0 ± 0%
Insert/krzysztofgb/gostree/1000_elements-12             1.001k ± 0%
Insert/ajwerner/orderstat/1000_elements-12               984.0 ± 0%
Insert/google/btree/1000_elements-12                    3.251k ± 0%
Insert/krzysztofgb/gostree/10000_elements-12            10.00k ± 0%
Insert/ajwerner/orderstat/10000_elements-12             9.987k ± 0%
Insert/google/btree/10000_elements-12                   32.10k ± 0%
```

</details>

<details>

<summary>Search</summary>

```text
                                                      │   sec/op    │
Search/krzysztofgb/gostree/100_elements-12              2.634µ ± 0%
Search/ajwerner/orderstat/100_elements-12               6.465µ ± 1%
Search/google/btree/100_elements-12                     935.1n ± 1%
Search/krzysztofgb/gostree/1000_elements-12             3.734µ ± 0%
Search/ajwerner/orderstat/1000_elements-12              9.991µ ± 0%
Search/google/btree/1000_elements-12                    8.821µ ± 1%
Search/krzysztofgb/gostree/10000_elements-12            6.475µ ± 0%
Search/ajwerner/orderstat/10000_elements-12             16.02µ ± 0%
Search/google/btree/10000_elements-12                   88.76µ ± 0%

                                                      │      B/op      │
Search/krzysztofgb/gostree/100_elements-12                0.000 ± 0%
Search/ajwerner/orderstat/100_elements-12                 647.0 ± 0%
Search/google/btree/100_elements-12                       984.0 ± 0%
Search/krzysztofgb/gostree/1000_elements-12               0.000 ± 0%
Search/ajwerner/orderstat/1000_elements-12                777.0 ± 0%
Search/google/btree/1000_elements-12                    7.922Ki ± 0%
Search/krzysztofgb/gostree/10000_elements-12              0.000 ± 0%
Search/ajwerner/orderstat/10000_elements-12               798.0 ± 0%
Search/google/btree/10000_elements-12                   78.28Ki ± 0%

                                                      │   allocs/op   │
Search/krzysztofgb/gostree/100_elements-12               0.000 ± 0%
Search/ajwerner/orderstat/100_elements-12                80.00 ± 1%
Search/google/btree/100_elements-12                      85.00 ± 0%
Search/krzysztofgb/gostree/1000_elements-12              0.000 ± 0%
Search/ajwerner/orderstat/1000_elements-12               97.00 ± 0%
Search/google/btree/1000_elements-12                     976.0 ± 0%
Search/krzysztofgb/gostree/10000_elements-12             0.000 ± 0%
Search/ajwerner/orderstat/10000_elements-12              99.00 ± 0%
Search/google/btree/10000_elements-12                   9.982k ± 0%
```

</details>

<details>

<summary>Delete</summary>

```text
                                                      │   sec/op    │
Delete/krzysztofgb/gostree/100_elements-12              5.399µ ± 0%
Delete/ajwerner/orderstat/100_elements-12               29.83µ ± 0%
Delete/google/btree/100_elements-12                     973.7n ± 1%
Delete/krzysztofgb/gostree/1000_elements-12             8.271µ ± 1%
Delete/ajwerner/orderstat/1000_elements-12              49.12µ ± 1%
Delete/google/btree/1000_elements-12                    9.477µ ± 0%
Delete/krzysztofgb/gostree/10000_elements-12            15.12µ ± 1%
Delete/ajwerner/orderstat/10000_elements-12             67.31µ ± 0%
Delete/google/btree/10000_elements-12                   94.75µ ± 0%

                                                      │      B/op      │
Delete/krzysztofgb/gostree/100_elements-12                0.000 ± 0%
Delete/ajwerner/orderstat/100_elements-12                 620.0 ± 0%
Delete/google/btree/100_elements-12                       920.0 ± 0%
Delete/krzysztofgb/gostree/1000_elements-12               0.000 ± 0%
Delete/ajwerner/orderstat/1000_elements-12                778.0 ± 0%
Delete/google/btree/1000_elements-12                    7.859Ki ± 0%
Delete/krzysztofgb/gostree/10000_elements-12              0.000 ± 0%
Delete/ajwerner/orderstat/10000_elements-12               799.0 ± 0%
Delete/google/btree/10000_elements-12                   78.25Ki ± 0%


                                                      │   allocs/op   │
Delete/krzysztofgb/gostree/100_elements-12               0.000 ± 0%
Delete/ajwerner/orderstat/100_elements-12                77.00 ± 1%
Delete/google/btree/100_elements-12                      80.00 ± 0%
Delete/krzysztofgb/gostree/1000_elements-12              0.000 ± 0%
Delete/ajwerner/orderstat/1000_elements-12               96.00 ± 0%
Delete/google/btree/1000_elements-12                     971.0 ± 0%
Delete/krzysztofgb/gostree/10000_elements-12             0.000 ± 0%
Delete/ajwerner/orderstat/10000_elements-12              99.00 ± 0%
Delete/google/btree/10000_elements-12                   9.981k ± 0%
```

</details>

<details>

<summary>Select</summary>

```text
                                                      │   sec/op    │
Select/krzysztofgb/gostree/100_elements-12              2.811µ ± 0%
Select/ajwerner/orderstat/100_elements-12               3.666µ ± 0%
Select/krzysztofgb/gostree/1000_elements-12             3.783µ ± 1%
Select/ajwerner/orderstat/1000_elements-12              5.449µ ± 0%
Select/krzysztofgb/gostree/10000_elements-12            6.829µ ± 1%
Select/ajwerner/orderstat/10000_elements-12             8.870µ ± 1%

                                                      │      B/op      │
Select/krzysztofgb/gostree/100_elements-12                0.000 ± 0%
Select/ajwerner/orderstat/100_elements-12                 0.000 ± 0%
Select/krzysztofgb/gostree/1000_elements-12               0.000 ± 0%
Select/ajwerner/orderstat/1000_elements-12                0.000 ± 0%
Select/krzysztofgb/gostree/10000_elements-12              0.000 ± 0%
Select/ajwerner/orderstat/10000_elements-12               0.000 ± 0%

                                                      │   allocs/op   │
Select/krzysztofgb/gostree/100_elements-12               0.000 ± 0%
Select/ajwerner/orderstat/100_elements-12                0.000 ± 0%
Select/krzysztofgb/gostree/1000_elements-12              0.000 ± 0%
Select/ajwerner/orderstat/1000_elements-12               0.000 ± 0%
Select/krzysztofgb/gostree/10000_elements-12             0.000 ± 0%
Select/ajwerner/orderstat/10000_elements-12              0.000 ± 0%
```

</details>

<details>

<summary>Rank</summary>

```text
                                                      │   sec/op    │
Rank/krzysztofgb/gostree/100_elements-12                2.908µ ± 0%
Rank/ajwerner/orderstat/100_elements-12                 7.827µ ± 0%
Rank/krzysztofgb/gostree/1000_elements-12               3.930µ ± 1%
Rank/ajwerner/orderstat/1000_elements-12                14.07µ ± 0%
Rank/krzysztofgb/gostree/10000_elements-12              6.542µ ± 1%
Rank/ajwerner/orderstat/10000_elements-12               22.56µ ± 0%

                                                      │      B/op      │
Rank/krzysztofgb/gostree/100_elements-12                  0.000 ± 0%
Rank/ajwerner/orderstat/100_elements-12                   647.0 ± 0%
Rank/krzysztofgb/gostree/1000_elements-12                 0.000 ± 0%
Rank/ajwerner/orderstat/1000_elements-12                  782.0 ± 0%
Rank/krzysztofgb/gostree/10000_elements-12                0.000 ± 0%
Rank/ajwerner/orderstat/10000_elements-12                 798.0 ± 0%

                                                      │   allocs/op   │
Rank/krzysztofgb/gostree/100_elements-12                 0.000 ± 0%
Rank/ajwerner/orderstat/100_elements-12                  80.00 ± 1%
Rank/krzysztofgb/gostree/1000_elements-12                0.000 ± 0%
Rank/ajwerner/orderstat/1000_elements-12                 97.00 ± 0%
Rank/krzysztofgb/gostree/10000_elements-12               0.000 ± 0%
Rank/ajwerner/orderstat/10000_elements-12                99.00 ± 0%
```

</details>

<details>

<summary>Mixed (Even Split) </summary>

```text
                                                      │   sec/op    │
MixedOperations/krzysztofgb/gostree/100_elements-12     4.942µ ± 2%
MixedOperations/ajwerner/orderstat/100_elements-12      14.61µ ± 0%
MixedOperations/google/btree/100_elements-12            9.688µ ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12    6.674µ ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12     21.86µ ± 0%
MixedOperations/google/btree/1000_elements-12           14.32µ ± 7%
MixedOperations/krzysztofgb/gostree/10000_elements-12   10.68µ ± 0%
MixedOperations/ajwerner/orderstat/10000_elements-12    34.15µ ± 0%
MixedOperations/google/btree/10000_elements-12          23.00µ ± 0%

                                                      │      B/op      │
MixedOperations/krzysztofgb/gostree/100_elements-12       960.0 ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12        452.0 ± 0%
MixedOperations/google/btree/100_elements-12            2.026Ki ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12      960.0 ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12       629.0 ± 0%
MixedOperations/google/btree/1000_elements-12           2.074Ki ± 0%
MixedOperations/krzysztofgb/gostree/10000_elements-12     960.0 ± 0%
MixedOperations/ajwerner/orderstat/10000_elements-12      639.0 ± 0%
MixedOperations/google/btree/10000_elements-12          2.702Ki ± 0%

                                                      │   allocs/op   │
MixedOperations/krzysztofgb/gostree/100_elements-12      20.00 ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12       56.00 ± 2%
MixedOperations/google/btree/100_elements-12             93.00 ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12     20.00 ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12      78.00 ± 0%
MixedOperations/google/btree/1000_elements-12            106.0 ± 0%
MixedOperations/krzysztofgb/gostree/10000_elements-12    20.00 ± 0%
MixedOperations/ajwerner/orderstat/10000_elements-12     79.00 ± 0%
MixedOperations/google/btree/10000_elements-12           124.0 ± 0%
```

</details>

## License

See [LICENSE](LICENSE) file for details.