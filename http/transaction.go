package http

import (
	"confact_client/conf"
	pb "confact_client/confact/proto"
	"confact_client/logs"
	"confact_client/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"strconv"
)

type ScanEntity struct {
	StartTs int64 `json:"start_ts"`
	EndTs   int64 `json:"end_ts"`
	Key     string `json:"key"`
	Type    pb.LogType `json:"type"`
}

type SetTypeEntity struct {
	RaftID int64 `json:"raft_id"`
	LogEntry pb.LogEntry `json:"log_entry"`
}

type DeleteLockEntity struct {
	Key   string `json:"key"`
	StartTs int64 `json:"start_ts"`
}


func TransactionGet(ctx *gin.Context){
	 key:=ctx.Query("key")
	 ts:=ctx.Query("ts")
	 startTs,_:=strconv.Atoi(ts)
	raftID := util.HashGet(key)
	logs.PrintInfo(raftID, raftID)
	status, resp, err := fasthttp.Get(nil, fmt.Sprintf("http://%s/get?key=%s&ts=%d", conf.JsonConf.RaftsHTTP[raftID], key,startTs))
	if err != nil {
		logs.PrintError(raftID, err.Error())
		ctx.JSON(500,err.Error())
		return
	}
	if status != fasthttp.StatusOK {
		logs.PrintError(raftID, "error")
		ctx.JSON(500,"error")
		return
	}
	var response *Response
	if err:=json.Unmarshal(resp, &response);err!=nil{
		logs.PrintError(raftID,err.Error())
		ctx.JSON(500,err.Error())
		return
	}
	ctx.JSON(200,response.Msg)
	return
}


func TransactionSet(ctx *gin.Context){
	entity:=&SetTypeEntity{}
        if err:=ctx.BindJSON(&entity);err!=nil{
        	ctx.JSON(500,err.Error())
		}
	raftID := util.HashSet(entity.LogEntry.Command.Key)
	logs.PrintInfo(raftID, raftID)
	fastReq := &fasthttp.Request{}
	fastReq.SetRequestURI(fmt.Sprintf("http://%s/setType", conf.JsonConf.RaftsHTTP[raftID]))
	fastReq.Header.SetMethod("POST")
	fastReq.Header.SetContentType("application/json")
	// 通过raft group找到leader节点
     entity.RaftID=raftID

	data, _ := json.Marshal(entity)
	fastReq.SetBody(data)
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}

	if err := client.Do(fastReq, resp); err != nil {
		logs.PrintError(raftID, err.Error())
		ctx.JSON(500,err.Error())
		return
	}
	var response *Response
	if err:=json.Unmarshal(resp.Body(), &response);err!=nil{
		logs.PrintError(raftID,err.Error())
		ctx.JSON(500,err.Error())
		return
	}
	ctx.JSON(200,response.Msg)
	return
}


func TransactionScan(ctx *gin.Context){
	entity:=&ScanEntity{}
	if err:=ctx.BindJSON(&entity);err!=nil{
		ctx.JSON(500,err.Error())
	}

	raftID := util.HashGet(entity.Key)
	logs.PrintInfo(raftID, raftID)
	fastReq := &fasthttp.Request{}
	fastReq.SetRequestURI(fmt.Sprintf("http://%s/transaction/scan", conf.JsonConf.RaftsHTTP[raftID]))
	fastReq.Header.SetMethod("POST")
	fastReq.Header.SetContentType("application/json")
	data, _ := json.Marshal(entity)
	fastReq.SetBody(data)
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}

	if err := client.Do(fastReq, resp); err != nil {
		logs.PrintError(raftID, err.Error())
		ctx.JSON(500,err.Error())
		return
	}
	var response *Response
	if err:=json.Unmarshal(resp.Body(), &response);err!=nil{
		logs.PrintError(raftID,err.Error())
		ctx.JSON(500,err.Error())
		return
	}
	ctx.JSON(200,response.Msg)
	return
}

