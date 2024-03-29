/*
Copyright 2023 HuanShiTianXin and The ZapLog Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,  // Debug 应是大量的，且通常在生产状态禁用
	"info":   zapcore.InfoLevel,   // Info 是默认的记录优先级
	"warn":   zapcore.WarnLevel,   // Warn 比 info 更重要,但不需要单独的人工审核
	"error":  zapcore.ErrorLevel,  // Error 是高优先级的,如果应用程序运行平稳，则不应生成任何错误级别的日志
	"dpanic": zapcore.DPanicLevel, // DPanic 特别重大的错误，在开发模式下引起 panic
	"panic":  zapcore.PanicLevel,  // Panic 记录消息后调用 panic
	"fatal":  zapcore.FatalLevel,  // Fatal 记录消息后调用 os.Exit(1)
}

// getLevel 获取日志级别
func getLevel(level string) zapcore.Level {
	if level, ok := levelMap[level]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// getEncoderConfig 获取编码器配置
func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	//更改时间编码,将时间编码改为可读的
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return encoderConfig
}

// getFileWriteSyncer 指定日志写入位置
func getFileWriteSyncer(f fileOut) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   f.filename,   // 日志文件的位置
		MaxSize:    f.maxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: f.maxBackups, // 保留旧文件的最大个数
		MaxAge:     f.maxAge,     // 保留旧文件的最大天数
		Compress:   f.compress,   // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// SliceWriter 自定义的 io.Writer，用于将日志输出保存到切片中
type SliceWriter struct {
	Logs *[]string
}

// Write 将日志内容写入 Logs 切片
func (w *SliceWriter) Write(p []byte) (n int, err error) {
	*w.Logs = append(*w.Logs, string(p))
	return len(p), nil
}

// getSliceWriteSyncer 指定日志写入位置为切片
func getSliceWriteSyncer(s *[]string) zapcore.WriteSyncer {
	sliceWriteLogger := &SliceWriter{Logs: s}
	return zapcore.AddSync(sliceWriteLogger)
}
