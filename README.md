# Laravel 设计模式 Go 协议定义

基于 Laravel 框架设计模式的 Go 语言协议定义项目。本项目提供了完整的接口定义，遵循 Laravel 的设计理念，同时充分利用 Go 语言的特性。

## 🎯 项目目标

本项目致力于将 Laravel 优秀的设计模式引入到 Go 生态系统中，提供：

- **标准化接口定义**：为 Go 应用提供统一的架构标准
- **Laravel 设计模式**：IoC 容器、门面模式、服务提供者等核心模式
- **强类型支持**：充分利用 Go 的类型系统和编译时检查
- **模块化设计**：每个功能模块独立成包，便于使用和维护
- **详细文档**：每个接口都有完整的使用示例和最佳实践

## 📦 包结构

项目采用模块化设计，将不同功能分成独立的包：

```
laraveldoc/
├── container/          # IoC 容器和依赖注入
├── facade/            # 门面模式和静态访问
├── application/       # 应用程序核心和生命周期
├── database/          # 基于 GORM 的数据库访问层
├── routing/           # HTTP 路由和请求处理
├── go.mod            # Go 模块定义
└── README.md         # 项目文档
```

## 🔧 包功能详解

### 1. Container 包 - IoC 容器系统

**路径**: `github.com/cnote0/laraveldoc/container`

提供强大的依赖注入容器功能：

**核心接口：**
- `Container` - IoC 容器核心接口，支持服务绑定、解析、生命周期管理
- `ServiceProvider` - 服务提供者模式，负责服务注册和引导
- `Resolver` - 依赖解析器，实现自动依赖注入
- `ContextualBinding` - 上下文绑定，支持基于条件的依赖注入

**使用示例：**
```go
import "github.com/cnote0/laraveldoc/container"

// 创建容器
container := NewContainer()

// 绑定服务
container.Bind("database", func(c container.Container) interface{} {
    return &Database{Host: "localhost"}
}, false)

// 单例绑定
container.Singleton("logger", func(c container.Container) interface{} {
    return &Logger{}
})

// 上下文绑定
container.When("UserController").Needs("Repository").Give("UserRepository")

// 解析服务
db, err := container.Make("database")
logger := container.MustMake("logger").(*Logger)
```

**主要特性：**
- 服务绑定和解析
- 单例模式支持
- 上下文绑定
- 自动依赖注入
- 循环依赖检测
- 服务标签管理

### 2. Facade 包 - 门面模式

**路径**: `github.com/cnote0/laraveldoc/facade`

提供优雅的静态访问接口：

**核心接口：**
- `Facade` - 门面核心接口，提供静态方法访问
- `StaticFacade` - 静态门面，支持反射调用
- `FacadeManager` - 门面管理器，统一管理所有门面
- `RealtimeFacade` - 实时门面，动态创建门面代理
- `MockInterface` - 测试模拟支持
- `FacadeMiddleware` - 门面中间件系统

**使用示例：**
```go
import "github.com/cnote0/laraveldoc/facade"

// 注册门面
manager := facade.NewFacadeManager(container)
manager.Register("DB", &DatabaseFacade{})
manager.Register("Cache", &CacheFacade{})
manager.Register("Log", &LoggerFacade{})

// 使用门面（通过生成的静态方法）
users := DB.Table("users").Where("active", true).Get()
Cache.Put("key", "value", time.Hour)
Log.Info("User login", userID)

// 实时门面
UserFacade := realtimeFacade.Create(&UserService{})
user, err := UserFacade.CreateUser(userData)

// 测试模拟
DB.Mock(mockDB)
defer DB.ClearMock()
```

**主要特性：**
- 静态方法调用
- 实时门面创建
- 测试模拟支持
- 方法调用中间件
- 门面生命周期管理

### 3. Application 包 - 应用程序核心

**路径**: `github.com/cnote0/laraveldoc/application`

应用程序的核心功能和生命周期管理：

**核心接口：**
- `Application` - 应用程序主接口，继承 Container
- `Kernel` - HTTP 内核，处理 HTTP 请求
- `ConsoleKernel` - 控制台内核，处理命令行
- `EventDispatcher` - 事件分发器，事件驱动架构
- `Config` - 配置管理器
- `LogManager` - 多通道日志系统
- `CacheManager` - 多驱动缓存系统

**使用示例：**
```go
import "github.com/cnote0/laraveldoc/application"

// 创建应用
app := application.New()

// 注册服务提供者
app.Register(&DatabaseServiceProvider{})
app.Register(&CacheServiceProvider{})

// 引导应用
err := app.Bootstrap()
if err != nil {
    log.Fatal(err)
}

// HTTP 处理
kernel := app.Make("http.kernel").(application.Kernel)
response := kernel.Handle(request)

// 控制台处理
console := app.Make("console.kernel").(application.ConsoleKernel)
code := console.HandleConsole(input, output)

// 事件系统
events := app.Make("events").(application.EventDispatcher)
events.Listen("user.created", userCreatedListener)
events.Dispatch("user.created", user)
```

**主要特性：**
- 应用生命周期管理
- 环境配置支持
- 事件驱动架构
- 命令行工具支持
- 多通道日志系统
- 多驱动缓存系统

### 4. Database 包 - 基于 GORM 的数据库访问层

**路径**: `github.com/cnote0/laraveldoc/database`

基于 gorm.io/gorm 设计的完整数据库访问层：

**核心接口：**
- `DB` - 核心数据库接口，兼容 GORM 设计模式
- `Model` - GORM 标准模型结构，支持软删除
- `EloquentBuilder` - Laravel 风格的 ORM 查询构建器
- `Migrator` - 数据库迁移器，支持表、列、索引管理
- `EloquentModel` - Laravel 风格模型接口
- `Relationship` - 关联关系基础接口

**使用示例：**
```go
import "github.com/cnote0/laraveldoc/database"

// 基础 GORM 操作
db := container.MustMake("database").(database.DB)

// 创建记录
user := &User{Name: "John", Email: "john@example.com"}
db.Create(user)

// 查询操作
var users []User
db.Where("active = ?", true).Find(&users)

// Laravel 风格的 Eloquent 查询
eloquent := container.MustMake("eloquent").(database.EloquentBuilder)
users := eloquent.Model(&User{}).Where("age", ">", 18).Get()

// 关联查询
users := eloquent.Model(&User{}).With("Profile", "Orders").Get()

// 数据库迁移
migrator := container.MustMake("migrator").(database.Migrator)
migrator.AutoMigrate(&User{}, &Profile{}, &Order{})
```

**主要特性：**
- 完整的 ORM 功能
- 强类型支持
- 软删除机制
- 关联关系管理
- 数据库迁移
- 查询构建器
- 事务管理

### 5. Routing 包 - HTTP 路由系统

**路径**: `github.com/cnote0/laraveldoc/routing`

完整的 HTTP 路由和请求处理系统：

**核心接口：**
- `Router` - 路由器接口，支持 RESTful 路由、分组、中间件
- `Route` - 单个路由定义，支持参数约束、中间件
- `RequestInterface` - HTTP 请求处理
- `ResponseInterface` - HTTP 响应处理
- `Middleware` - 中间件接口
- `UrlGenerator` - URL 生成和管理

**使用示例：**
```go
import "github.com/cnote0/laraveldoc/routing"

// 创建路由器
router := routing.NewRouter()

// 基础路由
router.Get("/users", userController.Index)
router.Post("/users", userController.Create)
router.Get("/users/{id}", userController.Show)
router.Put("/users/{id}", userController.Update)
router.Delete("/users/{id}", userController.Delete)

// 路由分组
api := router.Group("/api/v1")
api.Middleware("auth", "throttle")
api.Get("/profile", profileController.Show)

// 子域名路由
admin := router.Domain("admin.example.com")
admin.Get("/dashboard", adminController.Dashboard)

// 资源路由
router.Resource("/posts", postController)

// 路由分发
response, err := router.Dispatch(request)
```

**主要特性：**
- RESTful 路由支持
- 路由分组和中间件
- 路由参数和约束
- 子域名路由
- URL 生成和重定向
- 请求和响应处理

## 🌟 Laravel 设计模式分析

### IoC 容器 (Inversion of Control)

Laravel 的 IoC 容器是整个框架的核心，它管理类的依赖关系和生命周期：

**核心概念：**
- **服务绑定**：将抽象与具体实现绑定
- **依赖注入**：自动解析和注入依赖
- **服务提供者**：组织相关服务的注册逻辑
- **延迟加载**：按需解析服务

**优势：**
- 松耦合：组件之间不直接依赖
- 可测试性：轻松替换依赖进行测试
- 可扩展性：通过绑定替换实现

### 门面模式 (Facade Pattern)

门面提供了对底层复杂系统的简单接口：

**核心概念：**
- **静态访问**：通过静态方法调用服务
- **延迟解析**：运行时才解析实际服务
- **可测试性**：支持模拟和替换

**优势：**
- 简洁的 API
- 表达力强的代码
- 易于记忆的方法名

### 服务提供者模式

服务提供者负责向容器注册服务：

**核心概念：**
- **Register 阶段**：仅注册服务绑定
- **Boot 阶段**：执行需要其他服务的逻辑
- **延迟提供者**：按需加载

## 🚀 使用建议

### 通用建议

1. **保持接口的一致性**：实现时严格遵循接口定义
2. **注意并发安全**：在并发环境中使用时考虑线程安全
3. **合理使用上下文**：充分利用 Go 的 context 机制
4. **遵循错误处理约定**：使用 Go 风格的错误处理
5. **保持代码简洁**：避免过度设计，保持代码可读性

### Container 包使用建议

1. **使用强类型绑定**：尽量避免使用 interface{} 类型
2. **合理设计服务生命周期**：区分单例和多例服务
3. **避免循环依赖**：设计时注意依赖关系
4. **使用上下文绑定**：为不同场景提供不同实现

### Database 包使用建议

1. **使用强类型模型**：定义具体的结构体而非 map[string]interface{}
2. **实现软删除**：继承 Model 结构体或使用 DeletedAt 字段
3. **合理使用事务**：复杂操作要使用事务保证数据一致性
4. **优化查询性能**：使用 Preload 预加载关联，避免 N+1 查询
5. **设计良好的索引**：通过 Blueprint 接口设计合适的索引
6. **使用迁移管理**：通过 Migration 接口管理数据库结构变更

### Facade 包使用建议

1. **适度使用门面**：不是所有服务都需要门面
2. **测试时使用模拟**：充分利用 Mock 功能进行单元测试
3. **避免门面污染**：保持门面接口的简洁性

### Application 包使用建议

1. **合理组织服务提供者**：按功能模块组织提供者
2. **使用环境配置**：根据不同环境使用不同配置
3. **实现优雅关闭**：在应用关闭时清理资源

### Routing 包使用建议

1. **使用路由分组**：组织相关路由，应用通用中间件
2. **实现中间件**：利用中间件处理跨切面关注点
3. **参数验证**：在控制器中验证请求参数
4. **使用资源路由**：标准 CRUD 操作使用资源路由

## 🔍 GORM 集成特性

本项目的数据库协议基于 **gorm.io/gorm** 进行设计，具有以下优势：

1. **完整的ORM功能**：支持关联、钩子、预加载、事务等
2. **强类型支持**：使用具体类型而非泛型接口
3. **软删除**：内置 DeletedAt 类型，实现 sql.Scanner 和 driver.Valuer
4. **链式查询**：支持方法链式调用，提高代码可读性
5. **自动迁移**：支持数据库结构的自动迁移
6. **插件系统**：可扩展的插件架构
7. **多数据库支持**：支持 MySQL、PostgreSQL、SQLite、SQL Server
8. **连接池管理**：内置连接池和事务管理
9. **日志集成**：完整的SQL日志和性能监控
10. **测试友好**：支持干运行模式和模拟数据生成

## 🎨 Go 协议标准特点

1. **接口优先**：所有定义都基于接口，提供最大的灵活性
2. **组合优先**：通过接口组合实现复杂功能
3. **上下文支持**：全面支持 Go 的 context 机制
4. **错误处理**：遵循 Go 风格的错误处理约定
5. **并发安全**：考虑并发安全的设计
6. **强类型约束**：使用强类型而非 interface{}
7. **GORM兼容**：数据库层完全兼容 gorm.io/gorm 设计
8. **模块化设计**：每个功能模块独立成包，便于使用和维护

## 📝 示例用法

### 完整应用示例

```go
package main

import (
    "log"
    
    "github.com/cnote0/laraveldoc/container"
    "github.com/cnote0/laraveldoc/application"
    "github.com/cnote0/laraveldoc/database"
    "github.com/cnote0/laraveldoc/facade"
    "github.com/cnote0/laraveldoc/routing"
)

func main() {
    // 创建应用
    app := application.New()
    
    // 注册服务提供者
    app.Register(&DatabaseServiceProvider{})
    app.Register(&RouteServiceProvider{})
    app.Register(&LogServiceProvider{})
    
    // 引导应用
    if err := app.Bootstrap(); err != nil {
        log.Fatal("Failed to bootstrap application:", err)
    }
    
    // 启动服务提供者
    if err := app.BootProviders(); err != nil {
        log.Fatal("Failed to boot providers:", err)
    }
    
    // 获取路由器
    router := app.Make("router").(routing.Router)
    
    // 注册路由
    registerRoutes(router)
    
    // 启动HTTP服务器
    log.Println("Server starting on :8080")
    // http.ListenAndServe(":8080", router)
}

func registerRoutes(router routing.Router) {
    // API 路由组
    api := router.Group("/api/v1")
    api.Middleware("cors", "auth")
    
    // 用户相关路由
    users := api.Group("/users")
    users.Get("/", userController.Index)
    users.Post("/", userController.Create)
    users.Get("/{id}", userController.Show)
    users.Put("/{id}", userController.Update)
    users.Delete("/{id}", userController.Delete)
    
    // 管理员路由
    admin := router.Domain("admin.example.com")
    admin.Middleware("admin")
    admin.Get("/dashboard", adminController.Dashboard)
}
```

### 服务提供者示例

```go
type DatabaseServiceProvider struct{}

func (p *DatabaseServiceProvider) Register(container container.Container) error {
    // 注册数据库连接
    return container.Singleton("database", func(c container.Container) interface{} {
        config := c.MustMake("config").(application.Config)
        
        dsn := config.Get("database.dsn", "")
        db, err := gorm.Open(mysql.Open(dsn.(string)), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect database:", err)
        }
        
        return db
    })
}

func (p *DatabaseServiceProvider) Boot(container container.Container) error {
    // 执行数据库迁移
    db := container.MustMake("database").(database.DB)
    return db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

func (p *DatabaseServiceProvider) Provides() []string {
    return []string{"database"}
}

func (p *DatabaseServiceProvider) IsDeferred() bool {
    return false
}
```

## 📚 学习资源

- [Laravel 官方文档](https://laravel.com/docs)
- [GORM 官方文档](https://gorm.io/docs/)
- [Go 语言官方文档](https://golang.org/doc/)
- [设计模式](https://refactoring.guru/design-patterns)

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进本项目！

## 📄 许可证

本项目使用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。
