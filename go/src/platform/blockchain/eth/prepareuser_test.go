package eth

import (
	"testing"
	"time"
	"platform/om"
	"fmt"
)

// prepareuser.go 测试
func TestPrepareUser(t *testing.T) {
	ct := time.Now()
	pped := PrepareUser(1)
	om.Bigger(t, pped, 1, fmt.Sprintf("expected %d, prepared %d", 1, pped), "准备1个用户,耗时:", time.Since(ct))
	pped1 := PrepareUser(pped + 1)
	om.Bigger(t, pped1, pped, fmt.Sprintf("expected %d, prepared %d", pped + 1, pped1), "再多准备一个用户,耗时:", time.Since(ct))
	pped2 := PrepareUser(pped1 - 1)
	om.Bigger(t, pped1, pped, fmt.Sprintf("expected %d, prepared %d", pped1 - 1, pped2),
		"需要的用户数量少于实际已经准备的用户数量, 耗时:", time.Since(ct))
}
