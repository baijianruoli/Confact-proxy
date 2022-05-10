package conf

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gookit/color"
	"io/ioutil"
	"sync"
)

type TomlConfig struct {
	Nodes int64
	Replicate int64
	NodesHTTP []string
	MiddleWareRPC string
	MiddleWareHTTP string
}

type JsonConfig struct {
	Nodes int64 `json:"nodes"`
	Replicate int64 `json:"replicate"`
	NodesHTTP []string `json:"nodes_http"`
	NodesInfo map[int64]*NodeInfo `json:"nodes_info"`
	RaftsRPC  map[int64]string `json:"rafts_rpc"`
	RaftsHTTP  map[int64]string `json:"rafts_http"`
}

type NodeInfo struct {
	NodeID     int64 `json:"node_id"`
    NodeRafts  []*RaftInfo `json:"node_rafts"`

}

type RaftInfo struct {
	RaftID   int64 `json:"raft_id"`
	RaftRPC  string `json:"raft_rpc"`
	RaftList []int64 `json:"raft_list"`
}



var TomlConf TomlConfig
var JsonConf JsonConfig
var once sync.Once

// 生成multi-raft配置
func ConfigInit(){
	once.Do(func() {
		toml.DecodeFile("conf.toml", &TomlConf)
	})
	jsonConfig:=JsonConfig{
		Nodes: TomlConf.Nodes,
		Replicate: TomlConf.Replicate,
		NodesHTTP: TomlConf.NodesHTTP,
		NodesInfo: make(map[int64]*NodeInfo),
		RaftsRPC: make(map[int64]string),
		RaftsHTTP: make(map[int64]string),
	}
	atomicIndex:=int64(0)

	// 先初始化
	for i:=int64(0);i<TomlConf.Nodes;i++{
		jsonConfig.NodesInfo[i]=&NodeInfo{NodeID: i}
	}
	for i:=int64(0);i<TomlConf.Nodes;i++{

		idList:=make([]int64,0)
		for j:=atomicIndex;j<atomicIndex+TomlConf.Replicate;j++{
			idList=append(idList, j)
		}

		for j:=int64(0);j<TomlConf.Replicate;j++{
			raftInfo:=&RaftInfo{}
			raftInfo.RaftID=atomicIndex
			raftInfo.RaftList=idList
			raftInfo.RaftRPC=fmt.Sprint("localhost:",50010+atomicIndex)
			jsonConfig.NodesInfo[(i+j)%TomlConf.Nodes].NodeRafts=append(jsonConfig.NodesInfo[(i+j)%TomlConf.Nodes].NodeRafts, raftInfo)
		    jsonConfig.RaftsRPC[raftInfo.RaftID]=raftInfo.RaftRPC
		    jsonConfig.RaftsHTTP[raftInfo.RaftID]=TomlConf.NodesHTTP[(i+j)%TomlConf.Nodes]
			atomicIndex++
		}
	}
	JsonConf=jsonConfig
	data,err:=json.Marshal(jsonConfig)
	if err!=nil{
		color.Red.Println(err.Error())
	}
	if err:=ioutil.WriteFile("multi-raft.json",data,0644);err!=nil{
		color.Red.Println(err.Error())
	}

}
