package laraveldoc

import (
	"context"
	"time"
)

// Application 应用程序接口
// 定义应用程序的核心功能和生命周期管理
type Application interface {
	Container // 继承容器接口

	// Version 获取应用版本
	Version() string

	// Environment 获取当前环境
	Environment() string

	// IsEnvironment 检查是否为指定环境
	IsEnvironment(environments ...string) bool

	// IsProduction 是否为生产环境
	IsProduction() bool

	// IsDevelopment 是否为开发环境
	IsDevelopment() bool

	// IsDebug 是否开启调试模式
	IsDebug() bool

	// SetDebug 设置调试模式
	SetDebug(debug bool)

	// BasePath 获取应用根路径
	BasePath(path ...string) string

	// ConfigPath 获取配置路径
	ConfigPath(path ...string) string

	// StoragePath 获取存储路径
	StoragePath(path ...string) string

	// DatabasePath 获取数据库路径
	DatabasePath(path ...string) string

	// ResourcePath 获取资源路径
	ResourcePath(path ...string) string

	// PublicPath 获取公共路径
	PublicPath(path ...string) string

	// Bootstrap 启动应用程序
	Bootstrap() error

	// BootstrapWith 使用指定的引导程序启动
	BootstrapWith(bootstrappers []Bootstrapper) error

	// RegisterConfiguredProviders 注册配置的服务提供者
	RegisterConfiguredProviders() error

	// RegisterProvider 注册服务提供者
	RegisterProvider(provider ServiceProvider, force bool) ServiceProvider

	// GetProviders 获取已注册的服务提供者
	GetProviders(provider ServiceProvider) []ServiceProvider

	// BootProviders 启动所有服务提供者
	BootProviders() error

	// Terminate 终止应用程序
	Terminate() error

	// GetNamespace 获取应用命名空间
	GetNamespace() string

	// SetNamespace 设置应用命名空间
	SetNamespace(namespace string)
}

// Bootstrapper 启动程序接口
// 负责应用程序的各个启动阶段
type Bootstrapper interface {
	// Bootstrap 执行启动逻辑
	Bootstrap(app Application) error

	// Priority 获取启动优先级
	Priority() int

	// Name 获取启动程序名称
	Name() string
}

// Kernel 内核接口
// 定义应用程序内核功能
type Kernel interface {
	// Bootstrap 启动内核
	Bootstrap() error

	// Handle 处理请求
	Handle(request interface{}) (interface{}, error)

	// HandleWithContext 带上下文处理请求
	HandleWithContext(ctx context.Context, request interface{}) (interface{}, error)

	// Terminate 终止内核
	Terminate(request interface{}, response interface{}) error

	// GetApplication 获取应用程序实例
	GetApplication() Application

	// SetApplication 设置应用程序实例
	SetApplication(app Application)
}

// ConsoleKernel 控制台内核接口
// 处理命令行请求
type ConsoleKernel interface {
	Kernel

	// HandleConsole 处理控制台命令
	HandleConsole(input InputInterface, output OutputInterface) (int, error)

	// Call 调用命令
	Call(command string, parameters map[string]interface{}) (int, error)

	// Queue 队列命令
	Queue(command string, parameters map[string]interface{}) error

	// GetArtisan 获取Artisan实例
	GetArtisan() ArtisanInterface

	// SetArtisan 设置Artisan实例
	SetArtisan(artisan ArtisanInterface)
}

// InputInterface 输入接口
type InputInterface interface {
	// GetFirstArgument 获取第一个参数
	GetFirstArgument() string

	// HasParameterOption 检查是否有参数选项
	HasParameterOption(values []string, onlyParams bool) bool

	// GetParameterOption 获取参数选项
	GetParameterOption(values []string, defaultValue interface{}, onlyParams bool) interface{}

	// Bind 绑定输入定义
	Bind(definition InputDefinition) error

	// Validate 验证输入
	Validate() error

	// GetArguments 获取所有参数
	GetArguments() map[string]interface{}

	// GetArgument 获取指定参数
	GetArgument(name string) interface{}

	// SetArgument 设置参数
	SetArgument(name string, value interface{}) error

	// HasArgument 检查是否有参数
	HasArgument(name string) bool

	// GetOptions 获取所有选项
	GetOptions() map[string]interface{}

	// GetOption 获取指定选项
	GetOption(name string) interface{}

	// SetOption 设置选项
	SetOption(name string, value interface{}) error

	// HasOption 检查是否有选项
	HasOption(name string) bool

	// IsInteractive 是否为交互模式
	IsInteractive() bool

	// SetInteractive 设置交互模式
	SetInteractive(interactive bool)
}

// OutputInterface 输出接口
type OutputInterface interface {
	// Write 写入内容
	Write(messages []string, newline bool, verbosity int) error

	// WriteLine 写入一行
	WriteLine(message string, verbosity int) error

	// SetVerbosity 设置详细程度
	SetVerbosity(level int)

	// GetVerbosity 获取详细程度
	GetVerbosity() int

	// IsQuiet 是否为安静模式
	IsQuiet() bool

	// IsVerbose 是否为详细模式
	IsVerbose() bool

	// IsVeryVerbose 是否为非常详细模式
	IsVeryVerbose() bool

	// IsDebug 是否为调试模式
	IsDebug() bool

	// SetDecorated 设置装饰模式
	SetDecorated(decorated bool)

	// IsDecorated 是否为装饰模式
	IsDecorated() bool

	// SetFormatter 设置格式化器
	SetFormatter(formatter OutputFormatter)

	// GetFormatter 获取格式化器
	GetFormatter() OutputFormatter
}

// InputDefinition 输入定义
type InputDefinition interface {
	// SetDefinition 设置定义
	SetDefinition(definition []InputArgument) error

	// SetArguments 设置参数
	SetArguments(arguments []InputArgument) error

	// AddArguments 添加参数
	AddArguments(arguments []InputArgument) error

	// GetArgument 获取参数
	GetArgument(name string) (InputArgument, error)

	// HasArgument 检查是否有参数
	HasArgument(name string) bool

	// GetArguments 获取所有参数
	GetArguments() []InputArgument

	// GetArgumentCount 获取参数数量
	GetArgumentCount() int

	// GetArgumentRequiredCount 获取必需参数数量
	GetArgumentRequiredCount() int

	// GetArgumentDefaults 获取参数默认值
	GetArgumentDefaults() map[string]interface{}

	// SetOptions 设置选项
	SetOptions(options []InputOption) error

	// AddOptions 添加选项
	AddOptions(options []InputOption) error

	// GetOption 获取选项
	GetOption(name string) (InputOption, error)

	// HasOption 检查是否有选项
	HasOption(name string) bool

	// GetOptions 获取所有选项
	GetOptions() []InputOption

	// HasShortcut 检查是否有快捷方式
	HasShortcut(name string) bool

	// GetOptionForShortcut 通过快捷方式获取选项
	GetOptionForShortcut(shortcut string) (InputOption, error)

	// GetOptionDefaults 获取选项默认值
	GetOptionDefaults() map[string]interface{}
}

// InputArgument 输入参数
type InputArgument struct {
	// Name 参数名
	Name string

	// Mode 模式（必需、可选等）
	Mode int

	// Description 描述
	Description string

	// Default 默认值
	Default interface{}
}

// InputOption 输入选项
type InputOption struct {
	// Name 选项名
	Name string

	// Shortcut 快捷方式
	Shortcut string

	// Mode 模式
	Mode int

	// Description 描述
	Description string

	// Default 默认值
	Default interface{}
}

// OutputFormatter 输出格式化器
type OutputFormatter interface {
	// SetDecorated 设置装饰
	SetDecorated(decorated bool)

	// IsDecorated 是否装饰
	IsDecorated() bool

	// SetStyle 设置样式
	SetStyle(name string, style OutputFormatterStyle) error

	// HasStyle 检查是否有样式
	HasStyle(name string) bool

	// GetStyle 获取样式
	GetStyle(name string) OutputFormatterStyle

	// Format 格式化消息
	Format(message string) string
}

// OutputFormatterStyle 输出格式化样式
type OutputFormatterStyle interface {
	// SetForeground 设置前景色
	SetForeground(color string) error

	// SetBackground 设置背景色
	SetBackground(color string) error

	// SetOption 设置选项
	SetOption(option string) error

	// UnsetOption 取消选项
	UnsetOption(option string) error

	// SetOptions 设置多个选项
	SetOptions(options []string) error

	// Apply 应用样式
	Apply(text string) string
}

// ArtisanInterface Artisan 命令行接口
type ArtisanInterface interface {
	// Add 添加命令
	Add(command CommandInterface) CommandInterface

	// Get 获取命令
	Get(name string) (CommandInterface, error)

	// Has 检查是否有命令
	Has(name string) bool

	// GetNames 获取所有命令名
	GetNames() []string

	// All 获取所有命令
	All(namespace string) map[string]CommandInterface

	// GetAbbreviations 获取缩写
	GetAbbreviations() []string

	// Register 注册命令
	Register(name string) CommandInterface

	// AddCommands 添加多个命令
	AddCommands(commands []CommandInterface) error

	// Find 查找命令
	Find(name string) (CommandInterface, error)

	// FindNamespace 查找命名空间
	FindNamespace(namespace string) (string, error)

	// GetNamespaces 获取所有命名空间
	GetNamespaces() []string

	// ExtractNamespace 提取命名空间
	ExtractNamespace(name string, limit int) string

	// SetName 设置名称
	SetName(name string)

	// GetName 获取名称
	GetName() string

	// SetVersion 设置版本
	SetVersion(version string)

	// GetVersion 获取版本
	GetVersion() string

	// GetLongVersion 获取长版本
	GetLongVersion() string

	// SetAutoExit 设置自动退出
	SetAutoExit(autoExit bool)

	// SetCatchExceptions 设置捕获异常
	SetCatchExceptions(catchExceptions bool)

	// GetHelp 获取帮助
	GetHelp() string

	// Run 运行命令
	Run(input InputInterface, output OutputInterface) (int, error)

	// DoRun 执行运行
	DoRun(input InputInterface, output OutputInterface) (int, error)

	// SetHelperSet 设置帮助集
	SetHelperSet(helperSet HelperSetInterface)

	// GetHelperSet 获取帮助集
	GetHelperSet() HelperSetInterface

	// SetDefinition 设置定义
	SetDefinition(definition InputDefinition)

	// GetDefinition 获取定义
	GetDefinition() InputDefinition

	// GetDefaultInputDefinition 获取默认输入定义
	GetDefaultInputDefinition() InputDefinition

	// GetDefaultCommands 获取默认命令
	GetDefaultCommands() []CommandInterface

	// GetDefaultHelperSet 获取默认帮助集
	GetDefaultHelperSet() HelperSetInterface
}

// CommandInterface 命令接口
type CommandInterface interface {
	// Configure 配置命令
	Configure() error

	// Execute 执行命令
	Execute(input InputInterface, output OutputInterface) error

	// Interact 交互
	Interact(input InputInterface, output OutputInterface) error

	// Initialize 初始化
	Initialize(input InputInterface, output OutputInterface) error

	// Run 运行命令
	Run(input InputInterface, output OutputInterface) (int, error)

	// SetCode 设置代码
	SetCode(code func(InputInterface, OutputInterface) error) CommandInterface

	// MergeApplicationDefinition 合并应用定义
	MergeApplicationDefinition(mergeArgs bool) error

	// SetDefinition 设置定义
	SetDefinition(definition interface{}) error

	// GetDefinition 获取定义
	GetDefinition() InputDefinition

	// GetNativeDefinition 获取原生定义
	GetNativeDefinition() InputDefinition

	// AddArgument 添加参数
	AddArgument(name string, mode int, description string, defaultValue interface{}) CommandInterface

	// AddOption 添加选项
	AddOption(name string, shortcut string, mode int, description string, defaultValue interface{}) CommandInterface

	// SetName 设置名称
	SetName(name string) CommandInterface

	// GetName 获取名称
	GetName() string

	// SetHidden 设置隐藏
	SetHidden(hidden bool) CommandInterface

	// IsHidden 是否隐藏
	IsHidden() bool

	// SetDescription 设置描述
	SetDescription(description string) CommandInterface

	// GetDescription 获取描述
	GetDescription() string

	// SetHelp 设置帮助
	SetHelp(help string) CommandInterface

	// GetHelp 获取帮助
	GetHelp() string

	// GetProcessedHelp 获取处理过的帮助
	GetProcessedHelp() string

	// SetAliases 设置别名
	SetAliases(aliases []string) CommandInterface

	// GetAliases 获取别名
	GetAliases() []string

	// GetSynopsis 获取概要
	GetSynopsis(short bool) string

	// AddUsage 添加用法
	AddUsage(usage string) CommandInterface

	// GetUsages 获取用法
	GetUsages() []string

	// GetHelper 获取帮助器
	GetHelper(name string) HelperInterface

	// AsText 转换为文本
	AsText() string

	// AsXml 转换为XML
	AsXml(asDom bool) interface{}

	// SetApplication 设置应用
	SetApplication(application ArtisanInterface)

	// GetApplication 获取应用
	GetApplication() ArtisanInterface

	// IsEnabled 是否启用
	IsEnabled() bool
}

// HelperSetInterface 帮助集接口
type HelperSetInterface interface {
	// Set 设置帮助器
	Set(helper HelperInterface, alias string) error

	// Has 检查是否有帮助器
	Has(name string) bool

	// Get 获取帮助器
	Get(name string) HelperInterface

	// SetCommand 设置命令
	SetCommand(command CommandInterface)

	// GetCommand 获取命令
	GetCommand() CommandInterface
}

// HelperInterface 帮助器接口
type HelperInterface interface {
	// SetHelperSet 设置帮助集
	SetHelperSet(helperSet HelperSetInterface)

	// GetHelperSet 获取帮助集
	GetHelperSet() HelperSetInterface

	// GetName 获取名称
	GetName() string
}

// EventInterface 事件接口
type EventInterface interface {
	// IsPropagationStopped 是否停止传播
	IsPropagationStopped() bool

	// StopPropagation 停止传播
	StopPropagation()
}

// EventDispatcher 事件分发器接口
type EventDispatcher interface {
	// Dispatch 分发事件
	Dispatch(event interface{}, eventName string) interface{}

	// DispatchWithContext 带上下文分发事件
	DispatchWithContext(ctx context.Context, event interface{}, eventName string) interface{}

	// AddListener 添加监听器
	AddListener(eventName string, listener EventListener, priority int) error

	// AddSubscriber 添加订阅者
	AddSubscriber(subscriber EventSubscriber) error

	// RemoveListener 移除监听器
	RemoveListener(eventName string, listener EventListener) error

	// RemoveSubscriber 移除订阅者
	RemoveSubscriber(subscriber EventSubscriber) error

	// GetListeners 获取监听器
	GetListeners(eventName string) []EventListener

	// GetListenerPriority 获取监听器优先级
	GetListenerPriority(eventName string, listener EventListener) (int, error)

	// HasListeners 检查是否有监听器
	HasListeners(eventName string) bool
}

// EventListener 事件监听器
type EventListener func(event interface{}) error

// EventSubscriber 事件订阅者接口
type EventSubscriber interface {
	// GetSubscribedEvents 获取订阅的事件
	GetSubscribedEvents() map[string]interface{}
}

// Config 配置接口
type Config interface {
	// Get 获取配置值
	Get(key string, defaultValue interface{}) interface{}

	// Set 设置配置值
	Set(key string, value interface{}) error

	// Has 检查是否有配置
	Has(key string) bool

	// All 获取所有配置
	All() map[string]interface{}

	// OffsetExists 检查偏移是否存在
	OffsetExists(key string) bool

	// OffsetGet 获取偏移值
	OffsetGet(key string) interface{}

	// OffsetSet 设置偏移值
	OffsetSet(key string, value interface{})

	// OffsetUnset 取消偏移设置
	OffsetUnset(key string)

	// Prepend 前置值
	Prepend(key string, value interface{}) error

	// Push 推送值
	Push(key string, value interface{}) error
}

// LogManager 日志管理器接口
type LogManager interface {
	// Channel 获取日志通道
	Channel(name string) LoggerInterface

	// Driver 获取日志驱动
	Driver(driver string) LoggerInterface

	// Stack 创建日志栈
	Stack(channels []string, channel string) LoggerInterface

	// Build 构建日志器
	Build(config map[string]interface{}) LoggerInterface

	// GetDefaultDriver 获取默认驱动
	GetDefaultDriver() string

	// SetDefaultDriver 设置默认驱动
	SetDefaultDriver(name string)

	// Extend 扩展驱动
	Extend(driver string, callback func(Application, map[string]interface{}) LoggerInterface) LogManager

	// GetChannels 获取所有通道
	GetChannels() map[string]LoggerInterface
}

// LoggerInterface 日志器接口
type LoggerInterface interface {
	// Emergency 紧急日志
	Emergency(message string, context map[string]interface{}) error

	// Alert 警报日志
	Alert(message string, context map[string]interface{}) error

	// Critical 严重日志
	Critical(message string, context map[string]interface{}) error

	// Error 错误日志
	Error(message string, context map[string]interface{}) error

	// Warning 警告日志
	Warning(message string, context map[string]interface{}) error

	// Notice 通知日志
	Notice(message string, context map[string]interface{}) error

	// Info 信息日志
	Info(message string, context map[string]interface{}) error

	// Debug 调试日志
	Debug(message string, context map[string]interface{}) error

	// Log 记录日志
	Log(level string, message string, context map[string]interface{}) error

	// WithContext 带上下文
	WithContext(context map[string]interface{}) LoggerInterface
}

// CacheManager 缓存管理器接口
type CacheManager interface {
	// Store 获取缓存存储
	Store(name string) CacheStore

	// Driver 获取缓存驱动
	Driver(driver string) CacheStore

	// GetDefaultDriver 获取默认驱动
	GetDefaultDriver() string

	// SetDefaultDriver 设置默认驱动
	SetDefaultDriver(name string)

	// Extend 扩展驱动
	Extend(driver string, callback func(Application, map[string]interface{}) CacheStore) CacheManager

	// PurgeStores 清除存储
	PurgeStores() CacheManager
}

// CacheStore 缓存存储接口
type CacheStore interface {
	// Get 获取缓存
	Get(key string) (interface{}, error)

	// Many 获取多个缓存
	Many(keys []string) (map[string]interface{}, error)

	// Put 放置缓存
	Put(key string, value interface{}, ttl time.Duration) error

	// PutMany 放置多个缓存
	PutMany(values map[string]interface{}, ttl time.Duration) error

	// Add 添加缓存（如果不存在）
	Add(key string, value interface{}, ttl time.Duration) (bool, error)

	// Increment 增量
	Increment(key string, value int64) (int64, error)

	// Decrement 减量
	Decrement(key string, value int64) (int64, error)

	// Forever 永久缓存
	Forever(key string, value interface{}) error

	// Remember 记住缓存
	Remember(key string, ttl time.Duration, callback func() interface{}) (interface{}, error)

	// RememberForever 永久记住缓存
	RememberForever(key string, callback func() interface{}) (interface{}, error)

	// Forget 忘记缓存
	Forget(key string) (bool, error)

	// Flush 清空缓存
	Flush() (bool, error)

	// GetPrefix 获取前缀
	GetPrefix() string
}
