package main

import (
	"bufio"
	"net"
	"strconv"
	"time"
)

type Sender struct {
	IP         string
	Port       int
	Connection net.Conn
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

	sender.Connection = conn
}

func (sender *Sender) Close() {
	sender.Connection.Close()
}

func (sender *Sender) Send(request []byte) (response []byte) {
	// Send a TCP packet to the server
	_, err := sender.Connection.Write(request)
	if err != nil {
		LogError("Error sending request to real server:" + err.Error())
		panic(err)
	}

	reader := bufio.NewReader(sender.Connection)
	response, err = reader.ReadBytes('\n')
	if err != nil {
		LogError("Error receiving data:" + err.Error())
		panic(err)
	}

	return response
}
