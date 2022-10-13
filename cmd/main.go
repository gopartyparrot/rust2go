package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/pkg/profile"
	"github.com/stretchr/testify/assert"

	"github.com/gopartyparrot/rust2go"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	ctx := context.Background()

	start := time.Now()
	count := 0
	for {
		f := rand.Float64()
		u128bits, err := rust2go.F64ToFixBits(ctx, f)
		if err != nil {
			panic(err)
		}

		rf, err := rust2go.U128BitsToFix(ctx, u128bits)
		if err != nil {
			panic(err)
		}

		assert.InDelta(Assert{}, f, rf.InexactFloat64(), 0.0001, "should almost equal")

		count++
		if count%1e5 == 0 {
			fmt.Printf("%10d, %5d ops/sec\n", count, int(float64(count)/time.Since(start).Seconds()))
		}
	}
}

type Assert struct{}

func (a Assert) Errorf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
