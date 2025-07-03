package facade

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
