//  Copyright (c) 2020. The EFF Team Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

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
