zero:
	cargo build --release --target wasm32-unknown-unknown 
	cp target/wasm32-unknown-unknown/release/wasmlib.wasm ./wasmlib.wasm
	
test:
	@echo ===================rust   test====================
	cargo test
	@echo ===================golang test====================
	go test