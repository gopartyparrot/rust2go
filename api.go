package rust2go

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

var instance *Instance

var works = false

const warnMsg = "WARN (rust2go): U128BitsToFix return value may be wrong, please do not use this library!"

/*
**
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

	testCase := struct{ input, output string }{"1341723281558402372940796526592", "72734964837"} //known test case that fail at some machine
	testBig, _ := big.NewInt(0).SetString(testCase.input, 10)
	ret, err := inst.U128BitsToFix(context.Background(), testBig)
	if err != nil {
		panic(err)
	}
	works = ret.String() == testCase.output
	if !works {
		log.Println(warnMsg)
	}
}

func F64ToFixBits(ctx context.Context, f float64) (*big.Int, error) {
	if !works {
		log.Println(warnMsg)
	}
	return instance.F64ToFixBits(ctx, f)
}

func U128BitsToFix(ctx context.Context, b *big.Int) (decimal.Decimal, error) {
	if !works {
		log.Println(warnMsg)
	}
	return instance.U128BitsToFix(ctx, b)
}

func LoadU128BitsToFixCache(path string) (*FileCache[string, decimal.Decimal], error) {
	// line should in format: "u128bits:decimal"
	lineParser := func(line string) (string, decimal.Decimal, error) {
		arr := strings.Split(strings.TrimSpace(line), ":")
		if len(arr) != 2 {
			return "", decimal.Zero, fmt.Errorf("invalid line: %s", line)
		}
		d, err := decimal.NewFromString(arr[1])
		if err != nil {
			return "", decimal.Zero, fmt.Errorf("invalid decimal: %s", arr[1])
		}
		return arr[0], d, nil
	}
	return LoadFileCache(path, lineParser)
}

func CachedU128BitsToFix(ctx context.Context, cache *FileCache[string, decimal.Decimal], num *big.Int) (decimal.Decimal, error) {
	if cache == nil {
		return U128BitsToFix(ctx, num)
	}

	return cache.Load(num.String(), func(k string) (decimal.Decimal, error) {
		return instance.U128BitsToFix(ctx, num)
	}, func(k string, v decimal.Decimal) string {
		return fmt.Sprintf("%s:%s", k, v.String())
	})
}
