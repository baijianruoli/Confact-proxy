package http

import (
	"confact_client/conf"
	"confact_client/logs"
	"confact_client/util"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)


func RaftGet(key string,ts int64) (*Response,error){
	raftID := util.HashGet(key)
	logs.PrintInfo(raftID, raftID)
	status, resp, err := fasthttp.Get(nil, fmt.Sprintf("http://%s/get?key=%s&ts=%d", conf.JsonConf.RaftsHTTP[raftID], key,ts))
	if err != nil {
		logs.PrintError(raftID, err.Error())
		return nil,err
	}
	if status != fasthttp.StatusOK {
		logs.PrintError(raftID, "error")
		return nil,err
	}
	var response *Response
	if err:=json.Unmarshal(resp, &response);err!=nil{
		logs.PrintError(raftID,err.Error())
		return nil,err
	}
	return response,nil
}

