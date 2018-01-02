package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/qiniu/log"
	"github.com/weikaishio/chess/codec"
	"github.com/weikaishio/chess/game_slots/pb_game"
	"github.com/weikaishio/chess/pb/login"
)

const accountLoginUrl = "http://127.0.0.1:9090/login"
const gateUrl = "http://127.0.0.1:8886/"

func loginAccount() (loginResp login.LoginResp, success bool) {
	var req login.LoginReq
	req.Version = 1

	data, _ := proto.Marshal(&req)
	resp, err := http.Post(accountLoginUrl, "", bytes.NewReader(codec.EncryptWithLen(data)))
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("%s", err.Error())
		return
	}

	if err := proto.Unmarshal(codec.DecryptWithLen(buf), &loginResp); err != nil {
		log.Error("%s", err.Error())
		return
	}

	log.Info("account login resp:%s", loginResp.String())

	success = true
	return
}

func httpSpin(userid uint32, username, token string) bool {
	var req pb_game.UserThemeSpinC2S
	var resp pb_game.UserThemeSpinSuccessS2C

	req.Msgid = proto.Uint32(1)
	req.Gold = proto.Int64(2)
	req.Bet = proto.Int32(3)
	req.Isfree = proto.Uint32(4)
	req.Themeid = proto.Int32(5)

	var cg codec.ClientGame
	cg.Userid = userid
	cg.MsgBody, _ = proto.Marshal(&req)
	cg.Token = []byte(token)
	cg.Msgid = uint16(pb_game.MsgHeader_user_theme_spin_c2s)
	body, err := cg.Encode2Byte()
	if err != nil {
		log.Error("test fail:%v", err)
		return false
	}
	reqHttp, err := http.NewRequest("POST", gateUrl, bytes.NewReader(body))
	if err != nil {
		log.Error("httpSpin NewRequest err:%v", err)
		return false
	}
	client := &http.Client{}
	respBody, err := client.Do(reqHttp)
	if err != nil {
		log.Error("http.post er:%v", err)
		return false
	}
	if respBody.StatusCode != http.StatusOK {
		log.Warn("respBody.StatusCode:%d!=http.StatusOK", respBody.StatusCode)
		return false
	}
	var gc codec.GameClient
	if err := gc.DecodeFromReader(respBody.Body); err != nil {
		log.Error("receive test resp fail:%s", err.Error())
		return false
	}
	err = proto.Unmarshal(gc.MsgBody, &resp)
	fmt.Printf("receive httpSpin resp: %v,err:%v\n", resp, err)
	return true
}
func main() {
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f}
	codec.Init(key, iv) //todo:全局的加解密密钥，真商用，需要改成登录后获取的密钥

	loginResp, success := loginAccount()
	if !success {
		return
	}

	res := true
	for res {
		res = httpSpin(loginResp.Userid, loginResp.Username, loginResp.Token)
		time.Sleep(1 * time.Second)
	}
}
