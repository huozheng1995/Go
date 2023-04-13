package main

import (
	"bufio"
	"fmt"
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

func (sender *Sender) Open() error {
	destAddr := sender.IP + ":" + strconv.Itoa(sender.Port)
	conn, err := net.Dial("tcp", destAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return err
	}

	// Set keep-alive parameters
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)

	sender.Connection = conn
	return nil
}

func (sender *Sender) Close() {
	sender.Connection.Close()
}

func (sender *Sender) Send(request []byte) (response []byte, err error) {
	// Send a TCP packet to the server
	_, err = sender.Connection.Write(request)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return nil, err
	}

	reader := bufio.NewReader(sender.Connection)
	response, err = reader.ReadBytes('\n')
	_, err = sender.Connection.Read(response)
	if err != nil {
		fmt.Println("Error receiving data:", err)
		return nil, err
	}

	return response, nil
}
