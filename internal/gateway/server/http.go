package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	"im/api/pb"
	"im/internal/gateway/handler"
)

func NewHTTPServer(c *pb.GateWayConf, h *handler.WsHandler) *http.Server {

	router := mux.NewRouter()
	router.HandleFunc("/ws", h.Ws)

	httpSrv := http.NewServer(http.Address(c.Http.Addr))
	httpSrv.HandlePrefix("/", router)

	return httpSrv
}
