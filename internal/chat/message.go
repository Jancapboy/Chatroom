package chat

import "time"

// 消息类型常量
const (
	MsgTypeText    = "text"
	MsgType3DModel = "3d_model"
)

type Message struct {
	User     *User  `json:"user"`
	Type     string `json:"type"`
	Content  string `json:"message_content"`
	ModelUrl string `json:"model_url,omitempty"`
	SendTime int64  `json:"send_time"`
}

func NewMessage(user *User, content string) *Message {
	return &Message{
		User:     user,
		Type:     MsgTypeText,
		Content:  content,
		SendTime: time.Now().Unix(),
	}
}

func New3DMessage(user *User, content string, modelUrl string) *Message {
	return &Message{
		User:     user,
		Type:     MsgType3DModel,
		Content:  content,
		ModelUrl: modelUrl,
		SendTime: time.Now().Unix(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:     System,
		Type:     MsgTypeText,
		Content:  "欢迎 " + user.Nickname + " 加入聊天室",
		SendTime: time.Now().Unix(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:     System,
		Type:     MsgTypeText,
		Content:  user.Nickname + " 离开了聊天室",
		SendTime: time.Now().Unix(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:     System,
		Type:     MsgTypeText,
		Content:  content,
		SendTime: time.Now().Unix(),
	}
}
