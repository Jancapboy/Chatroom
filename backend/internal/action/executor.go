package action

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// MarkdownReportAction 生成Markdown报告
type MarkdownReportAction struct{}

func (a *MarkdownReportAction) Name() string        { return "markdown_report" }
func (a *MarkdownReportAction) Description() string { return "生成Markdown推演报告" }
func (a *MarkdownReportAction) Icon() string        { return "file-text" }

func (a *MarkdownReportAction) Validate(params map[string]interface{}) error {
	if _, ok := params["room_name"]; !ok {
		return fmt.Errorf("缺少参数: room_name")
	}
	if _, ok := params["messages"]; !ok {
		return fmt.Errorf("缺少参数: messages")
	}
	return nil
}

func (a *MarkdownReportAction) Execute(params map[string]interface{}) (string, error) {
	roomName := getString(params, "room_name", "未命名推演")
	topic := getString(params, "topic", "")
	conclusion := getString(params, "conclusion", "")
	messages := getStringSlice(params, "messages")
	
	// 生成报告内容
	content := a.generateReport(roomName, topic, conclusion, messages)
	
	// 写入文件
	reportsDir := "/tmp/chatroom-reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return "", fmt.Errorf("创建报告目录失败: %w", err)
	}
	
	filename := fmt.Sprintf("%s/%s_%s.md", reportsDir, 
		sanitizeFilename(roomName), 
		time.Now().Format("20060102_150405"))
	
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("写入报告失败: %w", err)
	}
	
	return fmt.Sprintf("报告已生成: %s", filename), nil
}

func (a *MarkdownReportAction) generateReport(roomName, topic, conclusion string, messages []string) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("# 推演报告: %s\n\n", roomName))
	sb.WriteString(fmt.Sprintf("**生成时间:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	
	if topic != "" {
		sb.WriteString(fmt.Sprintf("**推演主题:** %s\n\n", topic))
	}
	
	sb.WriteString("## 推演过程\n\n")
	for i, msg := range messages {
		sb.WriteString(fmt.Sprintf("### 记录 %d\n\n%s\n\n", i+1, msg))
	}
	
	if conclusion != "" {
		sb.WriteString("## 推演结论\n\n")
		sb.WriteString(conclusion)
		sb.WriteString("\n\n")
	}
	
	sb.WriteString("---\n\n")
	sb.WriteString("*由 ASI Chatroom 多智能体推演平台自动生成*\n")
	
	return sb.String()
}

// FeishuMessageAction 发送飞书消息（需要配置Token）
type FeishuMessageAction struct{}

func (a *FeishuMessageAction) Name() string        { return "feishu_message" }
func (a *FeishuMessageAction) Description() string { return "发送飞书消息通知" }
func (a *FeishuMessageAction) Icon() string        { return "send" }

func (a *FeishuMessageAction) Validate(params map[string]interface{}) error {
	if _, ok := params["content"]; !ok {
		return fmt.Errorf("缺少参数: content")
	}
	return nil
}

func (a *FeishuMessageAction) Execute(params map[string]interface{}) (string, error) {
	content := getString(params, "content", "")
	if content == "" {
		return "", fmt.Errorf("消息内容为空")
	}
	
	// TODO: 接入飞书API（需要Bot Token）
	// 目前返回模拟结果
	return fmt.Sprintf("飞书消息已发送（模拟）: %s", truncate(content, 50)), nil
}

// helper functions
func getString(params map[string]interface{}, key, defaultVal string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func getStringSlice(params map[string]interface{}, key string) []string {
	if v, ok := params[key]; ok {
		if arr, ok := v.([]string); ok {
			return arr
		}
		if arr, ok := v.([]interface{}); ok {
			var result []string
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_", "|", "_",
	)
	return replacer.Replace(name)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
