/*
 @Title
 @Description
 @Author  Leo
 @Update  2020/12/1 下午6:23
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func post(_type string, data string) (result []byte, code int, err error) {
	if callbackUrl=="" {
		return nil,0, fmt.Errorf("callback url didn't provide")
	}

	fields := url.Values{}

	fields.Add("type", _type)
	fields.Add("data", data)

	resp,err := http.PostForm(callbackUrl, fields)
	if err!=nil {
		return nil,0,err
	}

	defer func() {
		_=resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil,resp.StatusCode, nil
	}

	result,err = ioutil.ReadAll(resp.Body)

	if err!=nil {
		return nil,0,err
	}

	return result,http.StatusOK,nil
}