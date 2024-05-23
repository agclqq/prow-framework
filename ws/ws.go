package ws

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Msg struct {
	Type int
	Data []byte
}

type Ws struct {
	conn      *websocket.Conn
	rcvChan   chan *Msg //读取到的消息
	sendChan  chan *Msg //将此消息
	closeChan chan []byte
	closed    bool
	sync.Mutex
}

func NewWs(wsConn *websocket.Conn) *Ws {
	return &Ws{
		conn:      wsConn,
		rcvChan:   make(chan *Msg, 1024),
		sendChan:  make(chan *Msg, 1024),
		closeChan: make(chan []byte, 1),
	}
}

func (w *Ws) Run() {
	go w.Receive()
	go w.Send()
	go w.heartbeat()
}
func (w *Ws) heartbeat() {
	for {
		time.Sleep(2 * time.Second)
		err := w.conn.WriteMessage(websocket.PingMessage, []byte("ping"))
		if err != nil {
			fmt.Println("send heartbeat stop")
			return
		}
	}
}
func (w *Ws) Receive() {
	var err error
	for {
		msgType, data, err1 := w.conn.ReadMessage()
		if err1 != nil {
			err = err1
			goto ERROR
		}
		w.rcvChan <- &Msg{Type: msgType, Data: data}
	}
ERROR:
	w.close(err)
}

func (w *Ws) Send() {
	var err error
	for {
		select {
		case msg := <-w.sendChan:
			err = w.conn.WriteMessage(msg.Type, msg.Data)
			if err != nil {
				goto ERROR
			}
		}
	}
ERROR:
	w.close(err)
}

func (w *Ws) close(err error) error {
	w.Close()
	fmt.Println(err)
	return err
}

func (w *Ws) Close() {
	if w.closed {
		return
	}
	w.Lock()
	defer w.Unlock()
	close(w.sendChan)
	close(w.rcvChan)
	w.conn.Close()
	w.closed = true

}

type ConsumeFunc func(*Msg) error

func (w *Ws) ConsumeMsg(f ConsumeFunc) error {
	for {
		select {
		case msg, ok := <-w.rcvChan:
			if !ok {
				return errors.New("chan has been closed and can no longer consume")
			}
			err := f(msg)
			if err != nil {
				return err
			}
		}
	}
}
func (w *Ws) ProductMsg(msg *Msg) error {
	if w.closed {
		return errors.New("chan has been closed and can no longer consume")
	}
	w.sendChan <- msg
	return nil
}
