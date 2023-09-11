// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LogLevels = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"fatal": zap.FatalLevel,
	"panic": zap.PanicLevel,
}

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
}

// wrapper delegates all calls to the underlying zap.Logger
type wrapper struct {
	logger *zap.Logger
}

// Debug logs an debug msg with fields
func (l wrapper) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info msg with fields
func (l wrapper) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs an error msg with fields
func (l wrapper) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l wrapper) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l wrapper) With(fields ...zapcore.Field) Logger {
	return wrapper{logger: l.logger.With(fields...)}
}

func Start() (zap.AtomicLevel, Factory, error) {

	atom := zap.NewAtomicLevel()
	zaplog, err := zap.Config{
		Encoding:    "json",
		Level:       atom,
		OutputPaths: []string{"stdout"}, // You can change this to a file path if needed
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	}.Build()
	if err != nil {
		return zap.AtomicLevel{}, Factory{}, err
	}
	factory := NewFactory(zaplog)
	return atom, factory, nil
}
