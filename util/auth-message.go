package util

import (
	"crypto/sha256"
	"crypto/hmac"
	"time"
	"bytes"
	"math/rand"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"strconv"
	"github.com/goodnodes/Syeong_server/config"
	"github.com/goodnodes/Syeong_server/log"
	// "github.com/goodnodes/Syeong_server/model"
)

type SMS struct {
	Type          string `json:"type"`
	From          string `json:"from"`
	Content       string `json:"content"`
	Messages      []struct {
		To       string `json:"to"`
	} `json:"messages"`
}

var cfg = config.GetConfig("config/config.toml")
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
func SendMsg(pNum string) string {
	// Content에 랜덤한 4자리 숫자를 보내야 함.
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(8999) + 1000

	sms := SMS {
		Type : "SMS",
		From : "01083041394",
		Content : "셩 인증번호\n" + strconv.Itoa(random),
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
		logger.Error(err.Error())
		panic(err)
	}

	uri := "/sms/v2/services/" + serviceId + "/messages"
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", accessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", makeSignature("POST", uri))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	var respMap  map[string]interface{}

	json.Unmarshal(respBody, &respMap)
	
	return  respMap["requestId"].(string)
}


// requestId로 messageId를 가져오는 함수
func GetMsgId(requestId string) string {
	req, err := http.NewRequest("GET", "https://sens.apigw.ntruss.com/sms/v2/services/" + serviceId + "/messages?requestId=" + requestId, nil)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	uri := "/sms/v2/services/" + serviceId + "/messages?requestId=" + requestId
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", accessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", makeSignature("GET", uri))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	var messageData map[string]interface{}

	err = json.Unmarshal(respBody, &messageData)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// 메시지 id를 받아서 return해줌

	// 코드 설명 ->  messageData라는 map에 대하여 key는 string, value는 interface{}이다.
	// 이 때 value가 interface{}이기 때문에 다시 배열도 들어갈 수 있을 것이라 생각했다.
	// 그래서 messageData["messages"]의 value가 interface{}타입의 배열이라고 type assertion을 해 주었고
	// 해당 요소의 0번째 요소는 다시 map[string]interface{} 타입이다. 이를 위해 다시 0번째 요소에 대한 type assertion을 해주었으며
	// 그 map의 "messageId"라는 key의 value는 string 타입이기 때문에 다시 type assertion을 해주었다.
	if messageData["statusCode"] == "202" {
		messageId := messageData["messages"].([]interface{})[0].(map[string]interface{})["messageId"].(string)
		return messageId
	} else {
		logger.Error(err.Error())
		panic(err)
	}
	// return messageData.Messages[0].MessageId
}


// messageId로 message Content(내용)을 받아오는 함수
func GetMsgContent(messageId string) string {
	req, err := http.NewRequest("GET", "https://sens.apigw.ntruss.com/sms/v2/services/" + serviceId + "/messages/" + messageId, nil)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	uri := "/sms/v2/services/" + serviceId + "/messages/" + messageId
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", accessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", makeSignature("GET", uri))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	var result map[string]interface{}

	json.Unmarshal(respBody, &result)

	body := result["messages"].([]interface{})[0].(map[string]interface{})["content"].(string)
	return body
}


// 이어져야 할 로직
// 1. 메시지로 랜덤한 숫자를 보내고 body로 requestId와 현재 시간을 보낸다. Do
// 2. 이어지는 확인 라우트에서는 body로 숫자와 requestId, 아까 보낸 시간을 받는다. 이 때 현재 시간이 아까 보낸 시간보다 5분 이상이면 abort Do
// 3. requestId를 가지고 네이버클라우드 메시지 현재 상태 확인 api에 보낸다. 그리고 그 결과값인 messageId로 다시 메시지 확인 api에 보낸다. Do
// 4. 3번을 통해 받은 message body와 2번에서 받은 숫자가 같은지 확인한다. 다르다면 abort한다. 