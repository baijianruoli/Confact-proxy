package rpc

import "context"
import 	pb "confact_client/confact/proto"

type MiddleWare struct {

}


func (mw *MiddleWare) AppendEntries(ctx context.Context, args *pb.AppendEntriesArgs) (reply *pb.AppendEntriesReply, error error){
      return &pb.AppendEntriesReply{},nil
}

func (mw *MiddleWare) RequestVote(ctx context.Context, args *pb.RequestVoteArgs) (reply *pb.RequestVoteReply, error error) {
	return &pb.RequestVoteReply{},nil
}

func (mw *MiddleWare) HeartBeat(ctx context.Context,args *pb.HeartBeatArgs)(reply *pb.HeartBeatReply,err error){
	return &pb.HeartBeatReply{},nil
}