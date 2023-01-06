package util

import (
	"crypto/sha256"
	"crypto/hmac"
	"time"
	"bytes"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"strconv"
	"github.com/goodnodes/Syeong_server/config"
)

type SMS struct {
	Type          string `json:"type"`
	From          string `json:"from"`
	Content       string `json:"content"`
	Messages      []struct {
		To       string `json:"to"`
	} `json:"messages"`
}

var cfg = config.GetConfig("../config/config.toml")
var accessKey = cfg.SMS.Accesskey
var privateKey = cfg.SMS.Privatekey
var serviceId = cfg.SMS.Serviceid



// 인증 메시지에 필요한 시그니처를 만드는 함수
func makeSignature(method, uri string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	secretKeyBytes := []byte(privateKey)

	message := method + " " + uri + "\n" + timestamp + "\n" + accessKey
	messageBytes := []byte(message)

	h := hmac.New(sha256.New, secretKeyBytes)
	h.Write(messageBytes)
	signature := h.Sum(nil)

	signingKey := base64.StdEncoding.EncodeToString(signature)

	return signingKey
}


// 인증 메시지를 보내는 함수
func SendMsg(pNum string) {

	// Content에 랜덤한 숫자를 보내야 함.

	sms := SMS {
		Type : "SMS",
		From : "01064853201",
		Content : "test Body",
		Messages : []struct {
			To string `json:"to"`
		} {
			{
				To : pNum,
			},
		},
	}

	pBytes, _ := json.Marshal(sms)
	buff := bytes.NewBuffer(pBytes)

	req, err := http.NewRequest("POST", "https://sens.apigw.ntruss.com/sms/v2/services/" + serviceId + "/messages", buff)
	if err != nil {
		panic(err)
	}

}


// 이어져야 할 로직
// 1. 메시지로 랜덤한 숫자를 보내고 body로 requestId와 현재 시간을 보낸다.
// 2. 이어지는 확인 라우트에서는 body로 숫자와 requestId, 아까 보낸 시간을 받는다. 이 때 현재 시간이 아까 보낸 시간보다 5분 이상이면 abort
// 3. requestId를 가지고 네이버클라우드 메시지 현재 상태 확인 api에 보낸다. 그리고 그 결과값인 messageId로 다시 메시지 확인 api에 보낸다.
// 4. 3번을 통해 받은 message body와 2번에서 받은 숫자가 같은지 확인한다. 다르다면 abort한다.