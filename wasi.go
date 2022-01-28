package rust2go

import (
	_ "embed"

	"github.com/pkg/errors"
	"github.com/wasmerio/wasmer-go/wasmer"
)

//go:embed wasi.wasm
var wasiBytes []byte

var (
	store  *wasmer.Store
	module *wasmer.Module
)

func init() {
	store = wasmer.NewStore(wasmer.NewEngine())
	_module, err := wasmer.NewModule(store, wasiBytes)
	if err != nil {
		panic(err)
	}
	module = _module
}

func callStart(
	args ...string,
) (string, error) {
	builder := wasmer.NewWasiStateBuilder("wasi").
		CaptureStdout().
		CaptureStderr()
	for _, arg := range args {
		builder.Argument(arg)
	}
	wasiEnv, err := builder.Finalize()
	if err != nil {
		return "", err
	}
	importObject, err := wasiEnv.GenerateImportObject(store, module)
	if err != nil {
		return "", err
	}
	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		return "", err
	}

	start, err := instance.Exports.GetWasiStartFunction()
	if err != nil {
		return "", err
	}
	_, err = start()
	if err != nil {
		return "", errors.Wrap(err, "wasi start invoke failed")
	}

	errBytes := wasiEnv.ReadStderr()
	if len(errBytes) > 0 {
		return "", errors.Errorf("wasi stderr: %s", string(errBytes))
	}
	return string(wasiEnv.ReadStdout()), nil
}
