package rust2go

import (
	"context"
	"math/big"

	"github.com/shopspring/decimal"
)

var instance *Instance

func init() {
	inst, err := newInstance(context.Background())
	if err != nil {
		panic(err)
	}
	instance = inst
}

func F64ToFixBits(ctx context.Context, f float64) (*big.Int, error) {
	return instance.F64ToFixBits(ctx, f)
}

func U128BitsToFix(ctx context.Context, b *big.Int) (decimal.Decimal, error) {
	return instance.U128BitsToFix(ctx, b)
}
