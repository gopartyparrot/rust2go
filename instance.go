package rust2go

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const reInstantiateThreshold = 16366 //got this through TestReInstantiateThreshold

//go:embed wasmlib.wasm
var wasmBytes []byte

type Instance struct {
	enableReInstantiate bool //just keep this value true, unless you need to test it

	runtime wazero.Runtime
	mod     api.Module

	callCount int

	mu sync.Mutex
}

func newInstance(ctx context.Context) (*Instance, error) {
	rt, mod, err := createRuntime(ctx)
	if err != nil {
		return nil, err
	}
	return &Instance{enableReInstantiate: true, runtime: rt, mod: mod, callCount: 0}, nil
}

func createRuntime(ctx context.Context) (wazero.Runtime, api.Module, error) {
	rt := wazero.NewRuntime(ctx) // Create a new WebAssembly Runtime.
	// Instantiate a Go-defined module named "env" that exports a function to
	// log to the console.
	_, err := rt.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(logString).Export("log").
		Instantiate(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Instantiate a WebAssembly module that imports the "log" function defined
	// in "env" and exports "memory" and functions we'll use in this example.
	mod, err := rt.Instantiate(ctx, wasmBytes)
	if err != nil {
		return nil, nil, err
	}
	return rt, mod, nil
}

func (inst *Instance) resetRuntimeIfNeeded(ctx context.Context) (err error) {
	if !inst.enableReInstantiate {
		return nil
	}
	if inst.callCount < reInstantiateThreshold {
		return nil
	}

	if inst.runtime != nil {
		err = inst.runtime.Close(ctx)
		if err != nil {
			return errors.Wrap(err, "close current runtime failed")
		}
	}

	rt, mod, err := createRuntime(ctx)
	if err != nil {
		return err
	}
	inst.runtime = rt
	inst.mod = mod
	inst.callCount = 0
	return
}

func (inst *Instance) F64ToFixBits(ctx context.Context, f float64) (*big.Int, error) {
	ret, err := inst.processString(ctx, "f64_to_fix_bits", fmt.Sprintf("%f", f))
	if err != nil {
		return nil, err
	}
	b, ok := big.NewInt(0).SetString(ret, 10)
	if !ok {
		return nil, errors.Errorf("unable to parse result as u128: %s", b)
	}
	return b, nil
}

func (inst *Instance) U128BitsToFix(ctx context.Context, b *big.Int) (decimal.Decimal, error) {
	ret, err := inst.processString(ctx, "u128bits_to_fix", b.String())
	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromString(ret)
}

func (inst *Instance) processString(ctx context.Context, callFn string, input string) (string, error) {
	inst.mu.Lock()
	defer inst.mu.Unlock()

	err := inst.resetRuntimeIfNeeded(ctx)
	if err != nil {
		return "", errors.Wrap(err, "try reset runtime failed")
	}
	defer func() { inst.callCount++ }()

	inputSize := uint64(len(input))

	// Instead of an arbitrary memory offset, use Rust's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	allocateResults, err := inst.mod.ExportedFunction("allocate").Call(ctx, inputSize)
	if err != nil {
		return "", errors.Wrap(err, "allocate failed")
	}

	inputPtr := allocateResults[0]
	// This pointer was allocated by Rust, but owned by Go, So, we have to
	// deallocate it when finished
	defer inst.mod.ExportedFunction("deallocate").Call(ctx, inputPtr, inputSize)

	// The pointer is a linear memory offset, which is where we write the name.
	if !inst.mod.Memory().Write(uint32(inputPtr), []byte(input)) {
		return "", errors.Errorf("Memory.Write(%d, %d) out of range of memory size %d",
			inputPtr, inputSize, inst.mod.Memory().Size())
	}

	// Finally, we get the greeting message "greet" printed. This shows how to
	// read-back something allocated by Rust.
	retPtrs, err := inst.mod.ExportedFunction(callFn).Call(ctx, inputPtr, inputSize)
	if err != nil {
		return "", errors.Wrap(err, "call fn failed")
	}

	retPtr, retSize := uint32(retPtrs[0]>>32), uint32(retPtrs[0])

	// This pointer was allocated by Rust, but owned by Go, So, we have to
	// deallocate it when finished
	defer inst.mod.ExportedFunction("deallocate").Call(ctx, uint64(retPtr), uint64(retSize))

	// The pointer is a linear memory offset, which is where we write the name.
	bytes, ok := inst.mod.Memory().Read(retPtr, retSize)
	if !ok {
		return "", errors.Errorf("Memory.Read(%d, %d) out of range of memory size %d",
			retPtr, retSize, inst.mod.Memory().Size())
	}
	return string(bytes), nil
}

func logString(ctx context.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Fprintln(os.Stdout, string(buf)) //TODO can be configured to other Writer(may be use global var)
}
