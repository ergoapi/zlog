## zlog

```
go get -u github.com/ergoapi/zlog
```

## usage

```go
package main

import (
	"github.com/ergoapi/zlog"
	"go.uber.org/zap/zapcore"
)

func init() {
	cfg := zlog.Config{
		Simple:      true,
		HookFunc:    []func(entry zapcore.Entry) error{zlog.ExampleHook()},
		WriteLog:    false,
		WriteJSON:   false,
		WriteConfig: zlog.WriteConfig{},
		ServiceName: "example",
	}
	zlog.InitZlog(&cfg)
}

func main() {
	zlog.Debug("debug")
	zlog.Error("err")
	zlog.Fatal("fatal")
}
```