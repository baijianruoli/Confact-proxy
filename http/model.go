package http



type Req struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	RaftID int64      `json:"raft_id"`
}

type Response struct {
	Code int64 `json:"code"`
	Msg  interface{} `json:"msg"`
}



type ReqBatch struct {
	List  []*Req `json:"list"`
}
