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
// 包结构：
// - container_interface.go - Container 核心接口
// - service_provider.go - ServiceProvider 服务提供者接口
// - resolver.go - Resolver 依赖解析器接口
// - contextual_binding.go - ContextualBinding 上下文绑定接口
// - binding.go - Binding 绑定信息结构体
//
// 使用示例：
//
//	// 创建容器
//	container := NewContainer()
//
//	// 绑定服务
//	container.Bind("database", func(c Container) interface{} {
//		return &Database{Host: "localhost"}
//	}, false)
//
//	// 单例绑定
//	container.Singleton("logger", func(c Container) interface{} {
//		return &Logger{}
//	})
//
//	// 解析服务
//	db, err := container.Make("database")
//	logger := container.MustMake("logger").(*Logger)
//
//	// 注册服务提供者
//	provider := &DatabaseServiceProvider{}
//	container.RegisterProvider(provider)
//
//	// 上下文绑定
//	container.When("UserController").Needs("Repository").Give("UserRepository")
//	container.When("AdminController").Needs("Repository").Give("AdminRepository")
//
//	// 标签绑定
//	container.Tag([]string{"cache.redis", "cache.memory"}, "cache.drivers")
//	drivers := container.Tagged("cache.drivers")
package container
