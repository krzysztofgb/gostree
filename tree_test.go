package gostree

import (
	"cmp"
	"testing"
)

// Red-Black tree properties:
// 1. The root is black.
// 2. All leaves (NIL nodes) are black.
// 3. If a red node has children, both children are black (no two reds in a row).
// 4. Every path from a node to its descendant leaves has the same number of black nodes.
func checkRedBlackProperties[T cmp.Ordered](t *testing.T, tree *Tree[T]) {
	t.Helper()

	if tree.root != tree.nil && tree.root.color != BLACK {
		t.Error("Property 1 violated: root is not BLACK")
	}

	checkNoRedRedViolation(t, tree.root, tree.nil)

	blackHeight := -1
	checkBlackHeight(t, tree.root, tree.nil, 0, &blackHeight)
}

func checkNoRedRedViolation[T cmp.Ordered](t *testing.T, node, sentinel *Node[T]) {
	t.Helper()

	if node == sentinel {
		return
	}

	if node.color == RED {
		if node.left.color == RED {
			t.Errorf("Red-red violation: node %v has red left child", node.key)
		}
		if node.right.color == RED {
			t.Errorf("Red-red violation: node %v has red right child", node.key)
		}
	}

	checkNoRedRedViolation(t, node.left, sentinel)
	checkNoRedRedViolation(t, node.right, sentinel)
}

func checkBlackHeight[T cmp.Ordered](t *testing.T, node, sentinel *Node[T], currentBlackHeight int, blackHeight *int) {
	t.Helper()

	if node == sentinel {
		if *blackHeight == -1 {
			*blackHeight = currentBlackHeight
		} else if *blackHeight != currentBlackHeight {
			t.Errorf("Black height violation: expected %d, got %d", *blackHeight, currentBlackHeight)
		}
		return
	}

	if node.color == BLACK {
		currentBlackHeight++
	}

	checkBlackHeight(t, node.left, sentinel, currentBlackHeight, blackHeight)
	checkBlackHeight(t, node.right, sentinel, currentBlackHeight, blackHeight)
}

func verifySizes[T cmp.Ordered](t *testing.T, node, sentinel *Node[T]) int {
	t.Helper()

	if node == sentinel {
		return 0
	}

	leftSize := verifySizes(t, node.left, sentinel)
	rightSize := verifySizes(t, node.right, sentinel)
	expectedSize := leftSize + rightSize + 1

	if node.size != expectedSize {
		t.Errorf("Size mismatch at node %v: has %d, expected %d",
			node.key, node.size, expectedSize)
	}

	return expectedSize
}

func buildTree(values []int) *Tree[int] {
	tree := NewTree[int]()
	for _, v := range values {
		tree.Insert(v)
	}
	return tree
}

func TestNewTree(t *testing.T) {
	t.Run("creates_valid_empty_tree", func(t *testing.T) {
		tree := NewTree[int]()
		if tree == nil {
			t.Fatal("NewTree returned nil")
		}
		if tree.nil == nil {
			t.Fatal("sentinel node is nil")
		}
		if tree.nil.color != BLACK {
			t.Errorf("sentinel color = %v, want BLACK", tree.nil.color)
		}
		if tree.nil.size != 0 {
			t.Errorf("sentinel size = %d, want 0", tree.nil.size)
		}
		if tree.root != tree.nil {
			t.Error("root does not point to sentinel")
		}
	})

	t.Run("sentinel_is_self_referential", func(t *testing.T) {
		tree := NewTree[int]()
		if tree.nil.left != tree.nil || tree.nil.right != tree.nil || tree.nil.parent != tree.nil {
			t.Error("sentinel is not properly self-referential")
		}
	})

	t.Run("works_with_different_types", func(t *testing.T) {
		intTree := NewTree[int]()
		stringTree := NewTree[string]()
		floatTree := NewTree[float64]()

		trees := []interface{}{intTree, stringTree, floatTree}
		for i, tree := range trees {
			if tree == nil {
				t.Errorf("tree %d is nil", i)
			}
		}

		// Verify sentinel has zero value for each type
		if intTree.nil.key != 0 {
			t.Errorf("int sentinel key = %d, want 0", intTree.nil.key)
		}
		if stringTree.nil.key != "" {
			t.Errorf("string sentinel key = %q, want empty string", stringTree.nil.key)
		}
	})

	t.Run("multiple_trees_are_independent", func(t *testing.T) {
		tree1 := NewTree[int]()
		tree2 := NewTree[int]()

		if tree1 == tree2 || tree1.nil == tree2.nil {
			t.Error("trees are not independent")
		}
	})
}

func TestInsert(t *testing.T) {
	t.Run("single_element", func(t *testing.T) {
		tree := NewTree[int]()
		tree.Insert(10)

		if tree.root.key != 10 || tree.root.color != BLACK || tree.root.size != 1 {
			t.Error("root properties incorrect")
		}
		if tree.root.parent != tree.nil || tree.root.left != tree.nil || tree.root.right != tree.nil {
			t.Error("root links incorrect")
		}
	})

	t.Run("maintains_bst_property", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15})

		if tree.root.key != 10 {
			t.Errorf("root = %d, want 10", tree.root.key)
		}
		if tree.root.left.key != 5 || tree.root.right.key != 15 {
			t.Error("BST property violated")
		}
	})

	t.Run("updates_sizes_correctly", func(t *testing.T) {
		tree := NewTree[int]()

		tree.Insert(10)
		if tree.root.size != 1 {
			t.Errorf("after 1 insert: root.size = %d, want 1", tree.root.size)
		}

		tree.Insert(5)
		if tree.root.size != 2 {
			t.Errorf("after 2 inserts: root.size = %d, want 2", tree.root.size)
		}

		tree.Insert(15)
		if tree.root.size != 3 {
			t.Errorf("after 3 inserts: root.size = %d, want 3", tree.root.size)
		}

		verifySizes(t, tree.root, tree.nil)
	})

	t.Run("triggers_left_rotation", func(t *testing.T) {
		tree := buildTree([]int{10, 20, 30})

		// After rotation, 20 should be root
		if tree.root.key != 20 {
			t.Errorf("root = %d, want 20 after rotation", tree.root.key)
		}
		if tree.root.left.key != 10 || tree.root.right.key != 30 {
			t.Error("children incorrect after rotation")
		}

		checkRedBlackProperties(t, tree)
	})

	t.Run("triggers_right_rotation", func(t *testing.T) {
		tree := buildTree([]int{30, 20, 10})

		// After rotation, 20 should be root
		if tree.root.key != 20 {
			t.Errorf("root = %d, want 20 after rotation", tree.root.key)
		}
		if tree.root.left.key != 10 || tree.root.right.key != 30 {
			t.Error("children incorrect after rotation")
		}

		checkRedBlackProperties(t, tree)
	})

	t.Run("complex_insertions_maintain_properties", func(t *testing.T) {
		values := []int{13, 8, 17, 1, 11, 15, 25, 6, 22, 27}
		tree := NewTree[int]()

		for i, v := range values {
			tree.Insert(v)

			if tree.root.size != i+1 {
				t.Errorf("after inserting %d values: size = %d, want %d", i+1, tree.root.size, i+1)
			}

			checkRedBlackProperties(t, tree)
			verifySizes(t, tree.root, tree.nil)
		}
	})

	t.Run("handles_duplicates", func(t *testing.T) {
		tree := NewTree[int]()
		tree.Insert(10)
		tree.Insert(10) // duplicate
		tree.Insert(5)
		tree.Insert(10) // another duplicate
		tree.Insert(15)

		checkRedBlackProperties(t, tree)
		verifySizes(t, tree.root, tree.nil)
	})
}

func TestSearch(t *testing.T) {
	t.Run("empty_tree", func(t *testing.T) {
		tree := NewTree[int]()
		if tree.Search(10) {
			t.Error("Search should return false for empty tree")
		}
	})

	t.Run("single_element", func(t *testing.T) {
		tree := buildTree([]int{42})

		if !tree.Search(42) {
			t.Error("Search should find existing element")
		}
		if tree.Search(41) || tree.Search(43) {
			t.Error("Search should not find non-existing elements")
		}
	})

	t.Run("multiple_elements", func(t *testing.T) {
		tree := buildTree([]int{50, 30, 70, 20, 40, 60, 80})

		// existing elements
		existing := []int{50, 30, 70, 20, 40, 60, 80}
		for _, v := range existing {
			if !tree.Search(v) {
				t.Errorf("Search should find %d", v)
			}
		}

		// non-existing elements
		nonExisting := []int{10, 25, 35, 45, 55, 65, 75, 85}
		for _, v := range nonExisting {
			if tree.Search(v) {
				t.Errorf("Search should not find %d", v)
			}
		}
	})

	t.Run("search_after_rotations", func(t *testing.T) {
		// Insert values that cause rotations
		tree := buildTree([]int{1, 2, 3, 4, 5, 6, 7})

		// All values should still be searchable
		for i := 1; i <= 7; i++ {
			if !tree.Search(i) {
				t.Errorf("Search should find %d after rotations", i)
			}
		}
	})

	t.Run("search_doesnt_modify_tree", func(t *testing.T) {
		tree := buildTree([]int{50, 30, 70, 20, 40})

		sizeBefore := tree.root.size
		tree.Search(20)
		tree.Search(40)
		tree.Search(35) // non-existing
		sizeAfter := tree.root.size

		if sizeBefore != sizeAfter {
			t.Error("Search should not modify tree size")
		}
	})
}

func TestSelect(t *testing.T) {
	t.Run("empty_tree", func(t *testing.T) {
		tree := NewTree[int]()

		testCases := []int{-1, 0, 1}
		for _, idx := range testCases {
			if _, ok := tree.Select(idx); ok {
				t.Errorf("Select(%d) should return false for empty tree", idx)
			}
		}
	})

	t.Run("single_element", func(t *testing.T) {
		tree := buildTree([]int{42})

		val, ok := tree.Select(0)
		if !ok || val != 42 {
			t.Errorf("Select(0) = %d, %v; want 42, true", val, ok)
		}

		if _, ok := tree.Select(1); ok {
			t.Error("Select(1) should return false for single element tree")
		}

		if _, ok := tree.Select(-1); ok {
			t.Error("Select(-1) should return false")
		}
	})

	t.Run("returns_sorted_order", func(t *testing.T) {
		tree := buildTree([]int{30, 10, 50, 20, 40, 60, 70})

		expected := []int{10, 20, 30, 40, 50, 60, 70}
		for i, want := range expected {
			got, ok := tree.Select(i)
			if !ok || got != want {
				t.Errorf("Select(%d) = %d, %v; want %d, true", i, got, ok, want)
			}
		}
	})

	t.Run("boundary_conditions", func(t *testing.T) {
		tree := buildTree([]int{1, 2, 3, 4, 5})

		// First element
		val, ok := tree.Select(0)
		if !ok || val != 1 {
			t.Errorf("Select(0) = %d, %v; want 1, true", val, ok)
		}

		// Last element
		val, ok = tree.Select(4)
		if !ok || val != 5 {
			t.Errorf("Select(4) = %d, %v; want 5, true", val, ok)
		}

		// Out of bounds
		if _, ok := tree.Select(5); ok {
			t.Error("Select(5) should return false")
		}

		if _, ok := tree.Select(-1); ok {
			t.Error("Select(-1) should return false")
		}
	})

	t.Run("select_with_duplicates", func(t *testing.T) {
		tree := NewTree[int]()
		for _, v := range []int{5, 3, 7, 3, 5, 7} {
			tree.Insert(v)
		}

		// Expected order: 3, 3, 5, 5, 7, 7
		expected := []int{3, 3, 5, 5, 7, 7}
		for i, want := range expected {
			got, ok := tree.Select(i)
			if !ok || got != want {
				t.Errorf("Select(%d) = %d, %v; want %d, true", i, got, ok, want)
			}
		}
	})
}

func TestRank(t *testing.T) {
	t.Run("empty_tree", func(t *testing.T) {
		tree := NewTree[int]()

		rank := tree.Rank(10)
		if rank != 0 {
			t.Errorf("Rank in empty tree = %d, want 0", rank)
		}
	})

	t.Run("single_element", func(t *testing.T) {
		tree := buildTree([]int{42})

		testCases := []struct {
			key      int
			expected int
		}{
			{40, 0}, // less than root
			{42, 0}, // equal to root
			{50, 1}, // greater than root
		}

		for _, tc := range testCases {
			rank := tree.Rank(tc.key)
			if rank != tc.expected {
				t.Errorf("Rank(%d) = %d, want %d", tc.key, rank, tc.expected)
			}
		}
	})

	t.Run("multiple_elements", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15, 3, 7, 12, 20})

		// Test rank of existing elements
		testCases := []struct {
			key      int
			expected int
		}{
			{3, 0}, // minimum
			{5, 1},
			{7, 2},
			{10, 3},
			{12, 4},
			{15, 5},
			{20, 6}, // maximum
			{1, 0},  // less than min
			{25, 7}, // greater than max
			{8, 3},  // between elements
		}

		for _, tc := range testCases {
			rank := tree.Rank(tc.key)
			if rank != tc.expected {
				t.Errorf("Rank(%d) = %d, want %d", tc.key, rank, tc.expected)
			}
		}
	})

	t.Run("rank_with_duplicates", func(t *testing.T) {
		tree := NewTree[int]()
		for _, v := range []int{5, 3, 7, 3, 5, 7} {
			tree.Insert(v)
		}

		// Rank should return the rank of the first occurrence
		testCases := []struct {
			key      int
			expected int
		}{
			{3, 0},  // two 3's at positions 0,1
			{5, 2},  // two 5's at positions 2,3
			{7, 4},  // two 7's at positions 4,5
			{10, 6}, // all elements < 10
		}

		for _, tc := range testCases {
			rank := tree.Rank(tc.key)
			if rank != tc.expected {
				t.Errorf("Rank(%d) = %d, want %d", tc.key, rank, tc.expected)
			}
		}
	})

	t.Run("rank_select_consistency", func(t *testing.T) {
		tree := buildTree([]int{15, 6, 18, 3, 7, 17, 20, 2, 4, 13, 9})

		// For each element, verify that Select(Rank(x)) returns x
		values := []int{2, 3, 4, 6, 7, 9, 13, 15, 17, 18, 20}
		for _, v := range values {
			rank := tree.Rank(v)
			selected, ok := tree.Select(rank)
			if !ok || selected != v {
				t.Errorf("Select(Rank(%d)) = %d, want %d", v, selected, v)
			}
		}

		// For each position, verify that Rank(Select(i)) == i
		for i := 0; i < len(values); i++ {
			selected, ok := tree.Select(i)
			if !ok {
				t.Errorf("Select(%d) failed", i)
				continue
			}
			rank := tree.Rank(selected)
			if rank != i {
				t.Errorf("Rank(Select(%d)) = %d, want %d", i, rank, i)
			}
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("empty_tree", func(t *testing.T) {
		tree := NewTree[int]()

		if tree.Delete(10) {
			t.Error("Delete should return false for empty tree")
		}
	})

	t.Run("single_element", func(t *testing.T) {
		tree := buildTree([]int{42})

		if !tree.Delete(42) {
			t.Error("Delete should return true when deleting existing element")
		}

		if tree.root != tree.nil {
			t.Error("Tree should be empty after deleting only element")
		}

		if tree.Delete(42) {
			t.Error("Delete should return false when element no longer exists")
		}
	})

	t.Run("delete_leaf_nodes", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15, 3, 7, 12, 17})

		if !tree.Delete(3) {
			t.Error("Should delete leaf node 3")
		}
		if !tree.Delete(17) {
			t.Error("Should delete leaf node 17")
		}

		if tree.Search(3) || tree.Search(17) {
			t.Error("Deleted nodes should not be found")
		}

		if tree.root.size != 5 {
			t.Errorf("Root size = %d, want 5", tree.root.size)
		}

		checkRedBlackProperties(t, tree)
		verifySizes(t, tree.root, tree.nil)
	})

	t.Run("delete_node_with_one_child", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15, 12})

		if !tree.Delete(15) {
			t.Error("Should delete node with one child")
		}

		if tree.root.right.key != 12 {
			t.Errorf("Right child = %v, want 12", tree.root.right.key)
		}

		checkRedBlackProperties(t, tree)
		verifySizes(t, tree.root, tree.nil)
	})

	t.Run("delete_node_with_two_children", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15, 3, 7, 12, 17})

		if !tree.Delete(5) {
			t.Error("Should delete node with two children")
		}

		if tree.Search(5) {
			t.Error("5 should not exist after deletion")
		}

		for _, v := range []int{10, 3, 7, 15, 12, 17} {
			if !tree.Search(v) {
				t.Errorf("%d should still exist", v)
			}
		}

		checkRedBlackProperties(t, tree)
		verifySizes(t, tree.root, tree.nil)
	})

	t.Run("delete_root", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15})

		oldRootKey := tree.root.key
		if !tree.Delete(oldRootKey) {
			t.Error("Should delete root")
		}

		if tree.Search(oldRootKey) {
			t.Errorf("Old root %v should not exist", oldRootKey)
		}

		// New root should be one of the children
		if tree.root.key != 5 && tree.root.key != 15 {
			t.Errorf("New root = %v, expected 5 or 15", tree.root.key)
		}

		checkRedBlackProperties(t, tree)
		verifySizes(t, tree.root, tree.nil)
	})

	t.Run("complex_deletion_sequence", func(t *testing.T) {
		tree := buildTree([]int{7, 3, 18, 10, 22, 8, 11, 26, 2, 6, 13})

		deleteOrder := []int{6, 18, 11, 3, 10}
		expectedSizes := []int{10, 9, 8, 7, 6}

		for i, v := range deleteOrder {
			if !tree.Delete(v) {
				t.Errorf("Should delete %d", v)
			}

			if tree.root.size != expectedSizes[i] {
				t.Errorf("After deleting %d: size = %d, want %d", v, tree.root.size, expectedSizes[i])
			}

			checkRedBlackProperties(t, tree)
			verifySizes(t, tree.root, tree.nil)
		}
	})

	t.Run("delete_all_elements", func(t *testing.T) {
		values := []int{10, 5, 15, 3, 7, 12, 17}
		tree := buildTree(values)

		for _, v := range values {
			if !tree.Delete(v) {
				t.Errorf("Should delete %d", v)
			}
		}

		if tree.root != tree.nil {
			t.Error("Tree should be empty after deleting all elements")
		}
	})

	t.Run("delete_updates_order_statistics", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15, 3, 7, 12, 17, 1, 4, 6, 8})

		tree.Delete(7)

		expected := []int{1, 3, 4, 5, 6, 8, 10, 12, 15, 17}
		for i, want := range expected {
			got, ok := tree.Select(i)
			if !ok || got != want {
				t.Errorf("Select(%d) = %d, %v; want %d, true", i, got, ok, want)
			}
		}

		for i, v := range expected {
			rank := tree.Rank(v)
			if rank != i {
				t.Errorf("Rank(%d) = %d, want %d", v, rank, i)
			}
		}
	})
}

func TestIntegration(t *testing.T) {
	t.Run("large_dataset_operations", func(t *testing.T) {
		tree := NewTree[int]()
		n := 1000

		for i := 0; i < n; i += 2 {
			tree.Insert(i)
		}

		for i := 0; i < n; i += 2 {
			if !tree.Search(i) {
				t.Errorf("should find %d", i)
			}
		}

		for i := 1; i < n; i += 2 {
			if tree.Search(i) {
				t.Errorf("should not find %d", i)
			}
		}

		// Delete half the even numbers
		for i := 0; i < n/2; i += 2 {
			tree.Delete(i)
		}

		for i := 0; i < n/2; i += 2 {
			if tree.Search(i) {
				t.Errorf("should not find deleted %d", i)
			}
		}

		checkRedBlackProperties(t, tree)
		verifySizes(t, tree.root, tree.nil)
	})

	t.Run("stress_test_all_operations", func(t *testing.T) {
		tree := NewTree[int]()

		// Insert 100 random-ish values
		values := make([]int, 100)
		for i := range values {
			values[i] = (i * 7) % 100
			tree.Insert(values[i])
		}

		// Test rank/select consistency
		for i := 0; i < tree.root.size; i++ {
			val, ok := tree.Select(i)
			if !ok {
				t.Errorf("Select(%d) failed", i)
				continue
			}

			rank := tree.Rank(val)
			if rank > i {
				t.Errorf("Rank(Select(%d)) = %d, should be <= %d", i, rank, i)
			}
		}

		for i := 0; i < len(values); i += 3 {
			tree.Delete(values[i])
		}

		checkRedBlackProperties(t, tree)
		verifySizes(t, tree.root, tree.nil)
	})
}

func TestSize(t *testing.T) {
	t.Run("empty_tree", func(t *testing.T) {
		tree := NewTree[int]()
		if size := tree.Size(); size != 0 {
			t.Errorf("Size() = %d, want 0 for empty tree", size)
		}
	})

	t.Run("single_element", func(t *testing.T) {
		tree := NewTree[int]()
		tree.Insert(42)
		if size := tree.Size(); size != 1 {
			t.Errorf("Size() = %d, want 1", size)
		}
	})

	t.Run("multiple_elements", func(t *testing.T) {
		tree := NewTree[int]()
		elements := []int{10, 5, 15, 3, 7, 12, 17}
		
		for i, v := range elements {
			tree.Insert(v)
			expectedSize := i + 1
			if size := tree.Size(); size != expectedSize {
				t.Errorf("After inserting %d elements: Size() = %d, want %d", expectedSize, size, expectedSize)
			}
		}
	})

	t.Run("with_duplicates", func(t *testing.T) {
		tree := NewTree[int]()
		elements := []int{5, 3, 7, 3, 5, 7, 3}
		
		for i, v := range elements {
			tree.Insert(v)
			expectedSize := i + 1
			if size := tree.Size(); size != expectedSize {
				t.Errorf("After inserting element %d: Size() = %d, want %d", v, size, expectedSize)
			}
		}
	})

	t.Run("after_deletions", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15, 3, 7, 12, 17})
		initialSize := tree.Size()
		if initialSize != 7 {
			t.Errorf("Initial size = %d, want 7", initialSize)
		}

		// Delete some elements
		tree.Delete(3)
		if size := tree.Size(); size != 6 {
			t.Errorf("After deleting 3: Size() = %d, want 6", size)
		}

		tree.Delete(15)
		if size := tree.Size(); size != 5 {
			t.Errorf("After deleting 15: Size() = %d, want 5", size)
		}

		tree.Delete(7)
		if size := tree.Size(); size != 4 {
			t.Errorf("After deleting 7: Size() = %d, want 4", size)
		}
	})

	t.Run("delete_non_existing", func(t *testing.T) {
		tree := buildTree([]int{10, 5, 15})
		originalSize := tree.Size()
		
		// Try to delete non-existing element
		tree.Delete(20)
		if size := tree.Size(); size != originalSize {
			t.Errorf("Size after deleting non-existing element = %d, want %d", size, originalSize)
		}
	})

	t.Run("large_tree", func(t *testing.T) {
		tree := NewTree[int]()
		n := 1000
		
		for i := 0; i < n; i++ {
			tree.Insert(i)
		}
		
		if size := tree.Size(); size != n {
			t.Errorf("Size() = %d, want %d", size, n)
		}
		
		// Delete half the elements
		for i := 0; i < n/2; i++ {
			tree.Delete(i)
		}
		
		if size := tree.Size(); size != n/2 {
			t.Errorf("After deleting half: Size() = %d, want %d", size, n/2)
		}
	})

	t.Run("consistency_with_select", func(t *testing.T) {
		tree := buildTree([]int{30, 10, 50, 20, 40, 60, 70})
		size := tree.Size()
		
		// The last valid index for Select should be size-1
		if _, ok := tree.Select(size - 1); !ok {
			t.Errorf("Select(%d) should succeed for tree of size %d", size-1, size)
		}
		
		if _, ok := tree.Select(size); ok {
			t.Errorf("Select(%d) should fail for tree of size %d", size, size)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("specific_rotation_patterns", func(t *testing.T) {
		patterns := []struct {
			name   string
			values []int
		}{
			{"right_parent_case2", []int{10, 5, 15, 12, 20, 11}},
			{"left_parent_case2", []int{20, 15, 25, 10, 18, 19}},
			{"simple_right_rotation", []int{30, 20, 10}},
			{"complex_rebalancing", []int{50, 25, 75, 12, 37, 62, 87, 6, 18, 31, 43}},
		}

		for _, p := range patterns {
			t.Run(p.name, func(t *testing.T) {
				tree := buildTree(p.values)
				checkRedBlackProperties(t, tree)
				verifySizes(t, tree.root, tree.nil)
			})
		}
	})

	t.Run("deletion_fixup_patterns", func(t *testing.T) {
		// Build a larger tree to exercise deletion fixup cases
		tree := buildTree([]int{50, 25, 75, 12, 37, 62, 87, 6, 18, 31, 43, 56, 68, 81, 93,
			3, 9, 15, 21, 28, 34, 40, 46, 53, 59, 65, 71, 78, 84, 90, 96})

		// Delete in specific order to trigger different fixup scenarios
		deleteOrder := []int{3, 6, 9, 12, 15, 18, 21, 25}

		for _, v := range deleteOrder {
			tree.Delete(v)
			checkRedBlackProperties(t, tree)
		}
	})
}
