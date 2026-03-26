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
type luoguConfig struct {
	Username                         string `yaml:"username"`
	Password                         string `yaml:"password"`
	UserAgent                        string `yaml:"userAgent"`
	IsInServer                       bool   `yaml:"isInServer"`
	LuoguRecordsConcurrency          int    `yaml:"luoguRecordsConcurrency"`
	LuoguSolutionConcurrency         int    `yaml:"luoguSolutionConcurrency"`
	LuoguSubmissionDetailConcurrency int    `yaml:"luoguSubmissionDetailConcurrency"`
	LuoguTeamID                      int    `yaml:"luoguTeamID"`
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
	LuoguUpdateCookie      string `yaml:"luogu_update_cookie"`
	LuoguSubmissionDetail  string `yaml:"luogu_submission_detail"`
	LuoguSolution          string `yaml:"luogu_solution"`
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
		path = "/config/config.dev.yml"
	} else if os.Getenv("RUN_MODE") == "prod" {
		path = "/config/config.prod.yml"
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
