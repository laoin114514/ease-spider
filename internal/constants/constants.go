package constants

// 洛谷相关常量
const (
	LuoguStatusAccepted = 12
)

// Codeforces相关常量
const (
	CfMaxRecords = 50000
)

// 验证码识别相关常量
const (
	CaptchaLengthServer = 5
	CaptchaLengthLocal  = 6
	CaptchaStartIndex   = 0
)

// 时间相关常量
const (
	DefaultTimeoutSeconds = 60
	TimeZoneOffsetHours   = 8
)

// 难度等级映射
var LuoguDifficultyMap = []string{
	"grey", "red", "brown", "yellow", "green", "blue", "purple", "black",
}
