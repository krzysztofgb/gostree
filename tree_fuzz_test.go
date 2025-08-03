package gostree

import (
	"cmp"
	"testing"
)

// Operation codes for fuzzing
const (
	opInsert byte = iota
	opDelete
	opSelect
	opRank
	opSearch
)

// FuzzTree tests the tree with random operations
func FuzzTree(f *testing.F) {
	// Add seed corpus with various operation sequences
	f.Add([]byte{opInsert, 10, opInsert, 20, opInsert, 30})                                           // Simple insertions
	f.Add([]byte{opInsert, 50, opInsert, 10, opInsert, 90, opInsert, 20, opInsert, 30, opInsert, 40}) // Larger sequence
	f.Add([]byte{opInsert, 10, opInsert, 20, opInsert, 30, opDelete, 10})                             // Insert then delete
	f.Add([]byte{opInsert, 10, opInsert, 10, opInsert, 10})                                           // Duplicate keys
	f.Add([]byte{opDelete, 10, opDelete, 20, opDelete, 30})                                           // All deletes
	f.Add([]byte{opInsert, 10, opDelete, 10, opInsert, 10, opDelete, 10})                             // Alternating insert/delete
	f.Add([]byte{opSelect, 0, opSelect, 5, opSelect, 10})                                             // Select operations
	f.Add([]byte{opRank, 50, opRank, 10, opRank, 90})                                                 // Rank operations
	f.Add([]byte{opSearch, 10, opSearch, 20, opSearch, 30})                                           // Search operations
	f.Add([]byte{opInsert, 10, opSearch, 10, opDelete, 10, opSearch, 10})                             // Search existing and non-existing

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 2 {
			return
		}

		tree := NewTree[int]()

		// Track what we've inserted for validation
		elements := make(map[int]int) // value -> count

		// Process operations in pairs (operation, value)
		for i := 0; i < len(data)-1; i += 2 {
			op := data[i] % 5
			value := int(data[i+1])

			switch op {
			case opInsert:
				tree.Insert(value)
				elements[value]++

				checkRedBlackProperties(t, tree)
				verifyOrderStatisticProperties(t, tree, elements)
			case opDelete:
				beforeSize := tree.Size()
				deleted := tree.Delete(value)

				if elements[value] > 0 {
					if !deleted {
						t.Fatalf("Failed to delete existing element %d", value)
					}
					elements[value]--
					if elements[value] == 0 {
						delete(elements, value)
					}
				} else {
					if deleted {
						t.Fatalf("Successfully deleted non-existent element %d", value)
					}
				}

				afterSize := tree.Size()
				if deleted && afterSize != beforeSize-1 {
					t.Fatalf("Size not updated correctly after delete: before=%d, after=%d", beforeSize, afterSize)
				}
				if !deleted && afterSize != beforeSize {
					t.Fatalf("Size changed after failed delete: before=%d, after=%d", beforeSize, afterSize)
				}

				checkRedBlackProperties(t, tree)
				verifyOrderStatisticProperties(t, tree, elements)
			case opSelect:
				if tree.Size() > 0 {
					k := value % tree.Size()
					elem, ok := tree.Select(k)
					if !ok {
						t.Fatalf("Select(%d) failed on tree of size %d", k, tree.Size())
					}

					// Verify the selected element is correct
					// Note: Rank returns the position of the first occurrence of a value,
					// so with duplicates, rank <= k < rank + count(elem)
					rank := tree.Rank(elem)
					if rank > k {
						t.Fatalf("Select/Rank mismatch: Select(%d)=%d, but Rank(%d)=%d (rank > k)", k, elem, elem, rank)
					}
					// Verify that elem is at position k by checking elements before and after
					if k > 0 {
						prevElem, _ := tree.Select(k - 1)
						if prevElem > elem {
							t.Fatalf("Select returned wrong order: Select(%d)=%d > Select(%d)=%d", k-1, prevElem, k, elem)
						}
					}
					if k < tree.Size()-1 {
						nextElem, _ := tree.Select(k + 1)
						if nextElem < elem {
							t.Fatalf("Select returned wrong order: Select(%d)=%d < Select(%d)=%d", k, elem, k+1, nextElem)
						}
					}
				}
			case opRank:
				rank := tree.Rank(value)
				if rank < 0 || rank > tree.Size() {
					t.Fatalf("Rank(%d) returned invalid value %d for tree of size %d", value, rank, tree.Size())
				}

				// If the value exists, verify we can select it back
				if tree.Search(value) {
					elem, ok := tree.Select(rank)
					if !ok || elem > value {
						t.Fatalf("Rank/Select mismatch: Rank(%d)=%d, but Select(%d)=%d", value, rank, rank, elem)
					}
				}
			case opSearch:
				found := tree.Search(value)
				expected := elements[value] > 0
				if found != expected {
					t.Fatalf("Search(%d) returned %v, expected %v", value, found, expected)
				}
			}
		}

		checkRedBlackProperties(t, tree)
		verifyOrderStatisticProperties(t, tree, elements)
		verifyTreeIntegrity(t, tree)
	})
}

// verifyOrderStatisticProperties checks that size fields are correct
func verifyOrderStatisticProperties[T cmp.Ordered](t *testing.T, tree *Tree[T], elements map[int]int) {
	totalCount := 0
	for _, count := range elements {
		totalCount += count
	}

	if tree.Size() != totalCount {
		t.Fatalf("Tree size mismatch: tree.Size()=%d, expected=%d", tree.Size(), totalCount)
	}

	verifySizeFields(t, tree, tree.root, tree.nil)
}

// verifySizeFields recursively verifies that size fields are correct
func verifySizeFields[T cmp.Ordered](t *testing.T, tree *Tree[T], node, nil *Node[T]) int {
	if node == nil {
		return 0
	}

	leftSize := verifySizeFields(t, tree, node.left, nil)
	rightSize := verifySizeFields(t, tree, node.right, nil)
	expectedSize := leftSize + rightSize + 1

	if node.size != expectedSize {
		t.Fatalf("Size field mismatch at node: expected %d, got %d", expectedSize, node.size)
	}

	return expectedSize
}

// verifyTreeIntegrity performs additional integrity checks
func verifyTreeIntegrity[T cmp.Ordered](t *testing.T, tree *Tree[T]) {
	// Verify in-order traversal produces sorted sequence
	var values []T
	inOrderTraversal(tree, tree.root, tree.nil, &values)

	for i := 1; i < len(values); i++ {
		if values[i] < values[i-1] {
			t.Fatalf("Tree not in sorted order: %v < %v at positions %d, %d", values[i], values[i-1], i, i-1)
		}
	}

	// Verify Select returns elements in order
	for i := 0; i < tree.Size(); i++ {
		elem, ok := tree.Select(i)
		if !ok {
			t.Fatalf("Select(%d) failed", i)
		}
		if i < len(values) && elem != values[i] {
			t.Fatalf("Select(%d) returned %v, expected %v", i, elem, values[i])
		}
	}
}

// inOrderTraversal performs in-order traversal
func inOrderTraversal[T cmp.Ordered](tree *Tree[T], node, nil *Node[T], values *[]T) {
	if node == nil {
		return
	}

	inOrderTraversal(tree, node.left, nil, values)
	*values = append(*values, node.key)
	inOrderTraversal(tree, node.right, nil, values)
}
