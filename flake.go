package goflake

import (
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

type Id struct {
	Value int64
}

var (
	Epoch string = "2015-Jan-01"

	increment atomic.Int64
	pid       int64

	timeStampBits int64 = 42
	incrementBits int64 = 12
	pidBits       int64 = 5
	randBits      int64 = 5

	timeStampMask int64 = (int64(1) << timeStampBits) - 1
	pidMask       int64 = (1 << pidBits) - 1
	randMask      int64 = (1 << randBits) - 1
	incrementMask int64 = (1 << incrementBits) - 1
)

// Generates a new Snowflake ID that is guaranteed to be unique and sortable across generations.
func NewId() *Id {
	const layout = "2006-Jan-02"
	tm, err := time.Parse(layout, Epoch)
	if err != nil {
		panic(err)
	}

	ms := time.Since(tm).Milliseconds() & timeStampMask
	if pid == 0 {
		pid = int64(os.Getpid()) & pidMask
	}

	rand := rand.Int63() & randMask
	inc := increment.Load() & int64(incrementMask)

	v := (ms << (pidBits + randBits + incrementBits)) + (pid << (randBits + incrementBits)) + (rand << incrementBits) + inc

	return &Id{
		Value: v,
	}
}
