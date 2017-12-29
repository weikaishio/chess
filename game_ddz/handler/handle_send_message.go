package handler

import (
	"github.com/weikaishio/chess/common"
	"github.com/weikaishio/chess/game/server"
	"github.com/weikaishio/chess/game/session"
	"github.com/weikaishio/chess/game_ddz/pb_client"
	"github.com/weikaishio/chess/game_ddz/user"
	"github.com/weikaishio/chess/util/log"
	"github.com/golang/protobuf/proto"
)

func HandleSendMessage(userid uint32, connid uint32, msgBody []byte) {
	var req pb_client.SendMessageReq
	var resp pb_client.SendMessageResp

	if err := proto.Unmarshal(msgBody, &req); err != nil {
		log.Warn("unmarshal SendMessageReq fail:%s", err.Error())
		return
	}

	log.Info("receive SendMessageReq:%s", req.String())

	var result uint16
	result = common.ResultFail

	defer func() {
		exitFunc(userid, connid, MsgidSendMessageResp, result, &resp)
	}()

	ui := user.LoadUserInfo(req.Receiver, []int{user.FlagBasicInfo})
	if ui == nil {
		result = ResultFailUserNotExist
		return
	}

	if !session.Exist(req.Receiver) {
		result = ResultFailNotLogined
		return
	}

	var notify pb_client.MessageNotify
	notify.Sender = userid
	notify.Content = req.Content

	buf, _ := proto.Marshal(&notify)
	server.SendResp(req.Receiver, connid, MsgidMessageNotify, 0, buf)

	result = common.ResultSuccess
}
