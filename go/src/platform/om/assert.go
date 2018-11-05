package om
// Testing helpers for doozer.

import (
	"github.com/kr/pretty"
	"reflect"
	"testing"
	"runtime"
	"fmt"
)

func assert(t *testing.T, result bool, f func(), cd int) {
	if !result {
		_, file, line, _ := runtime.Caller(cd + 1)
		t.Logf("%s:%d", file, line)
		f()
	}
}

func equal(t *testing.T, exp, got interface{}, cd int, args ...interface{}) {
	fn := func() {
		for _, desc := range pretty.Diff(exp, got) {
			t.Logf(":", desc)
		}
		if len(args) > 0 {
			t.Logf(":", " -", fmt.Sprint(args...))
		}
	}
	result := reflect.DeepEqual(exp, got)
	if result {
		t.Log("Test OK:", args)
	}
	assert(t, result, fn, cd+1)
}

func Equal(t *testing.T, exp, got interface{}, args ...interface{}) {
	equal(t, exp, got, 1, args...)
}

func NotEqual(t *testing.T, exp, got interface{}, args ...interface{}) {
	fn := func() {
		t.Logf("!  Unexpected: <%#v>", exp)
		if len(args) > 0 {
			t.Logf("!", " -", fmt.Sprint(args...))
		}
	}
	result := !reflect.DeepEqual(exp, got)
	if result {
		t.Log("Test OK:", args)
	}
	assert(t, result, fn, 1)
}

func Bigger(t *testing.T, cmp1, cmp2 int, args ...interface{}) {
	fn := func() {
		t.Logf("Expect %d > %d ,isn't true", cmp1, cmp2)
	}
	result := cmp1 > cmp2
	if result {
		t.Log("Test OK:", args)
	}
	assert(t, result, fn, 1)
}
