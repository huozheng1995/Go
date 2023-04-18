package main

import (
	"net"
	"strconv"
	"time"
)

type Sender struct {
	IP   string
	Port int
	Conn net.Conn
}

func NewSender(ip string, port int) *Sender {
	return &Sender{
		IP:   ip,
		Port: port,
	}
}

func (sender *Sender) Open() {
	destAddr := sender.IP + ":" + strconv.Itoa(sender.Port)
	conn, err := net.Dial("tcp", destAddr)
	if err != nil {
		LogError("Error connecting:" + err.Error())
		panic(err)
	}

	// Set keep-alive parameters
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)

	sender.Conn = conn
}

func (sender *Sender) Close() {
	sender.Conn.Close()
}

func (sender *Sender) Send(request []byte) (response []byte, err error) {
	// Send a TCP packet to the server
	_, err = sender.Conn.Write(request)
	if err != nil {
		LogError("Error sending request to real server: " + err.Error())
		return nil, err
	}

	buffer := make([]byte, 4096)
	n, err := sender.Conn.Read(buffer)
	if err != nil {
		LogError("Error receiving response: " + err.Error())
		return nil, err
	}
	response = buffer[:n]
	return response, nil
}
