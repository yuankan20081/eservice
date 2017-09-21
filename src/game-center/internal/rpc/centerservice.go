package rpc

import (
	"game-center/internal/license"
	"game-share/centerservice"
	"github.com/go-xorm/xorm"
	"golang.org/x/net/context"
	"log"
)

type CenterService struct {
	db             *xorm.Engine
	licenseManager *license.Manager
}

func NewCenterService(db *xorm.Engine, m *license.Manager) *CenterService {
	return &CenterService{
		db:             db,
		licenseManager: m,
	}
}

func (cs *CenterService) AgentAuth(ctx context.Context, request *centerservice.AgentAuthRequest) (*centerservice.AgentAuthReply, error) {
	var reply = &centerservice.AgentAuthReply{
		Token:  request.Ticket,
		Server: "",
		Code:   centerservice.AgentAuthReply_FAIL,
	}

	if l, err := cs.licenseManager.Get(request.Ticket); err == nil {
		if l.BindIP == request.Ip {
			reply.Server = l.Server
			reply.Code = centerservice.AgentAuthReply_SUCCESS
		} else {
			log.Printf("do not match license bind ip<%s>, get<%s>\n", l.BindIP, request.Ip)
		}
	} else {
		log.Printf("invalid token<%s>, %s\n", request.Ticket, err)
	}

	return reply, nil
}
