# Google reCaptcha for golang

## 准备
- 在 [Google reCaptcha](https://www.google.com/recaptcha/admin "google reCaptcha admin") 网站获取 `siteKey` 和 `screctKey`

## 安装
- go get github.com/wq1019/reCaptcha-go

# Example
``` go
package main

import (
	"fmt"
	"github.com/wq1019/captcha"
	"log"
	"net/http"
)

var (
	sitekey   = "6LfR8YcUAAAAAK-QLovv4b8X40J-1oee_AtoNNSS"
	secretKey = "6LfR8YcUAAAAAOpQkKTJhFdq6eEZJhkFN0hxNoGQ"
)

func showHandler(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
	<!doctype html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Document</title>
		<script src='https://www.google.com/recaptcha/api.js?render=%s'></script>
		<script>
        	grecaptcha.ready(function () {
            	grecaptcha.execute('%s', {action: 'home'}).then(function (token) {
                    // Verify the token on the server.
                   	document.getElementById("response").value = token
            	});
			});
		</script>
	</head>
	<body>
		<form action="/verify" method="post">
			<input id="response" type="hidden" name="%s" value="">
			<input type="submit" value="提交">
		</form>
	</body>
	</html>
`, sitekey, sitekey, captcha.VerifyRespKey)

	// 输出 html 到 ResponseWriter
	_, err := w.Write([]byte(html))
	if err != nil {
		log.Fatal(err)
	}
}
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	reCaptcha := captcha.NewReCaptcha(sitekey, secretKey)
	isPass, err := reCaptcha.Verify(r)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if isPass == false {
		_, _ = w.Write([]byte("no"))
	}
	_, _ = w.Write([]byte("ok"))
}
func main() {
	http.HandleFunc("/", showHandler)
	http.HandleFunc("/verify", verifyHandler)

	log.Println("example link: http://captcha.local:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```