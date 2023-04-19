package main

import (
	"bytes"
	"myutil"
	"net"
	"strconv"
	"time"
)

type Mocker struct {
	IP          string
	Port        int
	PreSendSet  *myutil.Set
	PostSendSet *myutil.Set
	Listener    net.Listener
}

func NewMocker(ip string, port int) *Mocker {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		LogError("Error listening:" + err.Error())
		panic(err)
	}
	Log("Listener is started!")

	return &Mocker{
		IP:          ip,
		Port:        port,
		PreSendSet:  myutil.NewSet(),
		PostSendSet: myutil.NewSet(),
		Listener:    listener,
	}
}

func (m *Mocker) Start() {
	for {
		clientConn, err := m.Listener.Accept()
		if err != nil {
			LogError("Error accepting connection, error :" + err.Error())
			return
		}
		Log("Client socket is established!")

		serverConn, err := m.connectServer()
		if err != nil {
			LogError("Error connecting server, error: " + err.Error())
			return
		}
		Log("Server socket is established!")

		go m.handleClientSocket(clientConn, serverConn)
		go m.handleServerSocket(clientConn, serverConn)
	}
}

func (m *Mocker) connectServer() (net.Conn, error) {
	serverAddr := m.IP + ":" + strconv.Itoa(m.Port)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	// Set keep-alive parameters
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)

	return conn, nil
}

func (m *Mocker) handleClientSocket(clientConn net.Conn, serverConn net.Conn) {
	defer clientConn.Close()

	preElements := m.PreSendSet.Elements()
	for {
		buffer := make([]byte, 4096)
		n, err := clientConn.Read(buffer)
		if err != nil {
			LogError("Error reading from client: " + err.Error())
			return
		}
		Log("Read new client data")
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
			clientConn.Write(response)
			continue
		} else {
			Log("Response isn't found in PreSendSet, try to send request to real server")
			_, err = serverConn.Write(request)
			if err != nil {
				LogError("Error sending request to real server: " + err.Error())
				return
			}
		}
	}
}

func (m *Mocker) handleServerSocket(clientConn net.Conn, serverConn net.Conn) {
	defer serverConn.Close()

	postElements := m.PostSendSet.Elements()
	for {
		buffer := make([]byte, 4096)
		n, err := serverConn.Read(buffer)
		if err != nil {
			LogError("Error reading from server: " + err.Error())
			return
		}
		Log("Read new server data")
		response := buffer[:n]

		for _, item := range postElements {
			rr := item.(*ResLenResData)
			if rr.ResLen == len(response) {
				Log("Response is matched in PostSendSet!")
				response = *rr.ResData
				break
			}
		}
		clientConn.Write(response)
	}
}
