package logic

import (
	"context"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/zeromicro/cds/cmd/galaxy/internal/svc"
	"github.com/zeromicro/cds/cmd/galaxy/internal/types"
)

type AddLogic struct {
	ctx context.Context
	logx.Logger
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddLogic {
	return AddLogic{
		ctx:    ctx,
		Logger: logx.WithContext(ctx),
	}
	// TODO need set model here from svc
}

func (l *AddLogic) Add(req types.UserAddRequest) error {
	return nil
}
