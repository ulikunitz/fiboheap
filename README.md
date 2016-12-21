# fiboheap

The package provides a Go implementation of the Fibonacci heaps. It
requires the items to be stored in the heap to support a single Less
method. Beside extracing the minimum, which is a O(log n) operation, all
methods have time complexity O(1).

The standard library heap implementation is based on binary trees
requires a container type implementation supporting five methods. The
storing and extraction operation requires O(log n). It has however lower 
memory overhead than the Fibonacci Heap. This implementation requires an
internal node structure of 56 bytes per item.

## Install

```
go get -u github.com/ulikunitz/fiboheap
```

## Example

```go
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
```
