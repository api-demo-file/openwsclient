/*
 @Title
 @Description
 @Author  Leo
 @Update  2020/11/30 下午3:31
*/

package main

import (
	"fmt"
	"github.com/api-demo-file/openwsclient/client"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sigs = make(chan os.Signal)
	done = make(chan bool)

	clientObj *client.WebSocketClient
)

type CallbackImpl struct {}

func (cb *CallbackImpl) OnOpen() {
	fmt.Println("ws opened")
}

func (cb *CallbackImpl) OnClose() {
	fmt.Println("ws closed")
}

func (cb *CallbackImpl) OnMsg(t string, data []byte) {
	fmt.Println("received msg", t, string(data))
}

func main() {
	if len(os.Args)<2 {
		fmt.Println("usage: clienttest ${host}")
		os.Exit(1)
	}

	host := os.Args[1]

	cb := new(CallbackImpl)
	clientObj = client.NewWebSocketClient(host, cb)

	err := clientObj.Connect()
	if err!=nil {
		panic(err)
	}

	// serve
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // ctrl+c, kill, kill -2,
	go sigAwaiter()

	<-done
}

func sigAwaiter() {
	sig := <-sigs
	//fmt.Println(fmt.Sprintf("recv signal %s", sig.String()))
	client.InfoLog("recv signal %s", sig.String())

	clientObj.Close()

	time.Sleep(time.Second)

	done<-true
}