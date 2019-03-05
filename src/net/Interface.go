package net


type Transport interface {
	Listen()
	Read()
	Write()
	Close()
}

type DefTransport struct{
	Net  string
	Addr string
}

func (this *DefTransport) Listen(){

}

func (this *DefTransport) Read(){

}

func (this *DefTransport) Write(){

}

func (this *DefTransport) Close(){

}