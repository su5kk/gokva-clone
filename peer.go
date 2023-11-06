package main

type Peer struct {
	ID      string
	storage map[Key]StampedValue
}

func (p *Peer) Read(key Key) StampedValue {
	return p.storage[key]
}

func (p *Peer) Write(key Key, value StampedValue) (set bool) {
	old, ok := p.storage[key]
	if !ok || old.stamp.value < value.stamp.value {
		p.storage[key] = value
		set = true
	}

	return
}
