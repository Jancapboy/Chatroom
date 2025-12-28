// internal/setting/setting.go
package setting

import (
	"fmt"

	"github.com/spf13/viper"
)

// AI配置结构体
type AIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

type Setting struct {
	vp *viper.Viper
	AI AIConfig
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.AddConfigPath("configs/")
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	s := &Setting{vp: vp}

	// 读取AI配置
	if err := vp.UnmarshalKey("AI", &s.AI); err != nil {
		return nil, fmt.Errorf("解析AI配置失败: %w", err)
	}

	// 验证配置
	if s.AI.APIKey == "" {
		return nil, fmt.Errorf("AI API密钥未配置")
	}

	// 设置默认值
	if s.AI.Model == "" {
		s.AI.Model = "deepseek-chat"
	}
	if s.AI.BaseURL == "" {
		s.AI.BaseURL = "https://api.deepseek.com"
	}

	return s, nil
}

// 获取AI配置的方法
func (s *Setting) GetAIConfig() AIConfig {
	return s.AI
}
