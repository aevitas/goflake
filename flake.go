package goflake

import (
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

var increment int64

type Id struct {
	value int64
}

func (id *Id) Create() Id {
	return createInternal()
}

func (id *Id) String() string {
	return fmt.Sprint(id.value)
}

func createInternal() Id {
	const timestampBits = int64(42)
	const randBits = int64(5)
	const processIdBits = int64(5)
	const incrementBits = int64(12)

	const timestampMask = (int64(1) << timestampBits) - 1
	const randMask = (int64(1) << randBits) - 1
	const processIdMask = (int64(1) << processIdBits) - 1
	const incrementMask = (int64(1) << incrementBits) - 1

	ts := getElapsedMilliseconds()
	timestamp := ts & timestampMask

	rn := rand.NewSource(time.Now().UnixNano()).Int63()
	rand := rn & randMask

	pid := int64(os.Getpid())
	processId := pid & processIdMask

	atomic.AddInt64(&increment, 1)
	increment := increment & incrementMask

	val := (timestamp << (randBits + processIdBits + incrementBits)) + (rand << (processIdBits + incrementBits)) + (processId << incrementBits) + increment

	return Id{value: val}
}

func getElapsedMilliseconds() int64 {
	epoch := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)

	return epoch.UnixMilli()
}
