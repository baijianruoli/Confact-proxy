package main

import (
	"confact_client/conf"
	pb "confact_client/confact/proto"
	"confact_client/http"
	"confact_client/logs"
	"confact_client/util"
	"context"
	"github.com/gin-gonic/gin"
	"confact_client/middleware"
	"time"
)



func main() {
	conf.ConfigInit()
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/get", http.GetValue)
	r.GET("/getBatch", http.GetValueBatch)
	r.GET("/getBatchTest", http.GetValueBatchTest)
	r.POST("/set",http.SetValue)
	r.POST("/setBatch",http.SetValueBatch)
	r.GET("/setBatchTest",http.SetValueBatchTest)


	// transaction
	r.GET("/transaction/get",http.TransactionGet)
	r.POST("/transaction/set",http.TransactionSet)
	r.POST("/transaction/scan",http.TransactionScan)
	go r.Run(conf.TomlConf.MiddleWareHTTP) // listen and serve on 0.0.0.0:8080

    HeartBeat()

	//router := fasthttprouter.New()
	//// 不同的路由执行不同的处理函数
	//router.GET("/", Index)
	 //fasthttp.ListenAndServe(":7788",router.Handler)
}




func HeartBeat(){
    ticker:=time.NewTicker(500*time.Millisecond)
    util.LeaderMap=make(map[int64]int64)
    util.RaftFlag=make(map[int64]bool)
	for{
		<-ticker.C
		for k,_:=range conf.JsonConf.RaftsRPC{
			client:=util.GrpcClient(k)
			resp,err:=client.HeartBeat(context.TODO(),&pb.HeartBeatArgs{StartTs: time.Now().UnixNano()/1e6})
			if err!=nil{
				util.RaftFlag[k]=false
				if v,ok:=util.LeaderMap[k/conf.JsonConf.Replicate];ok&&v==k{
					util.LeaderMap[k/conf.JsonConf.Replicate]=-1
				}
				logs.PrintError(k,"心跳检测失败")
				continue
			}
			//logs.PrintInfo(k,"心跳检测成功：",resp.IsLeader)
			util.RaftFlag[k]=true
			// 找到哪一个raft group
			if resp.IsLeader{
				util.LeaderMap[k/conf.JsonConf.Replicate]=k
			}
		}
	}

}

//func Index(ctx *fasthttp.RequestCtx) {
//	fmt.Fprint(ctx, "Welcome")
//}