package utils

import (
	"testing"

	"spider/config"
)

func TestValidateTimerFrequency(t *testing.T) {
	old := config.AppConfig
	defer func() { config.AppConfig = old }()
	v := NewConfigValidator()

	// 合法取值 + 留空（留空表示用插件内置默认值）都应通过
	config.AppConfig = &config.Config{TimerFrequency: config.TimerFrequencyConfig{
		CfRecords:    "2m",
		Dingding:     "",
		NiukeRecords: "1h",
	}}
	if err := v.validateTimerFrequency(); err != nil {
		t.Fatalf("合法取值不应报错: %v", err)
	}

	// 非法取值必须在启动期被拦下
	config.AppConfig = &config.Config{TimerFrequency: config.TimerFrequencyConfig{CfRecords: "abc"}}
	if err := v.validateTimerFrequency(); err == nil {
		t.Fatal("非法取值应报错")
	}
}
