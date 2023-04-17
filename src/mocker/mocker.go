package main

import (
	"bufio"
	"bytes"
	"myutil"
	"net"
	"strconv"
)

type Mocker struct {
	MockSet  *myutil.Set
	Listener net.Listener
	Sender   *Sender
}

func NewMocker(ip string, port int) *Mocker {
	sender := NewSender(ip, port)
	sender.Open()
	Log("Sender socket is established!")

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		LogError("Error listening:" + err.Error())
		panic(err)
	}
	Log("Listener is started!")

	return &Mocker{
		MockSet:  myutil.NewSet(),
		Listener: listener,
		Sender:   sender,
	}
}

func (m *Mocker) Start() {
	for {
		conn, err := m.Listener.Accept()
		if err != nil {
			LogError("Error accepting connection:" + err.Error())
			panic(err)
		}

		go m.handleConnection(conn)
	}
}

func (m *Mocker) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	elements := m.MockSet.Elements()
	for {
		request, err := reader.ReadBytes('\n')
		if err != nil {
			LogError("Error reading request:" + err.Error())
			panic(err)
		}

		var response []byte
		for _, item := range elements {
			rr := item.(*ReqRes)
			if bytes.Equal(*rr.Request, request) {
				response = *rr.Response
				break
			}
		}

		if response != nil {
			Log("Response is found!")
			conn.Write(response)
		} else {
			Log("Response isn't found in Mocker, try to send request to real server")
			realResponse := m.Sender.Send(request)
			conn.Write(realResponse)
		}
	}
}
