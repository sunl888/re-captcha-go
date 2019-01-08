# Google reCaptcha package for golang

## Intro
`reCaptcha v3` 使网站所有者更简单有效的清除机器人所产生的有害流量，甚至不需要访问者证明他们是人类。这一变化使网站访问者更容易登录他们喜爱的网站，而不会浪费时间通过解决一个谜题或者识别一张图片来证明他们是真正的人。

## 准备
- 在 [Google reCaptcha](https://www.google.com/recaptcha/admin "google reCaptcha admin") 网站获取 `siteKey` 和 `screctKey`

## 安装
- go get https://github.com/wq1019/re-captcha-go

## Example
``` go
package main

import (
	"fmt"
	"github.com/wq1019/re-captcha-go"
	"log"
	"net/http"
)

var (
	sitekey   = "6LfR8YcUAAAAAK-QLovv4b8X40J-***********"
	secretKey = "6LfR8YcUAAAAAOpQkKTJhFdq6eE************"
)

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	reCaptcha := re_captcha_go.NewReCaptcha(secretKey)
	isOk, err := reCaptcha.Verify(r)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if isOk == false {
		_, _ = w.Write([]byte("no"))
	}
	_, _ = w.Write([]byte("ok"))
}
func main() {
	http.HandleFunc("/verify", verifyHandler)
	
	log.Println("example link: http://captcha.local:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```