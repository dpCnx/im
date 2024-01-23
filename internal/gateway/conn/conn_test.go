package conn

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnClient(t *testing.T) {
	h := http.Header{}
	h.Add("token", "1111")

	conn, _, err := websocket.DefaultDialer.Dial("ws://192.168.0.99:9990/ws?userId=111&platformId=1", h)
	if err != nil {
		t.Log(err)
	}

	fmt.Println(conn)
}
