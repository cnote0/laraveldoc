package laraveldoc

import (
	"database/sql"
)

// DatabaseManager 数据库管理器接口
type DatabaseManager interface {
	// Connection 获取数据库连接
	Connection(name string) ConnectionInterface

	// Reconnect 重新连接
	Reconnect(name string) ConnectionInterface

	// Disconnect 断开连接
	Disconnect(name string) error

	// GetDefaultConnection 获取默认连接名
	GetDefaultConnection() string

	// SetDefaultConnection 设置默认连接名
	SetDefaultConnection(name string)

	// Extend 扩展连接驱动
	Extend(name string, resolver func(map[string]interface{}, string) ConnectionInterface) DatabaseManager
}

// ConnectionInterface 数据库连接接口
type ConnectionInterface interface {
	// Table 创建查询构建器
	Table(table string) QueryBuilder

	// Select 执行select查询
	Select(query string, bindings []interface{}) ([]map[string]interface{}, error)

	// Insert 执行insert查询
	Insert(query string, bindings []interface{}) (sql.Result, error)

	// Update 执行update查询
	Update(query string, bindings []interface{}) (sql.Result, error)

	// Delete 执行delete查询
	Delete(query string, bindings []interface{}) (sql.Result, error)

	// Statement 执行语句
	Statement(query string, bindings []interface{}) (bool, error)

	// Transaction 事务处理
	Transaction(callback func(ConnectionInterface) error) error

	// BeginTransaction 开始事务
	BeginTransaction() error

	// Commit 提交事务
	Commit() error

	// Rollback 回滚事务
	Rollback() error

	// GetSchemaBuilder 获取模式构建器
	GetSchemaBuilder() SchemaBuilder

	// GetQueryGrammar 获取查询语法
	GetQueryGrammar() QueryGrammar

	// GetSchemaGrammar 获取模式语法
	GetSchemaGrammar() SchemaGrammar
}

// QueryBuilder 查询构建器接口
type QueryBuilder interface {
	// Select 选择字段
	Select(columns ...string) QueryBuilder

	// From 设置表名
	From(table string) QueryBuilder

	// Where 添加where条件
	Where(column string, operator string, value interface{}) QueryBuilder

	// OrWhere 添加or where条件
	OrWhere(column string, operator string, value interface{}) QueryBuilder

	// WhereIn 添加where in条件
	WhereIn(column string, values []interface{}) QueryBuilder

	// Join 添加join
	Join(table string, first string, operator string, second string) QueryBuilder

	// LeftJoin 添加left join
	LeftJoin(table string, first string, operator string, second string) QueryBuilder

	// OrderBy 添加排序
	OrderBy(column string, direction string) QueryBuilder

	// GroupBy 添加分组
	GroupBy(columns ...string) QueryBuilder

	// Having 添加having条件
	Having(column string, operator string, value interface{}) QueryBuilder

	// Limit 设置限制
	Limit(value int) QueryBuilder

	// Offset 设置偏移
	Offset(value int) QueryBuilder

	// Get 获取结果
	Get() ([]map[string]interface{}, error)

	// First 获取第一条记录
	First() (map[string]interface{}, error)

	// Count 获取计数
	Count() (int64, error)

	// Exists 检查是否存在
	Exists() (bool, error)

	// Insert 插入数据
	Insert(values map[string]interface{}) (sql.Result, error)

	// Update 更新数据
	Update(values map[string]interface{}) (sql.Result, error)

	// Delete 删除数据
	Delete() (sql.Result, error)

	// ToSQL 转换为SQL
	ToSQL() (string, []interface{})
}

// SchemaBuilder 模式构建器接口
type SchemaBuilder interface {
	// HasTable 检查表是否存在
	HasTable(table string) (bool, error)

	// HasColumn 检查列是否存在
	HasColumn(table string, column string) (bool, error)

	// GetColumnListing 获取列列表
	GetColumnListing(table string) ([]string, error)

	// Create 创建表
	Create(table string, callback func(Blueprint) error) error

	// Table 修改表
	Table(table string, callback func(Blueprint) error) error

	// Drop 删除表
	Drop(table string) error

	// DropIfExists 如果存在则删除表
	DropIfExists(table string) error

	// Rename 重命名表
	Rename(from string, to string) error
}

// Blueprint 表结构蓝图接口
type Blueprint interface {
	// ID 创建ID字段
	ID(column string) ColumnDefinition

	// String 创建字符串字段
	String(column string, length int) ColumnDefinition

	// Text 创建文本字段
	Text(column string) ColumnDefinition

	// Integer 创建整数字段
	Integer(column string) ColumnDefinition

	// BigInteger 创建大整数字段
	BigInteger(column string) ColumnDefinition

	// Float 创建浮点数字段
	Float(column string, total int, places int) ColumnDefinition

	// Boolean 创建布尔字段
	Boolean(column string) ColumnDefinition

	// Date 创建日期字段
	Date(column string) ColumnDefinition

	// DateTime 创建日期时间字段
	DateTime(column string) ColumnDefinition

	// Timestamp 创建时间戳字段
	Timestamp(column string) ColumnDefinition

	// Timestamps 创建创建和更新时间戳
	Timestamps() error

	// Primary 设置主键
	Primary(columns ...string) error

	// Index 创建索引
	Index(columns []string, name string) error

	// Unique 创建唯一索引
	Unique(columns []string, name string) error

	// Foreign 创建外键
	Foreign(columns []string) ForeignKeyDefinition

	// DropColumn 删除列
	DropColumn(columns ...string) error

	// DropPrimary 删除主键
	DropPrimary(index string) error

	// DropIndex 删除索引
	DropIndex(index string) error

	// DropUnique 删除唯一索引
	DropUnique(index string) error

	// DropForeign 删除外键
	DropForeign(index string) error
}

// ColumnDefinition 列定义接口
type ColumnDefinition interface {
	// Nullable 设置为可空
	Nullable() ColumnDefinition

	// Default 设置默认值
	Default(value interface{}) ColumnDefinition

	// Unsigned 设置为无符号
	Unsigned() ColumnDefinition

	// AutoIncrement 设置自动递增
	AutoIncrement() ColumnDefinition

	// Comment 设置注释
	Comment(comment string) ColumnDefinition

	// After 设置在某列之后
	After(column string) ColumnDefinition

	// First 设置为第一列
	First() ColumnDefinition
}

// ForeignKeyDefinition 外键定义接口
type ForeignKeyDefinition interface {
	// References 设置引用列
	References(columns ...string) ForeignKeyDefinition

	// On 设置引用表
	On(table string) ForeignKeyDefinition

	// OnDelete 设置删除时动作
	OnDelete(action string) ForeignKeyDefinition

	// OnUpdate 设置更新时动作
	OnUpdate(action string) ForeignKeyDefinition

	// Name 设置外键名称
	Name(name string) ForeignKeyDefinition
}

// QueryGrammar 查询语法接口
type QueryGrammar interface {
	// CompileSelect 编译select查询
	CompileSelect(query QueryBuilder) string

	// CompileInsert 编译insert查询
	CompileInsert(query QueryBuilder, values map[string]interface{}) string

	// CompileUpdate 编译update查询
	CompileUpdate(query QueryBuilder, values map[string]interface{}) string

	// CompileDelete 编译delete查询
	CompileDelete(query QueryBuilder) string

	// Parameter 参数化值
	Parameter(value interface{}) string

	// Wrap 包装标识符
	Wrap(value string) string
}

// SchemaGrammar 模式语法接口
type SchemaGrammar interface {
	// CompileCreate 编译创建表语句
	CompileCreate(blueprint Blueprint, command interface{}) string

	// CompileAdd 编译添加列语句
	CompileAdd(blueprint Blueprint, command interface{}) string

	// CompileDrop 编译删除表语句
	CompileDrop(blueprint Blueprint, command interface{}) string

	// CompileDropColumn 编译删除列语句
	CompileDropColumn(blueprint Blueprint, command interface{}) string

	// GetType 获取列类型
	GetType(column ColumnDefinition) string
}

// Model 模型接口
type Model interface {
	// GetTable 获取表名
	GetTable() string

	// SetTable 设置表名
	SetTable(table string) Model

	// GetKey 获取主键值
	GetKey() interface{}

	// GetKeyName 获取主键名
	GetKeyName() string

	// GetAttributes 获取属性
	GetAttributes() map[string]interface{}

	// SetAttribute 设置属性
	SetAttribute(key string, value interface{}) Model

	// GetAttribute 获取属性
	GetAttribute(key string) interface{}

	// Save 保存模型
	Save() error

	// Delete 删除模型
	Delete() error

	// Fresh 刷新模型
	Fresh() (Model, error)

	// Refresh 重新加载模型
	Refresh() Model

	// NewQuery 创建新查询
	NewQuery() QueryBuilder

	// GetConnection 获取连接
	GetConnection() ConnectionInterface
}

// EloquentBuilder Eloquent查询构建器接口
type EloquentBuilder interface {
	QueryBuilder

	// Find 根据主键查找
	Find(id interface{}) (Model, error)

	// FindMany 根据主键查找多个
	FindMany(ids []interface{}) ([]Model, error)

	// FindOrFail 查找或失败
	FindOrFail(id interface{}) (Model, error)

	// FirstOrFail 获取第一个或失败
	FirstOrFail() (Model, error)

	// Create 创建模型
	Create(attributes map[string]interface{}) (Model, error)

	// ForceCreate 强制创建模型
	ForceCreate(attributes map[string]interface{}) (Model, error)

	// UpdateOrCreate 更新或创建
	UpdateOrCreate(attributes map[string]interface{}, values map[string]interface{}) (Model, error)

	// FirstOrCreate 获取第一个或创建
	FirstOrCreate(attributes map[string]interface{}, values map[string]interface{}) (Model, error)

	// FirstOrNew 获取第一个或新建
	FirstOrNew(attributes map[string]interface{}, values map[string]interface{}) (Model, error)

	// With 预加载关联
	With(relations ...string) EloquentBuilder

	// WithCount 预加载关联计数
	WithCount(relations ...string) EloquentBuilder

	// Has 检查关联
	Has(relation string, operator string, count int) EloquentBuilder

	// WhereHas 根据关联条件查询
	WhereHas(relation string, callback func(QueryBuilder) QueryBuilder) EloquentBuilder

	// GetModel 获取模型实例
	GetModel() Model

	// SetModel 设置模型实例
	SetModel(model Model) EloquentBuilder
}

// Migration 迁移接口
type Migration interface {
	// Up 执行迁移
	Up() error

	// Down 回滚迁移
	Down() error

	// GetConnection 获取连接名
	GetConnection() string
}

// Migrator 迁移器接口
type Migrator interface {
	// Run 运行迁移
	Run(paths []string, options map[string]interface{}) error

	// Rollback 回滚迁移
	Rollback(paths []string, options map[string]interface{}) error

	// Reset 重置迁移
	Reset(paths []string, pretend bool) error

	// Refresh 刷新迁移
	Refresh(paths []string, options map[string]interface{}) error

	// Status 获取迁移状态
	Status(paths []string) ([]MigrationStatus, error)

	// GetRepository 获取迁移仓库
	GetRepository() MigrationRepository
}

// MigrationRepository 迁移仓库接口
type MigrationRepository interface {
	// GetRan 获取已运行的迁移
	GetRan() ([]string, error)

	// GetMigrations 获取迁移批次
	GetMigrations(steps int) ([]Migration, error)

	// GetLast 获取最后一批迁移
	GetLast() ([]Migration, error)

	// GetMigrationBatches 获取迁移批次映射
	GetMigrationBatches() (map[string]int, error)

	// Log 记录迁移
	Log(file string, batch int) error

	// Delete 删除迁移记录
	Delete(migration interface{}) error

	// GetNextBatchNumber 获取下一个批次号
	GetNextBatchNumber() (int, error)

	// CreateRepository 创建迁移表
	CreateRepository() error

	// RepositoryExists 检查迁移表是否存在
	RepositoryExists() (bool, error)
}

// MigrationStatus 迁移状态
type MigrationStatus struct {
	// Migration 迁移名称
	Migration string

	// Batch 批次号
	Batch int

	// Status 状态
	Status string
}

// Seeder 种子接口
type Seeder interface {
	// Run 运行种子
	Run() error
}
