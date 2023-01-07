package model

type LoginStruct struct {
	Pnum string `json:"pnum"`
	Pwd string `json:"pwd"`
}

type CheckSMSStruct struct {
	Code string `json:"code"`
	RequestId string `json:"requestid"`
	RequestTime int `json:"requesttime"`
}

type CheckSMSStatusStruct struct {
	RequestId string `json:"requestId"`
	StatusCode string `json:"statusCode"`
	StatusName string `json:"statusName"`
	Messages []struct {
		MessageId string `json:"messageId"`
		RequestTime interface{} `json:"requestTime"`
		ContentType string `json:"contentType"`
		CountryCode string `json:"countryCode"`
		From string `json:"from"`
		To string `json:"to"`
	} `json:"messages"`
}