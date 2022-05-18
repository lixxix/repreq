package repreq

import (
	"fmt"
	"net"
	"sync"
)

type Req struct {
	conn      net.Conn
	Connected bool // true if
	call      IMessage
}

type IMessage interface {
	OnMessage(buf []byte, req *Req)
}

type IRespone interface {
	OnRespone(buf []byte) []byte
}

func (r *Req) Connect(ip string) error {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return err
	}

	r.conn = conn
	r.Connected = true
	go r.Receive()
	return nil
}

func (r *Req) Receive() {
	defer func() {
		fmt.Println("服务器已经断开")
		ServerListen.unregister <- r
		r.conn.Close()
		r.Connected = false
	}()
	buff := make([]byte, 4096)
	for {
		sz, err := r.conn.Read(buff)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if r.call != nil {
			r.call.OnMessage(buff[0:sz], r)
		}
	}
}

func (r *Req) Send(data []byte) {
	r.conn.Write(data)
}

func (r *Req) Run() {
	go r.Receive()
}

func NewReqClient(conn net.Conn, call IMessage) *Req {
	return &Req{
		conn:      conn,
		call:      call,
		Connected: false,
	}
}

type Rep struct {
	Mutex       sync.Mutex // guards
	MutexReq    sync.Mutex // guards
	ResponeCall IRespone
	register    chan *Req
	unregister  chan *Req
	clients     map[*Req]bool
}

var ServerListen *Rep = nil

func ServerLink(conn net.Conn) {
	cli := NewReqClient(conn, ServerListen)
	cli.Run()
}

func (r *Rep) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		con, err := ln.Accept()
		if err != nil {
			return err
		}
		ServerLink(con)
	}
}

func (r *Rep) OnMessage(buf []byte, req *Req) {
	// 回复
	r.MutexReq.Lock()
	if r.ResponeCall != nil {
		req.Send(r.ResponeCall.OnRespone(buf))
	} else {
		rc := "recv:" + string(buf)
		req.Send([]byte(rc))
	}

	r.MutexReq.Unlock()
}

func (r *Rep) Run() {
	fmt.Println("request服务器开启")
	for {
		select {
		case client := <-r.register:
			r.Mutex.Lock()
			r.clients[client] = true
			r.Mutex.Unlock()
			fmt.Println("激活")
		case client := <-r.unregister:
			r.Mutex.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
			}
			r.Mutex.Unlock()
		}
	}
}

func RepServer(respone IRespone) *Rep {
	if ServerListen == nil {
		ServerListen = &Rep{
			Mutex:       sync.Mutex{},
			MutexReq:    sync.Mutex{},
			ResponeCall: respone,
			register:    make(chan *Req),
			unregister:  make(chan *Req),
			clients:     make(map[*Req]bool),
		}
	}
	return ServerListen
}
