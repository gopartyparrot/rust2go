rustgo

===

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

## FAQ

- Q: why wasi instead of wasm
  - A: easier to understand and don't need to use wasm memory api
- Q: why not https://github.com/bytecodealliance/wasmtime-go
  - A: wasmtime uses CGO to consume the C API of wasmtime project(which wrote in c++), it increase build size(at least 20M), and project only works on some specific platform.
