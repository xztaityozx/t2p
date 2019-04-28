package nutils

import (
	"fmt"
	"testing"
)

func TestIntMax(t *testing.T) {
	type TestData struct {
		desc   string
		in     []int
		expect int
	}
	tds := []TestData{
		{desc: "整数値のみ", in: []int{1, 2, 3}, expect: 3},
		{desc: "負数あり", in: []int{1, 2, -3}, expect: 2},
	}
	for _, v := range tds {
		got := IntMax(v.in...)
		if v.expect == got {
			t.Log(fmt.Sprintf("[OK] %s:", v.desc))
		} else {
			errMsg := fmt.Sprintf("expect = %d, got = %d\n", v.expect, got)
			t.Error(fmt.Sprintf("[NG] %s:\n%s", v.desc, errMsg))
		}
	}
}
