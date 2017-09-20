package rpc

import (
	"game-share/centerservice"
	"github.com/go-xorm/xorm"
	"golang.org/x/net/context"
)

type CenterService struct {
	db *xorm.Engine
}

func NewCenterService(db *xorm.Engine) *CenterService {
	return &CenterService{
		db: db,
	}
}

func (cs *CenterService) AgentAuth(ctx context.Context, request *centerservice.AgentAuthRequest) (*centerservice.AgentAuthReply, error) {
	var reply = &centerservice.AgentAuthReply{
		Token:  "",
		Server: "",
		Code:   centerservice.AgentAuthReply_FAIL,
	}



	return reply, nil
}
