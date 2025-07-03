package container

// ContextualBinding 上下文绑定接口
//
// ContextualBinding 允许根据不同的上下文（如类、方法或接口）
// 为同一个抽象提供不同的实现。这对于依赖注入中的条件绑定非常有用。
//
// 使用示例：
//
//	// 为不同的控制器提供不同的仓库实现
//	container.When("UserController").Needs("Repository").Give("UserRepository")
//	container.When("AdminController").Needs("Repository").Give("AdminRepository")
//
//	// 为测试环境提供模拟实现
//	container.When("*Test").Needs("PaymentGateway").Give(func(c Container) interface{} {
//		return &MockPaymentGateway{}
//	})
//
//	// 基于配置的条件绑定
//	container.When("ProductionService").Needs("Logger").GiveConfig("logging.production")
//	container.When("DebugService").Needs("Logger").GiveConfig("logging.debug")
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
