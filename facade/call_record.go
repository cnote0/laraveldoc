package facade

import (
	"context"
	"time"
)

// CallRecord 调用记录结构体
//
// CallRecord 记录了方法调用的详细信息，用于间谍模式和调用追踪。
type CallRecord struct {
	// Method 方法名
	Method string

	// Args 调用参数
	Args []interface{}

	// Result 返回值
	Result []interface{}

	// Timestamp 调用时间
	Timestamp time.Time

	// Context 调用上下文
	Context context.Context

	// Error 调用错误（如果有）
	Error error

	// Duration 执行时长
	Duration time.Duration

	// Metadata 元数据
	Metadata map[string]interface{}
}
