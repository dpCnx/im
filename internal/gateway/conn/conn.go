package conn

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"im/api/pb"
	"im/pkg/utils"
)

type Conn struct {
	ws         *websocket.Conn // websocket连接
	wsMutex    sync.Mutex      // WS写锁
	userId     int64           // 用户ID
	token      string
	platformId int
}

func NewConn(ws *websocket.Conn, userId int64, token string, platformId string) *Conn {
	return &Conn{
		ws:         ws,
		wsMutex:    sync.Mutex{},
		userId:     userId,
		token:      token,
		platformId: utils.StringToInt(platformId),
	}
}
func (c *Conn) GetWsConn() *websocket.Conn {
	return c.ws
}

func (c *Conn) Send(pt pb.PackageType, code int64, msg string, message proto.Message) error {
	var output = pb.Output{
		Type:    pt,
		Code:    code,
		Message: msg,
	}
	if message != nil {
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			return err
		}
		output.Data = msgBytes
	}

	outputBytes, err := proto.Marshal(&output)
	if err != nil {
		return err
	}

	if err = c.write(outputBytes); err != nil {
		return err
	}

	return nil
}

func (c *Conn) write(bytes []byte) error {
	c.wsMutex.Lock()
	defer c.wsMutex.Unlock()
	if err := c.ws.SetWriteDeadline(time.Now().Add(10 * time.Millisecond)); err != nil {
		return err
	}
	return c.ws.WriteMessage(websocket.BinaryMessage, bytes)
}
