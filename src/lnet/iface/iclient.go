package iface

type IClient interface {
	Connect() error
	Send(msg interface{})error
}
