package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/zeromicro/cds/cmd/dm/cmd/sync/config"
	"github.com/zeromicro/cds/cmd/galaxy/internal/clients"
	"github.com/zeromicro/cds/cmd/galaxy/internal/svc"
	"github.com/zeromicro/cds/cmd/galaxy/internal/types"
	"github.com/zeromicro/cds/pkg/strx"
	"github.com/zeromicro/cds/pkg/timex"
	"github.com/zeromicro/go-zero/core/logx"
)

type DmListLogic struct {
	ctx context.Context
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDmListLogic(ctx context.Context, svcCtx *svc.ServiceContext) DmListLogic {
	return DmListLogic{
		ctx:    ctx,
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
	// TODO need set model here from svc
}

func (l *DmListLogic) DmList(req types.ListRequest) (*types.DmListResponse, error) {
	cli := clients.NewDmClient(l.svcCtx.EtcdClient)
	sts, e := cli.All()
	if e != nil {
		logx.Error(e)
		return nil, e
	}
	stMap := make(map[int]*config.Status)
	for _, st := range sts {
		id, e := strconv.Atoi(st.ID)
		if e != nil {
			logx.Error(e)
			continue
		}
		stMap[id] = st
	}

	vs, e := l.svcCtx.DmModel.FindByDb(req.DbName, req.Page, req.Size)
	if e != nil {
		logx.Error(e)
		return nil, e
	}
	rp := &types.DmListResponse{
		DmList: make([]types.Dm, 0, len(vs)),
	}
	cnt, err := l.svcCtx.DmModel.GetCountByDb(req.DbName)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	for _, v := range vs {
		dm := types.Dm{
			ID:              v.ID,
			Name:            v.Name,
			SourceType:      v.SourceType,
			SourceDsn:       v.SourceDsn,
			SourceTable:     v.SourceTable,
			SourceQueryKey:  v.SourceQueryKey,
			TargetType:      v.TargetType,
			TargetDB:        v.TargetDB,
			TargetChProxy:   v.TargetChProxy,
			TargetTable:     v.TargetTable,
			WindowStartHour: v.WindowStartHour,
			WindowEndHour:   v.WindowEndHour,
			CreateTime:      v.CreateTime.Format(timex.StandardLayout),
		}
		if v.TargetShards != "" {
			var vs []string
			shards, e := strx.DecryptDsn(v.TargetShards)
			if e != nil {
				logx.Error(e)
				return nil, e
			}
			e = json.Unmarshal([]byte(shards), &vs)
			if e != nil {
				logx.Error(e)
				logx.Error(v.TargetShards)
				return nil, e
			}
			dm.TargetShards = strx.DeepSplit(vs, ",")
		}

		if job, ok := stMap[dm.ID]; ok {
			dm.Status = job.Status
			dm.Information = job.Information
			dm.UpdateTime = job.UpdateTime.Format(timex.StandardLayout)
			if time.Now().Unix()-job.UpdateTime.Unix() > 60 && dm.Status == "running" {
				dm.Status = "stopped"
				dm.Information = "任务超时"
				dm.UpdateTime = job.UpdateTime.Format(timex.StandardLayout)
			}

		} else {
			dm.Status = "未启动"
		}
		rp.DmList = append(rp.DmList, dm)
	}
	rp.PageAndSize = types.PageAndSize{
		Size: int(cnt),
		Page: req.Page,
	}
	return rp, nil
}
