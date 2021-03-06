package server

import (
	"encoding/binary"
	"time"

	"fmt"

	"github.com/weikaishio/chess/codec"
	"github.com/weikaishio/chess/common"
	"github.com/weikaishio/chess/game/session"
	"github.com/weikaishio/chess/util/log"
	"github.com/weikaishio/chess/util/services"
)

type handleFunc func(userid uint32, connid uint32, msgBody []byte)
type verifyFunc func(userid uint32, token string) bool

var requestQ chan codec.GateBackend = make(chan codec.GateBackend, 10000)
var handlers map[uint16]handleFunc = make(map[uint16]handleFunc)
var verifyHandle verifyFunc
var workerNum int
var loginReqMsgid uint16

func pushRequest(gb codec.GateBackend) {
	requestQ <- gb
}

//todo:瓶颈，有个monitorWorker 增加goroutinue
func workLoop() {
	for {
		gb := <-requestQ

		if gb.Msgid == common.MsgRoute {
			var cg codec.ClientGame
			if err := cg.Decode(gb.MsgBuf); err != nil {
				log.Warn("decode client game msg fail:%s", err.Error())
				SendResp(cg.Userid, gb.Connid, gb.Msgid, common.ResultFail, []byte(fmt.Sprintf("decode client game msg fail:%s", err.Error())))
				continue
			}

			f, ok := handlers[cg.Msgid]
			if !ok {
				log.Warn("find %d handler fail, handlers:%v", cg.Msgid, handlers)
				SendResp(cg.Userid, gb.Connid, gb.Msgid, common.ResultFail, []byte(fmt.Sprintf("find %d handler fail", cg.Msgid)))
				continue
			}

			if cg.Msgid != loginReqMsgid && !session.Exist(cg.Userid) {
				log.Info("user %d has not yet logined", cg.Userid)
				//continue
				//todo:auth or continue
				if !verifyHandle(cg.Userid, string(cg.Token)) {
					SendResp(cg.Userid, gb.Connid, gb.Msgid, common.ResultFail, []byte("verify fail"))
					continue
				}
			}
			log.Info("f(cg.Userid, gb.Connid, cg.MsgBody)")
			//todo:verify token
			f(cg.Userid, gb.Connid, cg.MsgBody)
		} else if gb.Msgid == common.MsgGateid {
			if len(gb.MsgBuf) == 4 {
				id := binary.LittleEndian.Uint32(gb.MsgBuf)
				common.SetGateid(id)

				log.Info("recv gateid:%d", id)
			}

		} else if gb.Msgid == common.MsgDisconnect {
			services.DelConnInfo(common.GetGateid(), gb.Connid)
			log.Info("recv disconnect:%d", gb.Connid)
		}
	}
}

func monitorWorker() {
	t := time.Second

	for {
		time.Sleep(t)

		qlen := len(requestQ)

		if qlen > 10 {
			go workLoop()
			workerNum++
			log.Warn("add work routine, workerNum = %d, queueLen = %d", workerNum, qlen)
			t = time.Millisecond * 10
		} else {
			t = time.Second
		}

		if workerNum > 10000 {
			log.Warn("monitorWorker exit")
			return
		}
	}
}

func RegisterHandler(msgid uint16, f handleFunc) {
	handlers[msgid] = f
}

func SetLoginReqMsgid(msgid uint16) {
	loginReqMsgid = msgid
}

func SetVeirfyHandler(verify verifyFunc) {
	verifyHandle = verify
}
