package main

import (
	"container/heap"
)

type farm struct {
	money     float32
	land_free int
	seeds     []int
	plants    PriorityQueue
}

// PQ code from https://cs.opensource.google/go/go/+/master:src/container/heap/example_pq_test.go
type Item struct {
	value    float32
	priority int // The priority of the item in the queue
	index    int // index is needed by update (from go docs)
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value float32, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func init_farm() *farm {
	const lands = 3
	const seeds = 1
	current := &farm{
		money:     0.0,
		land_free: lands,
		seeds:     make([]int, len(seed_type)),
	}

	for k := range current.seeds {
		current.seeds[k] = seeds
	}

	return current
}
