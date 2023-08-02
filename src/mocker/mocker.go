package main

import (
	"bytes"
	"fmt"
	"github.com/edward/mocker/logger"
	"myutil"
	"net"
	"strconv"
	"time"
)

type Mocker struct {
	MockedReqDataResData      *myutil.Set
	MockedReqDataResDataFiles *myutil.Set
	MockedResLenResData       *myutil.Set
	Settings                  Settings
}

func NewMocker(settings Settings) *Mocker {
	return &Mocker{
		MockedReqDataResData:      myutil.NewSet(),
		MockedReqDataResDataFiles: myutil.NewSet(),
		MockedResLenResData:       myutil.NewSet(),
		Settings:                  settings,
	}
}

func (m *Mocker) Start() {
	var listener net.Listener
	var err error
	if m.Settings.VirtualNetworkInterfaceMode {
		_, err = CreateNetworkInterface(m.Settings.ServerIP)
		if err != nil {
			logger.LogError("Failed to create Network Interface, error: " + err.Error())
			panic(err)
		}
		listener, err = net.Listen("tcp", m.Settings.ServerIP+":"+strconv.Itoa(m.Settings.MockerPort))
	} else {
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(m.Settings.MockerPort))
	}
	if err != nil {
		logger.LogError("Error listening, error: " + err.Error())
		panic(err)
	}
	logger.Log("Listening TCP connections on " + listener.Addr().String() + "...")

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			logger.LogError("Error accepting connection, error :" + err.Error())
			panic(err)
		}
		logger.Log(fmt.Sprintf("Client socket is established! Client: %v -> Local: %v", clientConn.RemoteAddr(), clientConn.LocalAddr()))

		var serverConn net.Conn
		if m.Settings.VirtualNetworkInterfaceMode {
			serverConn, err = m.connectServerByLocalInterface()
		} else {
			serverConn, err = m.connectServer()
		}
		if err != nil {
			logger.LogError("Error connecting server, error: " + err.Error())
			panic(err)
		}
		logger.Log(fmt.Sprintf("Server socket is established! Local: %v -> Server: %v", serverConn.LocalAddr(), serverConn.RemoteAddr()))

		go m.handleClientSocket(clientConn, serverConn)
		go m.handleServerSocket(clientConn, serverConn)
	}
}

func (m *Mocker) connectServer() (net.Conn, error) {
	serverAddr := m.Settings.ServerIP + ":" + strconv.Itoa(m.Settings.ServerPort)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)
	return conn, nil
}

func (m *Mocker) connectServerByLocalInterface() (net.Conn, error) {
	ip := m.Settings.LocalNetworkInterfaceAddress
	_, err := GetNetworkInterface(ip)
	if err != nil {
		logger.LogError("Cannot find local network interface, ip address: "+ip, err)
		return nil, err
	}

	localPort := m.Settings.ServerPort
	if localPort < 10000 {
		localPort += 10000
	}
	localAddr := &net.TCPAddr{IP: net.ParseIP(ip), Port: localPort}
	remoteAddr := &net.TCPAddr{IP: net.ParseIP(m.Settings.ServerIP), Port: m.Settings.ServerPort}
	dialer := net.Dialer{LocalAddr: localAddr}
	conn, err := dialer.Dial("tcp", remoteAddr.String())
	if err != nil {
		logger.LogError("Failed to connect to server through local network interface", err)
		return nil, err
	}

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)
	return conn, nil
}

func (m *Mocker) handleClientSocket(clientConn net.Conn, serverConn net.Conn) {
	defer clientConn.Close()

	logger.Log("Listening on client socket...")
	resDataSet := m.MockedReqDataResData.Elements()
	resDataFiles := m.MockedReqDataResDataFiles.Elements()
	buffer := make([]byte, 64*1024)
	for {
		n, err := clientConn.Read(buffer)
		if err != nil {
			logger.LogError("Error reading from client: " + err.Error())
			return
		}
		if n == cap(buffer) {
			logger.LogWarn("The length of data read from client is " + strconv.Itoa(n) + ", which has reached the capacity of the client read buffer. " +
				"The request data may have been fragmented, which could result in the MockedReqDataResData not matching.")
		}
		request := buffer[:n]
		logger.LogBytes("Read new client data, length: "+strconv.Itoa(n), request, m.Settings.PrintDetails)

		//handle MockedReqDataResData
		var response []byte
		for _, item := range resDataSet {
			rr := item.(*ReqDataResData)
			if bytes.Equal(*rr.ReqData, request) {
				response = *rr.ResData
				break
			}
		}
		if response != nil {
			logger.Log("Response is found in MockedReqDataResData!")
			clientConn.Write(response)
			continue
		}

		//handle MockedReqDataResDataFiles
		var fileUris []string
		for _, item := range resDataFiles {
			rr := item.(*ReqDataResDataFiles)
			if bytes.Equal(*rr.ReqData, request) {
				fileUris = rr.FileUris
				break
			}
		}
		if fileUris != nil {
			logger.Log("Response is found in MockedReqDataResDataFiles!")
			for _, fileUri := range fileUris {
				partResponse := HexFileToBytes(fileUri)
				clientConn.Write(partResponse)
			}
			continue
		}

		logger.Log("Response isn't found in MockedReqDataResData, try to send request to real server")
		_, err = serverConn.Write(request)
		if err != nil {
			logger.LogError("Error sending request to real server: " + err.Error())
			return
		}
	}
}

func (m *Mocker) handleServerSocket(clientConn net.Conn, serverConn net.Conn) {
	defer serverConn.Close()

	logger.Log("Listening on server socket...")
	postElements := m.MockedResLenResData.Elements()
	buffer := make([]byte, 512*1024)
	for {
		n, err := serverConn.Read(buffer)
		if err != nil {
			logger.LogError("Error reading from server: " + err.Error())
			return
		}
		if n == cap(buffer) {
			logger.LogWarn("The length of data read from server is " + strconv.Itoa(n) + ", which has reached the capacity of the server read buffer. " +
				"The response data may have been fragmented, which could result in the MockedResLenResData not matching.")
		}
		response := buffer[:n]
		logger.LogBytes("Read new server data, length: "+strconv.Itoa(n), response, m.Settings.PrintDetails)

		for _, item := range postElements {
			rr := item.(*ResLenResData)
			if rr.ResLen == len(response) {
				logger.Log("Response is matched in MockedResLenResData!")
				response = *rr.ResData
				break
			}
		}
		clientConn.Write(response)
	}
}
