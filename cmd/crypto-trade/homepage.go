package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
)

type TimeDataInput struct {
	Name string
	Time string
}

type TimeDataOutput struct {
	Result   string
	Text     string
	Time     string
	Duration string
}

func streamTime(timer *sse.Streamer) {
	fmt.Println("Streaming time  started")
	for serviceIsRunning {
		timer.SendString("", "time", time.Now().Format("02/01/2006, 15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

func getTime(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var data TimeDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		var responseData TimeDataOutput
		responseData.Result = "nok"
		responseData.Text = "problem with user json data"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	fmt.Println(data.Name)
	fmt.Println(data.Time)
	timer := time.Now()
	time.Sleep(1 * time.Second)
	end := time.Since(timer)
	fmt.Println("processing takes : " + end.String())
	var responseData TimeDataOutput
	responseData.Result = "ok"
	responseData.Text = "everything went smooth"
	responseData.Time = time.Now().Format("02/01/2006, 15:01:05")
	responseData.Duration = end.String()
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
}
