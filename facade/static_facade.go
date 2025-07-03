package facade

import (
	"context"
	"reflect"
)

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
