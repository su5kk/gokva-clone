package main

import (
	"context"
)

type Peer struct {
	ID      string
	storage map[Key]StampedValue
}

func NewPeer(id string) *Peer {
	return &Peer{
		ID:      id,
		storage: map[Key]StampedValue{},
	}
}

func (p *Peer) Read(key Key) StampedValue {
	return p.storage[key]
}

func (p *Peer) Write(ctx context.Context, key Key, value StampedValue) (set bool) {
	operationDone := make(chan struct{})

	go func() {
		old, ok := p.storage[key]
		if !ok || old.stamp.value < value.stamp.value {
			p.storage[key] = value
			set = true
		}
		close(operationDone)
	}()

	select {
	case <-ctx.Done():
		return false
	case <-operationDone:
		// operation completed: value may or may not have been written.
	}

	return
}
