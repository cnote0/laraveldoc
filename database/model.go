package database

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Model 基础模型结构体
//
// Model 提供了标准的数据库模型字段，包括主键、时间戳和软删除支持。
// 所有业务模型都应该嵌入此结构体。
//
// 使用示例：
//
//	type User struct {
//		database.Model
//		Name  string `json:"name"`
//		Email string `json:"email" gorm:"unique"`
//	}
//
//	type Product struct {
//		database.Model
//		Name        string          `json:"name"`
//		Price       decimal.Decimal `json:"price"`
//		CategoryID  uint            `json:"category_id"`
//	}
type Model struct {
	// ID 主键
	//
	// 使用 uint 类型的自增主键，符合 GORM 的默认约定。
	//
	// 示例：
	//   ID: 1, 2, 3, ...
	ID uint `gorm:"primarykey" json:"id"`

	// CreatedAt 创建时间
	//
	// 记录创建时的时间戳，GORM 会自动设置此字段。
	//
	// 示例：
	//   CreatedAt: 2023-01-01 12:00:00
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	//
	// 记录最后更新的时间戳，GORM 会在每次更新时自动设置。
	//
	// 示例：
	//   UpdatedAt: 2023-01-02 15:30:00
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt 软删除时间
	//
	// 使用自定义的 DeletedAt 类型，支持 JSON 序列化和数据库交互。
	// 当记录被软删除时，此字段会被设置为删除时间。
	//
	// 示例：
	//   DeletedAt: {Time: 2023-01-03 10:00:00, Valid: true}  // 已删除
	//   DeletedAt: {Time: time.Time{}, Valid: false}         // 未删除
	DeletedAt DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// DeletedAt 软删除时间字段
//
// DeletedAt 实现了 GORM 的软删除功能，同时支持 JSON 序列化。
// 它实现了 sql.Scanner 和 driver.Valuer 接口，可以与数据库交互。
//
// 使用示例：
//
//	var deletedAt DeletedAt
//
//	// 设置为已删除
//	deletedAt = DeletedAt{Time: time.Now(), Valid: true}
//
//	// 设置为未删除
//	deletedAt = DeletedAt{Valid: false}
//
//	// JSON 序列化
//	jsonData, _ := json.Marshal(deletedAt)
//	// 输出: "2023-01-01T12:00:00Z" 或 null
type DeletedAt struct {
	// Time 删除时间
	//
	// 记录被软删除的具体时间。
	Time time.Time

	// Valid 是否有效
	//
	// true: 记录已被软删除，Time 字段有效
	// false: 记录未被删除，Time 字段无效
	Valid bool
}

// Scan 实现 sql.Scanner 接口
//
// 将数据库中的值扫描到 DeletedAt 结构体中。
// 支持 NULL 值和时间戳的正确处理。
//
// 示例：
//
//	var dt DeletedAt
//	err := dt.Scan("2023-01-01 12:00:00")     // 设置时间
//	err := dt.Scan(nil)                       // 设置为 NULL
func (dt *DeletedAt) Scan(value interface{}) error {
	if value == nil {
		dt.Time, dt.Valid = time.Time{}, false
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		dt.Time, dt.Valid = v, true
		return nil
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		dt.Time, dt.Valid = t, true
		return nil
	}
	dt.Time, dt.Valid = time.Time{}, false
	return nil
}

// Value 实现 driver.Valuer 接口
//
// 将 DeletedAt 结构体转换为数据库可存储的值。
//
// 示例：
//
//	dt := DeletedAt{Time: time.Now(), Valid: true}
//	value, err := dt.Value()  // 返回 time.Time
//
//	dt := DeletedAt{Valid: false}
//	value, err := dt.Value()  // 返回 nil
func (dt DeletedAt) Value() (driver.Value, error) {
	if !dt.Valid {
		return nil, nil
	}
	return dt.Time, nil
}

// MarshalJSON 实现 JSON 序列化
//
// 将 DeletedAt 转换为 JSON 格式。
// 有效时间会序列化为 ISO8601 格式，无效时间序列化为 null。
//
// 示例：
//
//	dt := DeletedAt{Time: time.Now(), Valid: true}
//	json, _ := dt.MarshalJSON()  // "2023-01-01T12:00:00Z"
//
//	dt := DeletedAt{Valid: false}
//	json, _ := dt.MarshalJSON()  // null
func (dt DeletedAt) MarshalJSON() ([]byte, error) {
	if !dt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(dt.Time)
}

// UnmarshalJSON 实现 JSON 反序列化
//
// 从 JSON 数据反序列化到 DeletedAt 结构体。
// 支持 null 值和 ISO8601 时间格式。
//
// 示例：
//
//	var dt DeletedAt
//	dt.UnmarshalJSON([]byte(`"2023-01-01T12:00:00Z"`))  // 设置时间
//	dt.UnmarshalJSON([]byte(`null`))                    // 设置为无效
func (dt *DeletedAt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		dt.Valid = false
		return nil
	}
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	dt.Time = t
	dt.Valid = true
	return nil
}
