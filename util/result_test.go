package util

import (
	"fmt"
	"testing"
)

func TestResult_GetRootCause(t *testing.T) {
	res := MsgError("test", "this is error msg")
	res = res.With("outlayer 1")
	res.Code = -2
	res = res.With("outest")
	res.Code = -3
	fmt.Println(JsonStr(&res))

	root := res.GetRootCause()
	fmt.Println(JsonStr(&root))

}
