package handler

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"

	"github.com/weikaishio/chess/common"
	"github.com/weikaishio/chess/game_slots/pb_game"
	"github.com/weikaishio/chess/util/log"
	"github.com/weikaishio/chess/util/redis_cli"
)

func HandleVerifyToken(userid uint32, token string) bool {
	log.Info("HandleVerifyToken userId:%d,token:%s", userid, token)
	key := common.GenLoginInfoKey(userid)
	value, err := redis_cli.Get(key) //todo:使用user client或者rpc到server_login去校验
	if err != nil {
		log.Warn("VerifyToken redis_cli.Get(key) fail:%s", err.Error())
		return false
	}

	var loginInfo common.LoginInfo
	if err := json.Unmarshal([]byte(value), &loginInfo); err != nil {
		log.Warn("VerifyToken json.Unmarshal fail:%s", err.Error())
		return false
	}

	if token != loginInfo.Token {
		log.Info("token:%s != loginInfo.Token:%s", token, loginInfo.Token)
		return false
	}
	return true
}

func HandleSpin(userid uint32, connid uint32, msgBody []byte) {
	var req pb_game.UserThemeSpinC2S
	var resp pb_game.UserThemeSpinSuccessS2C

	if err := proto.Unmarshal(msgBody, &req); err != nil {
		log.Warn("unmarshal SendMessageReq fail:%s", err.Error())
		return
	}

	log.Info("receive UserThemeSpinC2S:%s", req.String())

	var result uint16
	result = common.ResultFail

	defer func() {
		exitFunc(userid, connid, 0, result, &resp)
	}()

	resp.Msgid = proto.Uint32(1011)
	symbolAry := make([]*pb_game.Symbol, 0)
	symbolAry = append(symbolAry, &pb_game.Symbol{
		X:        proto.Int32(1),
		Y:        proto.Int32(2),
		Symbolid: proto.Int32(3),
	})
	resp.Symbols = symbolAry

	eventAry := make([]*pb_game.Eventinfo, 0)
	coordinateAry := make([]*pb_game.Coordinate, 0)
	coordinateAry = append(coordinateAry, &pb_game.Coordinate{
		X: proto.Int32(1),
		Y: proto.Int32(2),
	})
	eventAry = append(eventAry, &pb_game.Eventinfo{
		Id:          proto.Int32(1),
		Num:         proto.Int32(2),
		Coordinates: coordinateAry,
	})
	resp.Eventinfo = eventAry
	lineAry := make([]*pb_game.Line, 0)
	lineAry = append(lineAry, &pb_game.Line{
		Lineid:      proto.Uint32(1),
		Coordinates: coordinateAry,
	})
	resp.Lines = lineAry
	resp.Replace = &pb_game.Replace{
		Replace1: symbolAry,
		Replace2: symbolAry,
		Replace3: symbolAry,
	}
	resp.Rewardinfo = &pb_game.Rewardinfo{
		Gold:   proto.Int64(1),
		Exp:    proto.Int64(2),
		Spin:   proto.Int32(3),
		Gameid: proto.Int32(4),
	}
	result = common.ResultSuccess
}
