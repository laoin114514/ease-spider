package utils

import (
	"fmt"
	"spider/config"
)

// ============================配置验证器===============================================
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

	// 验证 Luogu2Api SDK 配置
	if err := v.validateLuogu2ApiConfig(); err != nil {
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

// validateLuoguConfig 验证洛谷配置。
// 账号密码校验已随取 cookie 的定时服务一并删除：登录态现在由 Luogu2Api 服务的号池维护。
func (v *ConfigValidator) validateLuoguConfig() error {
	luogu := config.AppConfig.Luogu
	if luogu.UserAgent == "" {
		return fmt.Errorf("洛谷UserAgent不能为空")
	}
	if luogu.LuoguRecordsConcurrency <= 0 {
		return fmt.Errorf("洛谷提交记录并发数必须大于0")
	}
	return nil
}

// validateLuogu2ApiConfig 验证 Luogu2Api SDK 配置。
// 这两项是 SDK 的必填参数：地址为空或没有协议头、令牌为空时构造客户端就会失败，
// 提前在启动阶段拦下，避免定时任务每一轮都失败。
func (v *ConfigValidator) validateLuogu2ApiConfig() error {
	luogu2api := config.AppConfig.Luogu2Api
	if luogu2api.BaseURL == "" {
		return fmt.Errorf("luogu2api baseUrl不能为空")
	}
	if luogu2api.AdminToken == "" {
		return fmt.Errorf("luogu2api adminToken不能为空（需与Luogu2Api服务端的ADMIN_TOKEN一致）")
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
