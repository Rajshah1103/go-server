package handler

import (
	"github.com/Rajshah1103/go-server/utils"
	"net"
)

func Health(conn net.Conn, method, path string, headers map[string]string) {
	conn.Write([]byte(utils.BuildHTTPResponse("OK ✅")))
}
