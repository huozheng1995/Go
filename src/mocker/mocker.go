package main

import (
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
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		LogError("Error listening:" + err.Error())
		panic(err)
	}
	Log("Listener is started!")

	sender := NewSender(ip, port)

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
			return
		}
		Log("Source socket is established!")

		m.Sender.Open()
		Log("Sender socket is established!")

		go m.handleReqConn(conn)
		go m.handleResConn(conn)
	}
}

func (m *Mocker) handleReqConn(reqConn net.Conn) {
	defer reqConn.Close()

	for {
		buffer := make([]byte, 4096)
		n, err := reqConn.Read(buffer)
		if err != nil {
			LogError("Error reading request: " + err.Error())
			return
		}
		Log("Read new request data")
		request := buffer[:n]

		Log("Response isn't found in PreSendSet, try to send request to real server")
		_, err = m.Sender.Conn.Write(request)
		if err != nil {
			LogError("Error sending request to real server: " + err.Error())
			return
		}
	}
}
func (m *Mocker) handleResConn(reqConn net.Conn) {
	defer m.Sender.Conn.Close()

	for {
		buffer := make([]byte, 4096)
		n, err := m.Sender.Conn.Read(buffer)
		if err != nil {
			LogError("Error receiving response: " + err.Error())
			return
		}

		reqConn.Write(buffer[:n])
	}
}

func (m *Mocker) handleReqConn2(conn net.Conn) {
	defer conn.Close()

	m.Sender.Open()
	defer m.Sender.Conn.Close()
	Log("Sender socket is established!")

	serverHello := make([]byte, 4096)
	n, err := m.Sender.Conn.Read(serverHello)
	if err != nil {
		LogError("Error receiving ServerHello: " + err.Error())
		return
	}

	_, err = conn.Write(serverHello[:n])
	if err != nil {
		LogError("Error Hello : " + err.Error())
		return
	}

	preElements := m.PreSendSet.Elements()
	postElements := m.PostSendSet.Elements()
	for {
		buffer := make([]byte, 4096)
		n, err = conn.Read(buffer)
		if err != nil {
			LogError("Error reading request: " + err.Error())
			break
		}
		Log("Read new request data")
		request := buffer[:n]

		var response []byte
		for _, item := range preElements {
			rr := item.(*ReqDataResData)
			if bytes.Equal(*rr.ReqData, request) {
				response = *rr.ResData
				break
			}
		}
		if response != nil {
			Log("Response is found in PreSendSet!")
			conn.Write(response)
			continue
		}

		Log("Response isn't found in PreSendSet, try to send request to real server")
		response, err = m.Sender.Send(request)
		if err != nil {
			continue
		}
		Log("Got the response from the real server!")
		for _, item := range postElements {
			rr := item.(*ResLenResData)
			if rr.ResLen == len(response) {
				Log("Response is matched in PostSendSet!")
				response = *rr.ResData
				break
			}
		}
		conn.Write(response)
	}
}
