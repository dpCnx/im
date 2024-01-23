package handler

import (
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"im/api/pb"
	"im/internal/gateway/conn"
	"im/internal/gateway/data"
	"im/pkg/utils"
)

var ProviderSet = wire.NewSet(NewWsHandler)

type WsHandler struct {
	wsUpGrader  *websocket.Upgrader
	log         *log.Helper
	connManager *conn.Manager

	data *data.Data
}

func NewWsHandler(conf *pb.GateWayConf, connManager *conn.Manager, logger log.Logger, data *data.Data) *WsHandler {
	return &WsHandler{
		wsUpGrader: &websocket.Upgrader{
			ReadBufferSize:   int(conf.Ws.ReadBufferSize),
			WriteBufferSize:  int(conf.Ws.WriteBufferSize),
			HandshakeTimeout: time.Duration(conf.Ws.HandshakeTimeout) * time.Second,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		log:         log.NewHelper(logger),
		connManager: connManager,
		data:        data,
	}
}

func (ws *WsHandler) Ws(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("token")
	if token == "" {
		return
	}
	userId := r.URL.Query()["userId"][0]
	if userId == "" {
		return
	}
	platformId := r.URL.Query()["platformId"][0]

	c, err := ws.wsUpGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn := conn.NewConn(c, utils.StringToInt64(userId), token, platformId)
	ws.connManager.AddConn(conn)

	go ws.readMsg(conn)

}

func (ws *WsHandler) readMsg(c *conn.Conn) {
	for {
		messageType, msg, err := c.GetWsConn().ReadMessage()
		if messageType == websocket.PingMessage {
		}
		if err != nil {
			ws.connManager.DelConn(c)
			return
		}
		if messageType == websocket.CloseMessage {
			ws.connManager.DelConn(c)
			return
		}
		ws.msgParse(c, msg)
	}
}
