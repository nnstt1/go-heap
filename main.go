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

	l := h.List()
	for i, v := range l {
		fmt.Printf("[%d: %s]\n", i, v)
	}

}

func addHeap(value string, h *heap.Heap) {
	fmt.Printf("before h.Len(): %v\n", h.Len())
	h.Add(value)
	fmt.Printf("after  h.Len(): %v\n", h.Len())

	obj, exists, err := h.Get(value)
	if err != nil {
		fmt.Println("heap.Get error !")
		os.Exit(1)
	}
	if !exists {
		fmt.Printf("%s is not exists !\n", value)
		os.Exit(1)
	}
	fmt.Printf("value: %s\n", obj)
}
