<p align="center">
    <img src="./assets/gostree.png" style="border-radius: 20%" alt="gostree logo"/>
</p>

---

# gostree

[![Go Reference](https://pkg.go.dev/badge/github.com/krzysztofgb/gostree.svg)](https://pkg.go.dev/github.com/krzysztofgb/gostree)

A generic [order-statistic tree](https://en.wikipedia.org/wiki/Order_statistic_tree) tree implementation in Go.

## Installation

```bash
go get github.com/krzysztofgb/gostree
```

## Usage

```go
import "github.com/krzysztofgb/gostree"

// Create a new tree with a comparison function
tree := gostree.NewTree[int](func(a, b int) int {
    return a - b
})

// Insert elements
tree.Insert(5)
tree.Insert(3)
tree.Insert(7)
tree.Insert(1)
tree.Insert(9)

// Check if an element exists
exists := tree.Search(5)  // Returns true

// Get the k-th smallest element (0-indexed)
val, found := tree.Select(2)  // Returns 5, true (3rd smallest)

// Get the rank (number of elements less than x)
rank := tree.Rank(6)  // Returns 3 (elements 1, 3, 5 are less than 6)

// Delete an element
deleted := tree.Delete(3)  // Returns true

// Get the size
size := tree.Size()  // Returns 4
```

### Custom Types

You can use the tree with any type by providing an appropriate comparison function:

```go
type Person struct {
    Name string
    Age  int
}

// Create a tree sorted by age
tree := gostree.NewTree[Person](func(a, b Person) int {
    return a.Age - b.Age
})

// Create a tree sorted by name, then by age
tree := gostree.NewTree[Person](func(a, b Person) int {
    if a.Name < b.Name {
        return -1
    } else if a.Name > b.Name {
        return 1
    }
    // Names are equal, compare by age
    return a.Age - b.Age
})
```

## Concurrency Safety

**Write operations are NOT concurrent safe.**
The following methods modify the tree structure and require external synchronization when used concurrently:
- `Insert()`
- `Delete()`

**Read operations ARE concurrent safe.**
Multiple goroutines can safely call these methods simultaneously without external synchronization:
- `Search()`
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
Insert/krzysztofgb/gostree/100_elements-12              3.847µ ± 1%
Insert/ajwerner/orderstat/100_elements-12               10.47µ ± 0%
Insert/google/btree/100_elements-12                     13.58µ ± 1%
Insert/krzysztofgb/gostree/1000_elements-12             61.71µ ± 1%
Insert/ajwerner/orderstat/1000_elements-12              186.1µ ± 6%
Insert/google/btree/1000_elements-12                    213.3µ ± 0%
Insert/krzysztofgb/gostree/10000_elements-12            999.3µ ± 1%
Insert/ajwerner/orderstat/10000_elements-12             2.669m ± 0%
Insert/google/btree/10000_elements-12                   2.903m ± 1%

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
Search/krzysztofgb/gostree/100_elements-12              3.137µ ± 1%
Search/ajwerner/orderstat/100_elements-12               6.469µ ± 0%
Search/google/btree/100_elements-12                     936.2n ± 1%
Search/krzysztofgb/gostree/1000_elements-12             4.638µ ± 1%
Search/ajwerner/orderstat/1000_elements-12              10.08µ ± 0%
Search/google/btree/1000_elements-12                    8.829µ ± 0%
Search/krzysztofgb/gostree/10000_elements-12            8.021µ ± 0%
Search/ajwerner/orderstat/10000_elements-12             16.02µ ± 0%
Search/google/btree/10000_elements-12                   88.47µ ± 0%

                                                      │      B/op      │
Search/krzysztofgb/gostree/100_elements-12                0.000 ± 0%
Search/ajwerner/orderstat/100_elements-12                 647.0 ± 0%
Search/google/btree/100_elements-12                       984.0 ± 0%
Search/krzysztofgb/gostree/1000_elements-12               0.000 ± 0%
Search/ajwerner/orderstat/1000_elements-12                778.0 ± 0%
Search/google/btree/1000_elements-12                    7.930Ki ± 0%
Search/krzysztofgb/gostree/10000_elements-12              0.000 ± 0%
Search/ajwerner/orderstat/10000_elements-12               797.0 ± 0%
Search/google/btree/10000_elements-12                   78.23Ki ± 0%

                                                      │   allocs/op   │
Search/krzysztofgb/gostree/100_elements-12               0.000 ± 0%
Search/ajwerner/orderstat/100_elements-12                80.00 ± 1%
Search/google/btree/100_elements-12                      85.00 ± 0%
Search/krzysztofgb/gostree/1000_elements-12              0.000 ± 0%
Search/ajwerner/orderstat/1000_elements-12               97.00 ± 0%
Search/google/btree/1000_elements-12                     977.0 ± 0%
Search/krzysztofgb/gostree/10000_elements-12             0.000 ± 0%
Search/ajwerner/orderstat/10000_elements-12              99.00 ± 0%
Search/google/btree/10000_elements-12                   9.975k ± 0%
```

</details>

<details>

<summary>Delete</summary>

```text
                                                      │   sec/op    │
Delete/krzysztofgb/gostree/100_elements-12              6.092µ ± 0%
Delete/ajwerner/orderstat/100_elements-12               30.36µ ± 1%
Delete/google/btree/100_elements-12                     956.2n ± 0%
Delete/krzysztofgb/gostree/1000_elements-12             9.977µ ± 5%
Delete/ajwerner/orderstat/1000_elements-12              49.88µ ± 0%
Delete/google/btree/1000_elements-12                    9.574µ ± 1%
Delete/krzysztofgb/gostree/10000_elements-12            16.30µ ± 1%
Delete/ajwerner/orderstat/10000_elements-12             68.71µ ± 0%
Delete/google/btree/10000_elements-12                   94.88µ ± 1%

                                                      │      B/op      │
Delete/krzysztofgb/gostree/100_elements-12                0.000 ± 0%
Delete/ajwerner/orderstat/100_elements-12                 611.0 ± 0%
Delete/google/btree/100_elements-12                       912.0 ± 0%
Delete/krzysztofgb/gostree/1000_elements-12               0.000 ± 0%
Delete/ajwerner/orderstat/1000_elements-12                789.0 ± 0%
Delete/google/btree/1000_elements-12                    7.969Ki ± 0%
Delete/krzysztofgb/gostree/10000_elements-12              0.000 ± 0%
Delete/ajwerner/orderstat/10000_elements-12               799.0 ± 0%
Delete/google/btree/10000_elements-12                   78.20Ki ± 0%


                                                      │   allocs/op   │
Delete/krzysztofgb/gostree/100_elements-12               0.000 ± 0%
Delete/ajwerner/orderstat/100_elements-12                75.00 ± 1%
Delete/google/btree/100_elements-12                      79.00 ± 0%
Delete/krzysztofgb/gostree/1000_elements-12              0.000 ± 0%
Delete/ajwerner/orderstat/1000_elements-12               98.00 ± 0%
Delete/google/btree/1000_elements-12                     985.0 ± 0%
Delete/krzysztofgb/gostree/10000_elements-12             0.000 ± 0%
Delete/ajwerner/orderstat/10000_elements-12              99.00 ± 0%
Delete/google/btree/10000_elements-12                   9.974k ± 0%
```

</details>

<details>

<summary>Select</summary>

```text
                                                      │   sec/op    │
Select/krzysztofgb/gostree/100_elements-12              2.807µ ± 1%
Select/ajwerner/orderstat/100_elements-12               3.710µ ± 0%
Select/krzysztofgb/gostree/1000_elements-12             3.873µ ± 1%
Select/ajwerner/orderstat/1000_elements-12              5.461µ ± 1%
Select/krzysztofgb/gostree/10000_elements-12            6.917µ ± 0%
Select/ajwerner/orderstat/10000_elements-12             8.752µ ± 0%

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
Rank/krzysztofgb/gostree/100_elements-12                3.611µ ± 0%
Rank/ajwerner/orderstat/100_elements-12                 7.936µ ± 0%
Rank/krzysztofgb/gostree/1000_elements-12               5.093µ ± 0%
Rank/ajwerner/orderstat/1000_elements-12                14.09µ ± 0%
Rank/krzysztofgb/gostree/10000_elements-12              9.123µ ± 0%
Rank/ajwerner/orderstat/10000_elements-12               22.69µ ± 1%

                                                      │      B/op      │
Rank/krzysztofgb/gostree/100_elements-12                  0.000 ± 0%
Rank/ajwerner/orderstat/100_elements-12                   551.5 ± 0%
Rank/krzysztofgb/gostree/1000_elements-12                 0.000 ± 0%
Rank/ajwerner/orderstat/1000_elements-12                  777.0 ± 0%
Rank/krzysztofgb/gostree/10000_elements-12                0.000 ± 0%
Rank/ajwerner/orderstat/10000_elements-12                 797.0 ± 0%

                                                      │   allocs/op   │
Rank/krzysztofgb/gostree/100_elements-12                 0.000 ± 0%
Rank/ajwerner/orderstat/100_elements-12                  68.50 ± 1%
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
MixedOperations/krzysztofgb/gostree/100_elements-12     5.488µ ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12      14.49µ ± 3%
MixedOperations/google/btree/100_elements-12            9.957µ ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12    7.532µ ± 2%
MixedOperations/ajwerner/orderstat/1000_elements-12     21.80µ ± 3%
MixedOperations/google/btree/1000_elements-12           15.24µ ± 1%
MixedOperations/krzysztofgb/gostree/10000_elements-12   11.73µ ± 1%
MixedOperations/ajwerner/orderstat/10000_elements-12    33.75µ ± 0%
MixedOperations/google/btree/10000_elements-12          22.80µ ± 1%

                                                      │      B/op      │
MixedOperations/krzysztofgb/gostree/100_elements-12       960.0 ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12        458.0 ± 0%
MixedOperations/google/btree/100_elements-12            2.267Ki ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12      960.0 ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12       627.0 ± 0%
MixedOperations/google/btree/1000_elements-12           3.046Ki ± 0%
MixedOperations/krzysztofgb/gostree/10000_elements-12     960.0 ± 0%
MixedOperations/ajwerner/orderstat/10000_elements-12      639.0 ± 0%
MixedOperations/google/btree/10000_elements-12          2.662Ki ± 0%

                                                      │   allocs/op   │
MixedOperations/krzysztofgb/gostree/100_elements-12      20.00 ± 0%
MixedOperations/ajwerner/orderstat/100_elements-12       56.00 ± 0%
MixedOperations/google/btree/100_elements-12             98.00 ± 0%
MixedOperations/krzysztofgb/gostree/1000_elements-12     20.00 ± 0%
MixedOperations/ajwerner/orderstat/1000_elements-12      77.50 ± 1%
MixedOperations/google/btree/1000_elements-12            135.0 ± 0%
MixedOperations/krzysztofgb/gostree/10000_elements-12    20.00 ± 0%
MixedOperations/ajwerner/orderstat/10000_elements-12     79.00 ± 0%
MixedOperations/google/btree/10000_elements-12           123.0 ± 0%
```

</details>

## License

See [LICENSE](LICENSE) file for details.