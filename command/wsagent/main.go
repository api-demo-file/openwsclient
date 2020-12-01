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
	"gopkg.in/ini.v1"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	sigs = make(chan os.Signal)
	done = make(chan bool)

	helperInst *helper.Helper

	callbackUrl = ""
)

func main() {
	if len(os.Args)<2 {
		fmt.Println("usage: wsagent ${config path}")
		os.Exit(1)
	}

	var err error

	configPath := os.Args[1]

	conf,err := ini.Load(configPath)
	if err!=nil {
		panic(err)
	}

	// configs for log
	logLevel,err := conf.Section("log").Key("level").Int()
	if err!=nil {
		logLevel = 0
	}

	client.SetLogLevel(logLevel)

	// configs for connection
	host := conf.Section("app").Key("host").String()
	accessKey := conf.Section("app").Key("access_key").String()
	secret := conf.Section("app").Key("secret").String()

	helperInst,err = helper.NewHelper(host)

	if err!=nil {
		panic(err)
	}

	callbackUrl = conf.Section("agent").Key("callback_url").String()

	helperInst.OnPush(func(data []byte) {
		fmt.Println("OnPush ... ", string(data))
		_,_,err := post("push", string(data))
		if err!=nil {
			client.InfoLog("OnPush failed %s", err.Error())
		}
	})

	helperInst.OnError(func(i int32, s string) {
		fmt.Println("OnError ... ", i, s)
		_,_,err := post("push", fmt.Sprintf("code(%d) msg(%s)", i, s))
		if err!=nil {
			client.InfoLog("OnError failed %s", err.Error())
		}
	})

	subscribesConf := conf.Section("agent").Key("subscribes").String()
	client.InfoLog("subscribesConf %s", subscribesConf)
	if subscribesConf!="" {
		subscribeArr := strings.Split(subscribesConf, ",")
		for _,topic := range subscribeArr {
			client.InfoLog("topic %s", topic)
			helperInst.SetSubscribe(topic)
		}

	}

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