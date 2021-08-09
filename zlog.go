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

package zlog

import "go.uber.org/zap"

var (
	// Log log
	Log *zap.Logger
	// Zlog log sugar
	Zlog *zap.SugaredLogger
)

// InitZlog 初始化日志
func InitZlog(cfg *Config) {
	Log = zap.New(cfg.getCores()).WithOptions(cfg.debugMode()...)
	Zlog = Log.Sugar()
}

func Debug(f string, args ...interface{}) {
	Zlog.Debugf(f, args...)
}

func Info(f string, args ...interface{}) {
	Zlog.Infof(f, args...)
}

func Warn(f string, args ...interface{}) {
	Zlog.Warnf(f, args...)
}

func Error(f string, args ...interface{}) {
	Zlog.Errorf(f, args...)
}

func Panic(f string, args ...interface{}) {
	Zlog.Panicf(f, args...)
}

func Fatal(f string, args ...interface{}) {
	Zlog.Panicf(f, args...)
}
