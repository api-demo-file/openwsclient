/*
 @Title  ws 消息定义
 @Description  请填写文件描述（需要改）
 @Author  Leo  2020/4/21 2:55 下午
 @Update  Leo  2020/4/21 2:55 下午
*/

package client

const (

	// 规则：双字节无符号整数，最大 2**16 = 65535, 奇数请求，偶数响应

	// 心跳
	WSActionPing = "ping"
	WSActionPong = "pong"

	WSActionReq  = "req"
	WSActionSub  = "sub"
	WSActionPush = "push"

	RespCodeOK = 200
)

type Callback interface {
	OnMsg(string, []byte)
	OnOpen()
	OnClose()
}

type ProtocolBase struct {
	Action string `json:"action"`
}

// ping/pong
type HeartBeatPing struct {
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

type HeartBeatPong struct {
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

type ActionReq struct {
	Action string            `json:"action"`
	Ch     string            `json:"ch"`
	Params map[string]string `json:"params"`
}

type ActionPush struct {
	Action string                 `json:"action"`
	Ch     string                 `json:"ch"`
	Data   map[string]interface{} `json:"data"`
}

type ActionResp struct {
	Action string `json:"action"`
	Code   int32  `json:"code"`
	Ch     string `json:"ch"`
}

type ActionSub struct {
	Action string `json:"action"`
	Ch     string `json:"ch"`
}

//type UserSubResp struct {
//	Action string `json:"action"`
//	Code   int64  `json:"code"`
//	Ch     string `json:"ch"`
//}
