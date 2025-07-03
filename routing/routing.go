// Package routing 提供 Laravel 风格的路由系统协议定义
//
// 本包定义了完整的HTTP路由系统，包括路由注册、分发、中间件处理、
// 请求响应管理、URL生成等功能。设计基于 Laravel 的路由系统。
//
// 主要特性：
// - RESTful 路由支持
// - 路由分组和中间件
// - 路由参数和约束
// - 子域名路由
// - 路由缓存和优化
// - 请求和响应处理
// - URL 生成和重定向
//
// 使用示例：
//
//	// 创建路由器
//	router := routing.NewRouter()
//
//	// 基础路由
//	router.Get("/users", userController.Index)
//	router.Post("/users", userController.Create)
//	router.Get("/users/{id}", userController.Show)
//	router.Put("/users/{id}", userController.Update)
//	router.Delete("/users/{id}", userController.Delete)
//
//	// 路由分组
//	api := router.Group("/api/v1")
//	api.Middleware("auth", "throttle")
//	api.Get("/profile", profileController.Show)
//
//	// 子域名路由
//	admin := router.Domain("admin.example.com")
//	admin.Get("/dashboard", adminController.Dashboard)
//
//	// 资源路由
//	router.Resource("/posts", postController)
//
//	// 路由分发
//	response, err := router.Dispatch(request)
//
//	// URL 生成
//	urlGenerator := routing.NewUrlGenerator(router)
//	userUrl := urlGenerator.Route("user.show", map[string]interface{}{"id": 123})

package routing

import (
	"context"
)

// Router 路由器接口
//
// Router 是整个路由系统的核心，负责路由的注册、匹配和分发。
// 它支持多种HTTP方法、路由分组、中间件、参数约束等功能。
//
// 使用示例：
//
//	// 创建路由器实例
//	router := NewRouter()
//
//	// 注册基础路由
//	router.Get("/", homeHandler)
//	router.Post("/users", createUserHandler)
//	router.Get("/users/{id:[0-9]+}", getUserHandler)
//
//	// 路由分组
//	api := router.Group("/api")
//	api.Middleware("auth")
//	api.Get("/profile", profileHandler)
//
//	// 资源路由
//	router.Resource("/posts", postsController)
//
//	// 处理请求
//	response, err := router.Dispatch(request)
//	if err != nil {
//		// 处理错误
//	}

// Router 路由器接口
type Router interface {
	// Get 注册GET路由
	Get(uri string, action interface{}) Route

	// Post 注册POST路由
	Post(uri string, action interface{}) Route

	// Put 注册PUT路由
	Put(uri string, action interface{}) Route

	// Patch 注册PATCH路由
	Patch(uri string, action interface{}) Route

	// Delete 注册DELETE路由
	Delete(uri string, action interface{}) Route

	// Options 注册OPTIONS路由
	Options(uri string, action interface{}) Route

	// Any 注册任意方法路由
	Any(uri string, action interface{}) Route

	// Match 注册指定方法路由
	Match(methods []string, uri string, action interface{}) Route

	// Group 路由分组
	Group(attributes map[string]interface{}, callback func(Router)) Router

	// Prefix 设置路由前缀
	Prefix(prefix string) Router

	// Middleware 设置中间件
	Middleware(middleware ...string) Router

	// Namespace 设置命名空间
	Namespace(namespace string) Router

	// Name 设置路由名称前缀
	Name(name string) Router

	// Domain 设置域名
	Domain(domain string) Router

	// Where 设置路由参数约束
	Where(name string, expression string) Router

	// Resource 资源路由
	Resource(name string, controller string, options map[string]interface{}) RouteCollection

	// APIResource API资源路由
	APIResource(name string, controller string, options map[string]interface{}) RouteCollection

	// Redirect 重定向路由
	Redirect(uri string, destination string, status int) Route

	// PermanentRedirect 永久重定向路由
	PermanentRedirect(uri string, destination string) Route

	// View 视图路由
	View(uri string, view string, data map[string]interface{}) Route

	// Fallback 回退路由
	Fallback(action interface{}) Route

	// GetRoutes 获取所有路由
	GetRoutes() RouteCollection

	// Dispatch 分发请求
	Dispatch(request RequestInterface) ResponseInterface

	// DispatchToRoute 分发到路由
	DispatchToRoute(request RequestInterface) ResponseInterface
}

// Route 路由接口
type Route interface {
	// GetAction 获取动作
	GetAction() interface{}

	// SetAction 设置动作
	SetAction(action interface{}) Route

	// GetController 获取控制器
	GetController() string

	// GetControllerClass 获取控制器类
	GetControllerClass() string

	// GetControllerMethod 获取控制器方法
	GetControllerMethod() string

	// GetName 获取路由名称
	GetName() string

	// SetName 设置路由名称
	SetName(name string) Route

	// GetPrefix 获取前缀
	GetPrefix() string

	// SetPrefix 设置前缀
	SetPrefix(prefix string) Route

	// GetURI 获取URI
	GetURI() string

	// SetURI 设置URI
	SetURI(uri string) Route

	// GetMethods 获取HTTP方法
	GetMethods() []string

	// SetMethods 设置HTTP方法
	SetMethods(methods []string) Route

	// GetMiddleware 获取中间件
	GetMiddleware() []string

	// SetMiddleware 设置中间件
	SetMiddleware(middleware []string) Route

	// GetWhere 获取参数约束
	GetWhere() map[string]string

	// SetWhere 设置参数约束
	SetWhere(wheres map[string]string) Route

	// Where 添加参数约束
	Where(name string, expression string) Route

	// WhereNumber 数字约束
	WhereNumber(name string) Route

	// WhereAlpha 字母约束
	WhereAlpha(name string) Route

	// WhereAlphaNumeric 字母数字约束
	WhereAlphaNumeric(name string) Route

	// WhereUuid UUID约束
	WhereUuid(name string) Route

	// Middleware 添加中间件
	Middleware(middleware ...string) Route

	// WithoutMiddleware 排除中间件
	WithoutMiddleware(middleware ...string) Route

	// Name 设置名称
	Name(name string) Route

	// Domain 设置域名
	Domain(domain string) Route

	// Defaults 设置默认参数
	Defaults(key string, value interface{}) Route

	// Bind 绑定路由
	Bind(request RequestInterface) error

	// Matches 检查是否匹配请求
	Matches(request RequestInterface) bool

	// Run 运行路由
	Run() ResponseInterface

	// GetCompiled 获取编译后的路由
	GetCompiled() CompiledRoute

	// Compile 编译路由
	Compile() CompiledRoute
}

// RouteCollection 路由集合接口
type RouteCollection interface {
	// Add 添加路由
	Add(route Route) Route

	// Get 获取路由
	Get(method string, uri string) Route

	// HasNamedRoute 检查是否有命名路由
	HasNamedRoute(name string) bool

	// GetByName 根据名称获取路由
	GetByName(name string) Route

	// GetByAction 根据动作获取路由
	GetByAction(action string) Route

	// GetRoutes 获取所有路由
	GetRoutes() []Route

	// GetRoutesByMethod 根据方法获取路由
	GetRoutesByMethod() map[string][]Route

	// GetRoutesByName 根据名称获取路由
	GetRoutesByName() map[string]Route

	// Count 获取路由数量
	Count() int

	// Match 匹配请求
	Match(request RequestInterface) Route

	// Refresh 刷新路由
	Refresh()
}

// CompiledRoute 编译后的路由接口
type CompiledRoute interface {
	// GetRegex 获取正则表达式
	GetRegex() string

	// GetTokens 获取标记
	GetTokens() []interface{}

	// GetStaticPrefix 获取静态前缀
	GetStaticPrefix() string

	// GetVariables 获取变量
	GetVariables() []string
}

// RequestInterface 请求接口
type RequestInterface interface {
	// GetMethod 获取HTTP方法
	GetMethod() string

	// GetURI 获取URI
	GetURI() string

	// GetPath 获取路径
	GetPath() string

	// GetQuery 获取查询字符串
	GetQuery() string

	// GetHeaders 获取请求头
	GetHeaders() map[string][]string

	// GetHeader 获取指定请求头
	GetHeader(name string) string

	// HasHeader 检查是否有指定请求头
	HasHeader(name string) bool

	// GetInput 获取输入数据
	GetInput(key string, defaultValue interface{}) interface{}

	// All 获取所有输入数据
	All() map[string]interface{}

	// Has 检查是否有输入数据
	Has(key string) bool

	// File 获取上传文件
	File(key string) UploadedFile

	// HasFile 检查是否有上传文件
	HasFile(key string) bool

	// Cookie 获取Cookie
	Cookie(name string, defaultValue string) string

	// GetCookies 获取所有Cookie
	GetCookies() map[string]string

	// IP 获取客户端IP
	IP() string

	// UserAgent 获取用户代理
	UserAgent() string

	// GetRoute 获取路由
	GetRoute() Route

	// SetRoute 设置路由
	SetRoute(route Route)

	// GetRouteResolver 获取路由解析器
	GetRouteResolver() func() Route

	// SetRouteResolver 设置路由解析器
	SetRouteResolver(resolver func() Route)

	// Context 获取上下文
	Context() context.Context

	// WithContext 设置上下文
	WithContext(ctx context.Context) RequestInterface
}

// ResponseInterface 响应接口
type ResponseInterface interface {
	// GetContent 获取内容
	GetContent() string

	// SetContent 设置内容
	SetContent(content string) ResponseInterface

	// GetStatusCode 获取状态码
	GetStatusCode() int

	// SetStatusCode 设置状态码
	SetStatusCode(code int) ResponseInterface

	// GetHeaders 获取响应头
	GetHeaders() map[string][]string

	// SetHeader 设置响应头
	SetHeader(name string, value string) ResponseInterface

	// AddHeader 添加响应头
	AddHeader(name string, value string) ResponseInterface

	// RemoveHeader 移除响应头
	RemoveHeader(name string) ResponseInterface

	// WithCookie 设置Cookie
	WithCookie(cookie Cookie) ResponseInterface

	// WithoutCookie 移除Cookie
	WithoutCookie(name string) ResponseInterface

	// Send 发送响应
	Send() error

	// SendContent 发送内容
	SendContent() error

	// SendHeaders 发送响应头
	SendHeaders() error
}

// UploadedFile 上传文件接口
type UploadedFile interface {
	// GetClientOriginalName 获取原始文件名
	GetClientOriginalName() string

	// GetClientOriginalExtension 获取原始扩展名
	GetClientOriginalExtension() string

	// GetSize 获取文件大小
	GetSize() int64

	// GetMimeType 获取MIME类型
	GetMimeType() string

	// IsValid 检查文件是否有效
	IsValid() bool

	// Store 存储文件
	Store(path string, name string) (string, error)

	// StoreAs 存储文件并指定名称
	StoreAs(path string, name string) (string, error)

	// Move 移动文件
	Move(directory string, name string) error

	// GetPathname 获取路径名
	GetPathname() string

	// GetRealPath 获取真实路径
	GetRealPath() string
}

// Cookie Cookie接口
type Cookie interface {
	// GetName 获取名称
	GetName() string

	// GetValue 获取值
	GetValue() string

	// GetDomain 获取域
	GetDomain() string

	// GetPath 获取路径
	GetPath() string

	// GetExpiresTime 获取过期时间
	GetExpiresTime() int64

	// IsSecure 是否安全
	IsSecure() bool

	// IsHttpOnly 是否仅HTTP
	IsHttpOnly() bool

	// GetSameSite 获取SameSite
	GetSameSite() string
}

// Middleware 中间件接口
type Middleware interface {
	// Handle 处理请求
	Handle(request RequestInterface, next func(RequestInterface) ResponseInterface) ResponseInterface
}

// MiddlewareGroup 中间件组接口
type MiddlewareGroup interface {
	// GetMiddleware 获取中间件列表
	GetMiddleware() []string

	// AddMiddleware 添加中间件
	AddMiddleware(middleware string) MiddlewareGroup

	// PrependMiddleware 前置中间件
	PrependMiddleware(middleware string) MiddlewareGroup

	// RemoveMiddleware 移除中间件
	RemoveMiddleware(middleware string) MiddlewareGroup
}

// RouteModelBinding 路由模型绑定接口
type RouteModelBinding interface {
	// Bind 绑定模型
	Bind(value string, route Route) (interface{}, error)

	// GetBindingField 获取绑定字段
	GetBindingField() string
}

// UrlGenerator URL生成器接口
type UrlGenerator interface {
	// To 生成URL
	To(path string, parameters map[string]interface{}, secure bool) string

	// Route 生成命名路由URL
	Route(name string, parameters map[string]interface{}, absolute bool) string

	// Action 生成控制器动作URL
	Action(action string, parameters map[string]interface{}, absolute bool) string

	// Asset 生成资源URL
	Asset(path string, secure bool) string

	// SecureAsset 生成安全资源URL
	SecureAsset(path string) string

	// Current 获取当前URL
	Current() string

	// Full 获取完整URL
	Full() string

	// Previous 获取上一个URL
	Previous(fallback string) string

	// SetRequest 设置请求
	SetRequest(request RequestInterface)

	// GetRequest 获取请求
	GetRequest() RequestInterface
}

// RedirectResponse 重定向响应接口
type RedirectResponse interface {
	ResponseInterface

	// GetTargetUrl 获取目标URL
	GetTargetUrl() string

	// SetTargetUrl 设置目标URL
	SetTargetUrl(url string) RedirectResponse

	// With 设置会话数据
	With(key string, value interface{}) RedirectResponse

	// WithInput 设置输入数据
	WithInput(input map[string]interface{}) RedirectResponse

	// WithErrors 设置错误信息
	WithErrors(errors interface{}) RedirectResponse

	// WithCookies 设置Cookie
	WithCookies(cookies []Cookie) RedirectResponse

	// Away 外部重定向
	Away(url string) RedirectResponse

	// Back 返回上一页
	Back() RedirectResponse

	// Home 返回首页
	Home() RedirectResponse

	// Refresh 刷新当前页
	Refresh() RedirectResponse
}

// JsonResponse JSON响应接口
type JsonResponse interface {
	ResponseInterface

	// GetData 获取数据
	GetData() interface{}

	// SetData 设置数据
	SetData(data interface{}) JsonResponse

	// WithCallback 设置JSONP回调
	WithCallback(callback string) JsonResponse
}
