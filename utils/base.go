package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"user/connection"
	"user/domains/models"

	"github.com/devfeel/mapper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func PrintLog(location string, err interface{}) {
	logger := logrus.New()
	file, _ := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		logger.SetOutput(file)
		msg := fmt.Sprintf("%s : %v", location, err)
		logger.Info(msg)
	}
}

func PrintLogSukses(location string, err interface{}) {
	logger := logrus.New()
	file, _ := os.OpenFile("callback.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		logger.SetOutput(file)
		msg := fmt.Sprintf("%s : %v", location, err)
		logger.Info(msg)
	}
}

func IsNil(val interface{}) bool {
	if val == nil {
		return true
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		//use of IsNil method
		return reflect.ValueOf(val).IsNil()
	case reflect.Int, reflect.Float64, reflect.Uint, reflect.String:
		return false
	}
	if reflect.ValueOf(val).Kind() != reflect.Ptr && reflect.ValueOf(val).Len() == 0 {
		return true
	}
	return false
}

func AutoMap(from interface{}, to interface{}) error {
	jsonFrom, _ := json.Marshal(from)
	err := json.Unmarshal([]byte(string(jsonFrom)), to)
	return err
}

func MapRequest(ctx *gin.Context, request *models.BaseRequest, keys []string) *models.BaseRequest {
	valMap := make(map[string]interface{})
	for i := range keys {
		valMap[keys[i]] = ctx.Param(keys[i])
	}
	valMap["token"] = ctx.GetHeader("Authorization")

	// get query param
	queryParams := make(map[string]string)
	for k, v := range ctx.Request.URL.Query() {
		if len(v) == 1 && len(v[0]) != 0 {
			queryParams[k] = v[0]
		}
	}

	request.QueryParam = queryParams

	mapper.MapperMap(valMap, request)
	err := ctx.ShouldBindJSON(request.BodyData)
	PrintLog("Utils Map Request", err)

	return request
}

func CallAPI(httpMethod string, url string, jsonData interface{}, headers map[string]string, queryParams map[string]string) (*http.Response, error) {
	var request *http.Request = &http.Request{}

	if !IsNil(jsonData) {
		jsonValue, _ := json.Marshal(jsonData)
		request, _ = http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonValue))
	} else {
		request, _ = http.NewRequest(httpMethod, url, bytes.NewBuffer(nil))
	}

	if !IsNil(queryParams) {
		q := request.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	request.Header.Set("Content-Type", "application/json")
	for key, val := range headers {
		request.Header.Add(key, val)
	}
	response, err := connection.Client.Do(request)

	if err != nil {
		PrintLog("Utils CallAPI", err)
	}
	return response, err
}
