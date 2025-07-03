package facade

// SpyInterface 间谍接口
//
// SpyInterface 用于监控方法调用，记录调用信息但不改变行为。
// 主要用于测试中验证方法是否被正确调用。
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
// CallVerifier 提供方法调用的详细验证功能，
// 包括次数验证、参数验证等。
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

// MockInterface 模拟接口
//
// MockInterface 提供门面的模拟功能，用于测试中替换真实实现。
// 支持完全模拟、部分模拟和间谍模式。
//
// 使用示例：
//
//	// 在测试中模拟邮件发送
//	func TestUserRegistration(t *testing.T) {
//		// 模拟邮件门面
//		mockMailer := &MockMailer{}
//		mockMailer.On("Send", mock.Anything, mock.Anything).Return(nil)
//		Mail.Mock(mockMailer)
//		defer Mail.ClearMock()
//
//		// 执行测试
//		userService := &UserService{}
//		err := userService.Register("user@example.com")
//
//		// 验证
//		assert.NoError(t, err)
//		mockMailer.AssertExpectations(t)
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
