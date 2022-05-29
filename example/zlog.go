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
	"net/http"

	"github.com/ergoapi/util/exid"
	"github.com/ergoapi/zlog/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	cfg := zlog.Config{
		Simple:      false,
		HookFunc:    []func(entry zapcore.Entry) error{zlog.ExampleHook()},
		WriteLog:    true,
		WriteJSON:   true,
		WriteConfig: zlog.WriteConfig{},
		ServiceName: "example",
	}
	zlog.InitZlog(&cfg)
}

func AddTraceId() gin.HandlerFunc {
	return func(g *gin.Context) {
		traceId := g.GetHeader("traceId")
		if traceId == "" {
			traceId = exid.GenUUID()
		}
		ctx, log := zlog.GetLogger().AddCtx(g.Request.Context(), zap.Any("traceId", traceId))
		g.Request = g.Request.WithContext(ctx)
		log.Info("AddTraceId success")
		g.Next()
	}
}

// curl http://127.0.0.1:8888/test
func main() {
	g := gin.New()
	g.Use(AddTraceId())
	g.GET("/test", func(context *gin.Context) {
		log := zlog.GetLogger().GetCtx(context.Request.Context())
		log.Info("test")
		log.Debug("test")
		context.JSON(200, "success")
	})
	zlog.GetLogger().Info("example success")
	http.ListenAndServe(":8888", g)
}
