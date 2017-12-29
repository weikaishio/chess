package handler

import (
	"fmt"

	"github.com/weikaishio/chess/common"
	"github.com/weikaishio/chess/game/server"
)

func HandleEcho(userid uint32, connid uint32, msgBody []byte) {
	server.SendResp(userid, connid, MsgidEchoResp, common.ResultSuccess, []byte("hello every"))
	fmt.Printf("userid:%d,connid:%d,msgBody:%s\n", userid, connid, string(msgBody))
}
func HandleTest(userid uint32, connid uint32, msgBody []byte) {
	server.SendResp(userid, connid, MsgidTestResp, common.ResultSuccess, []byte("hello HandleTest"))
	fmt.Printf("userid:%d,connid:%d,msgBody:%s\n", userid, connid, string(msgBody))
}
