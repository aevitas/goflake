package goflake_test

import (
	"testing"

	"github.com/aevitas/goflake"
)

func TestNewId(t *testing.T) {
	id := goflake.NewId()

	if id.Value < 0 {
		t.Fail()
	}
}

func TestIdUniqueness(t *testing.T) {
	var ids []goflake.Id

	for i := 0; i < 100; i++ {
		ids = append(ids, *goflake.NewId())
	}

	for x := range ids {
		count := 0
		for y := range ids {
			if x == y {
				if count == 0 {
					count++
					continue
				}

				t.Fail()
			}
		}
	}
}

func BenchmarkNewIdPerf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		goflake.NewId()
	}
}
