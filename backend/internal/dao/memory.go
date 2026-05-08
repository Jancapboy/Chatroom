package dao

import (
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
)

// MemoryCreate 创建Agent记忆
func (d *Dao) MemoryCreate(memory *model.AgentMemory) *errcode.Error {
	return memory.Create(d.engine)
}

// MemoryListByAgent 获取Agent的所有记忆
func (d *Dao) MemoryListByAgent(agentID string, limit int) ([]model.AgentMemory, *errcode.Error) {
	return model.AgentMemory{}.ListByAgent(d.engine, agentID, limit)
}

// MemoryListByAgentAndTopic 获取Agent关于某话题的记忆
func (d *Dao) MemoryListByAgentAndTopic(agentID, topic string, limit int) ([]model.AgentMemory, *errcode.Error) {
	return model.AgentMemory{}.ListByAgentAndTopic(d.engine, agentID, topic, limit)
}

// MemoryDelete 删除记忆
func (d *Dao) MemoryDelete(id string) *errcode.Error {
	return model.AgentMemory{ID: id}.Delete(d.engine)
}
