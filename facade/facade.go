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
// 包结构：
// - facade_interface.go - Facade 核心接口
// - static_facade.go - StaticFacade 静态门面接口
// - facade_manager.go - FacadeManager 门面管理器接口
// - mock_interface.go - MockInterface、SpyInterface、CallVerifier 测试相关接口
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
//
//	// 测试模拟
//	mockMailer := &MockMailer{}
//	Mail.Mock(mockMailer)
//	defer Mail.ClearMock()
//
//	// 门面管理器
//	manager := NewFacadeManager()
//	manager.Register("Payment", &PaymentFacade{})
//	payment := manager.MustResolve("Payment")
package facade
