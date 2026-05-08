package dao

import (
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
)

// RoomCreate 创建房间
func (d *Dao) RoomCreate(room *model.Room) *errcode.Error {
	return room.Create(d.engine)
}

// RoomGet 获取房间详情
func (d *Dao) RoomGet(id string) (*model.Room, *errcode.Error) {
	return (&model.Room{ID: id}).Get(d.engine)
}

// RoomUpdate 更新房间
func (d *Dao) RoomUpdate(room *model.Room) *errcode.Error {
	return room.Update(d.engine)
}

// RoomDelete 删除房间
func (d *Dao) RoomDelete(id string) *errcode.Error {
	return model.Room{ID: id}.Delete(d.engine)
}

// RoomList 房间列表
func (d *Dao) RoomList(status string, page, pageSize int) ([]model.Room, int64, *errcode.Error) {
	return model.Room{}.List(d.engine, status, page, pageSize)
}
