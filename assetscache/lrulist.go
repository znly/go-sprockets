package assetscache

// assetLruListElement is an assetLruListElement of a linked list.
type assetLruListElement struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &all.root is both the next assetLruListElement of the last
	// list assetLruListElement (all.Back()) and the previous assetLruListElement of the first list
	// assetLruListElement (all.Front()).
	next, prev *assetLruListElement

	// The list to which this assetLruListElement belongs.
	list *assetLruList

	// The value stored with this assetLruListElement.
	Value *assetLruListEntry
}

// assetLruList represents a doubly linked list.
// The zero value for assetLruList is an empty list ready to use.
type assetLruList struct {
	root assetLruListElement // sentinel list assetLruListElement, only &root, root.prev, and root.next are used
	len  int                 // current list length excluding (this) sentinel assetLruListElement
}

// Next returns the next list assetLruListElement or nil.
func (alle *assetLruListElement) Next() *assetLruListElement {
	if p := alle.next; alle.list != nil && p != &alle.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list assetLruListElement or nil.
func (alle *assetLruListElement) Prev() *assetLruListElement {
	if p := alle.prev; alle.list != nil && p != &alle.list.root {
		return p
	}
	return nil
}

// Init initializes or clears list l.
func (all *assetLruList) Init() *assetLruList {
	all.root.next = &all.root
	all.root.prev = &all.root
	all.len = 0
	return all
}

// New returns an initialized list.
func newAssetLruList() *assetLruList { return new(assetLruList).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (all *assetLruList) Len() int {
	if all == nil {
		return 0
	}
	return all.len
}

// Front returns the first assetLruListElement of list l or nil.
func (all *assetLruList) Front() *assetLruListElement {
	if all.len == 0 {
		return nil
	}
	return all.root.next
}

// Back returns the last assetLruListElement of list l or nil.
func (all *assetLruList) Back() *assetLruListElement {
	if all.len == 0 {
		return nil
	}
	return all.root.prev
}

// lazyInit lazily initializes a zero assetLruList value.
func (all *assetLruList) lazyInit() {
	if all.root.next == nil {
		all.Init()
	}
}

// insert inserts e after at, increments all.len, and returns e.
func (all *assetLruList) insert(alle, at *assetLruListElement) *assetLruListElement {
	n := at.next
	at.next = alle
	alle.prev = at
	alle.next = n
	n.prev = alle
	alle.list = all
	all.len++
	return alle
}

// insertValue is a convenience wrapper for insert(&assetLruListElement{Value: v}, at).
func (all *assetLruList) insertValue(v *assetLruListEntry, at *assetLruListElement) *assetLruListElement {
	return all.insert(&assetLruListElement{Value: v}, at)
}

// remove removes e from its list, decrements all.len, and returns e.
func (all *assetLruList) remove(alle *assetLruListElement) *assetLruListElement {
	alle.prev.next = alle.next
	alle.next.prev = alle.prev
	alle.next = nil // avoid memory leaks
	alle.prev = nil // avoid memory leaks
	alle.list = nil
	all.len--
	return alle
}

// Remove removes e from l if e is an assetLruListElement of list l.
// It returns the assetLruListElement value alle.Value.
func (all *assetLruList) Remove(alle *assetLruListElement) *assetLruListEntry {
	if alle.list == all {
		// if alle.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero assetLruListElement) and all.remove will crash
		all.remove(alle)
	}
	return alle.Value
}

// PushFront inserts a new assetLruListElement e with value v at the front of list l and returns e.
func (all *assetLruList) PushFront(v *assetLruListEntry) *assetLruListElement {
	all.lazyInit()
	return all.insertValue(v, &all.root)
}

// PushFrontUniq inserts a new assetLruListElement e with value v at the front of list l and returns e.
func (all *assetLruList) PushFrontUniq(v *assetLruListEntry) *assetLruListElement {
	all.lazyInit()
	alle := all.Find(v)
	if alle == nil {
		alle = &assetLruListElement{Value: v}
	} else {
		all.remove(alle)
	}
	all.insert(alle, &all.root)
	return alle
}

// PushBack inserts a new assetLruListElement e with value v at the back of list l and returns e.
func (all *assetLruList) PushBack(v *assetLruListEntry) *assetLruListElement {
	all.lazyInit()
	return all.insertValue(v, all.root.prev)
}

// PushBackUniq inserts a new assetLruListElement e with value v at the back of list l and returns e.
func (all *assetLruList) PushBackUniq(v *assetLruListEntry) *assetLruListElement {
	all.lazyInit()
	all.lazyInit()
	e := all.Find(v)
	if e == nil {
		e = &assetLruListElement{Value: v}
	} else {
		all.remove(e)
	}
	all.insert(e, all.root.prev)
	return e
}

// InsertBefore inserts a new assetLruListElement e with value v immediately before mark and returns e.
// If mark is not an assetLruListElement of l, the list is not modified.
func (all *assetLruList) InsertBefore(v *assetLruListEntry, mark *assetLruListElement) *assetLruListElement {
	if mark.list != all {
		return nil
	}
	// see comment in assetLruList.Remove about initialization of l
	return all.insertValue(v, mark.prev)
}

// InsertAfter inserts a new assetLruListElement e with value v immediately after mark and returns e.
// If mark is not an assetLruListElement of l, the list is not modified.
func (all *assetLruList) InsertAfter(v *assetLruListEntry, mark *assetLruListElement) *assetLruListElement {
	if mark.list != all {
		return nil
	}
	// see comment in assetLruList.Remove about initialization of l
	return all.insertValue(v, mark)
}

// Find return the first assetLruListElement nil if v is not found
func (all *assetLruList) Find(v *assetLruListEntry) (alle *assetLruListElement) {
	if all.len == 0 {
		return nil
	}
	for alle = all.Front(); alle != nil && alle.Value != v; alle = alle.Next() {
	}
	return
}

// MoveToFront moves assetLruListElement e to the front of list l.
// If e is not an assetLruListElement of l, the list is not modified.
func (all *assetLruList) MoveToFront(alle *assetLruListElement) {
	if alle.list != all || all.root.next == alle {
		return
	}
	// see comment in assetLruList.Remove about initialization of l
	all.insert(all.remove(alle), &all.root)
}

// MoveToBack moves assetLruListElement e to the back of list l.
// If e is not an assetLruListElement of l, the list is not modified.
func (all *assetLruList) MoveToBack(alle *assetLruListElement) {
	if alle.list != all || all.root.prev == alle {
		return
	}
	// see comment in assetLruList.Remove about initialization of l
	all.insert(all.remove(alle), all.root.prev)
}

// MoveBefore moves assetLruListElement e to its new position before mark.
// If e or mark is not an assetLruListElement of l, or e == mark, the list is not modified.
func (all *assetLruList) MoveBefore(alle, mark *assetLruListElement) {
	if alle.list != all || alle == mark || mark.list != all {
		return
	}
	all.insert(all.remove(alle), mark.prev)
}

// MoveAfter moves assetLruListElement e to its new position after mark.
// If e or mark is not an assetLruListElement of l, or e == mark, the list is not modified.
func (all *assetLruList) MoveAfter(alle, mark *assetLruListElement) {
	if alle.list != all || alle == mark || mark.list != all {
		return
	}
	all.insert(all.remove(alle), mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.
func (all *assetLruList) PushBackList(other *assetLruList) {
	all.lazyInit()
	for i, alle := other.Len(), other.Front(); i > 0; i, alle = i-1, alle.Next() {
		all.insertValue(alle.Value, all.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.
func (all *assetLruList) PushFrontList(other *assetLruList) {
	all.lazyInit()
	for i, alle := other.Len(), other.Back(); i > 0; i, alle = i-1, alle.Prev() {
		all.insertValue(alle.Value, &all.root)
	}
}
