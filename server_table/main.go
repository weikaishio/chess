package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/weikaishio/chess/common"
	"github.com/weikaishio/chess/server_table/pkg"
	"github.com/weikaishio/chess/util/log"
	"github.com/weikaishio/chess/util/redis_cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s conf_path\n", os.Args[0])
		return
	}

	log.Info("server start, pid = %d", os.Getpid())

	if !pkg.Init(os.Args[1]) {
		return
	}

	listenPort := common.GetListenPort()

	redis_cli.Init(common.GetRedisAddr(), 1000)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(listenPort))
	if err != nil {
		log.Error("listen fail:%s", err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("accept fail:%s", err.Error())
			continue
		}

		go pkg.HandleConn(conn)
	}
}
