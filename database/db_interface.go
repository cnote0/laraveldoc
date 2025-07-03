package database

import (
	"context"
	"database/sql"
)

// DB 核心数据库接口，基于 GORM 的 DB 结构
//
// DB 是数据库操作的核心接口，提供了完整的 CRUD 操作、事务管理、
// 查询构建等功能。设计完全兼容 GORM 的 API。
type DB interface {
	// 数据库连接管理
	WithContext(ctx context.Context) DB
	Session(config *SessionConfig) DB
	Debug() DB
	DryRun() DB

	// 模型操作
	Model(value interface{}) DB
	Table(name string, args ...interface{}) DB
	Select(query interface{}, args ...interface{}) DB
	Omit(columns ...string) DB
	Where(query interface{}, args ...interface{}) DB
	Or(query interface{}, args ...interface{}) DB
	Not(query interface{}, args ...interface{}) DB

	// 创建操作
	Create(value interface{}) DB
	CreateInBatches(value interface{}, batchSize int) DB
	Save(value interface{}) DB

	// 查询操作
	Find(dest interface{}, conds ...interface{}) DB
	FindInBatches(dest interface{}, batchSize int, fc func(tx DB, batch int) error) DB
	First(dest interface{}, conds ...interface{}) DB
	Last(dest interface{}, conds ...interface{}) DB
	Take(dest interface{}, conds ...interface{}) DB
	FirstOrInit(dest interface{}, conds ...interface{}) DB
	FirstOrCreate(dest interface{}, conds ...interface{}) DB

	// 更新操作
	Update(column string, value interface{}) DB
	Updates(values interface{}) DB
	UpdateColumn(column string, value interface{}) DB
	UpdateColumns(values interface{}) DB

	// 删除操作
	Delete(value interface{}, conds ...interface{}) DB
	Unscoped() DB

	// 聚合操作
	Count(count *int64) DB
	Distinct(args ...interface{}) DB
	Group(name string) DB
	Having(query interface{}, args ...interface{}) DB
	Joins(query string, args ...interface{}) DB
	Preload(query string, args ...interface{}) DB

	// 分页和排序
	Limit(limit int) DB
	Offset(offset int) DB
	Order(value interface{}) DB

	// 原生查询
	Raw(sql string, values ...interface{}) DB
	Exec(sql string, values ...interface{}) DB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) DB
	ScanRows(rows *sql.Rows, dest interface{}) error
	Pluck(column string, dest interface{}) DB

	// 事务处理
	Begin(opts ...*sql.TxOptions) DB
	Commit() DB
	Rollback() DB
	SavePoint(name string) DB
	RollbackTo(name string) DB
	Transaction(fc func(tx DB) error, opts ...*sql.TxOptions) error

	// 关联操作
	Association(column string) Association

	// 数据库迁移
	AutoMigrate(dst ...interface{}) error
	Migrator() Migrator

	// 作用域和实例管理
	Scopes(funcs ...func(DB) DB) DB
	Attrs(attrs ...interface{}) DB
	Assign(attrs ...interface{}) DB

	// 设置和获取
	Set(key string, value interface{}) DB
	Get(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) DB
	InstanceGet(key string) (interface{}, bool)

	// 错误处理
	AddError(err error) error
	GetErrors() []error
	Error() error
	RowsAffected() int64

	// 数据库连接
	SqlDB() (*sql.DB, error)
	Close() error
}
