package util

import (
	"io/ioutil"
	"github.com/goodnodes/Syeong_server/model"
	"net/url"
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/goodnodes/Syeong_server/log"
)

var ID = cfg.GEO.Clientid
var PWD = cfg.GEO.Clientsecret

func GetGEO(address string) *model.GEO {
	url := "https://naveropenapi.apigw.ntruss.com/map-geocode/v2/geocode" + "?query=" + url.QueryEscape(address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	req.Header.Set("X-NCP-APIGW-API-KEY-ID", ID)
	req.Header.Set("X-NCP-APIGW-API-KEY", PWD)

	client := &http.Client{}
	resp, err := client.Do(req)
	ErrorHandler(err)
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	ErrorHandler(err)
	
	var data map[string]interface{}
	
	err = json.Unmarshal(respBody, &data)
	if err != nil{
		logger.Error(err.Error())
	}

	latitude := data["addresses"].([]interface{})[0].(map[string]interface{})["y"].(string)
	longitude := data["addresses"].([]interface{})[0].(map[string]interface{})["x"].(string)

	geo := model.GEO{
		Latitude : latitude,
		Longitude : longitude,
	}

	fmt.Println(geo)

	return &geo
}