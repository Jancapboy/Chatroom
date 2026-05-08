package dao

import (
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
)

// MessageCreate 创建消息
func (d *Dao) MessageCreate(msg *model.Message) *errcode.Error {
	return msg.Create(d.engine)
}

// MessageListByRoom 获取房间消息列表
func (d *Dao) MessageListByRoom(roomID string, round int, phase string, page, pageSize int) ([]model.Message, int64, *errcode.Error) {
	return model.Message{}.ListByRoom(d.engine, roomID, round, phase, page, pageSize)
}

// GetPhaseMessages 获取某Phase的消息
func (d *Dao) GetPhaseMessages(roomID, phase string, round int) ([]model.Message, *errcode.Error) {
	return model.Message{}.GetPhaseMessages(d.engine, roomID, phase, round)
}
