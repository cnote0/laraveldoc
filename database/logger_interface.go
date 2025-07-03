package database

import (
	"context"
	"time"
)

// LoggerInterface 数据库日志接口
//
// LoggerInterface 提供数据库操作的日志记录功能，与 GORM 的日志系统兼容。
type LoggerInterface interface {
	// LogMode 设置日志级别
	LogMode(level string) LoggerInterface

	// Info 记录信息日志
	Info(ctx context.Context, msg string, data ...interface{})

	// Warn 记录警告日志
	Warn(ctx context.Context, msg string, data ...interface{})

	// Error 记录错误日志
	Error(ctx context.Context, msg string, data ...interface{})

	// Trace 记录SQL执行轨迹
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
