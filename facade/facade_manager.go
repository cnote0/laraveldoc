package facade

// FacadeManager 门面管理器接口
//
// FacadeManager 负责注册、管理和解析门面实例。
// 提供统一的门面访问入口和生命周期管理。
//
// 使用示例：
//
//	manager := NewFacadeManager()
//
//	// 注册门面
//	manager.Register("Mail", &MailFacade{})
//	manager.Register("Queue", &QueueFacade{})
//	manager.Register("Cache", &CacheFacade{})
//
//	// 批量注册
//	facades := map[string]Facade{
//		"User":    &UserFacade{},
//		"Product": &ProductFacade{},
//		"Order":   &OrderFacade{},
//	}
//	manager.RegisterBatch(facades)
//
//	// 使用门面
//	mail := manager.MustResolve("Mail").(*MailFacade)
//	mail.Send("user@example.com", "Welcome!", template)
type FacadeManager interface {
	// Register 注册门面
	//
	// 将门面实例注册到管理器中，使其可以通过名称访问。
	//
	// 参数：
	//   name   - 门面名称（通常是简短的标识符）
	//   facade - 门面实例
	//
	// 示例：
	//   manager.Register("Mail", &MailFacade{})
	//   manager.Register("Queue", &QueueFacade{})
	Register(name string, facade Facade) error

	// RegisterBatch 批量注册门面
	//
	// 一次性注册多个门面，常用于应用初始化。
	//
	// 示例：
	//   facades := map[string]Facade{
	//       "User":    &UserFacade{},
	//       "Product": &ProductFacade{},
	//       "Order":   &OrderFacade{},
	//   }
	//   manager.RegisterBatch(facades)
	RegisterBatch(facades map[string]Facade) error

	// Resolve 解析门面
	//
	// 通过名称获取已注册的门面实例。
	//
	// 示例：
	//   mail, err := manager.Resolve("Mail")
	//   if err != nil {
	//       return err
	//   }
	//   mailFacade := mail.(*MailFacade)
	Resolve(name string) (Facade, error)

	// Has 检查门面是否存在
	//
	// 检查指定名称的门面是否已注册。
	//
	// 示例：
	//   if manager.Has("Redis") {
	//       redis := manager.MustResolve("Redis")
	//       // 使用 Redis 门面
	//   }
	Has(name string) bool

	// MustResolve 强制解析门面
	//
	// 类似 Resolve，但解析失败时会 panic。
	//
	// 示例：
	//   log := manager.MustResolve("Log").(*LogFacade)
	//   log.Info("Application started")
	MustResolve(name string) Facade

	// Remove 移除门面
	//
	// 从管理器中移除指定的门面。
	//
	// 示例：
	//   manager.Remove("DeprecatedService")
	Remove(name string) error

	// Clear 清空所有门面
	//
	// 移除所有已注册的门面，主要用于测试清理。
	//
	// 示例：
	//   defer manager.Clear() // 测试后清理
	Clear()

	// GetAll 获取所有门面
	//
	// 返回所有已注册的门面映射。
	//
	// 示例：
	//   facades := manager.GetAll()
	//   for name, facade := range facades {
	//       fmt.Printf("Facade: %s, Type: %T\n", name, facade)
	//   }
	GetAll() map[string]Facade

	// SetContainer 设置容器
	//
	// 设置门面管理器使用的依赖注入容器。
	//
	// 示例：
	//   manager.SetContainer(appContainer)
	SetContainer(container interface{})

	// GetContainer 获取容器
	//
	// 返回当前使用的容器实例。
	//
	// 示例：
	//   container := manager.GetContainer()
	GetContainer() interface{}
}
