# rustgo

call rust code from go

## Start up

```bash
# build wasm
make zero

# test
make test
```

## Usage

```bash
go get github.com/gopartyparrot/rust2go
```

```go
import (
  "fmt"
  "github.com/gopartyparrot/rust2go"
)

func main() {
  b, _ := rust2go.F64ToFixBits(ctx, 0.25)
  f, _ := rust2go.U128BitsToFix(ctx, big.NewInt(4611686018427387904))
}
```

## Benchmarks

```
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkFixFunctions
BenchmarkFixFunctions-12    	   68386	     42035 ns/op	   17977 B/op	     117 allocs/op
```

## FAQ

