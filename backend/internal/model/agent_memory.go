package model

import (
	"time"

	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"gorm.io/gorm"
)

// AgentMemory Agent长期记忆（跨房间学习）
type AgentMemory struct {
	ID         string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	AgentID    string    `json:"agent_id" gorm:"type:varchar(36);not null;index:idx_agent_topic"`
	RoomID     string    `json:"room_id" gorm:"type:varchar(36);index"`
	Topic      string    `json:"topic" gorm:"type:varchar(128);index:idx_agent_topic"` // 话题标签，如"微服务架构"
	Experience string    `json:"experience" gorm:"type:text"`                          // 经验文本
	Outcome    string    `json:"outcome" gorm:"type:varchar(32)"`                        // success / failure / neutral
	Confidence int       `json:"confidence"`                                             // 0-100
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (m AgentMemory) Create(db *gorm.DB) *errcode.Error {
	err := db.Create(&m).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (m AgentMemory) ListByAgent(db *gorm.DB, agentID string, limit int) ([]AgentMemory, *errcode.Error) {
	var memories []AgentMemory
	err := db.Where("agent_id = ?", agentID).Order("created_at DESC").Limit(limit).Find(&memories).Error
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return memories, nil
}

func (m AgentMemory) ListByAgentAndTopic(db *gorm.DB, agentID, topic string, limit int) ([]AgentMemory, *errcode.Error) {
	var memories []AgentMemory
	err := db.Where("agent_id = ? AND topic = ?", agentID, topic).Order("created_at DESC").Limit(limit).Find(&memories).Error
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return memories, nil
}

func (m AgentMemory) Delete(db *gorm.DB) *errcode.Error {
	err := db.Delete(&m).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}
