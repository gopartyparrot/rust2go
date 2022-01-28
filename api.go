package rust2go

import "fmt"

func ParseU64F64(u128 string) (string, error) {
	return callStart("parse-u64f64", fmt.Sprintf("--str=%s", u128))
}
