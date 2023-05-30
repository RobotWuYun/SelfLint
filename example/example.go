package example

import (
	"context"
	"github.com/gogf/gf/net/gtrace"
)

func myLog(format string, args ...interface{}) {
	ctx := context.Background()
	gtrace.NewSpan(ctx, "myLog")
	leafFunc(format, args...)
}

// ignore trace
func leafFunc(format string, args ...interface{}) {
	// do something
}
