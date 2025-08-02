package gostree

import (
	"testing"
	"math/rand"
	"github.com/ajwerner/orderstat"
)

var randGen *rand.Rand

func init() {
	randGen = rand.New(rand.NewSource(1337)) // Fixed seed for deterministic benchmarks
}

type orderstatInt int

func (a orderstatInt) Less(b orderstat.Item) bool {
	return a < b.(orderstatInt)
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
					case 0: // Insert
						tree.Insert(data[randGen.Intn(len(data))])
					case 1: // Search
						tree.Search(data[randGen.Intn(len(data))])
					case 2: // Select (k-th element)
						if tree.root.size > 0 {
							tree.Select(randGen.Intn(tree.root.size))
						}
					case 3: // Delete
						tree.Delete(data[randGen.Intn(len(data))])
					case 4: // Rank
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
					case 0: // Insert
						tree.ReplaceOrInsert(orderstatInt(data[randGen.Intn(len(data))]))
					case 1: // Search
						tree.Get(orderstatInt(data[randGen.Intn(len(data))]))
					case 2: // Select (k-th element)
						if tree.Len() > 0 {
							tree.Select(randGen.Intn(tree.Len()))
						}
					case 3: // Delete
						tree.Delete(orderstatInt(data[randGen.Intn(len(data))]))
					case 4: // Rank
						tree.Rank(orderstatInt(data[randGen.Intn(len(data))]))
					}
				}
			}
		})
	}
}
