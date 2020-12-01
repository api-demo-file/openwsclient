# openwsclient

长连接客户端      

包名 package name   
   
## 使用方法    

client 包：基础包，提供最基础连接，断线重连机制   
使用方法：参考 test/clienttest  

helper 包：扩展包，提供授权，订阅，自动管理断线重连后的授权订阅操作   
使用方法：参考 test/helper  

command/wsagent
代理程序，通过简单地配置，即可跨语言使用sdk   
wsagent 会将收到的 push 消息通过 http POST 转发到指定 url   

## 编译

本程序使用 go module 管理第三方包   

go mod download  
cd command/wsagent  
go build .  

## 运行

首先配置 config.ini

./wsagent ./config.ini

-----   

the client of bitz open WS  

package name  
github.com/api-demo-file/openwsclient   

## usage

package client: the basic program, provide function that can connect to WS server

package helper: the extra program, provide extra function includes authentication/subscribe and auto authentication/subscribe when reconnect.  

command/wsagent

The delegate program. It use sample config to cross program language to use the SDK.

THe program 'wsagent' will transfer the PUSH and ERROR msg to an URL using POST which was provided from config.ini.    

## compile  

It use 'go module' to manage packages.

go mod download   
cd command/wsagent    
go build .   

## run   

You need to configure config.ini first.  

./wsagent ./config.ini
