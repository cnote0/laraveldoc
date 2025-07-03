package container

import "context"

// Binding 绑定信息结构体
//
// Binding 存储服务绑定的所有信息，包括具体实现、
// 是否为单例、上下文信息等。
//
// 使用示例：
//
//	binding := &Binding{
//		Concrete: func(c Container) interface{} {
//			return &UserService{
//				DB: c.MustMake("database").(*Database),
//			}
//		},
//		Shared: true,
//		Context: map[string]interface{}{
//			"tags":        []string{"service", "user"},
//			"description": "User management service",
//		},
//		Dependencies: []string{"database", "logger"},
//		Alias:        []string{"user", "user.service"},
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
