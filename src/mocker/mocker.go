package main

import (
	"bytes"
	"fmt"
	"github.com/edward/mocker/logger"
	"myutil"
	"net"
	"path"
	"strconv"
	"time"
)

type Mocker struct {
	ReqDataResFiles *myutil.Set
	ResLenResFiles  *myutil.Set
	MockerConfig    MockerConfig
}

func NewMocker(config MockerConfig) *Mocker {
	reqDataResFiles := myutil.NewSet()
	for _, mockData := range config.MockDataGroup1 {
		reqData := HexFileToBytes(path.Join(config.MockDataLocation, mockData.RequestFile))
		reqDataResFiles.Add(&ReqDataResFiles{
			ReqData:  &reqData,
			FileUris: mockData.ResponseFiles,
		})
	}

	resLenResFiles := myutil.NewSet()
	for _, mockData := range config.MockDataGroup2 {
		resLenResFiles.Add(&ResLenResFiles{
			ResLen:   mockData.ResponseDataLength,
			FileUris: mockData.ResponseFiles,
		})
	}

	return &Mocker{
		ReqDataResFiles: reqDataResFiles,
		ResLenResFiles:  resLenResFiles,
		MockerConfig:    config,
	}
}

func (m *Mocker) Start() {
	var listener net.Listener
	var err error
	if m.MockerConfig.TunnelMode {
		_, err = CreateNetworkInterface(m.MockerConfig.ServerIP)
		if err != nil {
			logger.LogError("Failed to create Network Interface, error: " + err.Error())
			panic(err)
		}
		listener, err = net.Listen("tcp", m.MockerConfig.ServerIP+":"+strconv.Itoa(m.MockerConfig.MockerPort))
	} else {
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(m.MockerConfig.MockerPort))
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
		if m.MockerConfig.TunnelMode {
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
	serverAddr := m.MockerConfig.ServerIP + ":" + strconv.Itoa(m.MockerConfig.ServerPort)
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
	ip := m.MockerConfig.LocalNetworkInterfaceAddress
	_, err := GetNetworkInterface(ip)
	if err != nil {
		logger.LogError("Cannot find local network interface, ip address: "+ip, err)
		return nil, err
	}

	localAddr := &net.TCPAddr{IP: net.ParseIP(ip), Port: 0}
	remoteAddr := &net.TCPAddr{IP: net.ParseIP(m.MockerConfig.ServerIP), Port: m.MockerConfig.ServerPort}
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

	logger.Log("[C]Listening on client socket...")
	resDataFiles := m.ReqDataResFiles.Elements()
	buffer := make([]byte, 64*1024)
	for {
		n, err := clientConn.Read(buffer)
		if err != nil {
			logger.LogError("[C]Error reading from client: " + err.Error())
			return
		}
		if n == cap(buffer) {
			logger.LogWarn("[C]The length of data read from client is " + strconv.Itoa(n) + ", which has reached the capacity of the client read buffer. " +
				"The request data may have been fragmented, which could result in the MockedReqDataResData not matching.")
		}
		request := buffer[:n]
		logger.LogBytes("[C]Read new client data, length: "+strconv.Itoa(n), request, m.MockerConfig.PrintDetails)

		//handle ReqDataResFiles
		var fileUris []string
		for _, item := range resDataFiles {
			rr := item.(*ReqDataResFiles)
			if bytes.Equal(*rr.ReqData, request) {
				logger.Log("[C]Response is found in ReqDataResFiles!")
				fileUris = rr.FileUris
				break
			}
		}
		if fileUris != nil {
			for _, fileUri := range fileUris {
				partResponse := HexFileToBytes(path.Join(m.MockerConfig.MockDataLocation, fileUri))
				_, err = clientConn.Write(partResponse)
				if err != nil {
					logger.LogError("[C]Error sending response to client: " + err.Error())
					return
				}
			}
			continue
		}

		logger.Log("[C]Response isn't found in ReqDataResFiles, try to send request to server")
		_, err = serverConn.Write(request)
		if err != nil {
			logger.LogError("[C]Error sending request to server: " + err.Error())
			return
		}
	}
}

func (m *Mocker) handleServerSocket(clientConn net.Conn, serverConn net.Conn) {
	defer serverConn.Close()

	logger.Log("[S]Listening on server socket...")
	postElements := m.ResLenResFiles.Elements()
	buffer := make([]byte, 128*1024)
	for {
		n, err := serverConn.Read(buffer)
		if err != nil {
			logger.LogError("[S]Error reading from server: " + err.Error())
			return
		}
		if n == cap(buffer) {
			logger.LogWarn("[S]The length of data read from server is " + strconv.Itoa(n) + ", which has reached the capacity of the server read buffer. " +
				"The response data may have been fragmented, which could result in the ResLenResFiles not matching.")
		}
		response := buffer[:n]
		logger.LogBytes("[S]Read new server data, length: "+strconv.Itoa(n), response, m.MockerConfig.PrintDetails)

		//handle ResLenResFiles
		var fileUris []string
		for _, item := range postElements {
			rr := item.(*ResLenResFiles)
			if rr.ResLen == len(response) {
				logger.Log("[S]Response is matched in ResLenResFiles!")
				fileUris = rr.FileUris
				break
			}
		}
		if fileUris != nil {
			for _, fileUri := range fileUris {
				partResponse := HexFileToBytes(path.Join(m.MockerConfig.MockDataLocation, fileUri))
				_, err = clientConn.Write(partResponse)
				if err != nil {
					logger.LogError("[S]Error sending response to client: " + err.Error())
					return
				}
			}
			continue
		}

		_, err = clientConn.Write(response)
		if err != nil {
			logger.LogError("[S]Error sending response to client: " + err.Error())
			return
		}
	}
}
