package service

import (
	"context"

	"github.com/Jancapboy/Chatroom/backend/global"
	"github.com/Jancapboy/Chatroom/backend/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	return svc
}

func (svc Service) Dao() *dao.Dao {
	return svc.dao
}

