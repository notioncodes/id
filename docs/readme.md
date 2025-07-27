# Notion ID Utility

This package provides a parser for safely handling Notion IDs.

Notion IDs are a mix of hexadecimal and dashed formats (uuid v4). The parser
will normalize the ID to the dashed format if all conditions are met.

> [!NOTE]
> The parser is thread safe and can be reused for multiple calls in goroutines and
> also has a minimal cache to improve performance.

## Installation

```bash
go get github.com/notioncodes/id
```

## Usage

```go
import (
  "fmt"
  "log"
  "github.com/notioncodes/id"
)

func main() {
  // Notion likes to mix up the format of the ID so we support hexacdeimal
  // and dashed formats (uuid v4). The parser below will normalize the ID
  // to the dashed format if all goes well.
 potentialNotionID := "550e8400e29b41d4a716446655440000"

 // Create a new ID parser with a NoOpCache so that we don't have to
 // recreate the parser for each call down the line.
 //
 // Also, the parser is thread safe and can be reused for multiple calls in
 // goroutines.
 parser := id.NewIDParser(id.NewNoOpCache())

 // Parse the ID and get a normalized formatted ID as a UUID string.
 parsed, err := parser.Parse(potentialNotionID)

 // If the ID is invalid, we can't parse it safely.
 // We can see why the parse failed by checking the error.
 if err != nil {
  log.Fatal(err)
 }

 // The ID is valid, returns a normalized formatted ID as a UUID string:
 // "550e8400-e29b-41d4-a716-446655440000"
 fmt.Println(parsed)
}
```

## Performance

To run the benchmarks, run the following command:

```bash
go test -bench=. -benchmem
```

## Benchmark Results

```bash
goos: darwin
goarch: arm64
pkg: github.com/notioncodes/id
cpu: Apple M1 Max
BenchmarkValid/WithCache-10             130095438                9.194 ns/op           0 B/op          0 allocs/op
BenchmarkValid/WithoutCache-10          19246602                62.49 ns/op           64 B/op          2 allocs/op
BenchmarkInvalid/WithCache-10           263531785                4.562 ns/op           0 B/op          0 allocs/op
BenchmarkInvalid/WithoutCache-10        20465126                57.06 ns/op           48 B/op          2 allocs/op
PASS
ok      github.com/notioncodes/id       6.478s
```
