Goflake
=======

## Warning
This is a very early implementation of Discord-style snowflake IDs in Golang. It is roughly based on my [C# Snowflake implementation called FlakeID](https://github.com/aevitas/flakeid). The code in this repository is very much still a work in progress, and should not yet be used in any serious capacity.

## Snowflakes
Snowflake IDs were originally introduced by Twitter in 2010 as unique, decentralized IDs for Tweets. Their 8-byte size, ordered nature and guaranteed uniqueness make them ideal to use as resource identifiers. Since then, many applications at various scale have adopted Snowflake-esque identifiers.

This repository contains an implementation of decentralized, K-ordered Snowflake IDs based on the Discord Snowflake specification. The implementation heavily focuses on high-throughput, supporting upwards of 10.000 unique generations per second on commodity hardware.

## How it works

Every Snowflake fits in a 64-bit integer, consisting of various components that make it unique across generations.
The layout of the components that comprise a snowflake can be expressed as:

```
Timestamp                                   Thread Proc  Increment
111111111111111111111111111111111111111111  11111  11111 111111111111
64                                          22     17    12          0
```

The Timestamp component is represented as the milliseconds since the first second of 2015. Since we're using all 64 bits available, this epoch can be any point in time, as long as it's in the past. If the epoch is set to a point in time in the future, it may result in negative snowflakes being generated.

Where the original Discord reference mentions worker ID and process ID, we substitute these with the
thread and process ID respectively, as the combination of these two provide sufficient uniqueness, and they are
the closest we can get to the original specification within the .NET ecosystem.

The Increment component is a monotonically incrementing number, which is incremented every time a snowflake is generated.
This is in contrast with some other flake-ish implementations, which only increment the counter any time a snowflake is 
generated twice at the exact same instant in time. We believe Discord's implementation is more correct here,
as even two snowflakes that are generated at the exact same point in time will not be identical, because of their increments.

## Usage

Create a goflake using the `Create` method:

```golang
import (
	"github.com/aevitas/goflake"
)

func main() {
	id := Id.Create()
}
```