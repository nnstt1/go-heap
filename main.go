package main

import (
	"fmt"
	"os"

	"github.com/nnstt1/go-heap/heap"
)

func main() {
	keyFn := func(obj interface{}) (string, error) {
		return "key-" + obj.(string), nil
	}

	lessFn := func(podInfo1, podInfo2 interface{}) bool {
		return false
	}

	h := heap.New(keyFn, lessFn)
	addHeap("hoge", h)
	addHeap("fuga", h)
	addHeap("piyo", h)

	updateHeap("foo", h)

	peakHeap(h)
	peakHeap(h)
	peakHeap(h)

	deleteHeap("hoge", h)

	popHeap(h)
	popHeap(h)
	popHeap(h)

}

func addHeap(value string, h *heap.Heap) {
	fmt.Printf("Add(%v)\n", value)
	h.Add(value)

	_, exists, err := h.Get(value)
	if err != nil {
		fmt.Println("heap.Get error !")
		os.Exit(1)
	}
	if !exists {
		fmt.Printf("%s is not exists !\n", value)
		os.Exit(1)
	}
	list(h)
}

func updateHeap(value string, h *heap.Heap) {
	fmt.Printf("Update(%v)\n", value)
	h.Update(value)
	list(h)
}

func deleteHeap(value string, h *heap.Heap) {
	fmt.Printf("Delete(%v)\n", value)
	h.Delete(value)
	list(h)
}

func popHeap(h *heap.Heap) {
	obj, _ := h.Pop()
	fmt.Printf("Pop() => %s\n", obj.(string))
	list(h)
}

func peakHeap(h *heap.Heap) {
	obj := h.Peek()
	if obj == nil {
		fmt.Println("Peak() => not found")
	}
	fmt.Printf("Peak() => %s\n", obj.(string))
	list(h)
}

func list(h *heap.Heap) {
	l := h.List()
	for i, v := range l {
		fmt.Printf("[%d: %s]\n", i, v)
	}
}
