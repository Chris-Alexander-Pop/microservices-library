package avl

import (
	"sync"

	"golang.org/x/exp/constraints"
)

// Tree is an AVL Tree (self-balancing BST).
type Tree[K constraints.Ordered, V any] struct {
	Root *Node[K, V]
	mu   sync.RWMutex
}

type Node[K constraints.Ordered, V any] struct {
	Key    K
	Value  V
	Height int
	Left   *Node[K, V]
	Right  *Node[K, V]
}

func New[K constraints.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{}
}

func (t *Tree[K, V]) Put(Key K, Value V) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Root = insert(t.Root, Key, Value)
}

func (t *Tree[K, V]) Get(Key K) (V, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	n := search(t.Root, Key)
	if n == nil {
		var zero V
		return zero, false
	}
	return n.Value, true
}

func Height[K constraints.Ordered, V any](n *Node[K, V]) int {
	if n == nil {
		return 0
	}
	return n.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func updateHeight[K constraints.Ordered, V any](n *Node[K, V]) {
	n.Height = 1 + max(Height(n.Left), Height(n.Right))
}

func getBalance[K constraints.Ordered, V any](n *Node[K, V]) int {
	if n == nil {
		return 0
	}
	return Height(n.Left) - Height(n.Right)
}

func rightRotate[K constraints.Ordered, V any](y *Node[K, V]) *Node[K, V] {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	updateHeight(y)
	updateHeight(x)

	return x
}

func leftRotate[K constraints.Ordered, V any](x *Node[K, V]) *Node[K, V] {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	updateHeight(x)
	updateHeight(y)

	return y
}

func insert[K constraints.Ordered, V any](n *Node[K, V], Key K, Value V) *Node[K, V] {
	if n == nil {
		return &Node[K, V]{Key: Key, Value: Value, Height: 1}
	}

	if Key < n.Key {
		n.Left = insert(n.Left, Key, Value)
	} else if Key > n.Key {
		n.Right = insert(n.Right, Key, Value)
	} else {
		n.Value = Value // Update Value
		return n
	}

	updateHeight(n)
	balance := getBalance(n)

	// Left Left
	if balance > 1 && Key < n.Left.Key {
		return rightRotate(n)
	}
	// Right Right
	if balance < -1 && Key > n.Right.Key {
		return leftRotate(n)
	}
	// Left Right
	if balance > 1 && Key > n.Left.Key {
		n.Left = leftRotate(n.Left)
		return rightRotate(n)
	}
	// Right Left
	if balance < -1 && Key < n.Right.Key {
		n.Right = rightRotate(n.Right)
		return leftRotate(n)
	}

	return n
}

func search[K constraints.Ordered, V any](n *Node[K, V], Key K) *Node[K, V] {
	if n == nil {
		return nil
	}
	if Key < n.Key {
		return search(n.Left, Key)
	} else if Key > n.Key {
		return search(n.Right, Key)
	}
	return n
}
