package laraveldoc

import (
	"context"
	"reflect"
)

// Facade 门面接口
// 提供静态方式访问容器中的服务，简化API调用
type Facade interface {
	// GetFacadeAccessor 获取门面访问器
	// 返回在容器中注册的服务标识符
	GetFacadeAccessor() interface{}

	// SetFacadeContainer 设置门面使用的容器
	SetFacadeContainer(container Container)

	// GetFacadeContainer 获取门面使用的容器
	GetFacadeContainer() Container

	// ClearResolvedInstance 清除已解析的实例
	ClearResolvedInstance(name interface{})

	// ClearResolvedInstances 清除所有已解析的实例
	ClearResolvedInstances()

	// GetFacadeRoot 获取门面根对象
	GetFacadeRoot() (interface{}, error)
}

// StaticFacade 静态门面接口
// 提供静态方法调用支持
type StaticFacade interface {
	Facade

	// CallStatic 静态方法调用
	// method: 方法名
	// args: 方法参数
	CallStatic(method string, args ...interface{}) (interface{}, error)

	// CallStaticWithContext 带上下文的静态方法调用
	CallStaticWithContext(ctx context.Context, method string, args ...interface{}) (interface{}, error)
}

// FacadeManager 门面管理器接口
// 负责管理所有门面的注册和解析
type FacadeManager interface {
	// RegisterFacade 注册门面
	RegisterFacade(name string, facade Facade) error

	// GetFacade 获取门面实例
	GetFacade(name string) (Facade, error)

	// HasFacade 检查门面是否存在
	HasFacade(name string) bool

	// RemoveFacade 移除门面
	RemoveFacade(name string) error

	// GetRegisteredFacades 获取所有已注册的门面
	GetRegisteredFacades() map[string]Facade

	// SetContainer 设置容器
	SetContainer(container Container)

	// GetContainer 获取容器
	GetContainer() Container
}

// FacadeProxy 门面代理接口
// 提供动态代理功能
type FacadeProxy interface {
	// ProxyCall 代理方法调用
	ProxyCall(target interface{}, method string, args []interface{}) (interface{}, error)

	// ProxyCallWithContext 带上下文的代理方法调用
	ProxyCallWithContext(ctx context.Context, target interface{}, method string, args []interface{}) (interface{}, error)

	// GetProxyTarget 获取代理目标
	GetProxyTarget() interface{}

	// SetProxyTarget 设置代理目标
	SetProxyTarget(target interface{})
}

// FacadeAccessor 门面访问器
// 定义门面如何访问底层服务
type FacadeAccessor struct {
	// ServiceName 服务名称
	ServiceName string

	// ServiceType 服务类型
	ServiceType reflect.Type

	// Factory 工厂函数
	Factory func(Container) (interface{}, error)

	// Singleton 是否为单例
	Singleton bool
}

// RealtimeFacade 实时门面接口
// 支持实时更新的门面服务
type RealtimeFacade interface {
	Facade

	// Swap 交换底层服务实例
	Swap(instance interface{}) interface{}

	// SwapWithCallback 使用回调交换服务实例
	SwapWithCallback(callback func(interface{}) interface{}) interface{}

	// PartialMock 部分模拟
	PartialMock() MockInterface

	// Spy 创建间谍对象
	Spy() SpyInterface
}

// MockInterface 模拟接口
// 用于测试时的服务模拟
type MockInterface interface {
	// ShouldReceive 期望接收方法调用
	ShouldReceive(method string) ExpectationInterface

	// ShouldNotReceive 不应该接收方法调用
	ShouldNotReceive(method string) ExpectationInterface

	// AllowMockingNonExistentMethods 允许模拟不存在的方法
	AllowMockingNonExistentMethods(allow bool)

	// MockerGetExpectationCount 获取期望计数
	MockerGetExpectationCount() int
}

// SpyInterface 间谍接口
// 用于监控方法调用
type SpyInterface interface {
	// ShouldHaveReceived 应该已接收方法调用
	ShouldHaveReceived(method string, args ...interface{}) ExpectationInterface

	// ShouldNotHaveReceived 不应该已接收方法调用
	ShouldNotHaveReceived(method string, args ...interface{}) ExpectationInterface

	// GetCallHistory 获取调用历史
	GetCallHistory() []CallRecord
}

// ExpectationInterface 期望接口
// 定义方法调用期望
type ExpectationInterface interface {
	// With 指定参数
	With(args ...interface{}) ExpectationInterface

	// WithArgs 指定参数（别名）
	WithArgs(args ...interface{}) ExpectationInterface

	// WithAnyArgs 接受任意参数
	WithAnyArgs() ExpectationInterface

	// AndReturn 指定返回值
	AndReturn(values ...interface{}) ExpectationInterface

	// AndReturnUsing 使用函数指定返回值
	AndReturnUsing(callback func(...interface{}) interface{}) ExpectationInterface

	// AndThrow 抛出异常
	AndThrow(err error) ExpectationInterface

	// Times 指定调用次数
	Times(count int) ExpectationInterface

	// Once 调用一次
	Once() ExpectationInterface

	// Twice 调用两次
	Twice() ExpectationInterface

	// Never 永不调用
	Never() ExpectationInterface
}

// CallRecord 调用记录
type CallRecord struct {
	// Method 方法名
	Method string

	// Args 参数
	Args []interface{}

	// ReturnValues 返回值
	ReturnValues []interface{}

	// Error 错误
	Error error

	// Timestamp 时间戳
	Timestamp int64
}

// FacadeMiddleware 门面中间件接口
// 提供门面调用的中间件支持
type FacadeMiddleware interface {
	// Handle 处理门面调用
	Handle(call FacadeCall, next func(FacadeCall) (interface{}, error)) (interface{}, error)
}

// FacadeCall 门面调用信息
type FacadeCall struct {
	// FacadeName 门面名称
	FacadeName string

	// Method 方法名
	Method string

	// Args 参数
	Args []interface{}

	// Context 上下文
	Context context.Context

	// Target 目标服务
	Target interface{}
}
