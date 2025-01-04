package giutils

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
)

var admins []string

func init() {
	adminString := os.Getenv("ADMIN")
	if adminString != "" && adminString != " " {
		admins = strings.Split(strings.TrimSpace(adminString), ",")
	}
	//if os.Getenv("APP_ENV") == "local" {
	//	zaplog.Logger.Info("ADMIN:", zap.String("ADMIN", adminString))
	//}
}

// IsAdmin 检查给定字符串是否表示一个管理员身份。
// 该函数通过调用 contains 函数来判断给定的字符串是否存在于预定义的管理员列表中。
// 参数:
//
//	str: 待检查的字符串。
//
// 返回值:
//
//	如果给定的字符串存在于管理员列表中，则返回 true，否则返回 false。
func IsAdmin(str string) bool {
	// 使用 contains 函数检查给定字符串是否在管理员列表中。
	if contains(admins, str) {
		return true
	}
	return false
}

// contains 函数用于检查一个集合中是否包含特定的元素。
// 它使用了泛型类型 T，只要求 T 类型是可比较的。
// 这个函数通过遍历集合中的每个元素来实现，如果找到了匹配的元素，则返回 true，否则返回 false。
// 参数:
//
//	collection: 一个 T 类型的切片，代表要搜索的集合。
//	element: 一个 T 类型的元素，代表要查找的目标。
//
// 返回值:
//
//	bool: 如果集合中包含目标元素，则返回 true，否则返回 false。
func contains[T comparable](collection []T, element T) bool {
	// 遍历集合中的每个元素，与目标元素进行比较。
	for _, value := range collection {
		// 如果找到匹配的元素，立即返回 true。
		if value == element {
			return true
		}
	}
	// 如果遍历完集合后没有找到匹配的元素，返回 false。
	return false
}
