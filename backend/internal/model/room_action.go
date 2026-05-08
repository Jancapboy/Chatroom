package model

import (
	"time"

	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"gorm.io/gorm"
)

// RoomAction 推演执行操作记录
type RoomAction struct {
	ID          string                 `json:"id" gorm:"type:varchar(36);primaryKey"`
	RoomID      string                 `json:"room_id" gorm:"type:varchar(36);not null;index"`
	Type        string                 `json:"type" gorm:"type:varchar(32);not null"`      // action类型
	Status      string                 `json:"status" gorm:"type:varchar(32);default:'pending'"` // pending/approved/executed/failed
	Title       string                 `json:"title" gorm:"type:varchar(256)"`
	Description string                 `json:"description" gorm:"type:text"`
	Params      string                 `json:"params" gorm:"type:json"`                      // JSON参数
	Result      string                 `json:"result" gorm:"type:text"`                      // 执行结果
	CreatedBy   uint64                 `json:"created_by"`
	ApprovedBy  uint64                 `json:"approved_by"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	ExecutedAt  *time.Time             `json:"executed_at"`
}

func (a RoomAction) Create(db *gorm.DB) *errcode.Error {
	err := db.Create(&a).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (a RoomAction) Get(db *gorm.DB) (*RoomAction, *errcode.Error) {
	var action RoomAction
	err := db.Where("id = ?", a.ID).First(&action).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errcode.NotFound
	}
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return &action, nil
}

func (a RoomAction) ListByRoom(db *gorm.DB, roomID string) ([]RoomAction, *errcode.Error) {
	var actions []RoomAction
	err := db.Where("room_id = ?", roomID).Order("created_at DESC").Find(&actions).Error
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return actions, nil
}

func (a RoomAction) Update(db *gorm.DB) *errcode.Error {
	err := db.Save(&a).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (a RoomAction) Delete(db *gorm.DB) *errcode.Error {
	err := db.Delete(&a).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}
