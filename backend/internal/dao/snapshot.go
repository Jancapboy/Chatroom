package dao

import (
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/pkg/errcode"
)

// SnapshotCreate 创建快照
func (d *Dao) SnapshotCreate(snapshot *model.RoomSnapshot) *errcode.Error {
	return snapshot.Create(d.engine)
}

// SnapshotGet 获取快照
func (d *Dao) SnapshotGet(id string) (*model.RoomSnapshot, *errcode.Error) {
	return model.RoomSnapshot{ID: id}.Get(d.engine)
}

// SnapshotListByRoom 获取房间的所有快照
func (d *Dao) SnapshotListByRoom(roomID string) ([]model.RoomSnapshot, *errcode.Error) {
	return model.RoomSnapshot{}.ListByRoom(d.engine, roomID)
}

// SnapshotDelete 删除快照
func (d *Dao) SnapshotDelete(id string) *errcode.Error {
	return model.RoomSnapshot{ID: id}.Delete(d.engine)
}
