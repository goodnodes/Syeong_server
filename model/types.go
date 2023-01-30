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

type GEO struct {
	Latitude string `bson:"latitude" json:"latitude"`
	Longitude string `bson:"longitude" json:"longitude"`
}

type Tag struct {
	Key string
	Value int
}