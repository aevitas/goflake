# Goflake

This repository contains an implementation of Snowflake IDs, similar to [my implementation in C#](https://github.com/aevitas/flakeid/). Both are based on the [Discord Snowflake specification](https://discord.com/developers/docs/reference).

While the implementation in this repository is functional and reasonably performant, it is not currently considered ready for production scenarios.

## Usage

Install Goflake into any project using Go packages by running:

```
go get github.com/aevitas/goflake
```

Afterwards, simply generate an ID:

```go
id := goflake.NewId()
```

## Performance

Current performance on an M1 MacBook is very reasonable:
```
goos: darwin
goarch: arm64
pkg: github.com/aevitas/goflake
BenchmarkNewIdPerf-10    	22709946	        52.92 ns/op	       8 B/op	       1 allocs/op
PASS
```

Allowing for some ~20 million unique IDs to be generated every second.
