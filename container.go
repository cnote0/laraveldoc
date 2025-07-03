package laraveldoc

import (
	"context"
	"reflect"
)

// Container 服务容器接口
// 实现了IoC（控制反转）模式，负责管理服务的注册、解析和生命周期
type Container interface {
	// Bind 绑定服务到容器
	// abstract: 服务标识符（通常是接口类型或字符串）
	// concrete: 具体实现（可以是实例、工厂函数或类型）
	// shared: 是否为单例模式
	Bind(abstract interface{}, concrete interface{}, shared bool) error

	// BindSingleton 绑定单例服务
	BindSingleton(abstract interface{}, concrete interface{}) error

	// BindInstance 绑定已存在的实例
	BindInstance(abstract interface{}, instance interface{}) error

	// Resolve 解析服务实例
	// abstract: 服务标识符
	// 返回解析的服务实例
	Resolve(abstract interface{}) (interface{}, error)

	// ResolveWithContext 带上下文的服务解析
	ResolveWithContext(ctx context.Context, abstract interface{}) (interface{}, error)

	// Make 创建服务实例（类似Laravel的make方法）
	Make(abstract interface{}, parameters ...interface{}) (interface{}, error)

	// Has 检查服务是否已绑定
	Has(abstract interface{}) bool

	// Forget 移除服务绑定
	Forget(abstract interface{}) error

	// Flush 清空所有绑定
	Flush() error

	// GetBindings 获取所有绑定信息
	GetBindings() map[interface{}]Binding

	// Tag 为服务添加标签
	Tag(abstracts []interface{}, tag string) error

	// Tagged 获取标签下的所有服务
	Tagged(tag string) ([]interface{}, error)
}

// Binding 服务绑定信息
type Binding struct {
	// Abstract 抽象服务标识
	Abstract interface{}

	// Concrete 具体实现
	Concrete interface{}

	// Shared 是否为共享实例（单例）
	Shared bool

	// Instance 已解析的实例（仅用于单例）
	Instance interface{}

	// Tags 服务标签
	Tags []string

	// Dependencies 依赖关系
	Dependencies []interface{}
}

// ServiceProvider 服务提供者接口
// 负责向容器注册服务
type ServiceProvider interface {
	// Register 注册服务到容器
	Register(container Container) error

	// Boot 启动服务（在所有服务注册完成后调用）
	Boot(container Container) error

	// Provides 返回此提供者提供的服务标识符
	Provides() []interface{}

	// IsDeferred 是否为延迟加载的服务提供者
	IsDeferred() bool
}

// Resolver 依赖解析器接口
type Resolver interface {
	// ResolveDependencies 解析构造函数依赖
	ResolveDependencies(target reflect.Type, container Container) ([]reflect.Value, error)

	// ResolveMethodDependencies 解析方法依赖
	ResolveMethodDependencies(method reflect.Method, container Container, parameters []interface{}) ([]reflect.Value, error)
}

// ContextualBinding 上下文绑定接口
// 支持基于上下文的依赖注入
type ContextualBinding interface {
	// When 指定上下文条件
	When(context interface{}) ContextualBindingBuilder
}

// ContextualBindingBuilder 上下文绑定构建器
type ContextualBindingBuilder interface {
	// Needs 指定需要的依赖
	Needs(abstract interface{}) ContextualBindingGiven
}

// ContextualBindingGiven 上下文绑定给定接口
type ContextualBindingGiven interface {
	// Give 提供具体实现
	Give(concrete interface{}) error

	// GiveTagged 提供标签服务
	GiveTagged(tag string) error
}
