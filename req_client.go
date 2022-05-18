package repreq

import (
	"net"
	"time"
)

type ReqClient struct {
	conn      net.Conn
	connected bool
	timeOut   time.Duration
	recv      chan []byte
}

func (t *ReqClient) Connect(address string) error {
	if t.connected == false {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return err
		}
		t.conn = conn
		t.connected = true
		t.Run()
	}
	return nil
}

func (t *ReqClient) SetTimeOut(duration time.Duration) {
	t.timeOut = duration
}

func (t *ReqClient) Reading() {
	defer func() {
		t.conn.Close()
		t.connected = false
	}()
	tempbuff := make([]byte, 4096)
	for {
		sz, err := t.conn.Read(tempbuff)
		if err != nil {
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
