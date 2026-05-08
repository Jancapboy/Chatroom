package dao

import (
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
)

// ActionCreate 创建Action记录
func (d *Dao) ActionCreate(action *model.RoomAction) *errcode.Error {
	return action.Create(d.engine)
}

// ActionGet 获取Action
func (d *Dao) ActionGet(id string) (*model.RoomAction, *errcode.Error) {
	return model.RoomAction{ID: id}.Get(d.engine)
}

// ActionListByRoom 获取房间的所有Action
func (d *Dao) ActionListByRoom(roomID string) ([]model.RoomAction, *errcode.Error) {
	return model.RoomAction{}.ListByRoom(d.engine, roomID)
}

// ActionUpdate 更新Action
func (d *Dao) ActionUpdate(action *model.RoomAction) *errcode.Error {
	return action.Update(d.engine)
}

// ActionDelete 删除Action
func (d *Dao) ActionDelete(id string) *errcode.Error {
	return model.RoomAction{ID: id}.Delete(d.engine)
}
