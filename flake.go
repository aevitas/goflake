package goflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type ID int64

type Generator struct {
	mu        sync.Mutex
	lastTime  int64
	nodeID    int64
	increment int64
}

const (
	nodeIDBits    uint8 = 10
	sequenceBits  uint8 = 12
	maxNodeID     int64 = (1 << nodeIDBits) - 1
	incrementMask int64 = (1 << sequenceBits) - 1
	timeShift           = nodeIDBits + sequenceBits
	nodeShift           = sequenceBits
)

// Default is 2015-01-01T00:00:00Z, Discord's epoch
var Epoch int64 = 1420070400000

func NewGenerator(nodeID int64) (*Generator, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, fmt.Errorf("node ID must be between 0 and %d", maxNodeID)
	}
	return &Generator{
		nodeID: nodeID,
	}, nil
}

func (g *Generator) Next() (ID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// gets the diff between now and epoch, replacing the monotonic timer
	now := time.Now().UnixMilli() - Epoch

	// if we are generating a second ID within the same ms, increase increment by 1
	if now == g.lastTime {
		g.increment = (g.increment + 1) & incrementMask
		// increment will be 0 for 4096, because we AND the mask above, wait for next ms to generate id
		if g.increment == 0 {
			for now <= g.lastTime {
				now = time.Now().UnixMilli() - Epoch
			}
		}
	} else {
		g.increment = 0
	}

	// protect against clock drift
	if now < g.lastTime {
		return 0, errors.New("clock moved backwards, refusing to generate ID")
	}

	g.lastTime = now

	id := (now << timeShift) | (g.nodeID << nodeShift) | g.increment

	return ID(id), nil
}

// Time returns the timestamp portion of the ID in milliseconds since the Discord epoch.
func (id ID) Time() int64 {
	return (int64(id) >> timeShift)
}

// Node returns the node ID portion of the ID.
func (id ID) Node() int64 {
	return (int64(id) >> nodeShift) & maxNodeID
}

// Increment returns the sequence portion of the ID.
func (id ID) Increment() int64 {
	return int64(id) & incrementMask
}

// String returns the string representation of the ID.
func (id ID) String() string {
	return fmt.Sprintf("%d", id)
}

// Int64 returns the int64 representation of the ID.
func (id ID) Int64() int64 {
	return int64(id)
}
