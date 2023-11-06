package main

import (
	"fmt"
)

type (
	Key   string
	Value string
)

type WriteTs struct {
	value uint64
}

func (t *WriteTs) Inc() {
	t.value++
}

type StampedValue struct {
	stamp *WriteTs
	value Value
}

func main() {
	c := NewCoordinator(3)
	c.Set("a", "1")
	fmt.Println(c.Get("a"))
	for _, peer := range c.peers {
		fmt.Println(peer)
	}
}
