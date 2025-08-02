package gostree

import (
	"testing"
	"math/rand"
	"github.com/ajwerner/orderstat"
	"github.com/google/btree"
)

var randGen *rand.Rand

func init() {
	randGen = rand.New(rand.NewSource(1337)) // Fixed seed for deterministic benchmarks
}

type orderstatInt int

func (a orderstatInt) Less(b orderstat.Item) bool {
	return a < b.(orderstatInt)
}

type btreeInt int

func (b btreeInt) Less(c btree.Item) bool {
	return b < c.(btreeInt)
}

// generateRandomData creates a slice of random integers
func generateRandomData(n int) []int {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = randGen.Intn(n * 10)
	}
	return data
}

func BenchmarkInsert(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"100_elements", 100},
		{"1000_elements", 1000},
		{"10000_elements", 10000},
	}

	for _, bm := range benchmarks {
		data := generateRandomData(bm.size)

		b.Run("krzysztofgb/gostree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := NewTree[int]()
				for _, v := range data {
					tree.Insert(v)
				}
			}
		})

		b.Run("ajwerner/orderstat/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := orderstat.NewTree()
				for _, v := range data {
					tree.ReplaceOrInsert(orderstatInt(v))
				}
			}
		})

		b.Run("google/btree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := btree.New(2)
				for _, v := range data {
					tree.ReplaceOrInsert(btreeInt(v))
				}
			}
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"100_elements", 100},
		{"1000_elements", 1000},
		{"10000_elements", 10000},
	}

	for _, bm := range benchmarks {
		data := generateRandomData(bm.size)

		// Setup gostree
		gostreeTree := NewTree[int]()
		for _, v := range data {
			gostreeTree.Insert(v)
		}

		// Setup orderstat
		orderstatTree := orderstat.NewTree()
		for _, v := range data {
			orderstatTree.ReplaceOrInsert(orderstatInt(v))
		}

		b.Run("krzysztofgb/gostree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 100; j++ {
					gostreeTree.Search(data[randGen.Intn(len(data))])
				}
			}
		})

		b.Run("ajwerner/orderstat/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 100; j++ {
					orderstatTree.Get(orderstatInt(data[randGen.Intn(len(data))]))
				}
			}
		})

		b.Run("google/btree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := btree.New(2)
				for _, v := range data {
					tree.Get(btreeInt(v))
				}
			}
		})
	}
}

func BenchmarkSelect(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"100_elements", 100},
		{"1000_elements", 1000},
		{"10000_elements", 10000},
	}

	for _, bm := range benchmarks {
		data := generateRandomData(bm.size)

		// Setup gostree
		gostreeTree := NewTree[int]()
		for _, v := range data {
			gostreeTree.Insert(v)
		}

		// Setup orderstat
		orderstatTree := orderstat.NewTree()
		for _, v := range data {
			orderstatTree.ReplaceOrInsert(orderstatInt(v))
		}

		b.Run("krzysztofgb/gostree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 100; j++ {
					gostreeTree.Select(randGen.Intn(bm.size))
				}
			}
		})

		b.Run("ajwerner/orderstat/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 100; j++ {
					orderstatTree.Select(randGen.Intn(bm.size))
				}
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"100_elements", 100},
		{"1000_elements", 1000},
		{"10000_elements", 10000},
	}

	for _, bm := range benchmarks {
		data := generateRandomData(bm.size)

		b.Run("krzysztofgb/gostree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				tree := NewTree[int]()
				for _, v := range data {
					tree.Insert(v)
				}
				b.StartTimer()

				for j := 0; j < 100; j++ {
					tree.Delete(data[randGen.Intn(len(data))])
				}
			}
		})

		b.Run("ajwerner/orderstat/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				tree := orderstat.NewTree()
				for _, v := range data {
					tree.ReplaceOrInsert(orderstatInt(v))
				}
				b.StartTimer()

				for j := 0; j < 100; j++ {
					tree.Delete(orderstatInt(data[randGen.Intn(len(data))]))
				}
			}
		})

		b.Run("google/btree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := btree.New(2)
				for _, v := range data {
					tree.Delete(btreeInt(v))
				}
			}
		})
	}
}

func BenchmarkRank(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"100_elements", 100},
		{"1000_elements", 1000},
		{"10000_elements", 10000},
	}

	for _, bm := range benchmarks {
		data := generateRandomData(bm.size)

		// Setup gostree
		gostreeTree := NewTree[int]()
		for _, v := range data {
			gostreeTree.Insert(v)
		}

		// Setup orderstat
		orderstatTree := orderstat.NewTree()
		for _, v := range data {
			orderstatTree.ReplaceOrInsert(orderstatInt(v))
		}

		b.Run("krzysztofgb/gostree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 100; j++ {
					gostreeTree.Rank(data[randGen.Intn(len(data))])
				}
			}
		})

		b.Run("ajwerner/orderstat/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 100; j++ {
					orderstatTree.Rank(orderstatInt(data[randGen.Intn(len(data))]))
				}
			}
		})
	}
}

func BenchmarkMixedOperations(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"100_elements", 100},
		{"1000_elements", 1000},
		{"10000_elements", 10000},
	}

	for _, bm := range benchmarks {
		data := generateRandomData(bm.size)

		b.Run("krzysztofgb/gostree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				tree := NewTree[int]()
				// Pre-populate with initial data
				for _, v := range data[:bm.size/2] {
					tree.Insert(v)
				}
				b.StartTimer()

				// Mixed operations: 20% each of insert, search, select, delete, rank
				for j := 0; j < 100; j++ {
					switch j % 5 {
					case 0:
						tree.Insert(data[randGen.Intn(len(data))])
					case 1:
						tree.Search(data[randGen.Intn(len(data))])
					case 2:
						if tree.root.size > 0 {
							tree.Select(randGen.Intn(tree.root.size))
						}
					case 3:
						tree.Delete(data[randGen.Intn(len(data))])
					case 4:
						tree.Rank(data[randGen.Intn(len(data))])
					}
				}
			}
		})

		b.Run("ajwerner/orderstat/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				tree := orderstat.NewTree()
				// Pre-populate with initial data
				for _, v := range data[:bm.size/2] {
					tree.ReplaceOrInsert(orderstatInt(v))
				}
				b.StartTimer()

				// Mixed operations: 20% each of insert, search, select, delete, rank
				for j := 0; j < 100; j++ {
					switch j % 5 {
					case 0:
						tree.ReplaceOrInsert(orderstatInt(data[randGen.Intn(len(data))]))
					case 1:
						tree.Get(orderstatInt(data[randGen.Intn(len(data))]))
					case 2:
						if tree.Len() > 0 {
							tree.Select(randGen.Intn(tree.Len()))
						}
					case 3:
						tree.Delete(orderstatInt(data[randGen.Intn(len(data))]))
					case 4:
						tree.Rank(orderstatInt(data[randGen.Intn(len(data))]))
					}
				}
			}
		})

		b.Run("google/btree/"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				tree := btree.New(2)
				// Pre-populate with initial data
				for _, v := range data[:bm.size/2] {
					tree.ReplaceOrInsert(btreeInt(v))
				}
				b.StartTimer()

				// Mixed operations: 33.3% each of insert, get, delete
				for j := 0; j < 100; j++ {
					switch j % 3 {
					case 0:
						tree.ReplaceOrInsert(btreeInt(data[randGen.Intn(len(data))]))
					case 1:
						tree.Get(btreeInt(data[randGen.Intn(len(data))]))
					case 3:
						tree.Delete(btreeInt(data[randGen.Intn(len(data))]))
					}
				}
			}
		})
	}
}
