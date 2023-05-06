package main

import (
	"bytes"
	"myutil"
	"net"
	"strconv"
	"time"
)

type Mocker struct {
	MockerPort           int
	ServerIP             string
	ServerPort           int
	MockedReqDataResData *myutil.Set
	MockedResLenResData  *myutil.Set
	Listener             net.Listener
}

func NewMocker(mockerPort int, serverIp string, serverPort int) *Mocker {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(mockerPort))
	if err != nil {
		LogError("Error listening:" + err.Error())
		panic(err)
	}
	Log("Listener is started!")

	return &Mocker{
		MockerPort:           mockerPort,
		ServerIP:             serverIp,
		ServerPort:           serverPort,
		MockedReqDataResData: myutil.NewSet(),
		MockedResLenResData:  myutil.NewSet(),
		Listener:             listener,
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
	serverAddr := m.ServerIP + ":" + strconv.Itoa(m.ServerPort)
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

	preElements := m.MockedReqDataResData.Elements()
	buffer := make([]byte, 64*1024)
	for {
		n, err := clientConn.Read(buffer)
		if err != nil {
			LogError("Error reading from client: " + err.Error())
			return
		}
		Log("Read new client data, length: " + strconv.Itoa(n))
		if n == cap(buffer) {
			LogWarn("The length of data read from client is " + strconv.Itoa(n) + ", which has reached the capacity of the client read buffer. " +
				"The request data may have been fragmented, which could result in the MockedReqDataResData not matching.")
		}
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
			Log("Response is found in MockedReqDataResData!")
			clientConn.Write(response)
			continue
		} else {
			Log("Response isn't found in MockedReqDataResData, try to send request to real server")
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

	postElements := m.MockedResLenResData.Elements()
	buffer := make([]byte, 4*1024*1024)
	for {
		n, err := serverConn.Read(buffer)
		if err != nil {
			LogError("Error reading from server: " + err.Error())
			return
		}
		Log("Read new server data, length: " + strconv.Itoa(n))
		if n == cap(buffer) {
			LogWarn("The length of data read from server is " + strconv.Itoa(n) + ", which has reached the capacity of the server read buffer. " +
				"The response data may have been fragmented, which could result in the MockedResLenResData not matching.")
		}
		response := buffer[:n]

		for _, item := range postElements {
			rr := item.(*ResLenResData)
			if rr.ResLen == len(response) {
				Log("Response is matched in MockedResLenResData!")
				response = *rr.ResData
				break
			}
		}
		clientConn.Write(response)
	}
}
