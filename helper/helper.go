/*
 @Title
 @Description
 @Author  Leo
 @Update  2020/12/1 上午11:40
*/

package helper

import (
	"encoding/json"
	"github.com/api-demo-file/openwsclient/client"
)

type Helper struct {
	wsClient *client.WebSocketClient

	accessKey string
	secret string
	isLogin bool

	subscribeCHs map[string]string

	onPush func([]byte)
	onError func(code int32, msg string)
}

func NewHelper(connStr string) (*Helper ,error){
	helper := new(Helper)
	helper.subscribeCHs = make(map[string]string)
	helper.isLogin = false
	helper.wsClient = client.NewWebSocketClient(connStr, helper)
	helper.onPush = func(bytes []byte) {
		client.InfoLog("this is default CB, please provide it. Push %s", string(bytes))
	}
	helper.onError = func(code int32, msg string) {
		client.InfoLog("this is default CB, please provide it. Error(%d) %s", code, msg)
	}

	return helper,nil
}

func (helper *Helper) Close() {
	helper.wsClient.Close()
}

func (helper *Helper) OnOpen() {

	if helper.accessKey!="" {
		// authentication
		authParameters := client.GenerateAuthParameters(helper.accessKey, helper.secret)
		req := client.ActionReq{
			Action: client.WSActionReq,
			Ch:     "auth",
			Params: authParameters,
		}

		err := helper.wsClient.WriteJSON(req)
		if err!=nil {
			client.ErrLog("authentication failed %s", err.Error())
		}
	}
}

func (helper *Helper) OnMsg(s string, data []byte) {
	client.InfoLog("received new msg %s %s", s, string(data))

	switch s {
	case client.WSActionReq:
		helper.onReqReceived(data)
	case client.WSActionSub:
		helper.onSubReceived(data)
	case client.WSActionPush:
		helper.onPush(data)
	}
}

func (helper *Helper) OnClose() {}

func (helper *Helper) onSubReceived(data []byte) {
	resp:=new(client.ActionResp)
	err := json.Unmarshal(data, resp)
	if err!=nil {
		client.ErrLog("parse json %s failed %s", string(data), err.Error())
		return
	}

	if resp.Code == client.RespCodeOK {
		client.InfoLog("subscribe %s success", resp.Ch)
	}else{
		helper.onError(resp.Code, "subscribe "+resp.Ch+" failed")
	}
}

func (helper *Helper) onReqReceived(data []byte) {
	resp:=new(client.ActionResp)
	err := json.Unmarshal(data, resp)
	if err!=nil {
		client.ErrLog("parse json %s failed %s", string(data), err.Error())
		return
	}

	switch resp.Ch {
	case "auth":
		if resp.Code == client.RespCodeOK {
			helper.onAuthDone()
		}else{
			helper.onError(resp.Code, "authentication failed")
		}
	}

}

func (helper *Helper) onAuthDone() {
	// do subscribe
	for _,topic := range helper.subscribeCHs {
		subObj := &client.ActionSub{
			Action: client.WSActionSub,
			Ch:     topic,
		}
		err := helper.wsClient.WriteJSON(subObj)
		if err!=nil {
			client.ErrLog("subscribe %s failed %s", topic, err.Error())
			return
		}
	}
}

func (helper *Helper) SetSubscribe(topic string) {
	helper.subscribeCHs[topic] = topic
}

func (helper *Helper) OnPush(cb func([]byte)) {
	helper.onPush = cb
}

func (helper *Helper) OnError(cb func(int32,string)) {
	helper.onError = cb
}

//"ws://127.0.0.1:8120/ws/v1"
func (helper *Helper) RunWithAuth(accessKey, secret string) error {
	helper.accessKey = accessKey
	helper.secret = secret
	err := helper.wsClient.Connect()
	if err!=nil {
		return err
	}
	return nil
}