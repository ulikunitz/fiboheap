// Package fiboheap is an alternative implementation of a priority queue
// using Fibonacci Heaps. While the ExtracMin method requires O(log n)
// complexity all others have complexity O(n).
//
// The items stored in the heap must satisfy the Sortable interface
// requiring a Less function. Nothing more is needed.
//
// Please check the IntSortable example for the usage.
package fiboheap

// Sortable is the interface that values stored in the heap must
// support.
type Sortable interface {
	Less(b Sortable) bool
}

// Heap provides an Fibonacci Heap. It can be used without special
// initialization.
type Heap struct {
	// forest head contains root nodes; first child node is minimum node
	forest node
	// items provides the number of items stored.
	items int
}

// node represents a tree node in the Fibonacci Heap.
type node struct {
	item Sortable
	// siblings
	prev, next *node
	// head for children nodes
	first, last *node
	// number of children
	children int
}

// less compares the items stores in the nodes. It uses the Less method
// of the Sortable interface. It panics if nil values are handled.
func less(x *node, y *node) bool {
	return x.item.Less(y.item)
}

// appendChildren transfers all children of s at the end of the children in
// r. The node s will not contain any children after return.
func (r *node) appendChildren(s *node) {
	if s.first == nil {
		return
	}
	if r.last == nil {
		r.last = s.last
		r.first = s.first
	} else {
		r.last.next = s.first
		s.first.prev = r.last
		r.last = s.last
	}
	r.children += s.children
	s.first, s.last = nil, nil
	s.children = 0
}

// removeChild removes a child from parent node r.
func (r *node) removeChild(c *node) {
	if c.prev == nil {
		if r.first != c {
			panic("c is not a child of r")
		}
		r.first = c.next
	} else {
		c.prev.next = c.next
	}
	if c.next == nil {
		if r.last != c {
			panic("c is not a child of r")
		}
		r.last = c.prev
	} else {
		c.next.prev = c.prev
	}
	r.children--
	c.prev, c.next = nil, nil
}

// appendChild puts a child at the end of the children list in the
// parent r. The child must not have been stored in another tree.
func (r *node) appendChild(c *node) {
	if c.next != nil || c.prev != nil {
		panic("c is already a child")
	}
	if r.last == nil {
		r.first, r.last = c, c
	} else {
		c.prev = r.last
		r.last.next = c
		r.last = c
	}
	r.children++
}

// insertAtFront puts the child at the beginning of the list of parent.
// The child must not have been stored in a tree before.
func (r *node) insertAtFront(c *node) {
	if c.next != nil || c.prev != nil {
		panic("c is already a child")
	}
	if r.first == nil {
		r.first, r.last = c, c
	} else {
		c.next = r.first
		r.first.prev = c
		r.first = c
	}
	r.children++
}

// The rootSlice is used to store nodes with a specific index. It grows
// automatically if a node is stored at an index that doesn't fit the
// slice.
type rootSlice []*node

// Rerturns the node stored at index i. If the index is larger than the
// slice a nil value is returned.
func (r rootSlice) get(i int) *node {
	if i >= len(r) {
		return nil
	}
	return r[i]
}

// put stores node x at index i. If the index is larger than the current
// slice size sufficient space is obtained. If the required size is
// smaller than 32 the slice will get a size of 32 entries.
func (r *rootSlice) put(i int, x *node) {
	t := *r
	if i >= len(t) {
		if x == nil {
			return
		}
		c := i + 1
		if c < 32 {
			c = 32
		}
		t = make(rootSlice, c)
		copy(t, *r)
		*r = t
	}
	t[i] = x
}

// combine puts x and y in one tree. The tree with the smaller item will
// be the new root. The new root is returned.
func (r *node) combine(x, y *node) *node {
	// x and y are children of r
	if less(y, x) {
		x, y = y, x
	}
	r.removeChild(y)
	x.appendChild(y)
	return x
}

// restructureChildren ensures that are only logN entries in the child
// list, by combining nodes with the same number of childrens.
func (p *node) restructureChildren() {
	var a rootSlice
	x := p.first
	for x != nil {
		r := x
		x = x.next
		for {
			s := a.get(r.children)
			if s == nil {
				break
			}
			a.put(r.children, nil)
			r = p.combine(r, s)
		}
		a.put(r.children, r)
	}
}

// findMinChild finds the child with the minimum item of the children of
// the parent node.
func (r *node) findMinChild() *node {
	if r.first == nil {
		return nil
	}
	min := r.first
	for c := min.next; c != nil; c = c.next {
		if less(c, min) {
			min = c
		}
	}
	return min
}

// Len returns the number of Sortable items stored in the heap.
func (h *Heap) Len() int { return h.items }

// FindMin returns the minimum element. The running time is O(1).
func (h *Heap) FindMin() Sortable {
	if h.forest.first == nil {
		return nil
	}
	return h.forest.first.item
}

// Extract the minimum element. It requires O(log n) time.
func (h *Heap) ExtractMin() Sortable {
	if h.forest.first == nil {
		return nil
	}

	// Remove minimum item from roots and add the children to the
	// roots.
	min := h.forest.first
	h.forest.removeChild(min)
	h.forest.appendChildren(min)

	// Ensure that there are only log n roots in the forest.
	h.forest.restructureChildren()

	// Make new minimum first children of the forest.
	newMin := h.forest.findMinChild()
	if newMin != nil {
		h.forest.removeChild(newMin)
		h.forest.insertAtFront(newMin)
	}

	h.items--
	return min.item
}

// Insert puts a Sortable item into the heap. The running time for the
// operation is O(1).
func (h *Heap) Insert(x Sortable) {
	s := &node{item: x}
	if h.forest.first == nil || less(s, h.forest.first) {
		h.forest.insertAtFront(s)
	} else {
		h.forest.appendChild(s)
	}
	h.items++
}

// Merge adds the nodes from heap g to heap h. Heap g will be emptied.
// The running time is O(1).
func (h *Heap) Merge(g *Heap) {
	gmin := g.forest.first
	h.forest.appendChildren(&g.forest)
	h.items += g.items
	hmin := h.forest.first
	if gmin != nil && hmin != gmin && less(gmin, hmin) {
		h.forest.removeChild(gmin)
		h.forest.insertAtFront(gmin)
	}
	*g = Heap{}
}
