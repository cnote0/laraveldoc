package database

import "reflect"

// Migrator 数据库迁移接口
//
// Migrator 提供了数据库结构迁移的所有功能。
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
