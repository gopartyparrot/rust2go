# rustgo

call rust code from go

## How

- Build rust code to wasi
  - rust parse wasi args
  - rust print result(result/error) to stdout/stderr
- Call wasi from go, use wasm runtime(implemented by go)
  - read stdout/stderr as result/err

## Start up

```bash
# setup
make setup

# build wasm
make buildWasm

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
  fmt.Println(rust2go.ParseU64F64("2276328218695758774272"))
  // output:
  // 123.40000000000000568434 <nil>
}
```

## FAQ

- Q: why wasi instead of wasm
  - A: easier to understand and don't need to use wasm memory api
- Q: why not https://github.com/bytecodealliance/wasmtime-go
  - A: wasmtime uses CGO to consume the C API of wasmtime project(which wrote in rust), it increase build size(at least 20M), and project only works on some specific platform.
