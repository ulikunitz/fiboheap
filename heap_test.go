package fiboheap

import (
	"math/rand"
	"testing"
	"unsafe"
)

func TestHeap(t *testing.T) {
	var h Heap
	rand.Seed(13)
	const (
		batch1  = 1000
		extract = 523
		batch2  = 1234
	)
	for i := 0; i < batch1; i++ {
		n := rand.Intn(10000)
		h.Insert(IntSortable(n))
	}
	if h.Len() != batch1 {
		t.Fatalf("#1 h.Len() got %d; want %d", h.Len(), batch1)
	}
	m := int(h.ExtractMin().(IntSortable))
	for i := 1; i < extract; i++ {
		k := int(h.ExtractMin().(IntSortable))
		if k < m {
			t.Fatalf("extracted %d < %d", k, m)
		}
		m = k
	}
	if h.Len() != batch1-extract {
		t.Fatalf("#2 h.Len() got %d; want %d", h.Len(), batch1-extract)
	}
	for i := 0; i < batch2; i++ {
		n := rand.Intn(10000)
		h.Insert(IntSortable(n))
	}
	if h.Len() != batch1-extract+batch2 {
		t.Fatalf("#3 h.Len() got %d; want %d", h.Len(), batch1-extract)
	}
	m = int(h.ExtractMin().(IntSortable))
	for h.Len() > 0 {
		k := int(h.ExtractMin().(IntSortable))
		if k < m {
			t.Fatalf("extracted %d < %d", k, m)
		}
		m = k
	}
}

func TestNodeSize(t *testing.T) {
	var x node
	t.Logf("sizeof(node): %d bytes", unsafe.Sizeof(x))
}
