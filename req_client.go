package repreq

import (
	"fmt"
	"net"
	"time"
)

type ReqClient struct {
	conn          net.Conn
	connected     bool
	address       string
	fun_connected func()
	timeOut       time.Duration
	recv          chan []byte
}

func (t *ReqClient) Connect(address string, connected_func func()) error {
	if t.connected == false {
		t.address = address
		t.fun_connected = connected_func
		conn, err := net.Dial("tcp", address)
		if err != nil {
			//连接失败，退出并开启重连机制
			fmt.Println("DBServer link failed")
			go t.Reconnect()
			return err
		}
		t.conn = conn
		t.connected = true
		t.Run()
		fmt.Println("DBserver link success")
		if connected_func != nil {
			connected_func()
		}
	}
	return nil
}

func (t *ReqClient) IsConnected() bool {
	return t.connected
}

func (t *ReqClient) Reconnect() {
	time.Sleep(time.Second * 3)
	fmt.Println("dbserver reconnectinig..")
	t.Connect(t.address, t.fun_connected)
}

func (t *ReqClient) SetTimeOut(duration time.Duration) {
	t.timeOut = duration
}

func (t *ReqClient) Reading() {
	defer func() {
		t.conn.Close()
		t.connected = false
		t.Reconnect()
	}()
	tempbuff := make([]byte, 4096)
	for {
		sz, err := t.conn.Read(tempbuff)
		if err != nil {
			fmt.Println("error reading", err.Error())
			return
		}
		t.recv <- tempbuff[:sz]
	}
}

func (t *ReqClient) Request(buf []byte) []byte {
	t.conn.Write(buf)
	ticker := time.NewTicker(t.timeOut)
	for {
		select {
		case buf := <-t.recv:
			return buf
		case <-ticker.C:
			return []byte{}
		}
	}
}

func (t *ReqClient) Run() {
	go t.Reading()
}

func CreateReqClient() *ReqClient {
	return &ReqClient{
		connected: false,
		timeOut:   time.Second * 5,
		recv:      make(chan []byte),
	}
}
