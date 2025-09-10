package goflake

import (
	"sync"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	t.Run("ValidNodeID", func(t *testing.T) {
		_, err := NewGenerator(0)
		if err != nil {
			t.Errorf("expected no error for node ID 0, but got %v", err)
		}
		_, err = NewGenerator(maxNodeID)
		if err != nil {
			t.Errorf("expected no error for max node ID, but got %v", err)
		}
	})

	t.Run("InvalidNodeID", func(t *testing.T) {
		_, err := NewGenerator(-1)
		if err == nil {
			t.Error("expected an error for negative node ID, but got nil")
		}
		_, err = NewGenerator(maxNodeID + 1)
		if err == nil {
			t.Error("expected an error for oversized node ID, but got nil")
		}
	})
}

func TestNextIDUniqueness(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	const numIDs = 50000 // Generate a significant number of IDs
	ids := make(map[ID]bool, numIDs)

	for i := 0; i < numIDs; i++ {
		id, err := generator.Next()
		if err != nil {
			t.Fatalf("failed to generate ID: %v", err)
		}
		if ids[id] {
			t.Fatalf("duplicate ID generated: %d", id)
		}
		ids[id] = true
	}
}

func TestNextIDMonotonicity(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	var lastID ID
	for i := 0; i < 50000; i++ {
		id, err := generator.Next()
		if err != nil {
			t.Fatalf("failed to generate ID: %v", err)
		}
		if i > 0 && id <= lastID {
			t.Errorf("ID is not monotonic: current=%d, last=%d", id, lastID)
		}
		lastID = id
	}
}

func TestIDDecomposition(t *testing.T) {
	nodeID := int64(123)
	generator, err := NewGenerator(nodeID)
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	id, err := generator.Next()
	if err != nil {
		t.Fatalf("failed to generate ID: %v", err)
	}

	if id.Node() != nodeID {
		t.Errorf("expected node ID %d, but got %d", nodeID, id.Node())
	}

	// The increment should be 0 for the first ID in a new millisecond.
	if id.Increment() != 0 {
		t.Errorf("expected increment 0 for the first ID, but got %d", id.Increment())
	}

	// Check the timestamp. Allow for a small delta due to execution time.
	expectedTime := (id.Int64() >> timeShift)
	if id.Time() != expectedTime {
		t.Errorf("expected time %d, but got %d", expectedTime, id.Time())
	}

	// Generate a second ID immediately to test the increment
	id2, err := generator.Next()
	if err != nil {
		t.Fatalf("failed to generate second ID: %v", err)
	}

	// If they are in the same millisecond, the increment should be 1
	if id2.Time() == id.Time() {
		if id2.Increment() != 1 {
			t.Errorf("expected increment to be 1, but got %d", id2.Increment())
		}
	}
}

func TestConcurrentGeneration(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	const numGoroutines = 50
	const idsPerGoroutine = 2000
	totalIDs := numGoroutines * idsPerGoroutine

	var wg sync.WaitGroup
	idChan := make(chan ID, totalIDs)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := generator.Next()
				if err != nil {
					// Use t.Errorf from a goroutine as t.Fatalf will exit the wrong one
					t.Errorf("failed to generate ID: %v", err)
					return
				}
				idChan <- id
			}
		}()
	}

	wg.Wait()
	close(idChan)

	// Verify all collected IDs are unique
	ids := make(map[ID]bool, totalIDs)
	for id := range idChan {
		if ids[id] {
			t.Fatalf("found duplicate ID in concurrent generation: %d", id)
		}
		ids[id] = true
	}

	if len(ids) != totalIDs {
		t.Errorf("expected %d unique IDs, but got %d", totalIDs, len(ids))
	}
}

func BenchmarkSingleThread(b *testing.B) {
	generator, err := NewGenerator(1)
	if err != nil {
		b.Fatalf("failed to create generator: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = generator.Next()
	}
}

// BenchmarkConcurrent benchmarks ID generation from multiple goroutines.
func BenchmarkConcurrent(b *testing.B) {
	generator, err := NewGenerator(1)
	if err != nil {
		b.Fatalf("failed to create generator: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = generator.Next()
		}
	})
}
