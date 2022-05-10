package util

import (
	"confact_client/conf"
	pb "confact_client/confact/proto"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"sync"
	"unsafe"
)


var grpcClientMap sync.Map
var LeaderMap map[int64]int64
var RaftFlag map[int64]bool

func StringToByte(s string) []byte {
	var b []byte
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pByte := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pByte.Data = sh.Data
	pByte.Len = sh.Len
	pByte.Cap = sh.Len
	return b
}

func ByteToString(b []byte) string {
	var s string
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pByte := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Data = pByte.Data
	sh.Len = pByte.Len
	return s
}

func  GrpcClient(server int64) pb.RaftClient {
	client, ok :=grpcClientMap.Load(server)
	if ok {
		return client.(pb.RaftClient)
	} else {
		conn, err := grpc.Dial(conf.JsonConf.RaftsRPC[server], grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		c := pb.NewRaftClient(conn)
		grpcClientMap.Store(server, c)
		return c
	}
}