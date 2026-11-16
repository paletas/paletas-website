package concurrent

import (
	"sync"
)

type ConcurrentSlice struct {
	items []interface{}
	lock  *sync.RWMutex
}

type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

func NewConcurrentSlice() *ConcurrentSlice {
	return &ConcurrentSlice{
		items: make([]interface{}, 0),
		lock:  &sync.RWMutex{},
	}
}

func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.items = append(cs.items, item)
}

func (cs *ConcurrentSlice) Remove(index int) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.items = append(cs.items[:index], cs.items[index+1:]...)
}

func (cs *ConcurrentSlice) RemoveItem(item interface{}) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	for index, value := range cs.items {
		if value == item {
			cs.items = append(cs.items[:index], cs.items[index+1:]...)
		}
	}
}

func (cs *ConcurrentSlice) Get(index int) interface{} {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.items[index]
}

func (cs *ConcurrentSlice) GetAll() []interface{} {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.items
}

func (cs *ConcurrentSlice) GetByFilter(filter func(item interface{}) bool) []interface{} {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	var filteredItems []interface{}

	for _, item := range cs.items {
		if filter(item) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)

	f := func() {
		cs.lock.RLock()
		defer cs.lock.RUnlock()

		for index, item := range cs.items {
			c <- ConcurrentSliceItem{index, item}
		}
		close(c)
	}
	go f()

	return c
}

func (cs *ConcurrentSlice) Len() int {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return len(cs.items)
}
