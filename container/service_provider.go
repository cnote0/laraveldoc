package container

// ServiceProvider 服务提供者接口
//
// ServiceProvider 定义了服务注册和引导的标准接口。
// 服务提供者是向容器注册服务的主要方式，支持延迟加载和服务发现。
//
// 生命周期：
// 1. Register 阶段：注册所有服务绑定
// 2. Boot 阶段：执行服务初始化和配置
//
// 使用示例：
//
//	type DatabaseServiceProvider struct {
//		config *DatabaseConfig
//	}
//
//	func (p *DatabaseServiceProvider) Register(container Container) error {
//		// 注册数据库服务
//		return container.Singleton("database", func(c Container) interface{} {
//			return &Database{
//				Host:     p.config.Host,
//				Port:     p.config.Port,
//				Database: p.config.Database,
//			}
//		})
//	}
//
//	func (p *DatabaseServiceProvider) Boot(container Container) error {
//		// 初始化数据库连接
//		db := container.MustMake("database").(*Database)
//		return db.Connect()
//	}
//
//	func (p *DatabaseServiceProvider) Provides() []string {
//		return []string{"database", "db.connection", "db.migrator"}
//	}
//
//	func (p *DatabaseServiceProvider) IsDeferred() bool {
//		return false // 数据库服务应该立即加载
//	}
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
