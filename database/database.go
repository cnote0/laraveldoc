// Package database 提供基于 GORM 的 Laravel 风格数据库访问层协议定义
//
// 本包基于 gorm.io/gorm 设计，提供完整的 ORM 功能和数据库访问接口。
// 结合 Laravel 的 Eloquent ORM 设计模式，提供强类型的数据库操作。
//
// 主要特性：
// - 基于 GORM 的强类型数据库接口
// - Laravel 风格的 Eloquent 模型系统
// - 完整的关联关系支持
// - 数据库迁移和种子系统
// - 查询构建器和原生SQL支持
// - 事务管理和连接池
// - 软删除和模型工厂
//
// 使用示例：
//
//	// 基础 GORM 操作
//	db := container.MustMake("database").(database.DB)
//
//	// 创建记录
//	user := &User{Name: "John", Email: "john@example.com"}
//	db.Create(user)
//
//	// 查询操作
//	var users []User
//	db.Where("active = ?", true).Find(&users)
//
//	// Laravel 风格的 Eloquent 查询
//	eloquent := container.MustMake("eloquent").(database.EloquentBuilder)
//	users := eloquent.Model(&User{}).Where("age", ">", 18).Get()
//
//	// 关联查询
//	users := eloquent.Model(&User{}).With("Profile", "Orders").Get()
//
//	// 数据库迁移
//	migrator := container.MustMake("migrator").(database.Migrator)
//	migrator.AutoMigrate(&User{}, &Profile{}, &Order{})
package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"time"
)

// LoggerInterface 数据库日志接口
//
// LoggerInterface 提供数据库操作的日志记录功能，与 GORM 的日志系统兼容。
//
// 使用示例：
//
//	type CustomLogger struct {
//		level string
//	}
//
//	func (l *CustomLogger) LogMode(level string) LoggerInterface {
//		return &CustomLogger{level: level}
//	}
//
//	func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
//		log.Printf("[INFO] %s %v", msg, data)
//	}
//
//	func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
//		log.Printf("[WARN] %s %v", msg, data)
//	}
//
//	func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
//		log.Printf("[ERROR] %s %v", msg, data)
//	}
//
//	func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
//		elapsed := time.Since(begin)
//		sql, rows := fc()
//		log.Printf("[TRACE] %s [%v rows] took %v, error: %v", sql, rows, elapsed, err)
//	}
//
//	// 使用自定义日志器
//	config := &DatabaseConfig{
//		Logger: &CustomLogger{level: "info"},
//	}
type LoggerInterface interface {
	// LogMode 设置日志级别
	//
	// 参数：
	//   level - 日志级别 ("silent", "error", "warn", "info")
	//
	// 返回：
	//   LoggerInterface - 新的日志器实例
	//
	// 示例：
	//   logger := logger.LogMode("info")
	//   debugLogger := logger.LogMode("debug")
	LogMode(level string) LoggerInterface

	// Info 记录信息日志
	//
	// 用于记录一般信息，如配置信息、状态信息等。
	//
	// 示例：
	//   logger.Info(ctx, "Database connected successfully")
	//   logger.Info(ctx, "Using connection pool with %d max connections", maxConn)
	Info(ctx context.Context, msg string, data ...interface{})

	// Warn 记录警告日志
	//
	// 用于记录可能需要注意但不会导致错误的情况。
	//
	// 示例：
	//   logger.Warn(ctx, "Slow query detected")
	//   logger.Warn(ctx, "Connection pool nearly exhausted: %d/%d", used, max)
	Warn(ctx context.Context, msg string, data ...interface{})

	// Error 记录错误日志
	//
	// 用于记录错误信息，如SQL执行失败、连接问题等。
	//
	// 示例：
	//   logger.Error(ctx, "Failed to execute query: %v", err)
	//   logger.Error(ctx, "Database connection lost")
	Error(ctx context.Context, msg string, data ...interface{})

	// Trace 记录SQL执行轨迹
	//
	// 记录SQL语句的执行情况，包括执行时间、影响行数、错误信息等。
	// 这是GORM日志系统的核心方法。
	//
	// 参数：
	//   ctx           - 上下文
	//   begin         - 执行开始时间
	//   fc            - 获取SQL和影响行数的函数
	//   err           - 执行错误（如果有）
	//
	// 示例：
	//   logger.Trace(ctx, time.Now(), func() (string, int64) {
	//       return "SELECT * FROM users WHERE active = ?", 10
	//   }, nil)
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

// DB 核心数据库接口，基于 GORM 的 DB 结构
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

// SessionConfig 会话配置
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

// Model GORM 基础模型结构
type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// DeletedAt 软删除字段类型
type DeletedAt struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan 实现 sql.Scanner 接口
func (dt *DeletedAt) Scan(value interface{}) error {
	if value == nil {
		dt.Time, dt.Valid = time.Time{}, false
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		dt.Time, dt.Valid = v, true
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		dt.Time, dt.Valid = t, true
	}
	return nil
}

// Value 实现 driver.Valuer 接口
func (dt DeletedAt) Value() (driver.Value, error) {
	if !dt.Valid {
		return nil, nil
	}
	return dt.Time, nil
}

// MarshalJSON 实现 JSON 序列化
func (dt DeletedAt) MarshalJSON() ([]byte, error) {
	if !dt.Valid {
		return []byte("null"), nil
	}
	return dt.Time.MarshalJSON()
}

// UnmarshalJSON 实现 JSON 反序列化
func (dt *DeletedAt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		dt.Valid = false
		return nil
	}
	err := dt.Time.UnmarshalJSON(data)
	dt.Valid = err == nil
	return err
}

// Association 关联操作接口
type Association interface {
	Find(out interface{}, conds ...interface{}) error
	Append(values ...interface{}) error
	Replace(values ...interface{}) error
	Delete(values ...interface{}) error
	Clear() error
	Count() int64
}

// Migrator 数据库迁移接口
type Migrator interface {
	// 自动迁移
	AutoMigrate(dst ...interface{}) error

	// 数据库操作
	CurrentDatabase() string

	// 表操作
	CreateTable(dst ...interface{}) error
	DropTable(dst ...interface{}) error
	HasTable(dst interface{}) bool
	RenameTable(oldName, newName interface{}) error
	GetTables() ([]string, error)

	// 列操作
	AddColumn(dst interface{}, field string) error
	DropColumn(dst interface{}, field string) error
	AlterColumn(dst interface{}, field string) error
	HasColumn(dst interface{}, field string) bool
	RenameColumn(dst interface{}, oldName, field string) error
	ColumnTypes(dst interface{}) ([]ColumnType, error)

	// 视图操作
	CreateView(name string, option ViewOption) error
	DropView(name string) error

	// 约束操作
	CreateConstraint(dst interface{}, name string) error
	DropConstraint(dst interface{}, name string) error
	HasConstraint(dst interface{}, name string) bool

	// 索引操作
	CreateIndex(dst interface{}, name string) error
	DropIndex(dst interface{}, name string) error
	HasIndex(dst interface{}, name string) bool
	RenameIndex(dst interface{}, oldName, newName string) error
}

// ColumnType 列类型信息接口
type ColumnType interface {
	Name() string
	DatabaseTypeName() string
	Length() (length int64, ok bool)
	DecimalSize() (precision int64, scale int64, ok bool)
	Nullable() (nullable bool, ok bool)
	Unique() (unique bool, ok bool)
	ScanType() reflect.Type
	Comment() (value string, ok bool)
	DefaultValue() (value string, ok bool)
}

// ViewOption 视图创建选项
type ViewOption struct {
	Replace     bool
	CheckOption string
	Query       DB
}

// QueryBuilder 链式查询构建器
type QueryBuilder interface {
	// 条件构建
	Where(query interface{}, args ...interface{}) QueryBuilder
	Or(query interface{}, args ...interface{}) QueryBuilder
	Not(query interface{}, args ...interface{}) QueryBuilder

	// 字段选择
	Select(fields ...string) QueryBuilder
	Omit(fields ...string) QueryBuilder
	Distinct(args ...interface{}) QueryBuilder

	// 关联查询
	Joins(query string, args ...interface{}) QueryBuilder
	Preload(query string, args ...interface{}) QueryBuilder

	// 分组和排序
	Group(fields ...string) QueryBuilder
	Having(query interface{}, args ...interface{}) QueryBuilder
	Order(value interface{}) QueryBuilder

	// 分页
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder

	// 执行查询
	Find(dest interface{}) error
	First(dest interface{}) error
	Last(dest interface{}) error
	Take(dest interface{}) error
	Count(count *int64) error
	Pluck(column string, dest interface{}) error

	// 批量查询
	FindInBatches(dest interface{}, batchSize int, fc func(tx DB, batch int) error) error

	// 原生查询
	Raw(sql string, values ...interface{}) QueryBuilder
	Scan(dest interface{}) error

	// 作用域
	Scopes(funcs ...func(QueryBuilder) QueryBuilder) QueryBuilder
}

// Clause 查询子句接口
type Clause interface {
	Name() string
	Build(builder ClauseBuilder)
	MergeClause(clause *Clause)
}

// ClauseBuilder 子句构建器
type ClauseBuilder interface {
	WriteQuoted(field interface{})
	WriteByte(b byte) error
	WriteString(str string) (int, error)
	AddVar(vars ...interface{})
}

// Statement 查询语句构建
type Statement interface {
	Build(clauses ...string)
	Parse(value interface{}) error
	AddClause(clause Clause)
	Quote(field interface{}) string
	AddVar(vars ...interface{})
	SetContext(ctx context.Context)
	GetContext() context.Context
}

// EloquentModel Eloquent 模型接口
type EloquentModel interface {
	// 模型属性
	GetTable() string
	GetKey() interface{}
	GetKeyName() string
	GetAttributes() map[string]interface{}
	GetOriginal() map[string]interface{}
	GetDirty() map[string]interface{}
	GetChanges() map[string]interface{}

	// 属性操作
	GetAttribute(key string) interface{}
	SetAttribute(key string, value interface{})
	Fill(attributes map[string]interface{}) EloquentModel

	// 状态检查
	Exists() bool
	WasRecentlyCreated() bool
	IsDirty(attributes ...string) bool
	IsClean(attributes ...string) bool

	// 持久化操作
	Save(ctx context.Context) error
	Update(ctx context.Context, attributes map[string]interface{}) error
	Delete(ctx context.Context) error
	ForceDelete(ctx context.Context) error
	Restore(ctx context.Context) error

	// 关联操作
	Load(ctx context.Context, relations ...string) error
	LoadMissing(ctx context.Context, relations ...string) error

	// 时间戳
	Touch(ctx context.Context, attributes ...string) error
	GetCreatedAtColumn() string
	GetUpdatedAtColumn() string

	// 序列化
	ToMap() map[string]interface{}
	ToJSON() ([]byte, error)
	FromJSON(data []byte) error
}

// EloquentBuilder Eloquent 查询构建器
type EloquentBuilder interface {
	QueryBuilder

	// 模型实例
	Make() EloquentModel
	Create(attributes map[string]interface{}) (EloquentModel, error)
	FirstOrCreate(attributes map[string]interface{}, values ...map[string]interface{}) (EloquentModel, error)
	FirstOrNew(attributes map[string]interface{}, values ...map[string]interface{}) EloquentModel
	UpdateOrCreate(attributes map[string]interface{}, values map[string]interface{}) (EloquentModel, error)

	// 更新和删除
	UpdateModel(values map[string]interface{}) error
	DeleteModel() error
	ForceDelete() error
	Restore() error

	// 软删除
	WithTrashed() EloquentBuilder
	OnlyTrashed() EloquentBuilder
	WithoutTrashed() EloquentBuilder

	// 关联查询
	With(relations ...string) EloquentBuilder
	WithCount(relations ...string) EloquentBuilder
	Has(relation string, operator string, count int) EloquentBuilder
	WhereHas(relation string, callback func(QueryBuilder) QueryBuilder) EloquentBuilder

	// 关联关系
	BelongsTo(related EloquentModel, foreignKey ...string) BelongsTo
	HasOne(related EloquentModel, foreignKey ...string) HasOne
	HasMany(related EloquentModel, foreignKey ...string) HasMany
	BelongsToMany(related EloquentModel, table string, foreignPivotKey, relatedPivotKey string) BelongsToMany

	// 作用域
	Scope(name string, callback func(EloquentBuilder) EloquentBuilder) EloquentBuilder
	Global(callback func(EloquentBuilder) EloquentBuilder) EloquentBuilder
}

// Relationship 关联关系基础接口
type Relationship interface {
	AddConstraints()
	AddEagerConstraints(models []EloquentModel)
	InitRelation(models []EloquentModel, relation string) []EloquentModel
	Match(models []EloquentModel, results []EloquentModel, relation string) []EloquentModel
	GetResults() interface{}
}

// BelongsTo 属于关联
type BelongsTo interface {
	Relationship
	Associate(model EloquentModel) EloquentModel
	Dissociate() EloquentModel
	GetForeignKeyName() string
	GetOwnerKeyName() string
}

// HasOne 一对一关联
type HasOne interface {
	Relationship
	Make(attributes map[string]interface{}) EloquentModel
	Create(attributes map[string]interface{}) EloquentModel
	Save(model EloquentModel) EloquentModel
	Update(attributes map[string]interface{}) error
	GetForeignKeyName() string
	GetLocalKeyName() string
}

// HasMany 一对多关联
type HasMany interface {
	Relationship
	Make(attributes map[string]interface{}) EloquentModel
	Create(attributes map[string]interface{}) EloquentModel
	CreateMany(records []map[string]interface{}) []EloquentModel
	Save(model EloquentModel) EloquentModel
	SaveMany(models []EloquentModel) []EloquentModel
	Update(attributes map[string]interface{}) error
	Delete() error
	GetForeignKeyName() string
	GetLocalKeyName() string
}

// BelongsToMany 多对多关联
type BelongsToMany interface {
	Relationship
	Attach(id interface{}, attributes ...map[string]interface{}) error
	Detach(ids ...interface{}) error
	Sync(ids []interface{}) error
	SyncWithoutDetaching(ids []interface{}) error
	Toggle(ids []interface{}) error
	UpdateExistingPivot(id interface{}, attributes map[string]interface{}) error
	WherePivot(column string, operator string, value interface{}) BelongsToMany
	WithPivot(columns ...string) BelongsToMany
	WithTimestamps() BelongsToMany
	GetTable() string
	GetForeignPivotKeyName() string
	GetRelatedPivotKeyName() string
}

// Migration 数据库迁移接口
type Migration interface {
	Up(schema SchemaBuilder) error
	Down(schema SchemaBuilder) error
	GetName() string
	GetConnection() string
}

// SchemaBuilder 结构构建器接口
type SchemaBuilder interface {
	Create(table string, callback func(Blueprint)) error
	Table(table string, callback func(Blueprint)) error
	Rename(from, to string) error
	Drop(table string) error
	DropIfExists(table string) error
	HasTable(table string) bool
	HasColumn(table, column string) bool
	GetColumnListing(table string) ([]string, error)
	GetColumnType(table, column string) (string, error)
	Connection(name string) SchemaBuilder
}

// Blueprint 表结构蓝图
type Blueprint interface {
	// 主键和索引
	ID(column ...string) ColumnDefinition
	Primary(columns ...string) IndexDefinition
	Index(columns ...string) IndexDefinition
	UniqueIndex(columns ...string) IndexDefinition
	SpatialIndex(columns ...string) IndexDefinition
	ForeignKey(columns ...string) ForeignKeyDefinition

	// 基础列类型
	String(column string, length ...int) ColumnDefinition
	Text(column string, length ...string) ColumnDefinition
	Integer(column string, autoIncrement ...bool) ColumnDefinition
	BigInteger(column string, autoIncrement ...bool) ColumnDefinition
	SmallInteger(column string, autoIncrement ...bool) ColumnDefinition
	TinyInteger(column string, autoIncrement ...bool) ColumnDefinition
	UnsignedInteger(column string, autoIncrement ...bool) ColumnDefinition
	UnsignedBigInteger(column string, autoIncrement ...bool) ColumnDefinition
	UnsignedSmallInteger(column string, autoIncrement ...bool) ColumnDefinition
	UnsignedTinyInteger(column string, autoIncrement ...bool) ColumnDefinition

	// 浮点数类型
	Float(column string, precision, scale int) ColumnDefinition
	Double(column string, precision, scale int) ColumnDefinition
	Decimal(column string, precision, scale int) ColumnDefinition

	// 布尔和日期类型
	Boolean(column string) ColumnDefinition
	Date(column string) ColumnDefinition
	DateTime(column string, precision ...int) ColumnDefinition
	Time(column string, precision ...int) ColumnDefinition
	Timestamp(column string, precision ...int) ColumnDefinition

	// JSON 和二进制类型
	JSON(column string) ColumnDefinition
	JSONB(column string) ColumnDefinition
	Binary(column string) ColumnDefinition

	// UUID 和枚举类型
	UUID(column string) ColumnDefinition
	Enum(column string, values []string) ColumnDefinition

	// 时间戳
	Timestamps(precision ...int)
	SoftDeletes(column ...string)
	RememberToken()

	// 修改列
	DropColumn(columns ...string)
	RenameColumn(from, to string)
	DropIndex(indexName string)
	DropPrimary()
	DropForeign(indexName string)

	// 表选项
	Engine(engine string) Blueprint
	Charset(charset string) Blueprint
	Collation(collation string) Blueprint
	Comment(comment string) Blueprint
}

// ColumnDefinition 列定义接口
type ColumnDefinition interface {
	// 约束
	Nullable(nullable ...bool) ColumnDefinition
	Default(value interface{}) ColumnDefinition
	Unique() ColumnDefinition
	Primary() ColumnDefinition
	Index(indexName ...string) ColumnDefinition
	Comment(comment string) ColumnDefinition

	// 数值约束
	Unsigned() ColumnDefinition
	AutoIncrement() ColumnDefinition

	// 字符串约束
	Charset(charset string) ColumnDefinition
	Collation(collation string) ColumnDefinition

	// 时间约束
	UseCurrent() ColumnDefinition
	UseCurrentOnUpdate() ColumnDefinition

	// 位置
	After(column string) ColumnDefinition
	First() ColumnDefinition

	// 修改操作
	Change() ColumnDefinition

	// 外键
	References(column string) ForeignKeyDefinition
}

// IndexDefinition 索引定义接口
type IndexDefinition interface {
	Name(name string) IndexDefinition
	Unique() IndexDefinition
	Using(algorithm string) IndexDefinition
	Comment(comment string) IndexDefinition
}

// ForeignKeyDefinition 外键定义接口
type ForeignKeyDefinition interface {
	On(table string) ForeignKeyDefinition
	OnDelete(action string) ForeignKeyDefinition
	OnUpdate(action string) ForeignKeyDefinition
	Name(name string) ForeignKeyDefinition
}

// Seeder 数据库种子接口
type Seeder interface {
	Run(ctx context.Context) error
	GetName() string
	GetConnection() string
}

// Factory 模型工厂接口
type Factory interface {
	Make(attributes ...map[string]interface{}) EloquentModel
	Create(attributes ...map[string]interface{}) (EloquentModel, error)
	CreateMany(count int, attributes ...map[string]interface{}) ([]EloquentModel, error)
	State(state string) Factory
	Count(count int) Factory
	Sequence(callback func(int) map[string]interface{}) Factory
}

// DatabaseManager 数据库管理器接口
type DatabaseManager interface {
	Connection(name ...string) DB
	Reconnect(name ...string) DB
	Disconnect(name ...string) error
	GetConnections() map[string]DB
	GetDefaultConnection() string
	SetDefaultConnection(name string)
	Extend(name string, resolver func(DatabaseConfig) DB) DatabaseManager

	// 事务管理
	Transaction(callback func(DB) error, attempts ...int) error
	BeginTransaction() DB

	// 迁移管理
	GetMigrator() Migrator
	RunMigrations(path string) error
	RollbackMigrations(steps ...int) error
	GetMigrationRepository() MigrationRepository

	// 种子管理
	RunSeeders(seeders ...Seeder) error
	GetSeederRepository() SeederRepository
}

// MigrationRepository 迁移仓库接口
type MigrationRepository interface {
	GetRan() ([]string, error)
	GetMigrations(steps int) ([]Migration, error)
	GetMigrationBatches() (map[string]int, error)
	Log(file string, batch int) error
	Delete(migration string) error
	GetNextBatchNumber() (int, error)
	CreateRepository() error
	RepositoryExists() bool
}

// SeederRepository 种子仓库接口
type SeederRepository interface {
	GetRan() ([]string, error)
	Log(seeder string) error
	Delete(seeder string) error
	CreateRepository() error
	RepositoryExists() bool
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver                                   string            `json:"driver"`
	Host                                     string            `json:"host"`
	Port                                     int               `json:"port"`
	Database                                 string            `json:"database"`
	Username                                 string            `json:"username"`
	Password                                 string            `json:"password"`
	Charset                                  string            `json:"charset"`
	Collation                                string            `json:"collation"`
	Prefix                                   string            `json:"prefix"`
	Timezone                                 string            `json:"timezone"`
	DSN                                      string            `json:"dsn"`
	MaxOpenConnections                       int               `json:"max_open_connections"`
	MaxIdleConnections                       int               `json:"max_idle_connections"`
	MaxConnectionLifetime                    time.Duration     `json:"max_connection_lifetime"`
	SkipDefaultTransaction                   bool              `json:"skip_default_transaction"`
	DisableForeignKeyConstraintWhenMigrating bool              `json:"disable_foreign_key_constraint_when_migrating"`
	Logger                                   LoggerInterface   `json:"-"`
	SlowThreshold                            time.Duration     `json:"slow_threshold"`
	LogLevel                                 string            `json:"log_level"`
	IgnoreRecordNotFoundError                bool              `json:"ignore_record_not_found_error"`
	ParameterizedQueries                     bool              `json:"parameterized_queries"`
	PreparedStatements                       bool              `json:"prepared_statements"`
	DryRun                                   bool              `json:"dry_run"`
	Options                                  map[string]string `json:"options"`
}

// ConnectionPool 连接池接口
type ConnectionPool interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (ConnectionPool, error)
	Close() error
	Stats() sql.DBStats
}

// Plugin 数据库插件接口
type Plugin interface {
	Name() string
	Initialize(db DB) error
}

// Dialector 数据库方言接口
type Dialector interface {
	Name() string
	Initialize(db DB) error
	Migrator(db DB) Migrator
	DataTypeOf(field interface{}) string
	DefaultValueOf(field interface{}) interface{}
	BindVarTo(writer interface{}, stmt Statement, v interface{})
	QuoteTo(writer interface{}, str string)
	Explain(sql string, vars ...interface{}) string
}
