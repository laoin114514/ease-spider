package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Database       dbConfig             `yaml:"database"`
	Dingding       dingdingConfig       `yaml:"dingding"`
	Luogu          luoguConfig          `yaml:"luogu"`
	Luogu2Api      luogu2apiConfig      `yaml:"luogu2api"`
	Cf             cfConfig             `yaml:"cf"`
	TimerFrequency timerFrequencyConfig `yaml:"timerFrequency"`
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
type timerFrequencyConfig struct {
	CfOfficialContests     string `yaml:"cf_official_contests"`
	CfOfficialProblems     string `yaml:"cf_official_problems"`
	CfRecords              string `yaml:"cf_records"`
	CfTeamContests         string `yaml:"cf_team_contests"`
	CfTeamContestsProblems string `yaml:"cf_team_contests_problems"`
	Dingding               string `yaml:"dingding"`
	LuoguRecords           string `yaml:"luogu_records"`
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
