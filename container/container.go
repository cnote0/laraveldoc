// Package container 提供 Laravel 风格的 IoC 容器协议定义
//
// 本包实现了依赖注入容器的核心功能，包括服务绑定、解析、生命周期管理等。
// 设计基于 Laravel 的 IoC 容器模式，提供强类型的 Go 语言接口定义。
//
// 主要特性：
// - 服务绑定和解析
// - 单例模式支持
// - 上下文绑定
// - 依赖注入
// - 服务提供者模式
//
// 使用示例：
//
//	// 创建容器
//	container := NewContainer()
//
//	// 绑定服务
//	container.Bind("database", func(c Container) interface{} {
//		return &Database{Host: "localhost"}
//	})
//
//	// 单例绑定
//	container.Singleton("logger", func(c Container) interface{} {
//		return &Logger{}
//	})
//
//	// 解析服务
//	db, err := container.Make("database")
//	logger := container.MustMake("logger").(*Logger)
package container

import (
	"context"
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
	//       router.Use(middleware.(Middleware))
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

// ServiceProvider 服务提供者接口
//
// ServiceProvider 是 Laravel 应用程序的引导中心。所有核心服务都通过服务提供者绑定到容器中。
// 服务提供者定义了服务的注册和引导逻辑。
//
// 使用示例：
//
//	type DatabaseServiceProvider struct{}
//
//	func (p *DatabaseServiceProvider) Register(container Container) error {
//		// 注册数据库服务
//		return container.Singleton("database", func(c Container) interface{} {
//			config := c.MustMake("config").(*Config)
//			return &Database{
//				Host:     config.Database.Host,
//				Port:     config.Database.Port,
//				Username: config.Database.Username,
//				Password: config.Database.Password,
//			}
//		})
//	}
//
//	func (p *DatabaseServiceProvider) Boot(container Container) error {
//		// 服务引导逻辑
//		db := container.MustMake("database").(*Database)
//		return db.Connect()
//	}
//
//	// 注册服务提供者
//	app.Register(&DatabaseServiceProvider{})
type ServiceProvider interface {
	// Register 注册服务到容器
	//
	// 在此方法中应该只进行服务绑定，不应该访问其他服务，
	// 因为其他服务可能还未注册。
	//
	// 示例：
	//   func (p *CacheServiceProvider) Register(container Container) error {
	//       return container.Bind("cache", func(c Container) interface{} {
	//           return &RedisCache{Host: "localhost:6379"}
	//       }, true)
	//   }
	Register(container Container) error

	// Boot 引导服务
	//
	// 所有服务提供者的 Register 方法都被调用之后，才会调用 Boot 方法。
	// 在此方法中可以安全地访问其他已注册的服务。
	//
	// 示例：
	//   func (p *EventServiceProvider) Boot(container Container) error {
	//       dispatcher := container.MustMake("events").(*EventDispatcher)
	//       logger := container.MustMake("logger").(*Logger)
	//
	//       // 注册事件监听器
	//       dispatcher.Listen("user.created", func(event interface{}) {
	//           logger.Info("User created", event)
	//       })
	//
	//       return nil
	//   }
	Boot(container Container) error

	// Provides 返回此提供者提供的服务列表
	//
	// 用于延迟加载和服务发现。
	//
	// 示例：
	//   func (p *MailServiceProvider) Provides() []string {
	//       return []string{"mailer", "mail.manager", "mail.transport"}
	//   }
	Provides() []string

	// IsDeferred 返回此提供者是否为延迟提供者
	//
	// 延迟提供者只有在其提供的服务被请求时才会被加载。
	//
	// 示例：
	//   func (p *PaymentServiceProvider) IsDeferred() bool {
	//       return true // 只有在需要支付服务时才加载
	//   }
	IsDeferred() bool
}

// Resolver 依赖解析器接口
//
// Resolver 负责解析构造函数或方法的依赖关系，实现自动依赖注入。
//
// 使用示例：
//
//	type UserController struct {
//		userRepo *UserRepository
//		logger   *Logger
//	}
//
//	func NewUserController(userRepo *UserRepository, logger *Logger) *UserController {
//		return &UserController{
//			userRepo: userRepo,
//			logger:   logger,
//		}
//	}
//
//	// 解析器会自动注入依赖
//	controller, err := resolver.ResolveConstructor(reflect.TypeOf(&UserController{}), container)
type Resolver interface {
	// ResolveConstructor 解析构造函数依赖
	//
	// 分析构造函数的参数类型，从容器中解析对应的依赖并调用构造函数。
	//
	// 示例：
	//   type Service struct {
	//       db     *Database
	//       cache  *Cache
	//       logger *Logger
	//   }
	//
	//   func NewService(db *Database, cache *Cache, logger *Logger) *Service {
	//       return &Service{db: db, cache: cache, logger: logger}
	//   }
	//
	//   // 自动解析依赖并创建实例
	//   service, err := resolver.ResolveConstructor(reflect.TypeOf(&Service{}), container)
	ResolveConstructor(concrete reflect.Type, container Container) (interface{}, error)

	// ResolveMethod 解析方法依赖
	//
	// 分析方法的参数类型，从容器中解析对应的依赖并调用方法。
	//
	// 示例：
	//   func (s *UserService) SendWelcomeEmail(user *User, mailer *Mailer, logger *Logger) error {
	//       return mailer.Send(user.Email, "Welcome!", "welcome.html")
	//   }
	//
	//   // 自动注入依赖并调用方法
	//   result, err := resolver.ResolveMethod(service, "SendWelcomeEmail", container, params)
	ResolveMethod(instance interface{}, method string, container Container, parameters map[string]interface{}) ([]interface{}, error)

	// ResolveClosure 解析闭包依赖
	//
	// 分析闭包的参数类型，从容器中解析对应的依赖并调用闭包。
	//
	// 示例：
	//   handler := func(db *Database, cache *Cache) error {
	//       // 处理逻辑
	//       return nil
	//   }
	//
	//   result, err := resolver.ResolveClosure(handler, container, nil)
	ResolveClosure(closure interface{}, container Container, parameters map[string]interface{}) ([]interface{}, error)

	// GetDependencies 获取依赖列表
	//
	// 分析给定类型或函数的依赖关系，返回依赖类型列表。
	//
	// 示例：
	//   dependencies := resolver.GetDependencies(reflect.TypeOf(NewUserService))
	//   // dependencies: [*Database, *Logger, *Cache]
	GetDependencies(signature reflect.Type) ([]reflect.Type, error)
}

// ContextualBinding 上下文绑定接口
//
// ContextualBinding 允许根据使用上下文来绑定不同的服务实现。
// 这对于解决同一接口的不同实现场景非常有用。
//
// 使用示例：
//
//	// 为不同的控制器绑定不同的缓存实现
//	container.When("UserController").Needs("Cache").Give("RedisCache")
//	container.When("AdminController").Needs("Cache").Give("MemoryCache")
//
//	// 为特定类型绑定工厂函数
//	container.When("PaymentService").Needs("Gateway").Give(func(c Container) interface{} {
//		return &StripeGateway{ApiKey: os.Getenv("STRIPE_KEY")}
//	})
type ContextualBinding interface {
	// Needs 指定需要的依赖
	//
	// 参数：
	//   abstract - 依赖的抽象标识（通常是接口名或服务名）
	//
	// 示例：
	//   container.When("OrderService").Needs("PaymentGateway")
	Needs(abstract interface{}) ContextualBinding

	// Give 提供具体实现
	//
	// 参数：
	//   implementation - 具体实现（可以是服务名、工厂函数或实例）
	//
	// 示例：
	//   // 绑定到已注册的服务
	//   container.When("OrderService").Needs("PaymentGateway").Give("StripeGateway")
	//
	//   // 绑定到工厂函数
	//   container.When("TestService").Needs("Database").Give(func(c Container) interface{} {
	//       return &MockDatabase{}
	//   })
	//
	//   // 绑定到实例
	//   container.When("DevService").Needs("Logger").Give(&DebugLogger{Level: "debug"})
	Give(implementation interface{}) error

	// GiveTagged 绑定到带标签的服务集合
	//
	// 参数：
	//   tag - 服务标签
	//
	// 示例：
	//   container.When("NotificationService").Needs("Channels").GiveTagged("notification.channels")
	//
	//   // 此时 NotificationService 会接收到所有标记为 "notification.channels" 的服务
	GiveTagged(tag string) error

	// GiveConfig 根据配置绑定
	//
	// 参数：
	//   configKey - 配置键名
	//
	// 示例：
	//   container.When("DatabaseService").Needs("Driver").GiveConfig("database.default")
	//
	//   // 会根据配置项 "database.default" 的值来绑定相应的驱动
	GiveConfig(configKey string) error
}

// Binding 绑定信息结构
//
// Binding 包含了服务绑定的所有信息，包括具体实现、是否共享、上下文等。
//
// 使用示例：
//
//	binding := Binding{
//		Concrete: func(c Container) interface{} {
//			return &UserService{}
//		},
//		Shared: true,
//		Context: map[string]interface{}{
//			"scope": "request",
//			"tags":  []string{"service", "user"},
//		},
//	}
type Binding struct {
	// Concrete 具体实现
	//
	// 可以是：
	// - 工厂函数：func(Container) interface{}
	// - 结构体类型：reflect.Type
	// - 实例：任何对象
	//
	// 示例：
	//   // 工厂函数
	//   Concrete: func(c Container) interface{} {
	//       return &UserService{
	//           DB: c.MustMake("database").(*Database),
	//       }
	//   }
	//
	//   // 结构体类型
	//   Concrete: reflect.TypeOf(&UserService{})
	//
	//   // 实例
	//   Concrete: &UserService{}
	Concrete interface{}

	// Shared 是否为共享实例（单例）
	//
	// true:  在整个应用生命周期中只创建一个实例
	// false: 每次解析都创建新实例
	//
	// 示例：
	//   Shared: true  // 数据库连接应该是单例
	//   Shared: false // HTTP 请求对象应该是新实例
	Shared bool

	// Context 上下文信息
	//
	// 存储绑定的元数据，如标签、作用域、配置等。
	//
	// 示例：
	//   Context: map[string]interface{}{
	//       "tags":        []string{"cache", "redis"},
	//       "scope":       "singleton",
	//       "description": "Redis cache implementation",
	//       "version":     "1.0.0",
	//   }
	Context map[string]interface{}

	// Dependencies 依赖列表
	//
	// 记录此服务依赖的其他服务，用于依赖分析和循环依赖检测。
	//
	// 示例：
	//   Dependencies: []string{"database", "logger", "cache"}
	Dependencies []string

	// ResolvedAt 解析时间
	//
	// 记录服务首次解析的时间，用于调试和性能分析。
	ResolvedAt *context.Context

	// Alias 别名列表
	//
	// 此服务的所有别名。
	//
	// 示例：
	//   Alias: []string{"db", "database", "illuminate.database"}
	Alias []string
}
