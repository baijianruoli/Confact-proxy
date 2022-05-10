package util

import (
	"confact_client/conf"
	"crypto/md5"
	"math/rand"
)

func HashGet(key string) int64 {
	var ros int64
	var calc int64
	res := md5.Sum(StringToByte(key))
	for _, item := range res {
		ros += int64(item)
	}
	//找到对应的raft group
	raftGroup:=ros % conf.TomlConf.Nodes
    random:=rand.Intn(int(conf.TomlConf.Replicate))
    // total是raft
	raftGroupIndex:=raftGroup*conf.TomlConf.Replicate
    // calc是随机数
    calc=int64(random)
    // 随机加轮询的负载均衡算法
    // 可以继续优化
    for ok:=RaftFlag[calc+raftGroupIndex];!ok;ok=RaftFlag[calc+raftGroupIndex]{
    	calc=(calc+1)%conf.TomlConf.Replicate
    	if calc==int64(random){
    		return -1
		}
	}
	return calc+raftGroupIndex
}

func HashSet(key string) int64 {
	var ros int64
	res := md5.Sum(StringToByte(key))
	for _, item := range res {
		ros += int64(item)
	}
	raftGroup:=ros % conf.TomlConf.Nodes
	// 找到raft group的leader
	return LeaderMap[raftGroup]
}
