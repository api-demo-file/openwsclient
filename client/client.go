/*
 @Title
 @Description
 @Author  Leo
 @Update  2020/11/30 下午12:30
*/

package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
)

const (
	HeartbeatInterval = 20 // second
	ReconnectAfter    = 30 // second
)

type WebSocketClient struct {
	// client info
	connStr           string
	heartbeatInterval int64 // heartbeatInterval : Send Heartbeat command at every N seconds .
	reconnectAfter    int64 // reconnectDelay :  the connection will reconnect if don't receive PONG for N seconds.
	callbackObj       Callback

	// connect resources
	wsConn   *websocket.Conn
	httpResp *http.Response

	// the time will be updated when server response a PONG.
	PongAt           int64
	heartbeatCloseCh chan int32
}

// Create a WebSocketClient
// connStr : The url for program connect to. ex: wss://wsapi.bitz.ai/ws/v1
// msgCallback : WebSocketClient will invoke the function when a new msg received.
func NewWebSocketClient(connStr string, cb Callback) *WebSocketClient {
	ws := &WebSocketClient{
		connStr:           connStr,
		heartbeatInterval: HeartbeatInterval,
		reconnectAfter:    ReconnectAfter,
		callbackObj:       cb,
		heartbeatCloseCh:  make(chan int32, 10),
	}

	return ws
}

func (wsClient *WebSocketClient) Connect() error {
	err := wsClient.initSocket()
	if err != nil {
		return err
	}

	go wsClient.heartbeat()

	return nil
}

func (wsClient *WebSocketClient) Close() {
	err := wsClient.wsConn.Close()
	if err != nil {
		ErrLog("close connection failed %s", err.Error())
	}

	close(wsClient.heartbeatCloseCh) // close heartbeat task
}

func (wsClient *WebSocketClient) WriteJSON(v interface{}) error {
	return wsClient.wsConn.WriteJSON(v)
}

func (wsClient *WebSocketClient) initSocket() error {
	u := &url.URL{}

	u, err := u.Parse(wsClient.connStr)

	if err != nil {
		return err
	}

	DebugLog("url obj %+v", u)
	DebugLog(u.Scheme)
	DebugLog(u.Host)
	DebugLog(u.Path)

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	wsClient.PongAt = time.Now().Unix()
	wsClient.wsConn = conn
	wsClient.httpResp = resp

	go wsClient.readLoop()

	wsClient.callbackObj.OnOpen()
	InfoLog("websocket connected")
	return nil
}

func (wsClient *WebSocketClient) reconnect() {
	InfoLog("try to reconnect websocket")
	if wsClient.wsConn != nil {
		_=wsClient.wsConn.Close()
	}
	err := wsClient.initSocket()
	if err != nil {
		ErrLog("init socket failed %s", err.Error())
	}
}

func (wsClient *WebSocketClient) readLoop() {
	InfoLog("readLoop started")
	defer func() {
		wsClient.callbackObj.OnClose()
		InfoLog("readLoop exited")
	}()

	var err error

	for {
		var msg []byte
		_, msg, err = wsClient.wsConn.ReadMessage()
		if err != nil {
			ErrLog("ReadLoop ERROR: %s", err.Error())
			return
		}

		protocolType := new(ProtocolBase)
		err = json.Unmarshal(msg, protocolType)
		if err != nil {
			ErrLog("parse json [%s] failed: %s", string(msg), err.Error())
			continue
		}

		// handle with heartbeat
		if protocolType.Action == WSActionPong {
			wsClient.updateHeartbeat()
			continue
		}

		wsClient.callbackObj.OnMsg(protocolType.Action, msg)
	}
}

func (wsClient *WebSocketClient) updateHeartbeat() {
	wsClient.PongAt = time.Now().Unix()
}

func (wsClient *WebSocketClient) heartbeat() {
	InfoLog("heartbeat started")
	defer func() {
		InfoLog("heartbeat exited")
	}()

	duration := time.Second * time.Duration(wsClient.heartbeatInterval)

	for {
		select {
		case <-time.After(duration):
			now := time.Now().Unix()

			timeDiff := now-wsClient.PongAt
			DebugLog("now %d PongAt %d timeDiff %d", now, wsClient.PongAt, timeDiff)
			if timeDiff>=wsClient.reconnectAfter {
				wsClient.reconnect()
			}else{
				wsClient.sendHeartBeat()
			}

		case <-wsClient.heartbeatCloseCh:
			InfoLog("heartbeat will stop ")
			return
		}
	}

}

func (wsClient *WebSocketClient) sendHeartBeat() {
	heartbeatMsg := &HeartBeatPing{
		Action:    WSActionPing,
		Timestamp: time.Now().Unix(),
	}

	err := wsClient.wsConn.WriteJSON(heartbeatMsg)
	if err != nil {
		ErrLog("send Heartbeat failed: %s", err.Error())
	}else{
		InfoLog("send heartbeat done")
	}
}
