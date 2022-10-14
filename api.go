package rust2go

import (
	"context"
	"log"
	"math/big"

	"github.com/shopspring/decimal"
)

var instance *Instance

var works = false

/***
In practice, U128BitsToFix got wrong value some time.
IsWork() indicate whether it works for known test case.
If it return false, please do not use this library.
*/
func IsWork() bool {
	return works
}

func init() {
	inst, err := newInstance(context.Background())
	if err != nil {
		panic(err)
	}
	instance = inst

	testCase := struct{ input, output string }{"1341723281558402372940796526592", "72734964837"}
	testBig, _ := big.NewInt(0).SetString(testCase.input, 10)
	ret, err := U128BitsToFix(context.Background(), testBig)
	if err != nil {
		panic(err)
	}
	works = ret.String() == testCase.output
}

func F64ToFixBits(ctx context.Context, f float64) (*big.Int, error) {
	if !works {
		log.Println("WARN (rust2go): F64ToFixBits return value may be wrong, please do not use this library!")
	}
	return instance.F64ToFixBits(ctx, f)
}

func U128BitsToFix(ctx context.Context, b *big.Int) (decimal.Decimal, error) {
	if !works {
		log.Println("WARN (rust2go): U128BitsToFix return value may be wrong, please do not use this library!")
	}
	return instance.U128BitsToFix(ctx, b)
}
