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
