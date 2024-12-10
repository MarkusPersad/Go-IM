package database

import (
	"Go-IM/pkg/zaplog"
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	dbHost           = os.Getenv("DB_HOST")
	dbPort           = os.Getenv("DB_PORT")
	dbUser           = os.Getenv("DB_USER")
	dbPass           = os.Getenv("DB_PASS")
	dbName           = os.Getenv("DB_NAME")
	maxOpenConns, _  = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _  = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	dbMaxLifetime, _ = strconv.Atoi(os.Getenv("DB_CONN_MAX_LIFETIME"))
	valHost          = os.Getenv("VALKEY_CLIENT_HOST")
	valPort          = os.Getenv("VALKEY_CLIENT_PORT")
	valPass          = os.Getenv("VALKEY_CLIENT_PASS")
	dbInstance       Service
)

type Service interface {
	Health() map[string]string
	Close() error
	Set(id string, value string) error
	Get(id string, clear bool) string
	Verify(id, answer string, clear bool) bool
	GetDB() *gorm.DB
	GetValClient() valkey.Client
}

type service struct {
	db        *gorm.DB
	valClient valkey.Client
}

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)
	isSingularTable, _ := strconv.ParseBool(os.Getenv("DB_SINGULAR_TABLE"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: isSingularTable,
		},
	})
	if err != nil {
		zaplog.Logger.Error("Failed to connect to database.", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		zaplog.Logger.Error("Failed to connect to database.", zap.Error(err))
		os.Exit(2)
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbMaxLifetime) * time.Second)

	valClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valHost + ":" + valPort},
		Password:    valPass,
	})
	if err != nil {
		zaplog.Logger.Error("Failed to connect to valkey.", zap.Error(err))
		os.Exit(2)
	}
	dbInstance = &service{db: db, valClient: valClient}
	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stats := make(map[string]string)
	sqlDB, err := s.db.DB()
	if err != nil {
		stats["mariadb_status"] = "down"
		stats["mariadb_error"] = fmt.Sprintf("mariadb down: %v", err)
		zaplog.Logger.Fatal("mariadb down", zap.Error(err))
		return stats
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		stats["mariadb_status"] = "down"
		stats["mariadb_error"] = fmt.Sprintf("mariadb down: %v", err)
		zaplog.Logger.Fatal("mariadb down", zap.Error(err))
		return stats
	}
	stats["mariadb_status"] = "up"
	stats["mariadb_message"] = "It's healthy"

	dbStats := sqlDB.Stats()
	stats["mariadb_open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["mariadb_in_use"] = strconv.Itoa(dbStats.InUse)
	stats["mariadb_idle"] = strconv.Itoa(dbStats.Idle)
	stats["mariadb_wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["mariadb_wait_duration"] = dbStats.WaitDuration.String()
	stats["mariadb_max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["mariadb_max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["mariadb_message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["mariadb_message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["mariadb_message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["mariadb_message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}
	valResult := s.valClient.Do(ctx, s.valClient.B().Ping().Build())
	if valResult.Error() != nil {
		stats["valkey_status"] = "down"
		stats["valkey_error"] = fmt.Sprintf("valkey down: %v", valResult.Error())
		zaplog.Logger.Error("valkey down", zap.Error(valResult.Error()))
		return stats
	}
	stats["valkey_status"] = "up"
	stats["valkey_message"] = "It's healthy"
	valStatus := parseValkeyInfo(valResult.String())
	for k, v := range valStatus {
		stats[k] = v
	}
	return stats
}
func (s *service) Close() error {
	sqlDB, _ := s.db.DB()
	zaplog.Logger.Info("Closing database connection", zap.String("db", dbName))
	s.valClient.Close()
	zaplog.Logger.Info("Closing valkey connection")
	return sqlDB.Close()
}
func parseValkeyInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}

func (s *service) GetDB() *gorm.DB {
	return s.db
}

func (s *service) GetValClient() valkey.Client {
	return s.valClient
}
