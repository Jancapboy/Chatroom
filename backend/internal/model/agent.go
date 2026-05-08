package model

import (
	"time"

	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
	"gorm.io/gorm"
)

// RoomAgent 房间内智能体实例
type RoomAgent struct {
	ID           string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	RoomID       string    `json:"room_id" gorm:"type:varchar(36);not null"`
	TemplateID   *string   `json:"template_id" gorm:"type:varchar(36)"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Role         string    `json:"role" gorm:"type:varchar(32);not null"` // architect, risk_officer, strategist, analyst, executor
	Avatar       string    `json:"avatar" gorm:"type:varchar(255)"`
	Personality  string    `json:"personality" gorm:"type:text"`
	Expertise    string    `json:"expertise" gorm:"type:json"` // JSON array
	Model        string    `json:"model" gorm:"type:varchar(50);default:'deepseek-chat'"`
	SystemPrompt string    `json:"system_prompt" gorm:"type:text"`
	Energy       int       `json:"energy" gorm:"default:100"`
	Confidence   int       `json:"confidence" gorm:"default:50"`
	Stance       string    `json:"stance" gorm:"type:varchar(16);default:'neutral'"` // support, oppose, neutral
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
}

// AgentTemplate 预设智能体模板
type AgentTemplate struct {
	ID                    string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name                  string    `json:"name" gorm:"type:varchar(100);not null"`
	Role                  string    `json:"role" gorm:"type:varchar(32);not null"`
	Avatar                string    `json:"avatar" gorm:"type:varchar(255)"`
	Personality           string    `json:"personality" gorm:"type:text"`
	Expertise             string    `json:"expertise" gorm:"type:json"`
	DefaultModel          string    `json:"default_model" gorm:"type:varchar(50);default:'deepseek-chat'"`
	SystemPromptTemplate  string    `json:"system_prompt_template" gorm:"type:text"`
	SortOrder             int       `json:"sort_order" gorm:"default:0"`
	CreatedAt             time.Time `json:"created_at"`
}

func (a RoomAgent) Create(db *gorm.DB) *errcode.Error {
	err := db.Create(&a).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (a RoomAgent) ListByRoom(db *gorm.DB, roomID string) ([]RoomAgent, *errcode.Error) {
	var agents []RoomAgent
	err := db.Where("room_id = ? AND is_active = ?", roomID, true).Order("created_at ASC").Find(&agents).Error
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return agents, nil
}

func (a RoomAgent) Update(db *gorm.DB) *errcode.Error {
	err := db.Save(&a).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (t AgentTemplate) Create(db *gorm.DB) *errcode.Error {
	err := db.Create(&t).Error
	if err != nil {
		return errcode.Convert(err)
	}
	return nil
}

func (t AgentTemplate) List(db *gorm.DB) ([]AgentTemplate, *errcode.Error) {
	var templates []AgentTemplate
	err := db.Order("sort_order ASC, created_at ASC").Find(&templates).Error
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return templates, nil
}

func (t AgentTemplate) GetByID(db *gorm.DB, id string) (*AgentTemplate, *errcode.Error) {
	var template AgentTemplate
	err := db.Where("id = ?", id).First(&template).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errcode.NotFound
	}
	if err != nil {
		return nil, errcode.Convert(err)
	}
	return &template, nil
}
