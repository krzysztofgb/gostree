package gostree

type Color bool

const (
	RED   Color = false
	BLACK Color = true
)

// CompareFunc defines the comparison function for ordering elements.
// It should return:
//   - negative value if a < b
//   - zero if a == b
//   - positive value if a > b
type CompareFunc[T any] func(a, b T) int

type Node[T any] struct {
	key    T
	left   *Node[T]
	right  *Node[T]
	parent *Node[T]
	color  Color
	size   int // number of nodes in subtree rooted at this node
}

type Tree[T any] struct {
	root    *Node[T]
	nil     *Node[T] // sentinel node
	compare CompareFunc[T]
}

// getGrandparent returns the grandparent of the node
func (t *Tree[T]) getGrandparent(n *Node[T]) *Node[T] {
	if n.parent != t.nil {
		return n.parent.parent
	}

	return t.nil
}

// getSibling returns the sibling of the node
func (t *Tree[T]) getSibling(n *Node[T]) *Node[T] {
	if n.parent == t.nil {
		return t.nil
	}
	if n == n.parent.left {
		return n.parent.right
	}

	return n.parent.left
}

// getUncle returns the uncle (parent's sibling) of the node
func (t *Tree[T]) getUncle(n *Node[T]) *Node[T] {
	grandparent := t.getGrandparent(n)
	if grandparent == t.nil {
		return t.nil
	}
	if n.parent.isLeftChild() {
		return grandparent.right
	}

	return grandparent.left
}

// isLeftChild returns true if the node is a left child
func (n *Node[T]) isLeftChild() bool {
	return n.parent != nil && n == n.parent.left
}

// isRightChild returns true if the node is a right child
func (n *Node[T]) isRightChild() bool {
	return n.parent != nil && n == n.parent.right
}

// NewTree creates a new order-statistic tree.
func NewTree[T any](compare CompareFunc[T]) *Tree[T] {
	t := &Tree[T]{
		root:    nil,
		compare: compare,
		nil: &Node[T]{ // sentinel node
			key:    *new(T),
			left:   nil,
			right:  nil,
			parent: nil,
			color:  BLACK,
			size:   0,
		},
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
		if t.compare(key, current.key) < 0 {
			current = current.left
		} else {
			current = current.right
		}
	}

	// Insert the new node
	newNode.parent = parent
	if parent == t.nil {
		t.root = newNode
	} else if t.compare(newNode.key, parent.key) < 0 {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	// Fix red-black properties
	t.insertFixup(newNode)
}

// insertFixup maintains red-black tree properties after insertion
//
// The function handles violations where a RED node has a RED parent.
// There are 3 main cases (mirrored for left/right):
//
// Case 1: Uncle is RED
//
//	     G(B)                    G(R)
//	    /   \                  /   \
//	  P(R)   U(R)    =>      P(B)   U(B)
//	  /                      /
//	N(R)                   N(R)
//
// Case 2: Node is inner child (right child of left parent)
//
//	   G(B)                    G(B)
//	  /   \                  /   \
//	P(R)   U(B)    =>      N(R)   U(B)
//	  \                    /
//	   N(R)              P(R)
//
// Case 3: Node is outer child (left child of left parent)
//
//	     G(B)                    P(B)
//	    /   \                  /   \
//	  P(R)   U(B)    =>      N(R)   G(R)
//	  /                              \
//	N(R)                              U(B)
//
// Legend: G=Grandparent, P=Parent, N=NewNode, U=Uncle, (R)=RED, (B)=BLACK
func (t *Tree[T]) insertFixup(newNode *Node[T]) {
	for newNode.parent.color == RED {
		parent := newNode.parent
		grandparent := t.getGrandparent(newNode)

		if parent.isLeftChild() {
			uncle := t.getUncle(newNode)
			if uncle.color == RED {
				// Case 1: Uncle is RED - recolor and move up
				//     G(B)                G(R)
				//    /   \              /   \
				//  P(R)   U(R)  =>    P(B)   U(B)
				//  /                  /
				// N(R)               N(R)
				parent.color = BLACK
				uncle.color = BLACK
				grandparent.color = RED
				newNode = grandparent
			} else {
				if newNode.isRightChild() {
					// Case 2: Node is right child - rotate left to convert to case 3
					//     G(B)              G(B)
					//    /   \            /   \
					//  P(R)   U(B)  =>  N(R)   U(B)
					//    \              /
					//     N(R)        P(R)
					newNode = parent
					t.leftRotate(newNode)
				}
				// Case 3: Node is left child - rotate right and recolor
				//     G(B)              P(B)
				//    /   \            /   \
				//  P(R)   U(B)  =>  N(R)   G(R)
				//  /                        \
				// N(R)                      U(B)
				newNode.parent.color = BLACK
				grandparent.color = RED
				t.rightRotate(grandparent)
			}
		} else {
			// Mirror cases: parent is right child of grandparent
			uncle := t.getUncle(newNode)
			if uncle.color == RED {
				// Case 1: Uncle is RED - recolor and move up
				//     G(B)                G(R)
				//    /   \              /   \
				//  U(R)   P(R)  =>    U(B)   P(B)
				//           \                   \
				//            N(R)                N(R)
				parent.color = BLACK
				uncle.color = BLACK
				grandparent.color = RED
				newNode = grandparent
			} else {
				if newNode.isLeftChild() {
					// Case 2: Node is left child - rotate right to convert to case 3
					//     G(B)              G(B)
					//    /   \            /   \
					//  U(B)   P(R)  =>  U(B)   N(R)
					//         /                   \
					//       N(R)                  P(R)
					newNode = parent
					t.rightRotate(newNode)
				}
				// Case 3: Node is right child - rotate left and recolor
				//     G(B)              P(B)
				//    /   \            /   \
				//  U(B)   P(R)  =>  G(R)   N(R)
				//           \       /
				//            N(R)  U(B)
				newNode.parent.color = BLACK
				grandparent.color = RED
				t.leftRotate(grandparent)
			}
		}
	}
	t.root.color = BLACK
}

// leftRotate performs a left rotation on the given node
//
// Before:         After:
//
//	  x              y
//	 / \            / \
//	a   y    =>    x   c
//	   / \        / \
//	  b   c      a   b
//
// Where x = node, y = rightChild
// Parent relationships are updated accordingly
func (t *Tree[T]) leftRotate(node *Node[T]) {
	rightChild := node.right
	node.right = rightChild.left
	if rightChild.left != t.nil {
		rightChild.left.parent = node
	}
	rightChild.parent = node.parent
	if node.parent == t.nil {
		t.root = rightChild
	} else if node.isLeftChild() {
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
//
// Before:         After:
//
//	    y            x
//	   / \          / \
//	  x   c   =>   a   y
//	 / \              / \
//	a   b            b   c
//
// Where y = node, x = leftChild
// Parent relationships are updated accordingly
func (t *Tree[T]) rightRotate(node *Node[T]) {
	leftChild := node.left
	node.left = leftChild.right
	if leftChild.right != t.nil {
		leftChild.right.parent = node
	}
	leftChild.parent = node.parent
	if node.parent == t.nil {
		t.root = leftChild
	} else if node.isRightChild() {
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
	for current != t.nil {
		cmp := t.compare(key, current.key)
		if cmp == 0 {
			break
		} else if cmp < 0 {
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
// If there are duplicates of the key, it returns the rank of the leftmost occurrence.
func (t *Tree[T]) Rank(key T) int {
	rank := 0
	current := t.root

	for current != t.nil {
		if t.compare(key, current.key) <= 0 {
			// Key is less than or equal, go left
			current = current.left
		} else {
			// Key is greater, count this node and its left subtree
			rank += current.left.size + 1
			current = current.right
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
//
// Before:              After:
//
//	  P                   P
//	  |                   |
//	  U        =>         V
//	 / \                 / \
//	a   b              (V's subtree)
//
// Where P = parent of U, U = nodeToReplace, V = replacement
// This operation updates parent pointers but preserves V's children
func (t *Tree[T]) transplant(nodeToReplace, replacement *Node[T]) {
	if nodeToReplace.parent == t.nil {
		t.root = replacement
	} else if nodeToReplace.isLeftChild() {
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
//
// This function fixes violations when deleting a BLACK node.
// There are 4 main cases (mirrored for left/right):
//
// Case 1: Sibling is RED
//
//	   P(B)                    S(B)
//	  /   \                  /   \
//	N(B)   S(R)     =>     P(R)   SR(B)
//	      /   \           /   \
//	    SL(B) SR(B)     N(B)  SL(B)
//
// Case 2: Sibling and its children are BLACK
//
//	   P(?)                    P(?)
//	  /   \                  /   \
//	N(B)   S(B)     =>     N(B)   S(R)
//	      /   \                  /   \
//	    SL(B) SR(B)            SL(B) SR(B)
//
// Case 3: Sibling's far child is BLACK
//
//	   P(?)                    P(?)
//	  /   \                  /   \
//	N(B)   S(B)     =>     N(B)  SL(B)
//	      /   \                     \
//	    SL(R) SR(B)                  S(R)
//	                                   \
//	                                  SR(B)
//
// Case 4: Sibling's far child is RED
//
//	   P(?)                    S(?)
//	  /   \                  /   \
//	N(B)   S(B)     =>     P(B)   SR(B)
//	      /   \           /   \
//	    SL(?) SR(R)     N(B)  SL(?)
//
// Legend: P=Parent, N=Node, S=Sibling, SL=Sibling's Left, SR=Sibling's Right, (R)=RED, (B)=BLACK, (?)=Either color
func (t *Tree[T]) deleteFixup(node *Node[T]) {
	for node != t.root && node.color == BLACK {
		if node.isLeftChild() {
			sibling := t.getSibling(node)
			if sibling.color == RED {
				// Case 1: Sibling is RED - rotate left and recolor
				//    P(B)              S(B)
				//   /   \            /   \
				// N(B)   S(R)  =>  P(R)   SR(B)
				//       /   \      /   \
				//     SL(B) SR(B) N(B) SL(B)
				sibling.color = BLACK
				node.parent.color = RED
				t.leftRotate(node.parent)
				sibling = t.getSibling(node)
			}
			if sibling.left.color == BLACK && sibling.right.color == BLACK {
				// Case 2: Sibling's children are both BLACK - recolor sibling
				//    P(?)              P(?)
				//   /   \            /   \
				// N(B)   S(B)  =>  N(B)   S(R)
				//       /   \            /   \
				//     SL(B) SR(B)      SL(B) SR(B)
				sibling.color = RED
				node = node.parent
			} else {
				if sibling.right.color == BLACK {
					// Case 3: Sibling's right child is BLACK - rotate right and recolor
					//    P(?)              P(?)
					//   /   \            /   \
					// N(B)   S(B)  =>  N(B)  SL(B)
					//       /   \               \
					//     SL(R) SR(B)            S(R)
					//                              \
					//                             SR(B)
					sibling.left.color = BLACK
					sibling.color = RED
					t.rightRotate(sibling)
					sibling = t.getSibling(node)
				}
				// Case 4: Sibling's right child is RED - rotate left and recolor
				//    P(?)              S(?)
				//   /   \            /   \
				// N(B)   S(B)  =>  P(B)   SR(B)
				//       /   \      /   \
				//     SL(?) SR(R) N(B) SL(?)
				sibling.color = node.parent.color
				node.parent.color = BLACK
				sibling.right.color = BLACK
				t.leftRotate(node.parent)
				node = t.root
			}
		} else {
			// Mirror cases: node is right child
			sibling := t.getSibling(node)
			if sibling.color == RED {
				// Case 1: Sibling is RED - rotate right and recolor
				//       P(B)              S(B)
				//      /   \            /   \
				//    S(R)   N(B)  =>  SL(B)  P(R)
				//   /   \                   /   \
				// SL(B) SR(B)             SR(B) N(B)
				sibling.color = BLACK
				node.parent.color = RED
				t.rightRotate(node.parent)
				sibling = t.getSibling(node)
			}
			if sibling.right.color == BLACK && sibling.left.color == BLACK {
				// Case 2: Sibling's children are both BLACK - recolor sibling
				//       P(?)              P(?)
				//      /   \            /   \
				//    S(B)   N(B)  =>  S(R)   N(B)
				//   /   \            /   \
				// SL(B) SR(B)      SL(B) SR(B)
				sibling.color = RED
				node = node.parent
			} else {
				if sibling.left.color == BLACK {
					// Case 3: Sibling's left child is BLACK - rotate left and recolor
					//       P(?)              P(?)
					//      /   \            /   \
					//    S(B)   N(B)  =>  SR(B)  N(B)
					//   /   \            /
					// SL(B) SR(R)       S(R)
					//                  /
					//                SL(B)
					sibling.right.color = BLACK
					sibling.color = RED
					t.leftRotate(sibling)
					sibling = t.getSibling(node)
				}
				// Case 4: Sibling's left child is RED - rotate right and recolor
				//       P(?)              S(?)
				//      /   \            /   \
				//    S(B)   N(B)  =>  SL(B)  P(B)
				//   /   \                   /   \
				// SL(R) SR(?)             SR(?) N(B)
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

// Size returns the number of elements in the tree.
func (t *Tree[T]) Size() int {
	return t.root.size
}
