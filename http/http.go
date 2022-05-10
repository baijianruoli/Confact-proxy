package http

import (
	"confact_client/conf"
	"confact_client/logs"
	"confact_client/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"strconv"
	"sync"
	"time"
)

func GetValue(ctx *gin.Context) {
	key := ctx.Query("key")
    response,err:=RaftGet(key,time.Now().UnixNano()/1e6)
    if err!=nil{
    	ctx.JSON(501,err.Error())
    	return
	}
	ctx.JSON(200, response.Msg)
}

func GetValueBatch(ctx *gin.Context) {
	var (
		req ReqBatch
	)
	if err := ctx.BindJSON(&req); err != nil {
		logs.PrintError(0, err.Error())
		ctx.JSON(501, err.Error())
		return
	}
	raftNodeMap:=make(map[string][]*Req)
	// 分组
	for _,item:=range req.List{
		raftID := util.HashSet(item.Key)
		item.RaftID = raftID
		raftHttp:= conf.JsonConf.RaftsHTTP[raftID]
		if v,ok:=raftNodeMap[raftHttp];!ok{
			raftReqList:=make([]*Req,0)
			raftReqList=append(raftReqList, item)
			raftNodeMap[raftHttp]=raftReqList
		}else{
			v=append(v, item)
		}
	}
	var wg sync.WaitGroup
	for k,v:=range raftNodeMap{
		wg.Add(1)
		// 每个节点一个协程
		go func(k string,v []*Req) {
			fastReq := &fasthttp.Request{}
			fastReq.SetRequestURI(fmt.Sprintf("http://%s/getBatch", k))
			fastReq.Header.SetMethod("POST")
			fastReq.Header.SetContentType("application/json")
			data, _ := json.Marshal(v)
			fastReq.SetBody(data)
			resp := &fasthttp.Response{}
			client := &fasthttp.Client{}
			if err := client.Do(fastReq, resp); err != nil {
				logs.PrintError(1, err.Error())
				ctx.JSON(501, err.Error())
				return
			}
			wg.Done()
		}(k,v)
	}
	wg.Wait()
	ctx.JSON(200, "ok")
}

func GetValueBatchTest(ctx *gin.Context) {
	var (
		req ReqBatch
	)
	num:=ctx.Query("num")
	nums,_:=strconv.Atoi(num)
	req.List=make([]*Req,0)
	for i:=0;i<=nums;i++{
		req.List=append(req.List, &Req{Key: fmt.Sprintf("%d",i)})
	}
	raftNodeMap:=make(map[string][]*Req)
	// 分组
	for _,item:=range req.List{
		raftID := util.HashSet(item.Key)
		item.RaftID = raftID
		raftHttp:= conf.JsonConf.RaftsHTTP[raftID]
		if v,ok:=raftNodeMap[raftHttp];!ok{
			raftReqList:=make([]*Req,0)
			raftReqList=append(raftReqList, item)
			raftNodeMap[raftHttp]=raftReqList
		}else{
			v=append(v, item)
		}
	}
	now:=time.Now()
	var wg sync.WaitGroup
	for k,v:=range raftNodeMap{
		wg.Add(1)
		// 每个节点一个协程
		go func(k string,v []*Req) {
			fastReq := &fasthttp.Request{}
			fastReq.SetRequestURI(fmt.Sprintf("http://%s/getBatch", k))
			fastReq.Header.SetMethod("POST")
			fastReq.Header.SetContentType("application/json")
			data, _ := json.Marshal(v)
			fastReq.SetBody(data)
			resp := &fasthttp.Response{}
			client := &fasthttp.Client{}
			if err := client.Do(fastReq, resp); err != nil {
				logs.PrintError(1, err.Error())
				ctx.JSON(501, err.Error())
				return
			}
			wg.Done()
		}(k,v)
	}
	wg.Wait()
	ctx.JSON(200, time.Now().Sub(now).Seconds())
}

func SetValue(ctx *gin.Context) {
	var (
		req Req
	)
	if err := ctx.BindJSON(&req); err != nil {
		logs.PrintError(0, err.Error())
		ctx.JSON(501, err.Error())
		return
	}
	raftID := util.HashSet(req.Key)
	logs.PrintInfo(raftID, raftID)
	fastReq := &fasthttp.Request{}
	fastReq.SetRequestURI(fmt.Sprintf("http://%s/set", conf.JsonConf.RaftsHTTP[raftID]))
	fastReq.Header.SetMethod("POST")
	fastReq.Header.SetContentType("application/json")
	// 通过raft group找到leader节点
	req.RaftID = raftID
	data, _ := json.Marshal(req)
	fastReq.SetBody(data)
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}

	if err := client.Do(fastReq, resp); err != nil {
		logs.PrintError(raftID, err.Error())
		ctx.JSON(501, err.Error())
		return
	}
	var response Response
	if err:=json.Unmarshal(resp.Body(), &response);err!=nil{
		logs.PrintError(raftID,err.Error())
	}
	ctx.JSON(int(response.Code), response.Msg)
}

func SetValueBatch(ctx *gin.Context) {
	var (
		req ReqBatch
	)
	if err := ctx.BindJSON(&req); err != nil {
		logs.PrintError(0, err.Error())
		ctx.JSON(501, err.Error())
		return
	}
	for i:=0;i<=100000;i++{
		req.List=append(req.List, &Req{Key: fmt.Sprintf("%d",i),Value: 1234})
	}
	t:=time.Now()
	raftNodeMap:=make(map[string][]*Req)
	// 分组
	for _,item:=range req.List{
		raftID := util.HashSet(item.Key)
		item.RaftID = raftID
		raftHttp:= conf.JsonConf.RaftsHTTP[raftID]
		if _,ok:=raftNodeMap[raftHttp];!ok{
			raftReqList:=make([]*Req,0)
			raftReqList=append(raftReqList, item)
			raftNodeMap[raftHttp]=raftReqList
		}else{
			raftNodeMap[raftHttp]=append(raftNodeMap[raftHttp], item)
		}
	}
	var wg sync.WaitGroup
	for k,v:=range raftNodeMap{
		wg.Add(1)
		// 每个节点一个协程
		go func(k string,v []*Req) {
			fastReq := &fasthttp.Request{}
			fastReq.SetRequestURI(fmt.Sprintf("http://%s/setBatch", k))
			fastReq.Header.SetMethod("POST")
			fastReq.Header.SetContentType("application/json")
			data, _ := json.Marshal(v)
			fastReq.SetBody(data)
			resp := &fasthttp.Response{}
			client := &fasthttp.Client{}
			if err := client.Do(fastReq, resp); err != nil {
				logs.PrintError(1, err.Error())
				ctx.JSON(501, err.Error())
				return
			}
			wg.Done()
		}(k,v)
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(t).Seconds())
	ctx.JSON(200, "ok")
}


func SetValueBatchTest(ctx *gin.Context) {
	var (
		req ReqBatch
	)
	num:=ctx.Query("num")
	nums,_:=strconv.Atoi(num)
	req.List=make([]*Req,0)
	for i:=0;i<=nums;i++{
		req.List=append(req.List, &Req{Key: fmt.Sprintf("%d",i),Value: 1234})
	}
	t:=time.Now()
	raftNodeMap:=make(map[string][]*Req)
	// 分组
	for _,item:=range req.List{
		raftID := util.HashSet(item.Key)
		item.RaftID = raftID
		raftHttp:= conf.JsonConf.RaftsHTTP[raftID]
		if _,ok:=raftNodeMap[raftHttp];!ok{
			raftReqList:=make([]*Req,0)
			raftReqList=append(raftReqList, item)
			raftNodeMap[raftHttp]=raftReqList
		}else{
			raftNodeMap[raftHttp]=append(raftNodeMap[raftHttp], item)
		}
	}
	var wg sync.WaitGroup
    for k,v:=range raftNodeMap{
    	wg.Add(1)
    	// 每个节点一个协程
    	 go func(k string,v []*Req) {
			 fastReq := &fasthttp.Request{}
			 fastReq.SetRequestURI(fmt.Sprintf("http://%s/setBatch", k))
			 fastReq.Header.SetMethod("POST")
			 fastReq.Header.SetContentType("application/json")
			 data, _ := json.Marshal(v)
			 fastReq.SetBody(data)
			 resp := &fasthttp.Response{}
			 client := &fasthttp.Client{}
			 if err := client.Do(fastReq, resp); err != nil {
				 logs.PrintError(1, err.Error())
				 ctx.JSON(501, err.Error())
				 return
			 }
			 wg.Done()
		 }(k,v)
	}
	wg.Wait()
	ctx.JSON(200, time.Now().Sub(t).Seconds())
}
