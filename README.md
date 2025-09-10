# Goflake

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/goflake)](https://goreportcard.com/report/github.com/aevitas/goflake)
[![GoDoc](https://godoc.org/github.com/your-username/goflake?status.svg)](https://godoc.org/github.com/aevitas/goflake)

<p align="center"><img src="./assets/snowflake-96.png" width="96"></p>

`goflake` is a high-performance, thread-safe, and K-ordered unique ID generator written in idiomatic Go.

This implementation is based on the [Discord Snowflake specification](https://discord.com/developers/docs/reference). Snowflake IDs are 64-bit integers that are sortable, unique across a distributed system, and ideal for use as primary keys in databases. This package focuses on high throughput, correctness, and a simple, clean API.

## Features

-   **Discord-Compliant:** Follows the well-established bit layout for Snowflake IDs.
-   **High Performance:** Optimized for low-latency ID generation.
-   **Thread-Safe:** Safely generate IDs from multiple goroutines.
-   **Configurable Node ID:** Designed for distributed systems by allowing a configurable `NodeID`.
-   **Zero Dependencies:** A lightweight, standard library-only implementation.
-   **Idiomatic API:** Simple to use and integrate into any Go project.

## Installation

```sh
go get github.com/aevitas/goflake
```

## Quick Start

The core of the package is the `Generator`. You create a generator instance with a unique **node ID** and then call its `Next()` method.

```go
package main

import (
	"fmt"
	"log"
	"github.com/aevitas/goflake" // <-- Replace with your actual import path
)

func main() {
	// 1. Create a new generator.
	// The node ID must be unique for each running instance of your service.
	// It can be any number between 0 and 1023.
	nodeID := int64(1)
	generator, err := goflake.NewGenerator(nodeID)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// 2. Generate a new unique ID.
	id, err := generator.Next()
	if err != nil {
		log.Fatalf("Failed to generate ID: %v", err)
	}

	// 3. Use the ID.
	fmt.Printf("Generated ID: %d\n", id)

	// You can also decompose the ID to get its parts:
	fmt.Printf("  - Time:      %d (ms since epoch)\n", id.Time())
	fmt.Printf("  - Node ID:   %d\n", id.Node())
	fmt.Printf("  - Increment: %d\n", id.Increment())
}
```

## Anatomy of an ID

Every `goflake` ID is a 64-bit integer composed of three parts, ensuring uniqueness and sortability.

```
Timestamp                                   Node ID    Increment
111111111111111111111111111111111111111111  1111111111 111111111111
63                                          22         12          0
```

-   **Timestamp** (42 bits): The number of milliseconds that have passed since the configured **epoch**. This makes the IDs roughly time-sortable.
-   **Node ID** (10 bits): A unique identifier for the machine or process generating the ID. You are responsible for assigning a unique Node ID (0-1023) to each `Generator` instance to prevent collisions in a distributed system.
-   **Increment** (12 bits): A sequence number that increments for each ID generated within the same millisecond on the same node. This allows for up to **4096** unique IDs to be generated per millisecond, per node, some 4 million every second per node.

## Custom Epoch

The timestamp component is a delta from a predefined instant in time, known as the **epoch**. The default epoch is `2015-01-01T00:00:00Z`, matching Discord's implementation.

This default allows for valid ID generation for approximately 139 years. While you can change the epoch for your entire application, it is **critical** to use a single, consistent epoch across your whole system. Changing the epoch after IDs have already been generated will lead to collisions.

To change the epoch, modify the public `Epoch` variable before creating any generators:

```go
// Set a custom epoch to the Unix epoch (Jan 1, 1970)
goflake.Epoch = 0 // Milliseconds since standard Unix epoch
```

## A Warning for Web Clients (JavaScript/JSON)

⚠️ **Be careful when exposing 64-bit integer IDs to JavaScript clients.** ⚠️

JavaScript's `Number` type can only safely represent integers up to `Number.MAX_SAFE_INTEGER` (which is 2^53 - 1). A 64-bit Snowflake ID will exceed this limit, leading to precision loss and bugs.

For example, the ID `931124405369716748` might become `931124405369716700` in a browser.

**The recommended solution is to serialize the ID as a string in your JSON payloads.** You can easily achieve this using Go struct tags:

```go
import "github.com/aevitas/goflake"

type User struct {
    ID   goflake.ID `json:"id,string"` // <-- Make sure to mark the field as a string
    Name string     `json:"name"`
}
```

This will produce JSON like `{"id": "931124405369716748", "name": "Alice"}`, which can be safely handled by all JavaScript clients.

## Performance

This library is designed to be extremely fast. The implementation is thread safe. You can run the built-in benchmarks to see the performance on your own hardware.

To run the benchmarks:

```sh
go test -bench=.
```

| Benchmark               | Iterations | Time/Op (ns/op) |
| ----------------------- | ---------- | --------------- |
| `BenchmarkSingleThread` | `4923973`  | `242.9`         |
| `BenchmarkConcurrent`   | `4886612`  | `244.1`         |

## Contributions

Pull Requests are welcome! If you find an issue, please open a ticket. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](./LICENSE) file for details.
```