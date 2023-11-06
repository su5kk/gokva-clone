package main

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"
)

func NewCoordinator(peerCount int) *Coordinator {
	peers := make([]*Peer, 0, peerCount)
	for i := 0; i < peerCount; i++ {
		id := "peer-" + strconv.Itoa(i)
		peers = append(peers, NewPeer(id))
	}

	return &Coordinator{
		peers: peers,
		log:   log.Default(),
		ts:    &WriteTs{0},
	}
}

type Coordinator struct {
	peers []*Peer
	log   *log.Logger
	ts    *WriteTs
}

func (c *Coordinator) Set(key Key, value Value) error {
	c.ts.Inc()
	return c.set(key, StampedValue{
		value: value,
		stamp: c.ts,
	})
}

func (c *Coordinator) set(key Key, value StampedValue) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(c.peers))
	acks := make(chan bool, 1)

	for _, peer := range c.peers {
		go func(peer *Peer) {
			defer wg.Done()

			// Use a separate context with a timeout for each peer to ensure that
			// if a peer hangs, it doesn't affect the operation for other peers.
			peerCtx, peerCancel := context.WithTimeout(ctx, 5*time.Second)
			defer peerCancel()

			select {
			case acks <- peer.Write(peerCtx, key, value):
				// peer acknowledged write
			case <-peerCtx.Done():
				// peer timed out
			case <-ctx.Done():
				// acks are no longer required
			}
		}(peer)
	}

	go func() {
		wg.Wait()
		close(acks)
	}()

	err := c.awaitWrites(ctx, c.majority(), acks)
	if err != nil {
		return err
	}

	return nil
}

func (c *Coordinator) majority() uint64 {
	return uint64(len(c.peers)/2 + 1)
}

func (c *Coordinator) awaitWrites(ctx context.Context, threshold uint64, acks <-chan bool) error {
	count := uint64(0)

	for {
		select {
		case <-acks:
			count++
			if count >= threshold {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		default:
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// continue
		}
	}
}

func (c *Coordinator) Get(key Key) (Value, error) {
	reads := make(chan StampedValue, 1)

	var wg sync.WaitGroup
	wg.Add(len(c.peers))

	for _, peer := range c.peers {
		go func(peer *Peer) {
			defer wg.Done()
			reads <- peer.Read(key)
		}(peer)
	}

	go func() {
		wg.Wait()
		close(reads)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	actualReads, err := c.awaitReads(ctx, reads)
	if err != nil {
		return "", err
	}

	mostRecent := c.mostRecent(actualReads)

	if err = c.set(key, mostRecent); err != nil {
		return "", err
	}

	return mostRecent.value, nil
}

func (c *Coordinator) awaitReads(
	ctx context.Context,
	reads <-chan StampedValue,
) ([]StampedValue, error) {
	values := make([]StampedValue, 0)

	for {
		select {
		case <-ctx.Done():
			c.log.Println("shutting down")
			return values, ctx.Err()
		case read, ok := <-reads:
			if !ok {
				return values, nil
			}
			if read.stamp != nil {
				values = append(values, read)
			}
		default:
			if ctx.Err() != nil {
				c.log.Printf("shutting down: %v", ctx.Err())
				return values, ctx.Err()
			}
			// continue
		}
	}
}

func (c *Coordinator) mostRecent(reads []StampedValue) StampedValue {
	rec := reads[0]
	for _, read := range reads {
		if read.stamp.value > rec.stamp.value {
			rec = read
		}
	}
	return rec
}
