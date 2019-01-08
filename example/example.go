package main

import (
	"encoding/json"
	"fmt"
	"github.com/wq1019/re-captcha-go"
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
		<script src='https://cdn.bootcss.com/jquery/1.12.0/jquery.min.js'></script>
		<script>
        	grecaptcha.ready(function () {
            	grecaptcha.execute('%s', {action: 'homepage'}).then(function (token) {
                    // Verify the token on the server.
					$.ajax({
                    	url: "/verify",
                    	data: {
                        	"g-recaptcha-response": token
                    	},
                   		type: "POST",
                    	success: function (res) {
                     	   console.log(res);
                    	},
						error: function (res) {
							alert(res);
						}
                	});
            	});
			});
		</script>
	</head>
	<body>
		<h1 style="text-align:center;margin-top:30px;">Google reCaptcha test</h1>
	</body>
	</html>
`, sitekey, sitekey)

	// 输出 html 到 ResponseWriter
	_, err := w.Write([]byte(html))
	if err != nil {
		log.Fatal(err)
	}
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	reCaptcha := re_captcha_go.NewReCaptcha(sitekey, secretKey)
	isOk, err := reCaptcha.Verify(r)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if isOk == false {
		_, _ = w.Write([]byte("no"))
	}
	data := map[string]interface{}{
		"success": isOk,
	}
	resp, err := json.Marshal(data)
	if err != nil {
		_, _ = w.Write([]byte(re_captcha_go.JsonMarshalError(err).Error()))
		return
	}
	_, _ = w.Write(resp)
}

func main() {
	http.HandleFunc("/", showHandler)
	http.HandleFunc("/verify", verifyHandler)

	log.Println("example link: http://captcha.local:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
