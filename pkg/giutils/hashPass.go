package giutils

import (
	"Go-IM/pkg/zaplog"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

var cost int

func init() {
	if val, e := strconv.Atoi(os.Getenv("HASH_SALT")); e == nil {
		cost = val
	} else {
		cost = bcrypt.DefaultCost
	}
}

func GenerateHashPass(password string) string {

	hashPassword, e := bcrypt.GenerateFromPassword([]byte(password), cost)
	if e != nil {
		zaplog.Logger.Error("hash password failed", zap.Error(e))
		return ""
	}
	return string(hashPassword)
}

func CompareHashPassword(hashPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)) == nil
}
