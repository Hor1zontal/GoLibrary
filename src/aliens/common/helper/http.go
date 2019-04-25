/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/9
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package helper

import (
	"github.com/name5566/leaf/chanrpc"
	"aliens/log"
	"encoding/json"
	"strconv"
	"net/http"
)

const COMMAND_HTTP_PROXY int = 1

type HttpProxy struct {
	chanRpc *chanrpc.Server
	server *http.ServeMux
	handlers map[string]func(responseWriter http.ResponseWriter, request *http.Request)
}

func(this *HttpProxy) RegisterFunc(url string, handler func(responseWriter http.ResponseWriter, request *http.Request)) {
	this.handlers[url] = handler
	if this.server != nil {
		this.server.HandleFunc(url, this.ChanRpc)
	} else {
		http.HandleFunc(url, this.ChanRpc)
	}
}

func(this *HttpProxy) Init(chanRpc *chanrpc.Server) {
	this.InitMux(chanRpc, nil)
}

func(this *HttpProxy) InitMux(chanRpc *chanrpc.Server, server *http.ServeMux) {
	this.server = server
	this.handlers = make(map[string]func(responseWriter http.ResponseWriter, request *http.Request))
	this.chanRpc = chanRpc
	this.chanRpc.Register(COMMAND_HTTP_PROXY, this.Handle)
}

func(this *HttpProxy) ChanRpc(responseWriter http.ResponseWriter, request *http.Request) {
	this.chanRpc.Call0(COMMAND_HTTP_PROXY, responseWriter, request)
}

func(this *HttpProxy) Handle(args []interface{}) {
	responseWriter := args[0].(http.ResponseWriter)
	request := args[1].(*http.Request)

	requestPath := request.URL.Path
	//log.Debug("url %v", requestPath)
	handler := this.handlers[requestPath]
	if handler != nil {
		handler(responseWriter, request)
	}
}


type DataResponse struct {
	Code int 		`json:"code"`
	Data interface{}	`json:"data"`
}


func SendToClient(responseWriter http.ResponseWriter, content string) {
	_, err := responseWriter.Write([]byte(content))
	if err != nil {
		log.Debug(err.Error())
	}
}

func GetResponse(code int) string {
	return string("{\"code\":" + strconv.Itoa(code) + "}")
}


func GetDataResponse(code int, data interface{}) string {
	result, _ := json.Marshal(&DataResponse{Code:code, Data:data})
	return string(result)
}



