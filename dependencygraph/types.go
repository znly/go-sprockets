package dependencygraph

// Element is an Element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next Element of the last
	// list Element (l.Back()) and the previous Element of the first list
	// Element (l.Front()).
	next, prev *Element

	// The list to which this Element belongs.
	list *List

	// The value stored with this Element.
	Value *Node
}

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list Element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel Element
}

type Node struct {
	path string
	edge *List
}

type Graph struct {
	nodeList map[string]*Node
}
