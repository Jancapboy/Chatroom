// global/setting.go
package global

import (
	"github.com/Jancapboy/Chatroom/internal/setting"
)

var (
	ServerSettings   *setting.ServerSetting
	DatabaseSettings *setting.DatabaseSetting
	JWTSettings      *setting.JWTSetting
	ChatroomSettings *setting.ChatroomSetting
	AISettings       *setting.AISettingS // 新增 AI 配置
)
