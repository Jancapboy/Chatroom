package model

import (
	"time"

	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"gorm.io/gorm"
)

// Message 房间消息模型
type Message struct {
	ID           string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	RoomID       string    `json:"room_id" gorm:"type:varchar(36);not null;index:idx_room_round"`
	SenderID     string    `json:"sender_id" gorm:"type:varchar(36);not null"`
	SenderType   string    `json:"sender_type" gorm:"type:varchar(16);not null"` // agent, user, system
	SenderName   string    `json:"sender_name" gorm:"type:varchar(100)"`
	SenderAvatar string    `json:"sender_avatar" gorm:"type:varchar(255)"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	MsgType      string    `json:"msg_type" gorm:"type:varchar(32);default:'text'"` // text, decision, consensus, phase_change, fork_notice, system
	Phase        string    `json:"phase" gorm:"type:varchar(32)"`
	Round        int       `json:"round" gorm:"default:1;index:idx_room_round"`
	Metadata     string    `json:"metadata" gorm:"type:json"` // JSON string: {confidence, stance, votes}
	CreatedAt    time.Time `json:"created_at"`
}

func (m Message) Create(db *gorm.DB) *errcode.Error {
	err := db.Create(&m).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (m Message) ListByRoom(db *gorm.DB, roomID string, round int, phase string, page, pageSize int) ([]Message, int64, *errcode.Error) {
	var messages []Message
	var total int64

	query := db.Model(&Message{}).Where("room_id = ?", roomID)
	if round > 0 {
		query = query.Where("round = ?", round)
	}
	if phase != "" {
		query = query.Where("phase = ?", phase)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errcode.Convert(err)
	}

	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&messages).Error
	if err != nil {
		return nil, 0, errcode.Convert(err)
	}
	return messages, total, nil
}

func (m Message) GetPhaseMessages(db *gorm.DB, roomID, phase string, round int) ([]Message, *errcode.Error) {
	var messages []Message
	err := db.Where("room_id = ? AND phase = ? AND round = ?", roomID, phase, round).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return messages, nil
}
