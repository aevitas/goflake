package goflake

import (
	"log"
	"testing"
	"time"
)

func TestCreate(*testing.T) {
	var id Id

	start := time.Now()

	i := 0
	for i < 1000 {
		id = id.Create()

		i += 1
	}

	elapsed := time.Since(start)
	log.Printf("1000 ids took %s", elapsed)
}
