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
// 包结构：
// - db_interface.go - DB 核心数据库接口
// - model.go - Model 基础模型和 DeletedAt 软删除结构体
// - logger_interface.go - LoggerInterface 日志接口
// - session_config.go - SessionConfig 会话配置
// - association.go - Association 关联接口
// - migrator.go - Migrator 迁移器接口
// - query_builder.go - QueryBuilder 查询构建器接口
// - eloquent.go - EloquentModel 和 EloquentBuilder 接口
// - relationships.go - 各种关联关系接口
// - migration.go - Migration 和 SchemaBuilder 迁移相关接口
// - factory.go - Factory 工厂接口
// - manager.go - DatabaseManager 数据库管理器接口
// - config.go - DatabaseConfig 配置结构体
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
//
//	// 事务处理
//	err := db.Transaction(func(tx database.DB) error {
//		if err := tx.Create(&user1).Error(); err != nil {
//			return err
//		}
//		if err := tx.Create(&user2).Error(); err != nil {
//			return err
//		}
//		return nil
//	})
//
//	// 软删除模型定义
//	type User struct {
//		database.Model
//		Name  string `json:"name"`
//		Email string `json:"email" gorm:"unique"`
//	}
package database
