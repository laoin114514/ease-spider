package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Database       dbConfig             `yaml:"database"`
	Dingding       dingdingConfig       `yaml:"dingding"`
	Luogu          luoguConfig          `yaml:"luogu"`
	Luogu2Api      luogu2apiConfig      `yaml:"luogu2api"`
	Cf             cfConfig             `yaml:"cf"`
	TimerFrequency TimerFrequencyConfig `yaml:"timerFrequency"`
	DebugConfig    debugConfig          `yaml:"debug"`
}
type dbConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
}
type dingdingConfig struct {
	AppKey    string `yaml:"appKey"`
	AppSecret string `yaml:"appSecret"`
	WeekRange int    `yaml:"weekRange"`
}

// luoguConfig 洛谷抓取相关配置
type luoguConfig struct {
	// UserAgent 通用 User-Agent：internal/utils/request.go 的所有请求都会带上它
	UserAgent               string `yaml:"userAgent"`
	LuoguRecordsConcurrency int    `yaml:"luoguRecordsConcurrency"`
}

// luogu2apiConfig 是 Luogu2Api 服务（pkg/luogu2api SDK 的服务端）的连接配置
type luogu2apiConfig struct {
	// BaseURL 服务地址，必须带 http:// 或 https:// 前缀，例如 http://127.0.0.1:8080
	BaseURL string `yaml:"baseUrl"`
	// AdminToken 需与服务端的 ADMIN_TOKEN 一致，随 X-Admin-Token 请求头发送
	AdminToken string `yaml:"adminToken"`
}
type cfConfig struct {
	GroupCode                        string `yaml:"groupCode"`
	ManagerAccount                   string `yaml:"managerAccount"`
	CfRecordsConcurrency             int    `yaml:"cfRecordsConcurrency"`
	CfTeamContestProblemsConcurrency int    `yaml:"cfTeamContestProblemsConcurrency"`
}

// TimerFrequencyConfig 各定时任务的执行周期。
//
// 取值形如 30s / 5m / 2h / 1d，不带单位（或单位无法识别）时按分钟处理；
// 留空表示不配置，插件会退回自己的默认间隔。
type TimerFrequencyConfig struct {
	CfOfficialContests     string `yaml:"cf_official_contests"`
	CfOfficialProblems     string `yaml:"cf_official_problems"`
	CfRecords              string `yaml:"cf_records"`
	CfTeamContests         string `yaml:"cf_team_contests"`
	CfTeamContestsProblems string `yaml:"cf_team_contests_problems"`
	Dingding               string `yaml:"dingding"`
	LuoguRecords           string `yaml:"luogu_records"`
	NiukeRecords           string `yaml:"niuke_records"`
}

// Entry 是 timerFrequency 的一项（键名 + 原始取值），用于启动期校验与文档。
type Entry struct {
	Key   string
	Value string
}

// Entries 按固定顺序列出所有定时任务配置项。
func (t TimerFrequencyConfig) Entries() []Entry {
	return []Entry{
		{Key: "cf_official_contests", Value: t.CfOfficialContests},
		{Key: "cf_official_problems", Value: t.CfOfficialProblems},
		{Key: "cf_records", Value: t.CfRecords},
		{Key: "cf_team_contests", Value: t.CfTeamContests},
		{Key: "cf_team_contests_problems", Value: t.CfTeamContestsProblems},
		{Key: "dingding", Value: t.Dingding},
		{Key: "luogu_records", Value: t.LuoguRecords},
		{Key: "niuke_records", Value: t.NiukeRecords},
	}
}

// TimerFrequency 返回定时任务配置；配置尚未初始化时返回零值，
// 此时所有插件都会退回各自的默认间隔。
func TimerFrequency() TimerFrequencyConfig {
	if AppConfig == nil {
		return TimerFrequencyConfig{}
	}
	return AppConfig.TimerFrequency
}

// ParseInterval 解析 timerFrequency 的取值：
// 支持 30s / 5m / 2h / 1d；不带单位或单位无法识别时按分钟处理（与历史实现一致）。
// ok 为 false 表示取值为空或完全无法解析。
func ParseInterval(raw string) (time.Duration, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, false
	}
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	n, err := strconv.Atoi(s[:i])
	if err != nil || n <= 0 {
		return 0, false
	}
	switch s[i:] {
	case "s":
		return time.Duration(n) * time.Second, true
	case "m":
		return time.Duration(n) * time.Minute, true
	case "h":
		return time.Duration(n) * time.Hour, true
	case "d":
		return time.Duration(n) * 24 * time.Hour, true
	default:
		return time.Duration(n) * time.Minute, true
	}
}

// IntervalOr 解析 raw，取值非法或为空时返回 def。
func IntervalOr(raw string, def time.Duration) time.Duration {
	if d, ok := ParseInterval(raw); ok {
		return d
	}
	return def
}

type debugConfig struct {
	All bool `yaml:"all"`
}

var AppConfig *Config

func Init() error {
	err := godotenv.Load("config/.env")
	if err != nil {
		return err
	}
	path := ""
	if os.Getenv("RUN_MODE") == "dev" {
		path = "config/config.dev.yml"
	} else if os.Getenv("RUN_MODE") == "prod" {
		path = "config/config.prod.yml"
	} else {
		return fmt.Errorf("环境变量RUN_MODE必须为dev或prod")
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &AppConfig)
	if err != nil {
		return err
	}
	return nil
}
