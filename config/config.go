package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Database       dbConfig             `yaml:"database"`
	Dingding       dingdingConfig       `yaml:"dingding"`
	Luogu          luoguConfig          `yaml:"luogu"`
	Cf             cfConfig             `yaml:"cf"`
	TimerFrequency timerFrequencyConfig `yaml:"timerFrequency"`
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
}
type luoguConfig struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	UserAgent  string `yaml:"userAgent"`
	IsInServer bool   `yaml:"isInServer"`
}
type cfConfig struct {
	GroupCode string `yaml:"groupCode"`
}
type timerFrequencyConfig struct {
	CfOfficialContests     string `yaml:"cf_official_contests"`
	CfOfficialProblems     string `yaml:"cf_official_problems"`
	CfRecords              string `yaml:"cf_records"`
	CfTeamContests         string `yaml:"cf_team_contests"`
	CfTeamContestsProblems string `yaml:"cf_team_contests_problems"`
	Dingding               string `yaml:"dingding"`
	LuoguRecords           string `yaml:"luogu_records"`
	LuoguUpdateCookie      string `yaml:"luogu_update_cookie"`
}

var AppConfig *Config

func Init() error {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &AppConfig)
	if err != nil {
		return err
	}
	return nil
}
