package lnet

import (
	"errors"
	"go.uber.org/zap"
	"lnet/iface"
	"sync"
)

/*
	连接管理模块
*/
type TransportManager struct {
	transports    map[uint32]iface.ITransport //管理的连接信息
	transportLock sync.RWMutex                //读写连接的读写锁
}

/*
	创建一个连接管理
*/
func NewTransportManager() *TransportManager {
	return &TransportManager{
		transports: make(map[uint32]iface.ITransport),
	}
}

//添加连接
func (this *TransportManager) Add(trans iface.ITransport) {
	//保护共享资源Map 加写锁
	this.transportLock.Lock()
	defer this.transportLock.Unlock()

	//将conn连接添加到TransportManager中
	this.transports[trans.GetId()] = trans
	Logger.Info("connection add to ConnManager successfully", zap.Int32("num", this.Len()))
}

//删除连接
func (this *TransportManager) Remove(trans iface.ITransport) {
	//保护共享资源Map 加写锁
	this.transportLock.Lock()
	defer this.transportLock.Unlock()

	//删除连接信息
	delete(this.transports, trans.GetId())

	Logger.Info("connection Remove successfully", zap.Uint32("Id", trans.GetId()), zap.Int32("num", this.Len()))
}

//利用ConnID获取链接
func (this *TransportManager) Get(id uint32) (iface.ITransport, error) {
	//保护共享资源Map 加读锁
	this.transportLock.RLock()
	defer this.transportLock.RUnlock()

	if trans, ok := this.transports[id]; ok {
		return trans, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接长度
func (this *TransportManager) Len() int32 {
	return int32(len(this.transports))
}

//清除并停止所有连接
func (this *TransportManager) ClearTransport() {
	//保护共享资源Map 加写锁
	this.transportLock.Lock()
	defer this.transportLock.Unlock()

	//停止并删除全部的连接信息
	for connID, trans := range this.transports {
		//停止
		trans.Close()
		//删除
		delete(this.transports, connID)
	}

	Logger.Info("Clear All transports successfully", zap.Int32("num", this.Len()))
}
