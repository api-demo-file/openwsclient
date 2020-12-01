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
	"github.com/api-demo-file/openwsclient/helper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sigs = make(chan os.Signal)
	done = make(chan bool)

	helperInst *helper.Helper
)

func main() {
	if len(os.Args)<4 {
		fmt.Println("usage: helpertest ${host} ${accessKey} ${secret}")
		os.Exit(1)
	}

	host := os.Args[1]
	accessKey := os.Args[2]
	secret := os.Args[3]

	client.SetLogLevel(client.LogLevelInfo)

	var err error
	helperInst,err = helper.NewHelper(host)

	if err!=nil {
		panic(err)
	}

	helperInst.OnPush(func(data []byte) {
		fmt.Println("OnPush ... ", string(data))
	})

	helperInst.OnError(func(i int32, s string) {
		fmt.Println("OnError ... ", i, s)
	})

	helperInst.SetSubscribe("orders#ethusdt")
	err = helperInst.RunWithAuth(accessKey, secret)

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
	client.InfoLog("recv signal %s", sig.String())

	helperInst.Close()

	time.Sleep(time.Second)

	done<-true
}