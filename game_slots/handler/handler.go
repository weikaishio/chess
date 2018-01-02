package handler

import (
	"github.com/golang/protobuf/proto"
	"github.com/weikaishio/chess/common"
	"github.com/weikaishio/chess/game/server"
	"github.com/weikaishio/chess/game_slots/pb_game"
)

func init() {
	server.RegisterHandler(uint16(pb_game.MsgHeader_user_theme_spin_c2s), HandleSpin)
	server.SetVeirfyHandler(HandleVerifyToken)

}
func exitFunc(userid, connid uint32, msgid uint16, result uint16, resp proto.Message) {
	if result != common.ResultSuccess || resp == nil {
		server.SendResp(userid, 0, msgid, result, nil)
		return
	}

	buf, _ := proto.Marshal(resp)
	server.SendResp(userid, connid, msgid, result, buf)
}
