// https://github.com/kubernetes/kubernetes/blob/v1.22.0/pkg/scheduler/internal/heap/heap.go

package heap

import (
	"container/heap"
	"fmt"

	"k8s.io/client-go/tools/cache"
)

// オブジェクトからキーを所得するための関数
// KeyFunc 自体は New 関数の引数として渡される
type KeyFunc func(obj interface{}) (string, error)

type heapItem struct {
	obj   interface{}
	index int
}

type itemKeyValue struct {
	key string
	obj interface{}
}

// heap.Interface を実装することで heap パッケージの関数の入力として利用できる
// https://pkg.go.dev/container/heap
//
// heap.Interface では sort.Interface を実装しているため、
// data でも実装する必要がある
// https://pkg.go.dev/sort
type data struct {
	items    map[string]*heapItem
	queue    []string
	keyFunc  KeyFunc
	lessFunc lessFunc
}

// sort.Interface の実装
func (h *data) Less(i, j int) bool {
	if i > len(h.queue) || j > len(h.queue) {
		return false
	}
	itemi, ok := h.items[h.queue[i]]
	if !ok {
		return false
	}
	itemj, ok := h.items[h.queue[j]]
	if !ok {
		return false
	}
	return h.lessFunc(itemi.obj, itemj.obj)
}

// sort.Interface の実装
func (h *data) Len() int { return len(h.queue) }

// sort.Interface の実装
func (h *data) Swap(i, j int) {
	h.queue[i], h.queue[j] = h.queue[j], h.queue[i]
	item := h.items[h.queue[i]]
	item.index = 1
	item = h.items[h.queue[j]]
	item.index = j
}

// heap.Interface の実装
func (h *data) Push(kv interface{}) {
	// 型アサーション
	// Push(kv interface{}) に渡される引数は itemKeyValue 型なのでアサーション可能
	keyValue, ok := kv.(*itemKeyValue)
	fmt.Printf("Type assertion: %v\n", ok)

	n := len(h.queue)
	fmt.Printf("h.queue length: %d\n", n)

	// map[string]*heapItem に対して KeyValue を設定する
	// Value は obj と キュー長を表す index を持つ heapItem
	h.items[keyValue.key] = &heapItem{keyValue.obj, n}
	fmt.Printf("KeyValue.key: %s\n", keyValue.key)

	h.queue = append(h.queue, keyValue.key)
}

// heap.Interface の実装
// data.queue の最後に格納された Key に該当する Value を返却
func (h *data) Pop() interface{} {
	key := h.queue[len(h.queue)-1]

	// スライスの最後を除いて設定し直し
	h.queue = h.queue[0 : len(h.queue)-1]

	// map[string]*heapItem から Value を取得
	item, ok := h.items[key]
	if !ok {
		return nil
	}

	// builtin.delete(m map[Type]Type1, key Type) で map から指定の KeyValue を削除
	delete(h.items, key)
	return item.obj
}

func (h *data) Peek() interface{} {
	if len(h.queue) > 0 {
		return h.items[h.queue[0]].obj
	}
	return nil
}

// 元ソースでは metricRecorder も定義しているが、
// Heap 実装のみをターゲットにしているため除外
type Heap struct {
	data *data
}

func (h *Heap) Add(obj interface{}) error {
	key, err := h.data.keyFunc(obj)
	if err != nil {
		return cache.KeyError{Obj: obj, Err: err}
	}
	if _, exists := h.data.items[key]; exists {
		h.data.items[key].obj = obj
		heap.Fix(h.data, h.data.items[key].index)
	} else {
		// container/heap.Push(h Interface, x interface{}) が呼ばれて
		// container/heap.Push(x interface{}) を実装している
		// h.Push(kv interface{}) に繋がっている
		heap.Push(h.data, &itemKeyValue{key, obj})
	}

	return nil
}

func (h *Heap) AddIfNotPresent(obj interface{}) error {
	key, err := h.data.keyFunc(obj)
	if err != nil {
		return cache.KeyError{Obj: obj, Err: err}
	}
	if _, exists := h.data.items[key]; !exists {
		heap.Push(h.data, &itemKeyValue{key, obj})
	}
	return nil
}

func (h *Heap) Update(obj interface{}) error {
	return h.Add(obj)
}

func (h *Heap) Delete(obj interface{}) error {
	key, err := h.data.keyFunc(obj)
	if err != nil {
		return cache.KeyError{Obj: obj, Err: err}
	}
	if item, ok := h.data.items[key]; ok {
		heap.Remove(h.data, item.index)
		return nil
	}
	return fmt.Errorf("object not found")
}

// data.Peak() を実行
func (h *Heap) Peek() interface{} {
	return h.data.Peek()
}

func (h *Heap) Pop() (interface{}, error) {
	obj := heap.Pop(h.data)
	if obj != nil {
		return obj, nil
	}
	return nil, fmt.Errorf("object was removed from heap data")
}

func (h *Heap) Get(obj interface{}) (interface{}, bool, error) {
	key, err := h.data.keyFunc(obj)
	if err != nil {
		return nil, false, cache.KeyError{Obj: obj, Err: err}
	}
	return h.GetByKey(key)
}

func (h *Heap) GetByKey(key string) (interface{}, bool, error) {
	item, exists := h.data.items[key]
	if !exists {
		return nil, false, nil
	}
	return item.obj, true, nil
}

func (h *Heap) List() []interface{} {
	list := make([]interface{}, 0, len(h.data.items))
	for _, item := range h.data.items {
		list = append(list, item.obj)
	}
	return list
}

func (h *Heap) Len() int {
	return len(h.data.queue)
}

func New(keyFn KeyFunc, lessFn lessFunc) *Heap {
	//return NewWithRecorder(keyFn, lessFn)
	return &Heap{
		data: &data{
			items:    map[string]*heapItem{},
			queue:    []string{},
			keyFunc:  keyFn,
			lessFunc: lessFn,
		},
	}
}

// metricRecoder は実装しないためコメントアウト
// func NewWithRecorder(keyFn KeyFunc, lessFn lessFunc) *Heap {
// 	return &Heap{
// 		data: &data{
// 			items:    map[string]*heapItem{},
// 			queue:    []string{},
// 			keyFunc:  keyFn,
// 			lessFunc: lessFn,
// 		},
// 	}
// }

type lessFunc = func(item1, item2 interface{}) bool
