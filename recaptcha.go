package re_captcha_go

import (
	"encoding/json"
	"fmt"
	"github.com/wq1019/re-captcha-go/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	requestTimeout = time.Second * 10
	verifyRespKey  = "g-recaptcha-response"
	reCaptchaApi   = "https://www.google.com/recaptcha/api/siteverify"
)

type reCaptcha struct {
	Sitekey   string // Use this in the HTML code your site serves to users.
	SecretKey string // Use this for communication between your server and Google.
}

// Google reCaptcha response
type reCaptchaResponse struct {
	Score          float32   `json:"score,omitempty"`            // the score for this request (0.0 - 1.0)
	Action         string    `json:"action,omitempty"`           // the action name for this request (important to verify)
	Success        bool      `json:"success"`                    // whether this request was a valid reCAPTCHA token for your site
	Hostname       string    `json:"hostname,omitempty"`         // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes     []string  `json:"error-codes,omitempty"`      // errors
	ChallengeTs    time.Time `json:"challenge_ts"`               // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	ApkPackageName string    `json:"apk_package_name,omitempty"` // the package name of the app where the reCAPTCHA was solved
}

func (r *reCaptcha) Verify(request *http.Request) (bool, error) {
	var (
		err           error
		verifyRespVal = request.FormValue(verifyRespKey)
	)
	// 判断前端有没有传 verifyRespKey
	if verifyRespVal == "" {
		err = errors.BadRequest(fmt.Sprintf("找不到 %s 字段.", verifyRespKey))
		return false, err
	}
	// 请求 api
	client := &http.Client{
		Timeout: requestTimeout,
	}
	resp, err := client.PostForm(reCaptchaApi, url.Values{
		"secret":   {r.SecretKey},
		"response": {verifyRespVal},
		"remoteip": {request.RemoteAddr},
	})
	if err != nil {
		err = errors.BadRequest(err.Error(), err)
		return false, err
	}
	defer resp.Body.Close()
	// 读取 body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.BodyReadError(err)
		return false, err
	}
	// 解析 body 到 struct
	unmarshalResp := &reCaptchaResponse{}
	err = json.Unmarshal(body, unmarshalResp)
	if err != nil {
		err = errors.JsonUnmarshalError(err)
		return false, err
	}
	// 判断有没有错误返回
	if len(unmarshalResp.ErrorCodes) > 0 {
		err = errors.ReCaptchaVerifyError(unmarshalResp.ErrorCodes)
		return false, err
	}
	return unmarshalResp.Success, nil
}

func NewReCaptcha(siteKey, secretKey string) *reCaptcha {
	return &reCaptcha{
		Sitekey:   siteKey,
		SecretKey: secretKey,
	}
}
