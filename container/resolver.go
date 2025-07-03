package container

import "reflect"

// Resolver 依赖解析器接口
//
// Resolver 负责分析类型依赖关系并自动注入。
// 它可以解析构造函数参数、方法参数和闭包参数的依赖关系。
//
// 使用示例：
//
//	type UserService struct {
//		DB     *Database `inject:"database"`
//		Logger *Logger   `inject:"logger"`
//		Cache  *Cache    `inject:"cache"`
//	}
//
//	func NewUserService(db *Database, logger *Logger, cache *Cache) *UserService {
//		return &UserService{DB: db, Logger: logger, Cache: cache}
//	}
//
//	// 使用解析器自动注入
//	resolver := container.MustMake("resolver").(Resolver)
//	service, err := resolver.ResolveConstructor(reflect.TypeOf(&UserService{}), container)
//	userService := service.(*UserService)
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
