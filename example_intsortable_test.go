package fiboheap

import "fmt"

type IntSortable int

func (s IntSortable) Less(r Sortable) bool {
	t := r.(IntSortable)
	return s < t
}

func Example_intSortable() {
	var h Heap
	for _, k := range []int{2, 1, 5, 3} {
		h.Insert(IntSortable(k))
	}
	fmt.Printf("minimum: %d\n", h.FindMin().(IntSortable))
	for h.Len() > 0 {
		fmt.Printf("%d ", h.ExtractMin().(IntSortable))
	}
	// Output:
	// minimum: 1
	// 1 2 3 5
}
