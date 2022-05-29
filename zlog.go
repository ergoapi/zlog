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

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var (
	// Log log
	Log *Logger
)

type Logger struct {
	*zap.Logger

	CtxKey string
}

// InitZlog 初始化日志
func InitZlog(cfg *Config) {
	zlog := zap.New(cfg.getCores()).WithOptions(cfg.debugMode()...)
	Log = &Logger{}
	Log.Logger = zlog
	defer Log.Logger.Sync()
}

func GetLogger() *Logger {
	if Log == nil {
		fmt.Println("zlog not init")
		return nil
	}
	return Log
}

func (l *Logger) GetCtx(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.Logger
}

func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.Logger
}

func (l *Logger) AddCtx(ctx context.Context, field ...zap.Field) (context.Context, *zap.Logger) {
	log := l.With(field...)
	ctx = context.WithValue(ctx, l.CtxKey, log)
	return ctx, log
}
