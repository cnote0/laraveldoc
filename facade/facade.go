// Package facade 提供 Laravel 风格的门面（Facade）模式协议定义
//
// 门面模式为复杂的子系统提供一个简单的接口，使得访问底层服务变得更加简洁。
// 本包实现了静态方法调用、实时门面、模拟支持等功能。
//
// 主要特性：
// - 静态方法调用支持
// - 门面自动解析
// - 实时门面（Real-time Facades）
// - 测试模拟支持
// - 门面中间件
// - 门面管理器
//
// 使用示例：
//
//	// 注册门面
//	facade.Register("DB", "database")
//	facade.Register("Cache", "cache")
//	facade.Register("Log", "logger")
//
//	// 使用门面（通过生成的静态方法）
//	users := DB.Table("users").Where("active", true).Get()
//	Cache.Put("key", "value", time.Hour)
//	Log.Info("User login", userID)
//
//	// 实时门面
//	userService := facade.Real(&UserService{})
//	userService.CreateUser(userData)
package facade

import (
	"context"
	"reflect"
	"time"
)

// Facade 门面核心接口
//
// Facade 提供了访问底层服务的简化接口，通过静态方法调用的方式
// 隐藏了复杂的服务解析和依赖注入逻辑。
//
// 使用示例：
//
//	// 定义数据库门面
//	type DatabaseFacade struct{}
//
//	func (f *DatabaseFacade) GetFacadeAccessor() string {
//		return "database" // 容器中的服务标识
//	}
//
//	func (f *DatabaseFacade) ShouldPreventStaleExecution() bool {
//		return false // 是否阻止陈旧执行
//	}
//
//	// 使用门面
//	DB := &DatabaseFacade{}
//	users := DB.Table("users").Get() // 通过生成的代理方法调用
//
//	// 注册门面
//	manager.Register("DB", DB)
type Facade interface {
	// GetFacadeAccessor 获取门面访问器
	//
	// 返回在容器中注册的服务标识符，门面将通过此标识符
	// 从容器中解析实际的服务实例。
	//
	// 示例：
	//   func (f *CacheFacade) GetFacadeAccessor() string {
	//       return "cache.store" // 对应容器中的缓存服务
	//   }
	GetFacadeAccessor() string

	// ShouldPreventStaleExecution 是否阻止陈旧执行
	//
	// 当返回 true 时，门面会在每次调用时重新解析服务实例，
	// 确保总是使用最新的服务。适用于在测试中经常更换实现的场景。
	//
	// 示例：
	//   func (f *TestFacade) ShouldPreventStaleExecution() bool {
	//       return true // 测试环境下总是使用最新实例
	//   }
	ShouldPreventStaleExecution() bool

	// GetFacadeRoot 获取门面根实例
	//
	// 从容器中解析并返回实际的服务实例。此方法通常由
	// 门面管理器调用，不需要手动调用。
	//
	// 示例：
	//   instance, err := facade.GetFacadeRoot()
	//   if err != nil {
	//       return err
	//   }
	//   database := instance.(*Database)
	GetFacadeRoot() (interface{}, error)

	// ClearResolvedInstance 清除已解析的实例
	//
	// 清除缓存的服务实例，下次调用时会重新从容器解析。
	// 主要用于测试或需要重置服务状态的场景。
	//
	// 示例：
	//   facade.ClearResolvedInstance() // 清除缓存
	//   newInstance := facade.GetFacadeRoot() // 重新解析
	ClearResolvedInstance()

	// SetFacadeApplication 设置应用实例
	//
	// 设置门面使用的应用容器实例，门面将从此容器中解析服务。
	//
	// 示例：
	//   facade.SetFacadeApplication(app)
	SetFacadeApplication(app interface{})

	// GetFacadeApplication 获取应用实例
	//
	// 返回门面当前使用的应用容器实例。
	//
	// 示例：
	//   app := facade.GetFacadeApplication()
	//   container := app.(Container)
	GetFacadeApplication() interface{}
}

// StaticFacade 静态门面接口
//
// StaticFacade 支持类似静态方法的调用方式，提供更简洁的API。
// 通过代码生成或反射机制，可以实现直接调用底层服务方法的效果。
//
// 使用示例：
//
//	// 生成的静态门面代码
//	type DBFacade struct {
//		*StaticFacadeImpl
//	}
//
//	func (f *DBFacade) Table(name string) QueryBuilder {
//		result := f.CallMethod("Table", []interface{}{name})
//		return result[0].(QueryBuilder)
//	}
//
//	func (f *DBFacade) Select(columns ...string) QueryBuilder {
//		result := f.CallMethod("Select", []interface{}{columns})
//		return result[0].(QueryBuilder)
//	}
//
//	// 使用方式
//	DB := &DBFacade{}
//	users := DB.Table("users").Select("id", "name").Where("active", true).Get()
type StaticFacade interface {
	Facade

	// CallMethod 调用底层服务的方法
	//
	// 通过反射调用底层服务的指定方法，并传递参数。
	// 这是静态门面实现的核心方法。
	//
	// 参数：
	//   methodName - 要调用的方法名
	//   args      - 方法参数列表
	//
	// 返回：
	//   []interface{} - 方法返回值列表
	//   error        - 调用错误
	//
	// 示例：
	//   // 调用数据库的 table 方法
	//   result, err := facade.CallMethod("Table", []interface{}{"users"})
	//   if err != nil {
	//       return err
	//   }
	//   queryBuilder := result[0].(QueryBuilder)
	CallMethod(methodName string, args []interface{}) ([]interface{}, error)

	// CallMethodWithContext 带上下文调用方法
	//
	// 类似 CallMethod，但支持传递上下文信息。
	//
	// 示例：
	//   ctx := context.WithTimeout(context.Background(), 5*time.Second)
	//   result, err := facade.CallMethodWithContext(ctx, "LongRunningQuery", args)
	CallMethodWithContext(ctx context.Context, methodName string, args []interface{}) ([]interface{}, error)

	// HasMethod 检查方法是否存在
	//
	// 检查底层服务是否有指定的方法。
	//
	// 示例：
	//   if facade.HasMethod("Transaction") {
	//       // 支持事务操作
	//   }
	HasMethod(methodName string) bool

	// GetMethodSignature 获取方法签名
	//
	// 返回指定方法的签名信息，包括参数类型和返回类型。
	//
	// 示例：
	//   signature, err := facade.GetMethodSignature("Create")
	//   // signature: func(data map[string]interface{}) (*Model, error)
	GetMethodSignature(methodName string) (reflect.Type, error)
}

// FacadeManager 门面管理器接口
//
// FacadeManager 负责管理应用中的所有门面，提供门面的注册、解析和生命周期管理。
//
// 使用示例：
//
//	manager := NewFacadeManager(container)
//
//	// 注册门面
//	manager.Register("DB", &DatabaseFacade{})
//	manager.Register("Cache", &CacheFacade{})
//	manager.Register("Log", &LoggerFacade{})
//
//	// 解析门面
//	db, err := manager.Resolve("DB")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 批量注册
//	facades := map[string]Facade{
//		"Auth":    &AuthFacade{},
//		"Config":  &ConfigFacade{},
//		"Event":   &EventFacade{},
//	}
//	manager.RegisterBatch(facades)
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

// RealtimeFacade 实时门面接口
//
// RealtimeFacade 是 Laravel 的实时门面功能，允许将任何类转换为门面，
// 无需预先定义门面类。这提供了极大的灵活性。
//
// 使用示例：
//
//	// 定义普通服务类
//	type UserService struct {
//		repository *UserRepository
//	}
//
//	func (s *UserService) CreateUser(data map[string]interface{}) (*User, error) {
//		// 实现逻辑
//	}
//
//	func (s *UserService) GetUser(id int) (*User, error) {
//		// 实现逻辑
//	}
//
//	// 创建实时门面
//	UserFacade := realtimeFacade.Create(&UserService{})
//
//	// 像门面一样使用
//	user, err := UserFacade.CreateUser(userData)
//	existingUser, err := UserFacade.GetUser(123)
type RealtimeFacade interface {
	// Create 创建实时门面
	//
	// 将任何对象转换为门面，支持调用该对象的所有公共方法。
	//
	// 参数：
	//   target - 目标对象（任何结构体实例）
	//
	// 返回：
	//   interface{} - 门面代理对象
	//
	// 示例：
	//   service := &PaymentService{gateway: stripe}
	//   PaymentFacade := realtimeFacade.Create(service)
	//
	//   // 现在可以像门面一样调用
	//   result := PaymentFacade.ProcessPayment(amount, card)
	Create(target interface{}) interface{}

	// CreateWithContainer 创建带容器的实时门面
	//
	// 创建实时门面，并指定依赖注入容器用于解析依赖。
	//
	// 示例：
	//   service := &OrderService{} // 依赖将通过容器注入
	//   OrderFacade := realtimeFacade.CreateWithContainer(service, container)
	CreateWithContainer(target interface{}, container interface{}) interface{}

	// Wrap 包装现有实例
	//
	// 包装一个已存在的实例，使其具有门面的调用方式。
	//
	// 示例：
	//   logger := app.Make("logger").(*Logger)
	//   LogFacade := realtimeFacade.Wrap(logger)
	//   LogFacade.Info("Message logged via facade")
	Wrap(instance interface{}) interface{}

	// IsRealtimeFacade 检查是否为实时门面
	//
	// 检查给定对象是否为实时门面创建的代理对象。
	//
	// 示例：
	//   if realtimeFacade.IsRealtimeFacade(someObject) {
	//       // 这是一个实时门面
	//   }
	IsRealtimeFacade(obj interface{}) bool

	// GetTarget 获取目标对象
	//
	// 从实时门面代理对象中获取原始的目标对象。
	//
	// 示例：
	//   target := realtimeFacade.GetTarget(facadeProxy)
	//   originalService := target.(*UserService)
	GetTarget(facadeProxy interface{}) interface{}
}

// MockInterface 门面模拟接口
//
// MockInterface 为门面提供测试模拟功能，允许在测试中替换门面的实现。
//
// 使用示例：
//
//	// 在测试中模拟门面
//	func TestUserCreation(t *testing.T) {
//		// 创建模拟对象
//		mockDB := &MockDatabase{}
//		mockDB.On("Create", mock.Anything).Return(&User{ID: 1}, nil)
//
//		// 替换门面实现
//		DB.Mock(mockDB)
//		defer DB.ClearMock()
//
//		// 执行测试
//		user, err := userService.CreateUser(userData)
//		assert.NoError(t, err)
//		assert.Equal(t, 1, user.ID)
//
//		// 验证模拟调用
//		mockDB.AssertCalled(t, "Create", userData)
//	}
type MockInterface interface {
	// Mock 设置模拟实现
	//
	// 将门面的底层实现替换为模拟对象，常用于单元测试。
	//
	// 参数：
	//   mock - 模拟对象，应该实现与原服务相同的接口
	//
	// 示例：
	//   mockMailer := &MockMailer{}
	//   mockMailer.On("Send", mock.Anything).Return(nil)
	//   Mail.Mock(mockMailer)
	Mock(mock interface{})

	// PartialMock 部分模拟
	//
	// 只模拟指定的方法，其他方法仍使用原始实现。
	//
	// 参数：
	//   methods - 要模拟的方法映射
	//
	// 示例：
	//   Cache.PartialMock(map[string]interface{}{
	//       "Get": func(key string) interface{} {
	//           return "mocked_value"
	//       },
	//       "Put": func(key string, value interface{}, ttl time.Duration) error {
	//           return nil
	//       },
	//   })
	PartialMock(methods map[string]interface{})

	// Spy 创建间谍对象
	//
	// 创建间谍对象，可以监控方法调用但不改变行为。
	//
	// 示例：
	//   spy := Log.Spy()
	//
	//   // 执行一些操作
	//   userService.CreateUser(userData)
	//
	//   // 验证日志调用
	//   spy.ShouldHaveReceived("Info").Times(1)
	//   spy.ShouldHaveReceived("Error").Times(0)
	Spy() SpyInterface

	// ClearMock 清除模拟
	//
	// 恢复门面的原始实现，移除所有模拟设置。
	//
	// 示例：
	//   defer facade.ClearMock() // 测试结束后清理
	ClearMock()

	// IsMocked 检查是否已模拟
	//
	// 检查门面是否当前处于模拟状态。
	//
	// 示例：
	//   if DB.IsMocked() {
	//       // 当前使用的是模拟数据库
	//   }
	IsMocked() bool

	// GetMock 获取模拟对象
	//
	// 返回当前设置的模拟对象。
	//
	// 示例：
	//   mock := Cache.GetMock()
	//   mockCache := mock.(*MockCache)
	//   mockCache.AssertExpectations(t)
	GetMock() interface{}
}

// SpyInterface 间谍接口
//
// SpyInterface 用于监控和验证方法调用，而不改变原有行为。
//
// 使用示例：
//
//	spy := Mail.Spy()
//
//	// 执行业务逻辑
//	orderService.ProcessOrder(order)
//
//	// 验证邮件发送
//	spy.ShouldHaveReceived("Send").Times(2)
//	spy.ShouldHaveReceived("Send").With("order-confirmation", order.User.Email)
type SpyInterface interface {
	// ShouldHaveReceived 验证方法调用
	//
	// 验证指定方法是否被调用。
	//
	// 参数：
	//   methodName - 方法名
	//
	// 返回：
	//   CallVerifier - 调用验证器，可以进一步验证调用次数和参数
	//
	// 示例：
	//   spy.ShouldHaveReceived("ProcessPayment")
	//   spy.ShouldHaveReceived("SendEmail").Times(3)
	//   spy.ShouldHaveReceived("LogError").With("payment failed")
	ShouldHaveReceived(methodName string) CallVerifier

	// ShouldNotHaveReceived 验证方法未被调用
	//
	// 验证指定方法不应该被调用。
	//
	// 示例：
	//   spy.ShouldNotHaveReceived("SendSMS") // 应该没有发送短信
	ShouldNotHaveReceived(methodName string) CallVerifier

	// GetCallCount 获取调用次数
	//
	// 返回指定方法的调用次数。
	//
	// 示例：
	//   count := spy.GetCallCount("Log")
	//   assert.Equal(t, 5, count)
	GetCallCount(methodName string) int

	// GetCalls 获取所有调用记录
	//
	// 返回所有方法调用的详细记录。
	//
	// 示例：
	//   calls := spy.GetCalls()
	//   for _, call := range calls {
	//       fmt.Printf("Method: %s, Args: %v\n", call.Method, call.Args)
	//   }
	GetCalls() []CallRecord

	// Reset 重置间谍
	//
	// 清除所有调用记录，重新开始监控。
	//
	// 示例：
	//   spy.Reset() // 清除之前的调用记录
	Reset()
}

// CallVerifier 调用验证器接口
//
// CallVerifier 用于验证方法调用的次数和参数。
type CallVerifier interface {
	// Times 验证调用次数
	//
	// 验证方法被调用的确切次数。
	//
	// 示例：
	//   spy.ShouldHaveReceived("Save").Times(1)
	//   spy.ShouldHaveReceived("Log").Times(0) // 等同于 ShouldNotHaveReceived
	Times(count int) CallVerifier

	// AtLeast 验证最少调用次数
	//
	// 验证方法至少被调用指定次数。
	//
	// 示例：
	//   spy.ShouldHaveReceived("Validate").AtLeast(1)
	AtLeast(count int) CallVerifier

	// AtMost 验证最多调用次数
	//
	// 验证方法最多被调用指定次数。
	//
	// 示例：
	//   spy.ShouldHaveReceived("RetryOperation").AtMost(3)
	AtMost(count int) CallVerifier

	// With 验证调用参数
	//
	// 验证方法是否使用指定参数调用。
	//
	// 示例：
	//   spy.ShouldHaveReceived("SendEmail").With("user@example.com", "Welcome")
	//   spy.ShouldHaveReceived("UpdateUser").With(userID, mock.AnythingOfType("map[string]interface{}"))
	With(args ...interface{}) CallVerifier

	// WithArgs 验证参数匹配函数
	//
	// 使用自定义函数验证参数。
	//
	// 示例：
	//   spy.ShouldHaveReceived("ProcessOrder").WithArgs(func(args []interface{}) bool {
	//       order := args[0].(*Order)
	//       return order.Amount > 100
	//   })
	WithArgs(matcher func([]interface{}) bool) CallVerifier
}

// CallRecord 调用记录结构
//
// CallRecord 记录了方法调用的详细信息。
type CallRecord struct {
	// Method 方法名
	Method string

	// Args 调用参数
	Args []interface{}

	// Result 返回值
	Result []interface{}

	// Timestamp 调用时间
	Timestamp time.Time

	// Context 调用上下文
	Context context.Context

	// Error 调用错误（如果有）
	Error error
}

// FacadeMiddleware 门面中间件接口
//
// FacadeMiddleware 允许在门面方法调用前后执行额外逻辑，
// 如日志记录、性能监控、权限检查等。
//
// 使用示例：
//
//	type LoggingMiddleware struct{}
//
//	func (m *LoggingMiddleware) Handle(call *MethodCall, next func(*MethodCall) ([]interface{}, error)) ([]interface{}, error) {
//		start := time.Now()
//		defer func() {
//			duration := time.Since(start)
//			log.Printf("Method %s took %v", call.Method, duration)
//		}()
//
//		log.Printf("Calling method: %s with args: %v", call.Method, call.Args)
//		return next(call)
//	}
//
//	// 注册中间件
//	DB.UseMiddleware(&LoggingMiddleware{})
type FacadeMiddleware interface {
	// Handle 处理方法调用
	//
	// 在方法调用的生命周期中执行，可以在调用前后添加逻辑。
	//
	// 参数：
	//   call - 方法调用信息
	//   next - 下一个中间件或实际方法调用
	//
	// 返回：
	//   []interface{} - 方法返回值
	//   error        - 执行错误
	//
	// 示例：
	//   func (m *AuthMiddleware) Handle(call *MethodCall, next func(*MethodCall) ([]interface{}, error)) ([]interface{}, error) {
	//       if !m.isAuthorized(call.Context) {
	//           return nil, errors.New("unauthorized")
	//       }
	//       return next(call)
	//   }
	Handle(call *MethodCall, next func(*MethodCall) ([]interface{}, error)) ([]interface{}, error)

	// Priority 中间件优先级
	//
	// 返回中间件的执行优先级，数值越小优先级越高。
	//
	// 示例：
	//   func (m *AuthMiddleware) Priority() int {
	//       return 100 // 认证中间件应该较早执行
	//   }
	//
	//   func (m *LoggingMiddleware) Priority() int {
	//       return 900 // 日志中间件可以较晚执行
	//   }
	Priority() int
}

// MethodCall 方法调用信息结构
//
// MethodCall 包含了门面方法调用的所有信息。
type MethodCall struct {
	// Facade 门面实例
	Facade Facade

	// Method 方法名
	Method string

	// Args 调用参数
	Args []interface{}

	// Context 调用上下文
	Context context.Context

	// Metadata 元数据
	//
	// 存储调用相关的额外信息，如用户ID、请求ID等。
	//
	// 示例：
	//   Metadata: map[string]interface{}{
	//       "user_id":    123,
	//       "request_id": "req-abc-123",
	//       "ip_address": "192.168.1.1",
	//   }
	Metadata map[string]interface{}

	// StartTime 开始时间
	StartTime time.Time
}
