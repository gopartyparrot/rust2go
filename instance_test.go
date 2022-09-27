package rust2go

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestExample(t *testing.T) {
	b, _ := F64ToFixBits(ctx, 0.25)
	f, _ := U128BitsToFix(ctx, big.NewInt(4611686018427387904))
	fmt.Println(f, b)
}

func TestInstance(t *testing.T) {
	f := rand.Float64()

	u128bits, err := F64ToFixBits(ctx, f)
	require.NoError(t, err)

	f64, err := U128BitsToFix(ctx, u128bits)
	require.NoError(t, err)

	deltaEq(t, f, f64.InexactFloat64())

	//u128::max ==> 340282366920938463463374607431768211455
	max_exceeded := big.NewInt(0)
	max_exceeded, ok := max_exceeded.SetString("34028236692093846346337460743176821145500000", 10)
	require.True(t, ok)

	_, err = U128BitsToFix(ctx, max_exceeded)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ERR: invalid u128")
}

func BenchmarkInstantiateRuntime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inst, err := newInstance(ctx)
		require.NoError(b, err)

		_, err = inst.F64ToFixBits(ctx, rand.Float64())
		require.NoError(b, err)

		require.NoError(b, inst.runtime.Close(ctx))
	}
}

func BenchmarkFixFunctions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := rand.Float64()
		u128bits, err := F64ToFixBits(ctx, f)
		require.NoError(b, err)

		_, err = U128BitsToFix(ctx, u128bits)
		require.NoError(b, err)
	}
}

// this test find a reInstantiateThreshold
// TODO should test in different env, e.g. server env, docker env
func TestReInstantiateThreshold(t *testing.T) {
	inst, _ := newInstance(ctx)
	inst.enableReInstantiate = false

	count := 0
	for {
		f := rand.Float64()
		u128bits, err := inst.F64ToFixBits(ctx, f)
		if err != nil {
			break
		}

		f64, err := inst.U128BitsToFix(ctx, u128bits)
		if err != nil {
			break
		}

		deltaEq(t, f, f64.InexactFloat64())
		count++
	}

	t.Log("reInstantiateThreshold", count)
}

// simple test: how many ops per second
func TestSimpleSpeed(t *testing.T) {
	if testing.Short() {
		return
	}

	count := 0
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
outer:
	for {
		select {
		case <-ctx.Done():
			break outer
		default:
		}
		f := rand.Float64()
		u128bits, err := F64ToFixBits(ctx, f)
		require.NoError(t, err)

		f64, err := U128BitsToFix(ctx, u128bits)
		require.NoError(t, err)

		deltaEq(t, f, f64.InexactFloat64())
		count++
	}
	t.Logf("total count: %d, ops/sec: %d\n", count, 2*count/5)
}

// should works in goroutines, since mutex used
func TestConcurrency(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU()-2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

		outer:
			for {
				select {
				case <-ctx.Done():
					break outer
				default:
				}

				f := rand.Float64()
				u128bits, err := F64ToFixBits(ctx, f)
				require.NoError(t, err)

				f64, err := U128BitsToFix(ctx, u128bits)
				require.NoError(t, err)
				deltaEq(t, f, f64.InexactFloat64())
			}
		}()
	}

	wg.Wait()
}

func deltaEq(t *testing.T, f0, f1 float64) {
	require.InDelta(t, f0, f1, 0.001, "not nearly equal: %f %f", f0, f1)
}
