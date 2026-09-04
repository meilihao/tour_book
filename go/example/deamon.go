package example

import (
	"log/slog"
	"runtime/debug"
	"time"
)

func Fn() {
	for {
		// do something
	}
}

func Start() {
	Go(Fn, time.Second)
}

// Go 守护进程
func Go(f func(), retryInterval ...time.Duration) {
	d := time.Duration(0)
	if len(retryInterval) != 0 {
		d = retryInterval[0]
	}
	go func() {
		defer func() {
			// 异常退出重启进程
			if r := recover(); r != nil {
				slog.Error("Recovered in Go", "error", r, "stack", debug.Stack())
				if d > 0 {
					time.Sleep(d)
				}
				Go(f, d)
			}
		}()
		f()
	}()
}
