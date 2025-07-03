package database

// Association 关联操作接口
//
// Association 提供了模型关联关系的操作方法。
type Association interface {
	Find(out interface{}, conds ...interface{}) error
	Append(values ...interface{}) error
	Replace(values ...interface{}) error
	Delete(values ...interface{}) error
	Clear() error
	Count() int64
}
