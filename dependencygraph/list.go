package dependencygraph

// Next returns the next list Element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list Element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func NewList() *List { return new(List).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List) Len() int { return l.len }

// Front returns the first Element of list l or nil.
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last Element of list l or nil.
func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *List) insert(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *List) insertValue(v *Node, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *List) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Remove removes e from l if e is an Element of list l.
// It returns the Element value e.Value.
func (l *List) Remove(e *Element) *Node {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new Element e with value v at the front of list l and returns e.
func (l *List) PushFront(v *Node) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushFrontUniq inserts a new Element e with value v at the front of list l and returns e.
func (l *List) PushFrontUniq(v *Node) *Element {
	l.lazyInit()
	e := l.Find(v)
	if e == nil {
		e = &Element{Value: v}
	} else {
		l.remove(e)
	}
	l.insert(e, &l.root)
	return e
}

// PushBack inserts a new Element e with value v at the back of list l and returns e.
func (l *List) PushBack(v *Node) *Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// PushBackUniq inserts a new Element e with value v at the back of list l and returns e.
func (l *List) PushBackUniq(v *Node) *Element {
	l.lazyInit()
	l.lazyInit()
	e := l.Find(v)
	if e == nil {
		e = &Element{Value: v}
	} else {
		l.remove(e)
	}
	l.insert(e, l.root.prev)
	return e
}

// InsertBefore inserts a new Element e with value v immediately before mark and returns e.
// If mark is not an Element of l, the list is not modified.
func (l *List) InsertBefore(v *Node, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new Element e with value v immediately after mark and returns e.
// If mark is not an Element of l, the list is not modified.
func (l *List) InsertAfter(v *Node, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// Find return the first Element nil if v is not found
func (l *List) Find(v *Node) (e *Element) {
	if l.len == 0 {
		return nil
	}
	for e = l.Front(); e != nil && e.Value != v; e = e.Next() {
	}
	return
}

// MoveToFront moves Element e to the front of list l.
// If e is not an Element of l, the list is not modified.
func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.insert(l.remove(e), &l.root)
}

// MoveToBack moves Element e to the back of list l.
// If e is not an Element of l, the list is not modified.
func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.insert(l.remove(e), l.root.prev)
}

// MoveBefore moves Element e to its new position before mark.
// If e or mark is not an Element of l, or e == mark, the list is not modified.
func (l *List) MoveBefore(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.insert(l.remove(e), mark.prev)
}

// MoveAfter moves Element e to its new position after mark.
// If e or mark is not an Element of l, or e == mark, the list is not modified.
func (l *List) MoveAfter(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.insert(l.remove(e), mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.
func (l *List) PushBackList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.
func (l *List) PushFrontList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

func (l *List) String() (ret string) {
	for e := l.Front(); e != nil; e = e.Next() {
		ret += "\n" + e.Value.path
	}
	return
}
