package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/zeromicro/cds/cmd/galaxy/internal/logic"
	"github.com/zeromicro/cds/cmd/galaxy/internal/svc"
	"github.com/zeromicro/cds/cmd/galaxy/internal/types"
)

func getUserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetUserInfoLogic(r.Context(), ctx)
		resp, err := l.GetUserInfo(req)
		formatFullResponse(resp, err, w, r)
	}
}
