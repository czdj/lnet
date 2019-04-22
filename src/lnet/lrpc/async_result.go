package lrpc

import (
	"errors"
	"fmt"
	"lnet/lutils"
	"proto/pb"
	"sync"
	"time"
)

type AsyncResult struct {
	uid    string
	result chan *pb.RpcRspInfo
}

var G_asyncResult *AsyncResultMgr = NewAsyncResultMgr()

func NewAsyncResult(uid string) *AsyncResult {
	return &AsyncResult{
		uid:    uid,
		result: make(chan *pb.RpcRspInfo, 1),
	}
}

func (this *AsyncResult) GetUid() string {
	return this.uid
}

func (this *AsyncResult) SetResult(data *pb.RpcRspInfo) {
	this.result <- data
}

func (this *AsyncResult) GetResult(timeout time.Duration) (*pb.RpcRspInfo, error) {
	select {
	case <-time.After(timeout):
		close(this.result)
		return &pb.RpcRspInfo{}, errors.New(fmt.Sprintf("GetResult AsyncResult: timeout %s", this.uid))
	case result := <-this.result:
		return result, nil
	}
	return &pb.RpcRspInfo{}, errors.New("GetResult AsyncResult error. reason: no")
}

type AsyncResultMgr struct {
	results map[string]*AsyncResult
	sync.RWMutex
}

func NewAsyncResultMgr() *AsyncResultMgr {
	return &AsyncResultMgr{
		results: make(map[string]*AsyncResult, 0),
	}
}

func (this *AsyncResultMgr) Add() *AsyncResult {
	this.Lock()
	defer this.Unlock()

	r := NewAsyncResult(lutils.GetUUIDStr())
	this.results[r.GetUid()] = r

	return r
}

func (this *AsyncResultMgr) Remove(uid string) {
	this.Lock()
	defer this.Unlock()

	delete(this.results, uid)
}

func (this *AsyncResultMgr) GetAsyncResult(uid string) (*AsyncResult, error) {
	this.RLock()
	defer this.RUnlock()

	r, ok := this.results[uid]
	if ok {
		return r, nil
	} else {
		return nil, errors.New("not found AsyncResult")
	}
}

func (this *AsyncResultMgr) FillAsyncResult(uid string, data *pb.RpcRspInfo) error {
	r, err := this.GetAsyncResult(uid)
	if err == nil {
		this.Remove(uid)
		r.SetResult(data)
		return nil
	} else {
		return err
	}
}
