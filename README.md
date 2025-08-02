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

## Performance

### Benchmark Results

Benchmarks were run against:
* [ajwerner/orderstat](https://github.com/ajwerner/orderstat)
* [google/btree](https://github.com/google/btree)

Note, benchmarks for "Rank" and "Select" were omitted for `google/btree`, as they are not supported by that implementation.

<details>

<summary>Insert</summary>

```text
                                             │   sec/op    │
Insert/krzysztofgb/gostree/100_elements-12     3.330µ ± 1%
Insert/ajwerner/orderstat/100_elements-12      10.70µ ± 2%
Insert/google/btree/100_elements-12            14.30µ ± 3%
Insert/krzysztofgb/gostree/1000_elements-12    54.48µ ± 3%
Insert/ajwerner/orderstat/1000_elements-12     189.3µ ± 1%
Insert/google/btree/1000_elements-12           220.2µ ± 2%
Insert/krzysztofgb/gostree/10000_elements-12   915.5µ ± 1%
Insert/ajwerner/orderstat/10000_elements-12    2.703m ± 0%
Insert/google/btree/10000_elements-12          3.013m ± 1%

                                             │     B/op     │
Insert/krzysztofgb/gostree/100_elements-12     4.734Ki ± 0%
Insert/ajwerner/orderstat/100_elements-12      8.023Ki ± 0%
Insert/google/btree/100_elements-12            8.914Ki ± 0%
Insert/krzysztofgb/gostree/1000_elements-12    46.92Ki ± 0%
Insert/ajwerner/orderstat/1000_elements-12     71.13Ki ± 0%
Insert/google/btree/1000_elements-12           94.31Ki ± 0%
Insert/krzysztofgb/gostree/10000_elements-12   468.8Ki ± 0%
Insert/ajwerner/orderstat/10000_elements-12    1.076Mi ± 0%
Insert/google/btree/10000_elements-12          920.8Ki ± 0%

                                             │  allocs/op  │
Insert/krzysztofgb/gostree/100_elements-12      101.0 ± 0%
Insert/ajwerner/orderstat/100_elements-12       71.00 ± 0%
Insert/google/btree/100_elements-12             282.0 ± 0%
Insert/krzysztofgb/gostree/1000_elements-12    1.001k ± 0%
Insert/ajwerner/orderstat/1000_elements-12      984.0 ± 0%
Insert/google/btree/1000_elements-12           3.251k ± 0%
Insert/krzysztofgb/gostree/10000_elements-12   10.00k ± 0%
Insert/ajwerner/orderstat/10000_elements-12    9.987k ± 0%
Insert/google/btree/10000_elements-12          32.10k ± 0%
```

</details>

<details>

<summary>Search</summary>

```text
                                             │   sec/op    │
Search/krzysztofgb/gostree/100_elements-12     2.670µ ± 1%
Search/ajwerner/orderstat/100_elements-12      6.448µ ± 1%
Search/google/btree/100_elements-12            987.1n ± 0%
Search/krzysztofgb/gostree/1000_elements-12    3.760µ ± 0%
Search/ajwerner/orderstat/1000_elements-12     10.10µ ± 0%
Search/google/btree/1000_elements-12           9.089µ ± 0%
Search/krzysztofgb/gostree/10000_elements-12   6.520µ ± 1%
Search/ajwerner/orderstat/10000_elements-12    16.18µ ± 0%
Search/google/btree/10000_elements-12          91.92µ ± 0%

                                             │     B/op     │
Search/krzysztofgb/gostree/100_elements-12       0.000 ± 0%
Search/ajwerner/orderstat/100_elements-12        536.0 ± 0%
Search/google/btree/100_elements-12              984.0 ± 0%
Search/krzysztofgb/gostree/1000_elements-12      0.000 ± 0%
Search/ajwerner/orderstat/1000_elements-12       778.0 ± 0%
Search/google/btree/1000_elements-12           7.820Ki ± 0%
Search/krzysztofgb/gostree/10000_elements-12     0.000 ± 0%
Search/ajwerner/orderstat/10000_elements-12      797.0 ± 0%
Search/google/btree/10000_elements-12          78.23Ki ± 0%

                                             │  allocs/op   │
Search/krzysztofgb/gostree/100_elements-12      0.000 ± 0%
Search/ajwerner/orderstat/100_elements-12       67.00 ± 1%
Search/google/btree/100_elements-12             85.00 ± 0%
Search/krzysztofgb/gostree/1000_elements-12     0.000 ± 0%
Search/ajwerner/orderstat/1000_elements-12      97.00 ± 0%
Search/google/btree/1000_elements-12            963.0 ± 0%
Search/krzysztofgb/gostree/10000_elements-12    0.000 ± 0%
Search/ajwerner/orderstat/10000_elements-12     99.00 ± 0%
Search/google/btree/10000_elements-12          9.976k ± 0%
```

</details>

<details>

<summary>Select</summary>

```text
                                             │   sec/op    │
Select/krzysztofgb/gostree/100_elements-12     2.764µ ± 0%
Select/ajwerner/orderstat/100_elements-12      3.620µ ± 1%
Select/krzysztofgb/gostree/1000_elements-12    3.814µ ± 1%
Select/ajwerner/orderstat/1000_elements-12     5.485µ ± 0%
Select/krzysztofgb/gostree/10000_elements-12   6.912µ ± 1%
Select/ajwerner/orderstat/10000_elements-12    8.822µ ± 0%

                                             │     B/op     │
Select/krzysztofgb/gostree/100_elements-12     0.000 ± 0%
Select/ajwerner/orderstat/100_elements-12      0.000 ± 0%
Select/krzysztofgb/gostree/1000_elements-12    0.000 ± 0%
Select/ajwerner/orderstat/1000_elements-12     0.000 ± 0%
Select/krzysztofgb/gostree/10000_elements-12   0.000 ± 0%
Select/ajwerner/orderstat/10000_elements-12    0.000 ± 0%

                                             │  allocs/op   │
Select/krzysztofgb/gostree/100_elements-12     0.000 ± 0%
Select/ajwerner/orderstat/100_elements-12      0.000 ± 0%
Select/krzysztofgb/gostree/1000_elements-12    0.000 ± 0%
Select/ajwerner/orderstat/1000_elements-12     0.000 ± 0%
Select/krzysztofgb/gostree/10000_elements-12   0.000 ± 0%
Select/ajwerner/orderstat/10000_elements-12    0.000 ± 0%
```

</details>

<details>

<summary>Rank</summary>

```text
                                           │   sec/op    │
Rank/krzysztofgb/gostree/100_elements-12     2.839µ ± 0%
Rank/ajwerner/orderstat/100_elements-12      7.993µ ± 0%
Rank/krzysztofgb/gostree/1000_elements-12    3.870µ ± 0%
Rank/ajwerner/orderstat/1000_elements-12     14.31µ ± 0%
Rank/krzysztofgb/gostree/10000_elements-12   6.433µ ± 0%
Rank/ajwerner/orderstat/10000_elements-12    23.11µ ± 1%

                                           │     B/op     │
Rank/krzysztofgb/gostree/100_elements-12     0.000 ± 0%
Rank/ajwerner/orderstat/100_elements-12      536.0 ± 0%
Rank/krzysztofgb/gostree/1000_elements-12    0.000 ± 0%
Rank/ajwerner/orderstat/1000_elements-12     780.0 ± 0%
Rank/krzysztofgb/gostree/10000_elements-12   0.000 ± 0%
Rank/ajwerner/orderstat/10000_elements-12    797.0 ± 0%

                                           │  allocs/op   │
Rank/krzysztofgb/gostree/100_elements-12     0.000 ± 0%
Rank/ajwerner/orderstat/100_elements-12      67.00 ± 1%
Rank/krzysztofgb/gostree/1000_elements-12    0.000 ± 0%
Rank/ajwerner/orderstat/1000_elements-12     97.00 ± 0%
Rank/krzysztofgb/gostree/10000_elements-12   0.000 ± 0%
Rank/ajwerner/orderstat/10000_elements-12    99.00 ± 0%
```

</details>

<details>

<summary>Delete</summary>

```text
                                             │   sec/op    │
Delete/krzysztofgb/gostree/100_elements-12     5.481µ ± 0%
Delete/ajwerner/orderstat/100_elements-12      30.71µ ± 0%
Delete/google/btree/100_elements-12            1.056µ ± 0%
Delete/krzysztofgb/gostree/1000_elements-12    8.453µ ± 1%
Delete/ajwerner/orderstat/1000_elements-12     48.18µ ± 2%
Delete/google/btree/1000_elements-12           9.820µ ± 1%
Delete/krzysztofgb/gostree/10000_elements-12   15.13µ ± 4%
Delete/ajwerner/orderstat/10000_elements-12    68.96µ ± 0%
Delete/google/btree/10000_elements-12          98.01µ ± 0%

                                             │     B/op     │
Delete/krzysztofgb/gostree/100_elements-12       0.000 ± 0%
Delete/ajwerner/orderstat/100_elements-12        539.5 ± 0%
Delete/google/btree/100_elements-12              992.0 ± 0%
Delete/krzysztofgb/gostree/1000_elements-12      0.000 ± 0%
Delete/ajwerner/orderstat/1000_elements-12       777.0 ± 0%
Delete/google/btree/1000_elements-12           7.883Ki ± 0%
Delete/krzysztofgb/gostree/10000_elements-12     0.000 ± 0%
Delete/ajwerner/orderstat/10000_elements-12      799.0 ± 0%
Delete/google/btree/10000_elements-12          78.21Ki ± 0%

                                             │  allocs/op   │
Delete/krzysztofgb/gostree/100_elements-12     0.000 ± 0%
Delete/ajwerner/orderstat/100_elements-12      66.50 ± 1%
Delete/google/btree/100_elements-12            89.00 ± 0%
Delete/krzysztofgb/gostree/1000_elements-12    0.000 ± 0%
Delete/ajwerner/orderstat/1000_elements-12     96.00 ± 0%
Delete/google/btree/1000_elements-12           974.0 ± 0%
Delete/krzysztofgb/gostree/10000_elements-12   0.000 ± 0%
Delete/ajwerner/orderstat/10000_elements-12    99.00 ± 0%
Delete/google/btree/10000_elements-12         9.976k ± 0%
```

</details>

<details>

<summary>Mixed (Even Split) </summary>

```text
                                                      │   sec/op    │
MixedOperations/krzysztofgb/gostree/100_elements-12     4.970µ ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12      14.63µ ± 1%
MixedOperations/google/btree/100_elements-12            10.37µ ± 1%
MixedOperations/krzysztofgb/gostree/1000_elements-12    6.719µ ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12     21.23µ ± 0%
MixedOperations/google/btree/1000_elements-12           15.37µ ± 0%
MixedOperations/krzysztofgb/gostree/10000_elements-12   10.53µ ± 1%
MixedOperations/ajwerner/orderstat/10000_elements-12    33.90µ ± 0%
MixedOperations/google/btree/10000_elements-12          24.04µ ± 1%

                                                      │     B/op     │
MixedOperations/krzysztofgb/gostree/100_elements-12       960.0 ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12        432.0 ± 0%
MixedOperations/google/btree/100_elements-12            2.497Ki ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12      960.0 ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12       626.0 ± 0%
MixedOperations/google/btree/1000_elements-12           2.844Ki ± 0%
MixedOperations/krzysztofgb/gostree/10000_elements-12     960.0 ± 0%
MixedOperations/ajwerner/orderstat/10000_elements-12      639.0 ± 0%
MixedOperations/google/btree/10000_elements-12          3.009Ki ± 0%

                                                      │ allocs/op  │
MixedOperations/krzysztofgb/gostree/100_elements-12     20.00 ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12      53.00 ± 0%
MixedOperations/google/btree/100_elements-12            102.0 ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12    20.00 ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12     77.00 ± 0%
MixedOperations/google/btree/1000_elements-12           129.0 ± 0%
MixedOperations/krzysztofgb/gostree/10000_elements-12   20.00 ± 0%
MixedOperations/ajwerner/orderstat/10000_elements-12    79.00 ± 0%
MixedOperations/google/btree/10000_elements-12          134.0 ± 1%
```

</details>

## License

See [LICENSE](LICENSE) file for details.