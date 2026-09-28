package getuserrecords

import (
	"testing"
	"time"

	"spider/config"
)

func TestMetaIntervalComesFromConfig(t *testing.T) {
	old := config.AppConfig
	defer func() { config.AppConfig = old }()
	config.AppConfig = &config.Config{TimerFrequency: config.TimerFrequencyConfig{LuoguRecords: "3m"}}

	if got := (&GetUserRecords{}).Meta().Interval; got != 3*time.Minute {
		t.Fatalf("Meta().Interval = %v, want 3m（应来自 timerFrequency.luogu_records）", got)
	}
}

func TestMetaIntervalFallsBackToDefault(t *testing.T) {
	old := config.AppConfig
	defer func() { config.AppConfig = old }()
	config.AppConfig = &config.Config{}

	if got := (&GetUserRecords{}).Meta().Interval; got != 5*time.Minute {
		t.Fatalf("未配置时 Meta().Interval = %v, want 5m（内置默认值）", got)
	}
}
