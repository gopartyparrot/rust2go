setup:
	rustup target add wasm32-wasi
buildWasm:
	cargo build --target wasm32-wasi --release
	cp target/wasm32-wasi/release/rust2go.wasm ./wasi.wasm
test:
	@echo ===================rust   test====================
	cargo test
	@echo ===================golang test====================
	go test