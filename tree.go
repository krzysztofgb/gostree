package gostree

import "cmp"

type Color bool

const (
	RED   Color = false
	BLACK Color = true
)

type Node[T cmp.Ordered] struct {
	key    T
	left   *Node[T]
	right  *Node[T]
	parent *Node[T]
	color  Color
	size   int // number of nodes in subtree rooted at this node
}

type Tree[T cmp.Ordered] struct {
	root *Node[T]
	nil  *Node[T] // sentinel node
}

// NewTree creates a new order-statistic red-black tree.
func NewTree[T cmp.Ordered]() *Tree[T] {
	t := &Tree[T]{}
	// Create sentinel node
	t.nil = &Node[T]{
		color: BLACK,
		size:  0,
	}
	// Make sentinel self-referential
	t.nil.left = t.nil
	t.nil.right = t.nil
	t.nil.parent = t.nil
	// Initialize root to sentinel
	t.root = t.nil
	return t
}

// Insert adds a new key to the red-black tree
// and maintains the red-black properties.
func (t *Tree[T]) Insert(key T) {
	newNode := &Node[T]{
		key:    key,
		left:   t.nil,
		right:  t.nil,
		parent: t.nil,
		color:  RED,
		size:   1,
	}

	parent := t.nil
	current := t.root

	// Find insertion position
	for current != t.nil {
		parent = current
		// Update size on the path down
		current.size++
		if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}

	// Insert the new node
	newNode.parent = parent
	if parent == t.nil {
		t.root = newNode
	} else if newNode.key < parent.key {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	// Fix red-black properties
	t.insertFixup(newNode)
}

// insertFixup maintains red-black tree properties after insertion
func (t *Tree[T]) insertFixup(newNode *Node[T]) {
	for newNode.parent.color == RED {
		parent := newNode.parent
		grandparent := parent.parent

		if parent == grandparent.left {
			uncle := grandparent.right
			if uncle.color == RED {
				// Case 1: Uncle is RED - recolor and move up
				parent.color = BLACK
				uncle.color = BLACK
				grandparent.color = RED
				newNode = grandparent
			} else {
				if newNode == parent.right {
					// Case 2: Node is right child - rotate left to convert to case 3
					newNode = parent
					t.leftRotate(newNode)
				}
				// Case 3: Node is left child - rotate right and recolor
				newNode.parent.color = BLACK
				newNode.parent.parent.color = RED
				t.rightRotate(newNode.parent.parent)
			}
		} else {
			// Mirror cases: parent is right child of grandparent
			uncle := grandparent.left
			if uncle.color == RED {
				// Case 1: Uncle is RED - recolor and move up
				parent.color = BLACK
				uncle.color = BLACK
				grandparent.color = RED
				newNode = grandparent
			} else {
				if newNode == parent.left {
					// Case 2: Node is left child - rotate right to convert to case 3
					newNode = parent
					t.rightRotate(newNode)
				}
				// Case 3: Node is right child - rotate left and recolor
				newNode.parent.color = BLACK
				newNode.parent.parent.color = RED
				t.leftRotate(newNode.parent.parent)
			}
		}
	}
	t.root.color = BLACK
}

// leftRotate performs a left rotation on the given node
func (t *Tree[T]) leftRotate(node *Node[T]) {
	rightChild := node.right
	node.right = rightChild.left
	if rightChild.left != t.nil {
		rightChild.left.parent = node
	}
	rightChild.parent = node.parent
	if node.parent == t.nil {
		t.root = rightChild
	} else if node == node.parent.left {
		node.parent.left = rightChild
	} else {
		node.parent.right = rightChild
	}
	rightChild.left = node
	node.parent = rightChild

	// Update sizes
	node.size = node.left.size + node.right.size + 1
	rightChild.size = rightChild.left.size + rightChild.right.size + 1
}

// rightRotate performs a right rotation on the given node
func (t *Tree[T]) rightRotate(node *Node[T]) {
	leftChild := node.left
	node.left = leftChild.right
	if leftChild.right != t.nil {
		leftChild.right.parent = node
	}
	leftChild.parent = node.parent
	if node.parent == t.nil {
		t.root = leftChild
	} else if node == node.parent.right {
		node.parent.right = leftChild
	} else {
		node.parent.left = leftChild
	}
	leftChild.right = node
	node.parent = leftChild

	// Update sizes
	node.size = node.left.size + node.right.size + 1
	leftChild.size = leftChild.left.size + leftChild.right.size + 1
}

// Search checks if a key exists in the tree.
// It returns true if the key is found, false otherwise.
func (t *Tree[T]) Search(key T) bool {
	return t.search(key) != t.nil
}

func (t *Tree[T]) search(key T) *Node[T] {
	current := t.root
	for current != t.nil && key != current.key {
		if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}
	return current
}

// Select returns the k-th smallest element (0-indexed).
func (t *Tree[T]) Select(k int) (T, bool) {
	var zero T
	if k < 0 || k >= t.root.size {
		return zero, false
	}

	node := t.selectNode(t.root, k)
	return node.key, true
}

func (t *Tree[T]) selectNode(current *Node[T], k int) *Node[T] {
	for current != t.nil {
		leftSize := current.left.size
		if k < leftSize {
			current = current.left
		} else if k == leftSize {
			return current
		} else {
			k -= leftSize + 1
			current = current.right
		}
	}
	return current
}

// Rank returns the number of elements less than the given key.
func (t *Tree[T]) Rank(key T) int {
	rank := 0
	current := t.root

	for current != t.nil {
		if key < current.key {
			current = current.left
		} else if key > current.key {
			rank += current.left.size + 1
			current = current.right
		} else {
			// key == current.key
			rank += current.left.size
			break
		}
	}

	return rank
}

// Delete removes one occurrence of a key from the tree.
func (t *Tree[T]) Delete(key T) bool {
	nodeToDelete := t.search(key)
	if nodeToDelete == t.nil {
		return false
	}

	t.deleteNode(nodeToDelete)
	return true
}

func (t *Tree[T]) deleteNode(nodeToDelete *Node[T]) {
	nodeActuallyDeleted := nodeToDelete
	originalColor := nodeActuallyDeleted.color
	var replacementNode *Node[T]

	if nodeToDelete.left == t.nil {
		// Node has no left child
		replacementNode = nodeToDelete.right
		t.transplant(nodeToDelete, nodeToDelete.right)
	} else if nodeToDelete.right == t.nil {
		// Node has no right child
		replacementNode = nodeToDelete.left
		t.transplant(nodeToDelete, nodeToDelete.left)
	} else {
		// Node has two children - find successor
		nodeActuallyDeleted = t.minimum(nodeToDelete.right)
		originalColor = nodeActuallyDeleted.color
		replacementNode = nodeActuallyDeleted.right

		if nodeActuallyDeleted.parent == nodeToDelete {
			replacementNode.parent = nodeActuallyDeleted
		} else {
			t.transplant(nodeActuallyDeleted, nodeActuallyDeleted.right)
			nodeActuallyDeleted.right = nodeToDelete.right
			nodeActuallyDeleted.right.parent = nodeActuallyDeleted
		}

		t.transplant(nodeToDelete, nodeActuallyDeleted)
		nodeActuallyDeleted.left = nodeToDelete.left
		nodeActuallyDeleted.left.parent = nodeActuallyDeleted
		nodeActuallyDeleted.color = nodeToDelete.color
	}

	// Update sizes from the deletion point upward
	t.updateSizeUpward(replacementNode.parent)

	if originalColor == BLACK {
		t.deleteFixup(replacementNode)
	}
}

// transplant replaces subtree rooted at nodeToReplace with subtree rooted at replacement
func (t *Tree[T]) transplant(nodeToReplace, replacement *Node[T]) {
	if nodeToReplace.parent == t.nil {
		t.root = replacement
	} else if nodeToReplace == nodeToReplace.parent.left {
		nodeToReplace.parent.left = replacement
	} else {
		nodeToReplace.parent.right = replacement
	}
	replacement.parent = nodeToReplace.parent
}

// minimum returns the node with minimum key in subtree rooted at the given node
func (t *Tree[T]) minimum(node *Node[T]) *Node[T] {
	for node.left != t.nil {
		node = node.left
	}
	return node
}

// updateSizeUpward recalculates sizes from node to root
func (t *Tree[T]) updateSizeUpward(node *Node[T]) {
	for node != t.nil {
		node.size = node.left.size + node.right.size + 1
		node = node.parent
	}
}

// deleteFixup maintains red-black tree properties after deletion
func (t *Tree[T]) deleteFixup(node *Node[T]) {
	for node != t.root && node.color == BLACK {
		if node == node.parent.left {
			sibling := node.parent.right
			if sibling.color == RED {
				// Case 1: Sibling is RED - rotate left and recolor
				sibling.color = BLACK
				node.parent.color = RED
				t.leftRotate(node.parent)
				sibling = node.parent.right
			}
			if sibling.left.color == BLACK && sibling.right.color == BLACK {
				// Case 2: Sibling's children are both BLACK - recolor sibling
				sibling.color = RED
				node = node.parent
			} else {
				if sibling.right.color == BLACK {
					// Case 3: Sibling's right child is BLACK - rotate right and recolor
					sibling.left.color = BLACK
					sibling.color = RED
					t.rightRotate(sibling)
					sibling = node.parent.right
				}
				// Case 4: Sibling's right child is RED - rotate left and recolor
				sibling.color = node.parent.color
				node.parent.color = BLACK
				sibling.right.color = BLACK
				t.leftRotate(node.parent)
				node = t.root
			}
		} else {
			// Mirror cases: node is right child
			sibling := node.parent.left
			if sibling.color == RED {
				// Case 1: Sibling is RED - rotate right and recolor
				sibling.color = BLACK
				node.parent.color = RED
				t.rightRotate(node.parent)
				sibling = node.parent.left
			}
			if sibling.right.color == BLACK && sibling.left.color == BLACK {
				// Case 2: Sibling's children are both BLACK - recolor sibling
				sibling.color = RED
				node = node.parent
			} else {
				if sibling.left.color == BLACK {
					// Case 3: Sibling's left child is BLACK - rotate left and recolor
					sibling.right.color = BLACK
					sibling.color = RED
					t.leftRotate(sibling)
					sibling = node.parent.left
				}
				// Case 4: Sibling's left child is RED - rotate right and recolor
				sibling.color = node.parent.color
				node.parent.color = BLACK
				sibling.left.color = BLACK
				t.rightRotate(node.parent)
				node = t.root
			}
		}
	}
	node.color = BLACK
}
