package main

import (
	"bufio"
	"bytes"
	"myutil"
	"net"
	"strconv"
)

type Mocker struct {
	PreSendSet  *myutil.Set
	PostSendSet *myutil.Set
	Listener    net.Listener
	Sender      *Sender
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
		PreSendSet:  myutil.NewSet(),
		PostSendSet: myutil.NewSet(),
		Listener:    listener,
		Sender:      sender,
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
	preElements := m.PreSendSet.Elements()
	postElements := m.PostSendSet.Elements()
	for {
		request, err := reader.ReadBytes('\n')
		if err != nil {
			LogError("Error reading request:" + err.Error())
			panic(err)
		}

		var response []byte
		for _, item := range preElements {
			rr := item.(*ReqDataResData)
			if bytes.Equal(*rr.ReqData, request) {
				Log("Response is found in PreSendSet!")
				response = *rr.ResData
				break
			}
		}

		if response != nil {
			conn.Write(response)
		} else {
			Log("Response isn't found in Mocker, try to send request to real server")
			response = m.Sender.Send(request)
			for _, item := range postElements {
				rr := item.(*ResLenResData)
				if rr.ResLen == len(response) {
					Log("Response is matched PostSendSet!")
					response = *rr.ResData
					break
				}
			}

			conn.Write(response)
		}
	}
}
