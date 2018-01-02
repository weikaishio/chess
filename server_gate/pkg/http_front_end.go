package pkg

import (
	"net/http"
	"time"

	"io/ioutil"

	"fmt"

	"github.com/weikaishio/chess/common"
	"github.com/weikaishio/chess/server_gate/connid"
	"github.com/weikaishio/chess/util/log"
)

type httpWriterWithBlock struct {
	http.ResponseWriter
	finishSig chan struct{}
}

func httpRespErr(id uint32, w http.ResponseWriter, err string) {
	defer func() {
		if id > 0 {
			connid.Release(id)
			delConn(id)
		}
	}()
	//defer sendBackendMsg(id, common.MsgDisconnect, nil, nil) here todo?
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err))
}
func DoHttpFrontEnd(w http.ResponseWriter, r *http.Request) {
	log.Info("connection from %s", r.RemoteAddr)

	id := connid.Get()
	if id == connid.InvalidId {
		log.Warn("connid exhaust")
		httpRespErr(0, w, "connid exhaust")
		return
	}
	hw := httpWriterWithBlock{w, make(chan struct{})}
	log.Info("connid:%d", id)
	putConn(id, hw)

	//_, password, ok := r.BasicAuth()
	//if !ok && password == "" {
	//	httpRespErr(0, w, "auth fail")
	//	return
	//}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll(r.Body) err:%v,body:%v\n", err, body)
		httpRespErr(0, w, "params invalid")
		return
	}

	incRecvMsgCounter()
	sendBackendMsg(id, common.MsgRoute, body)
	select {
	case <-time.Tick(15 * time.Second):
		log.Info("timeout http 400")
		httpRespErr(0, w, "handle timeout")
		break
	case <-hw.finishSig:
		break
	}
}
