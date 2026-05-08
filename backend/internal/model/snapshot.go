package model

import (
	"time"

	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"gorm.io/gorm"
)

// RoomSnapshot 推演状态快照
type RoomSnapshot struct {
	ID             string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	RoomID         string    `json:"room_id" gorm:"type:varchar(36);not null;index"`
	Round          int       `json:"round" gorm:"not null"`
	Phase          string    `json:"phase" gorm:"type:varchar(32)"`
	ConsensusScore int       `json:"consensus_score"`
	AgentStates    string    `json:"agent_states" gorm:"type:json"`    // JSON: [{agent_id, stance, confidence, energy}]
	KeyDecisions   string    `json:"key_decisions" gorm:"type:json"`   // JSON: [{content, stance, agent_id}]
	TriggerReason  string    `json:"trigger_reason" gorm:"type:text"`  // 为什么创建这个快照
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (s RoomSnapshot) Create(db *gorm.DB) *errcode.Error {
	err := db.Create(&s).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (s RoomSnapshot) Get(db *gorm.DB) (*RoomSnapshot, *errcode.Error) {
	var snapshot RoomSnapshot
	err := db.Where("id = ?", s.ID).First(&snapshot).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errcode.NotFound
	}
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return &snapshot, nil
}

func (s RoomSnapshot) ListByRoom(db *gorm.DB, roomID string) ([]RoomSnapshot, *errcode.Error) {
	var snapshots []RoomSnapshot
	err := db.Where("room_id = ?", roomID).Order("round ASC, created_at ASC").Find(&snapshots).Error
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return snapshots, nil
}

func (s RoomSnapshot) Delete(db *gorm.DB) *errcode.Error {
	err := db.Delete(&s).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}
