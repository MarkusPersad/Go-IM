package database

import (
	"Go-IM/pkg/zaplog"
	"go.uber.org/zap"
)

func (s *service) SetAndTime(key, value string, timeout int64) error {
	return s.valClient.Do(ctx, s.valClient.B().Setex().Key(key).Seconds(timeout).Value(value).Build()).Error()
}
func (s *service) GetValue(key string) string {
	result := s.valClient.Do(ctx, s.valClient.B().Get().Key(key).Build())
	if result.Error() != nil {
		return ""
	}
	val, e := result.ToString()
	if e != nil {
		zaplog.Logger.Error("valkey get error", zap.Error(e))
		return ""
	}
	return val
}

// DelValue 删除指定键的值。
// key: 需要删除的键。
// 返回: 如果删除操作失败，返回错误；否则返回nil。
func (s *service) DelValue(key string) error {
	// 执行删除操作，并检查错误。
	if e := s.valClient.Do(ctx, s.valClient.B().Del().Key(key).Build()).Error(); e != nil {
		// 如果发生错误，记录错误日志并返回错误。
		zaplog.Logger.Error("valkey delete error", zap.Error(e))
		return e
	}
	// 删除操作成功，返回nil。
	return nil
}
