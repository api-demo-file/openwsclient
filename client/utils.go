/*
 @Title
 @Description
 @Author  Leo
 @Update  2020/11/30 下午5:41
*/

package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"
)

var (
	logLevel = 0

	LogLevelDebug = 0
	LogLevelInfo = 1
	LogLevelError = 2
	LogLevelNone = 100
)

func SetLogLevel(l int) {
	logLevel = l
}

func ErrLog(_fmt string, args... interface{}) {
	if logLevel>LogLevelError {
		return
	}
	writeLog("ERROR", _fmt, args...)
}

func InfoLog(_fmt string, args... interface{}) {
	if logLevel>LogLevelInfo {
		return
	}

	writeLog("INFO", _fmt, args...)
}

func DebugLog(_fmt string, args... interface{}) {
	if logLevel>LogLevelDebug {
		return
	}

	writeLog("DEBUG", _fmt, args...)
}

func writeLog(level string, _fmt string, args... interface{}) {
	n := time.Now().Format("2006-01-02_15:04:05")
	msg := fmt.Sprintf(_fmt, args...)
	caller := GetCaller(2)
	fmt.Println(fmt.Sprintf("%s\t%s\t%s\t%s", n, level, caller, msg))
}

func GetCaller(skip int) string {
	_,file,line,_ := runtime.Caller(skip+1)
	file = file[strings.LastIndex(file, "/")+1:]
	return fmt.Sprintf("%s:%d", file, line)
}

func GenerateAuthParameters(accessKey,secret string) map[string]string {
	parameters:=make(map[string]string)
	parameters["accessKey"] = accessKey
	parameters["signatureMethod"] = "HmacSHA256"
	parameters["signatureVersion"] = "1.0"
	parameters["timestamp"] = time.Now().In(time.UTC).Format("2006-01-02T15:04:05")
	var preStrBuf bytes.Buffer
	paramKeys := make([]string, 0)
	for k,_ := range parameters {
		paramKeys = append(paramKeys, k)
	}
	sort.Strings(paramKeys)
	for _,k := range paramKeys {
		v := parameters[k]
		if preStrBuf.Len()>0 {
			preStrBuf.WriteRune('&')
		}
		preStrBuf.WriteString(k)
		preStrBuf.WriteRune('=')
		preStrBuf.WriteString(v)
	}
	h:=hmac.New(sha256.New, []byte(secret))
	h.Write(preStrBuf.Bytes())
	parameters["signature"] = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return parameters
}