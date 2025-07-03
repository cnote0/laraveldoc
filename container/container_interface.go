package container

import (
	"reflect"
)

// Container IoC 容器核心接口
//
// Container 是依赖注入容器的核心，负责管理服务的绑定、解析和生命周期。
// 它提供了多种绑定方式和解析策略，支持单例模式和上下文绑定。
//
// 使用示例：
//
//	// 基础绑定
//	container.Bind("config", func(c Container) interface{} {
//		return &Config{Env: "production"}
//	})
//
//	// 单例绑定 - 全局唯一实例
//	container.Singleton("database", func(c Container) interface{} {
//		return &Database{Pool: createPool()}
//	})
//
//	// 实例绑定 - 绑定已存在的实例
//	container.Instance("app", &Application{Name: "MyApp"})
//
//	// 条件绑定 - 基于上下文的绑定
//	container.When("UserController").Needs("Repository").Give("UserRepository")
//
//	// 解析服务
//	config, err := container.Make("config")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 强制解析（panic on error）
//	db := container.MustMake("database").(*Database)
//
//	// 标签绑定和解析
//	container.Tag([]string{"cache.redis", "cache.memory"}, "cache.drivers")
//	drivers := container.Tagged("cache.drivers")
type Container interface {
	// Bind 绑定服务到容器
	//
	// 参数：
	//   abstract - 服务标识符（字符串或类型）
	//   concrete - 服务实现（可以是工厂函数、结构体或实例）
	//   shared   - 是否为单例
	//
	// 示例：
	//   container.Bind("mailer", func(c Container) interface{} {
	//       return &SMTPMailer{Host: "smtp.gmail.com"}
	//   }, false)
	Bind(abstract interface{}, concrete interface{}, shared bool) error

	// BindIf 条件绑定 - 仅在服务未绑定时执行绑定
	//
	// 用于避免重复绑定，特别适用于库或包的初始化。
	//
	// 示例：
	//   container.BindIf("logger", func(c Container) interface{} {
	//       return &DefaultLogger{}
	//   }, true)
	BindIf(abstract interface{}, concrete interface{}, shared bool) error

	// Singleton 绑定单例服务
	//
	// 单例服务在整个应用生命周期中只会被创建一次。
	//
	// 示例：
	//   container.Singleton("database", func(c Container) interface{} {
	//       return &Database{
	//           Connection: createConnection(),
	//           Pool: createPool(),
	//       }
	//   })
	Singleton(abstract interface{}, concrete interface{}) error

	// Instance 绑定已存在的实例
	//
	// 直接绑定一个已创建的实例，该实例将作为单例使用。
	//
	// 示例：
	//   app := &Application{Name: "MyApp", Version: "1.0.0"}
	//   container.Instance("app", app)
	Instance(abstract interface{}, instance interface{}) error

	// Make 解析服务
	//
	// 从容器中解析指定的服务，如果服务未绑定或解析失败，返回错误。
	//
	// 示例：
	//   mailer, err := container.Make("mailer")
	//   if err != nil {
	//       return err
	//   }
	//   m := mailer.(*SMTPMailer)
	Make(abstract interface{}) (interface{}, error)

	// MustMake 强制解析服务
	//
	// 类似 Make，但解析失败时会 panic。适用于确定服务已绑定的场景。
	//
	// 示例：
	//   logger := container.MustMake("logger").(*Logger)
	//   logger.Info("Application started")
	MustMake(abstract interface{}) interface{}

	// MakeWith 带参数解析服务
	//
	// 解析服务时传递额外参数，这些参数可以在工厂函数中使用。
	//
	// 示例：
	//   container.Bind("user.repository", func(c Container, params map[string]interface{}) interface{} {
	//       db := c.MustMake("database").(*Database)
	//       tableName := params["table"].(string)
	//       return &UserRepository{DB: db, Table: tableName}
	//   })
	//
	//   repo, err := container.MakeWith("user.repository", map[string]interface{}{
	//       "table": "users",
	//   })
	MakeWith(abstract interface{}, parameters map[string]interface{}) (interface{}, error)

	// Bound 检查服务是否已绑定
	//
	// 示例：
	//   if container.Bound("cache") {
	//       cache := container.MustMake("cache")
	//       // 使用缓存
	//   }
	Bound(abstract interface{}) bool

	// Resolved 检查服务是否已解析
	//
	// 示例：
	//   if container.Resolved("database") {
	//       // 数据库连接已建立
	//   }
	Resolved(abstract interface{}) bool

	// Alias 为服务创建别名
	//
	// 示例：
	//   container.Bind("illuminate.database", databaseFactory, true)
	//   container.Alias("illuminate.database", "db")
	//
	//   // 现在可以通过别名访问
	//   db := container.MustMake("db").(*Database)
	Alias(abstract interface{}, alias interface{}) error

	// Tag 为服务添加标签
	//
	// 标签允许将相关服务分组，便于批量解析。
	//
	// 示例：
	//   container.Tag([]string{"cache.redis", "cache.file", "cache.memory"}, "cache.stores")
	//
	//   stores := container.Tagged("cache.stores")
	//   for _, store := range stores {
	//       // 初始化每个缓存存储
	//   }
	Tag(abstracts []interface{}, tag string) error

	// Tagged 获取带有指定标签的所有服务
	//
	// 示例：
	//   middlewares := container.Tagged("middleware")
	//   for _, middleware := range middlewares {
	//       // 注册中间件
	//   }
	Tagged(tag string) []interface{}

	// When 开始上下文绑定
	//
	// 上下文绑定允许根据依赖关系的上下文来解析不同的实现。
	//
	// 示例：
	//   // 为 UserController 提供特定的仓库实现
	//   container.When("UserController").Needs("Repository").Give("UserRepository")
	//
	//   // 为 AdminController 提供不同的仓库实现
	//   container.When("AdminController").Needs("Repository").Give("AdminRepository")
	When(concrete interface{}) ContextualBinding

	// Call 调用方法并注入依赖
	//
	// 自动解析方法参数的依赖关系并调用方法。
	//
	// 示例：
	//   type UserService struct{}
	//   func (s *UserService) CreateUser(repo *UserRepository, logger *Logger) error {
	//       // 实现逻辑
	//   }
	//
	//   result, err := container.Call(&UserService{}, "CreateUser", nil)
	Call(instance interface{}, method string, parameters map[string]interface{}) ([]interface{}, error)

	// Build 构建实例
	//
	// 根据给定的类型构建实例，自动注入依赖关系。
	//
	// 示例：
	//   type UserController struct {
	//       Repository *UserRepository `inject:""`
	//       Logger     *Logger         `inject:""`
	//   }
	//
	//   controller, err := container.Build(reflect.TypeOf(&UserController{}))
	Build(concrete reflect.Type) (interface{}, error)

	// Flush 清空容器
	//
	// 清除所有绑定、别名和已解析的实例。主要用于测试。
	//
	// 示例：
	//   defer container.Flush() // 测试后清理
	Flush()

	// GetBindings 获取所有绑定
	GetBindings() map[interface{}]Binding

	// IsShared 检查服务是否为单例
	IsShared(abstract interface{}) bool

	// Extend 扩展已绑定的服务
	//
	// 允许修改已绑定服务的实现，常用于装饰者模式。
	//
	// 示例：
	//   container.Extend("logger", func(service interface{}, c Container) interface{} {
	//       baseLogger := service.(*Logger)
	//       return &TimestampLogger{Logger: baseLogger}
	//   })
	Extend(abstract interface{}, closure func(interface{}, Container) interface{}) error
}
