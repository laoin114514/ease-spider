package utils

import (
	"fmt"
	"spider/config"
)

// ConfigValidator 配置验证器
type ConfigValidator struct{}

// NewConfigValidator 创建配置验证器
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// ValidateConfig 验证配置
func (v *ConfigValidator) ValidateConfig() error {
	if config.AppConfig == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 验证数据库配置
	if err := v.validateDatabaseConfig(); err != nil {
		return err
	}

	// 验证洛谷配置
	if err := v.validateLuoguConfig(); err != nil {
		return err
	}

	// 验证钉钉配置
	if err := v.validateDingdingConfig(); err != nil {
		return err
	}

	return nil
}

// validateDatabaseConfig 验证数据库配置
func (v *ConfigValidator) validateDatabaseConfig() error {
	db := config.AppConfig.Database
	if db.Host == "" {
		return fmt.Errorf("数据库主机地址不能为空")
	}
	if db.Port == "" {
		return fmt.Errorf("数据库端口不能为空")
	}
	if db.User == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if db.Password == "" {
		return fmt.Errorf("数据库密码不能为空")
	}
	if db.DbName == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	return nil
}

// validateLuoguConfig 验证洛谷配置
func (v *ConfigValidator) validateLuoguConfig() error {
	luogu := config.AppConfig.Luogu
	if luogu.Username == "" {
		return fmt.Errorf("洛谷用户名不能为空")
	}
	if luogu.Password == "" {
		return fmt.Errorf("洛谷密码不能为空")
	}
	if luogu.UserAgent == "" {
		return fmt.Errorf("洛谷UserAgent不能为空")
	}
	return nil
}

// validateDingdingConfig 验证钉钉配置
func (v *ConfigValidator) validateDingdingConfig() error {
	dingding := config.AppConfig.Dingding
	if dingding.AppKey == "" {
		return fmt.Errorf("钉钉AppKey不能为空")
	}
	if dingding.AppSecret == "" {
		return fmt.Errorf("钉钉AppSecret不能为空")
	}
	return nil
}
