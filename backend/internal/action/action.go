package action

import (
	"time"
)

// Action 执行操作接口
type Action interface {
	Name() string                                              // 操作名称
	Description() string                                     // 描述
	Icon() string                                            // 图标标识
	Validate(params map[string]interface{}) error            // 参数校验
	Execute(params map[string]interface{}) (string, error)   // 执行操作，返回结果
}

// ActionRegistry Action注册表
var ActionRegistry = make(map[string]Action)

// RegisterAction 注册Action
func RegisterAction(action Action) {
	ActionRegistry[action.Name()] = action
}

// GetAction 获取Action
func GetAction(name string) (Action, bool) {
	action, ok := ActionRegistry[name]
	return action, ok
}

// ListActions 列出所有可用Action
func ListActions() []ActionInfo {
	var list []ActionInfo
	for _, action := range ActionRegistry {
		list = append(list, ActionInfo{
			Name:        action.Name(),
			Description: action.Description(),
			Icon:        action.Icon(),
		})
	}
	return list
}

// ActionInfo Action信息（用于前端展示）
type ActionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// ActionDraft Action草稿（推演完成后生成）
type ActionDraft struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`        // action名称
	Title       string                 `json:"title"`       // 标题
	Description string                 `json:"description"` // 描述
	Params      map[string]interface{} `json:"params"`      // 执行参数
	Status      string                 `json:"status"`      // pending / approved / executed / failed
	CreatedAt   time.Time              `json:"created_at"`
}

func init() {
	// 注册内置Actions
	RegisterAction(&MarkdownReportAction{})
	RegisterAction(&FeishuMessageAction{})
}

// contains helper
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsSub(s, substr))
}

func containsSub(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
