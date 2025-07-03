package database

import (
	"context"
	"time"
)

// SessionConfig 会话配置结构体
//
// SessionConfig 用于配置数据库会话的行为，包括事务模式、钩子设置等。
type SessionConfig struct {
	DryRun                   bool
	PrepareStmt              bool
	NewDB                    bool
	SkipHooks                bool
	SkipDefaultTransaction   bool
	DisableNestedTransaction bool
	AllowGlobalUpdate        bool
	FullSaveAssociations     bool
	QueryFields              bool
	Context                  context.Context
	Logger                   LoggerInterface
	NowFunc                  func() time.Time
	CreateBatchSize          int
}
