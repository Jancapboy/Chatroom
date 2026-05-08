package model

import (
	"time"

	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"gorm.io/gorm"
)

// Room 推演房间模型
type Room struct {
	ID            string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null"`
	Topic         string    `json:"topic" gorm:"type:text"`
	Description   string    `json:"description" gorm:"type:text"`
	Status        string    `json:"status" gorm:"type:varchar(32);default:'preparing'"` // preparing, running, paused, completed, archived
	TemplateID    *string   `json:"template_id" gorm:"type:varchar(36)"`
	CurrentPhase  string    `json:"current_phase" gorm:"type:varchar(32);default:'info_gathering'"`
	CurrentRound  int       `json:"current_round" gorm:"default:1"`
	MaxRounds     int       `json:"max_rounds" gorm:"default:10"`
	ConsensusScore int      `json:"consensus_score" gorm:"default:0"`
	ForkedFrom    *string   `json:"forked_from" gorm:"type:varchar(36)"`
	CreatedBy     uint64    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	Agents   []RoomAgent `json:"agents,omitempty" gorm:"foreignKey:RoomID"`
	Messages []Message   `json:"messages,omitempty" gorm:"foreignKey:RoomID"`
}

func (r *Room) Create(db *gorm.DB) *errcode.Error {
	now := time.Now()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = now
	}
	err := db.Create(r).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (r *Room) Get(db *gorm.DB) (*Room, *errcode.Error) {
	var room Room
	err := db.Where("id = ?", r.ID).Preload("Agents").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(100)
	}).First(&room).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errcode.NotFound
	}
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return &room, nil
}

func (r *Room) Update(db *gorm.DB) *errcode.Error {
	r.UpdatedAt = time.Now()
	err := db.Save(r).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (r Room) Delete(db *gorm.DB) *errcode.Error {
	err := db.Delete(&r).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (r Room) List(db *gorm.DB, status string, page, pageSize int) ([]Room, int64, *errcode.Error) {
	var rooms []Room
	var total int64

	query := db.Model(&Room{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errcode.Convert(err)
	}

	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rooms).Error
	if err != nil {
		return nil, 0, errcode.Convert(err)
	}
	return rooms, total, nil
}
