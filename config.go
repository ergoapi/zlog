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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ergoapi/util/zos"
	"github.com/ergoapi/util/ztime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultMaxSize  = 50 // MB
	defaultBackups  = 3  // 备份
	defaultMaxAge   = 7  // days
	defaultCompress = true
)

// Config 配置
type Config struct {
	Simple      bool                              // 简易模式
	HookFunc    []func(entry zapcore.Entry) error // hook
	WriteLog    bool                              // 写日志
	WriteJSON   bool                              // json
	WriteConfig WriteConfig
	ServiceName string
}

// WriteConfig 写日志配置
type WriteConfig struct {
	LogPath    string // 日志路径
	MaxSize    int    // 日志大小
	MaxBackups int    // 备份
	MaxAge     int    // 天数
}

func (c *Config) debugMode() []zap.Option {
	var cfgopts []zap.Option
	if !c.Simple {
		cfgopts = append(cfgopts, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	}
	if len(c.HookFunc) != 0 {
		cfgopts = append(cfgopts, zap.Hooks(c.HookFunc...))
		// cfgopts = append(cfgopts, zap.Hooks(warnHook()))
	}
	return cfgopts
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getEncoder 文件encoder
func (c *Config) getEncoder(enablejson bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if enablejson {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// logfilesyncer log syncer
func (c *Config) logfilesyncer(lvl ...string) zapcore.WriteSyncer {
	var wss []zapcore.WriteSyncer
	if c.Simple {
		// 不分级别，写一个文件
		wss = append(wss, c.getLogWriter())
	} else {
		// 分级别，写不同文件
		if len(lvl) == 0 {
			lvl = append(lvl, "debug")
		}
		wss = append(wss, c.getLogWriter(lvl[0]))
	}
	return zapcore.NewMultiWriteSyncer(wss...)
}

// consolesyncer console打出日志
func (c *Config) consolesyncer() zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
}

func (c *Config) svcname() string {
	if c.ServiceName == "" || len(c.ServiceName) == 0 {
		return "unknow"
	}
	return c.ServiceName
}

func (c *Config) getLogConfig() *WriteConfig {
	logcfg := c.WriteConfig
	if logcfg.LogPath == "" {
		if zos.IsMacOS() {
			logcfg.LogPath = fmt.Sprintf("/tmp/%v/%v", c.svcname(), ztime.GetToday())
		} else {
			logcfg.LogPath = fmt.Sprintf("/var/log/%v/%v", c.svcname(), ztime.GetToday())
		}
	}
	if logcfg.MaxAge <= defaultMaxAge {
		logcfg.MaxAge = defaultMaxAge
	}
	if logcfg.MaxBackups <= defaultBackups {
		logcfg.MaxBackups = defaultBackups
	}
	if logcfg.MaxSize <= defaultMaxSize {
		logcfg.MaxSize = defaultMaxSize
	}
	return &logcfg
}

func (c *Config) getLogWriter(loglevel ...string) zapcore.WriteSyncer {
	logcfg := c.getLogConfig()
	var logpath string
	if len(loglevel) > 0 {
		logpath = fmt.Sprintf("%v/%v.log", logcfg.LogPath, loglevel[0])
	} else {
		logpath = fmt.Sprintf("%v.log", logcfg.LogPath)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    logcfg.MaxSize,
		MaxBackups: logcfg.MaxBackups,
		MaxAge:     logcfg.MaxAge,
		Compress:   defaultCompress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func (c *Config) getCores() zapcore.Core {
	var cors []zapcore.Core
	// level debug custom err
	debugPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.InfoLevel
	})
	customPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.InfoLevel && level < zap.WarnLevel
	})
	warnPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zap.WarnLevel
	})
	errPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	defaultPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.DebugLevel
	})

	consolecoder := c.getEncoder(false)
	consoleCore := zapcore.NewCore(consolecoder, zapcore.NewMultiWriteSyncer(c.consolesyncer()), defaultPriority)
	cors = append(cors, consoleCore)
	if c.WriteLog {
		// 输出文件
		filecoder := c.getEncoder(c.WriteJSON)
		warnCore := zapcore.NewCore(filecoder, zapcore.NewMultiWriteSyncer(c.logfilesyncer("warn")), warnPriority)
		if c.Simple {
			simpleCore := zapcore.NewCore(filecoder, zapcore.NewMultiWriteSyncer(c.logfilesyncer("default")), defaultPriority)
			cors = append(cors, warnCore, simpleCore)
		} else {
			debugCore := zapcore.NewCore(filecoder, zapcore.NewMultiWriteSyncer(c.logfilesyncer("debug")), debugPriority)
			customCore := zapcore.NewCore(filecoder, zapcore.NewMultiWriteSyncer(c.logfilesyncer("custom")), customPriority)
			errCore := zapcore.NewCore(filecoder, zapcore.NewMultiWriteSyncer(c.logfilesyncer("err")), errPriority)
			cors = append(cors, warnCore, debugCore, customCore, errCore)
		}
	}
	return zapcore.NewTee(cors...)
}

func ExampleHook() func(entry zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		if entry.Level < zapcore.WarnLevel {
			return nil
		}
		log.Printf("example hook msg: %v", entry.Message)
		return nil
	}
}
