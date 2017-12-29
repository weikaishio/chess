package handler

import (
	"io"

	"github.com/weikaishio/chess/pb/center"
	"github.com/weikaishio/chess/server_center/conn_info"
	"github.com/golang/protobuf/proto"
)

var addConnInfoResp *center.AddConnInfoResp = &center.AddConnInfoResp{}

func HandleAddConnInfo(client io.Writer, req proto.Message) error {

	addConnInfoReq, ok := req.(*center.AddConnInfoReq)
	if !ok || addConnInfoReq.Info == nil {
		return nil
	}

	oldUserid, isNew := conn_info.Add(*(addConnInfoReq.Info))
	if oldUserid != 0 && oldUserid != addConnInfoReq.Info.Userid {
		sendDelConnInfoNotify(addConnInfoReq.Info, nil)
	}

	if isNew {
		sendNewConnInfoNotify(addConnInfoReq.Info, client)
	}

	return sendResp(client, addConnInfoResp)
}
