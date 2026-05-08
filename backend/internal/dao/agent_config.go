package dao

import (
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
)

// AgentCreate 创建房间内Agent
func (d *Dao) AgentCreate(agent *model.RoomAgent) *errcode.Error {
	return agent.Create(d.engine)
}

// AgentListByRoom 获取房间内的Agent列表
func (d *Dao) AgentListByRoom(roomID string) ([]model.RoomAgent, *errcode.Error) {
	return model.RoomAgent{}.ListByRoom(d.engine, roomID)
}

// AgentUpdate 更新Agent状态
func (d *Dao) AgentUpdate(agent *model.RoomAgent) *errcode.Error {
	return agent.Update(d.engine)
}

// TemplateCreate 创建Agent模板
func (d *Dao) TemplateCreate(template *model.AgentTemplate) *errcode.Error {
	return template.Create(d.engine)
}

// TemplateList 获取所有Agent模板
func (d *Dao) TemplateList() ([]model.AgentTemplate, *errcode.Error) {
	return model.AgentTemplate{}.List(d.engine)
}

// TemplateGetByID 根据ID获取模板
func (d *Dao) TemplateGetByID(id string) (*model.AgentTemplate, *errcode.Error) {
	return model.AgentTemplate{}.GetByID(d.engine, id)
}
