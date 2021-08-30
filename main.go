package main

import (
	"fmt"

	"github.com/nnstt1/go-heap/heap"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func main() {
	keyFn := func(obj interface{}) (string, error) {
		return cache.MetaNamespaceKeyFunc(obj.(*framework.QueuedPodInfo).Pod)
	}

	lessFn := func(podInfo1, podInfo2 interface{}) bool {
		return false
	}

	h := heap.New(keyFn, lessFn)
	fmt.Printf("h.Len(): %v\n", h.Len())

}
