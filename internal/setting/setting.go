// internal/setting/setting.go
package setting

import (
	"fmt"
	//"time"

	"github.com/spf13/viper"
)

// 1. 添加 AI 配置结构体
type AISettingS struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

// 2. 修改 Setting 结构体（如果还没有的话）
type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.AddConfigPath("configs/")
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	return &Setting{vp}, nil
}

// 3. 添加读取 AI 配置的方法

// 4. 添加 GetAIConfig 方法
func (s *Setting) GetAIConfig() (*AISettingS, error) {
	aiConfig := &AISettingS{}
	if err := s.ReadSection("AI", aiConfig); err != nil {
		return nil, err
	}

	// 设置默认值
	if aiConfig.Model == "" {
		aiConfig.Model = "deepseek-chat"
	}
	if aiConfig.BaseURL == "" {
		aiConfig.BaseURL = "https://api.deepseek.com"
	}

	return aiConfig, nil
}
