package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"myutil"
	"myutil/file"
	"net"
	"path"
	"strconv"
	"time"
)

type MockerConn struct {
	ConnId     int
	ClientConn net.Conn
	ServerConn net.Conn
	isClosed   bool
}

func (m *MockerConn) Close() {
	if !m.isClosed {
		m.ClientConn.Close()
		m.ServerConn.Close()
		m.isClosed = true
	}
}

type ReqDataResFiles struct {
	ReqData  *[]byte
	FileUris []string
}

type ResLenResFiles struct {
	ResLen   int
	FileUris []string
}

type Mocker struct {
	ReqDataResFiles *myutil.Set
	ResLenResFiles  *myutil.Set
	MockerConfig    *MockerConfig
}

func NewMocker(config *MockerConfig) *Mocker {
	reqDataResFiles := myutil.NewSet()
	for _, mockData := range config.MockDataGroup1 {
		var reqData []byte
		func() {
			fileUri := path.Join(config.MockDataLocation, mockData.RequestFile)
			hexFile, err := file.NewHexFile(fileUri)
			defer hexFile.Close()
			if err != nil {
				Logger.LogError("Main", "Error opening file: "+fileUri+", error: "+err.Error())
				panic(err)
			}
			reqData, err = hexFile.ReadAll()
			if err != nil {
				Logger.LogError("Main", "Error reading file: "+fileUri+", error: "+err.Error())
				panic(err)
			}
		}()

		fileUris := make([]string, len(mockData.ResponseFiles))
		for i, fileName := range mockData.ResponseFiles {
			fileUris[i] = path.Join(config.MockDataLocation, fileName)
		}
		reqDataResFiles.Add(&ReqDataResFiles{
			ReqData:  &reqData,
			FileUris: fileUris,
		})
	}

	resLenResFiles := myutil.NewSet()
	for _, mockData := range config.MockDataGroup2 {
		fileUris := make([]string, len(mockData.ResponseFiles))
		for i, fileName := range mockData.ResponseFiles {
			fileUris[i] = path.Join(config.MockDataLocation, fileName)
		}
		resLenResFiles.Add(&ResLenResFiles{
			ResLen:   mockData.ResponseDataLength,
			FileUris: fileUris,
		})
	}

	ip, err := GetIPByDomain(config.ServerIP)
	if err != nil {
		Logger.LogError("Main", "Failed to get IP, domain:"+config.ServerIP+", error: "+err.Error())
		panic(err)
	}
	if ip == nil {
		err = errors.New("There is no IP for domain, name: " + config.ServerIP)
		Logger.LogError("Main", err.Error())
		panic(err)
	}
	config.ServerIP = ip.String()

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
		if m.MockerConfig.CreateNetworkInterfaceManual {
			err = CreateInterfaceManual(m.MockerConfig.ServerIP)
		} else {
			err = CreateInterface(m.MockerConfig.ServerIP)
		}
		if err != nil {
			Logger.LogError("Main", "Failed to create Network Interface, error: "+err.Error())
			panic(err)
		}
		listener, err = net.Listen("tcp", m.MockerConfig.ServerIP+":"+strconv.Itoa(m.MockerConfig.MockerPort))
	} else {
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(m.MockerConfig.MockerPort))
	}
	if err != nil {
		Logger.LogError("Main", "Error listening, error: "+err.Error())
		panic(err)
	}
	Logger.Log("Main", "Listening TCP connections on "+listener.Addr().String()+"...")

	connId := 0
	for {
		clientConn, err := listener.Accept()
		connId++
		if err != nil {
			Logger.LogError("Main", "Error accepting connection, error :"+err.Error())
			panic(err)
		}
		Logger.Log("Main", fmt.Sprintf("Client socket is established! Client: %v -> Local: %v", clientConn.RemoteAddr(), clientConn.LocalAddr()))

		var serverConn net.Conn
		if m.MockerConfig.TunnelMode {
			serverConn, err = m.connectServerByLocalInterface()
		} else {
			serverConn, err = m.connectServer()
		}
		if err != nil {
			Logger.LogError("Main", "Error connecting server, error: "+err.Error())
			panic(err)
		}
		Logger.Log("Main", fmt.Sprintf("Server socket is established! Local: %v -> Server: %v", serverConn.LocalAddr(), serverConn.RemoteAddr()))

		mockerConn := MockerConn{
			ConnId:     connId,
			ClientConn: clientConn,
			ServerConn: serverConn,
		}
		go m.handleClientSocket(&mockerConn)
		go m.handleServerSocket(&mockerConn)
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
	_, err := GetInterfaceByIP(ip)
	if err != nil {
		Logger.LogError("Main", "Cannot find local network interface, ip address: "+ip+", error: "+err.Error())
		return nil, err
	}

	localAddr := &net.TCPAddr{IP: net.ParseIP(ip), Port: 0}
	remoteAddr := &net.TCPAddr{IP: net.ParseIP(m.MockerConfig.ServerIP), Port: m.MockerConfig.ServerPort}
	dialer := net.Dialer{LocalAddr: localAddr}
	conn, err := dialer.Dial("tcp", remoteAddr.String())
	if err != nil {
		Logger.LogError("Main", "Failed to connect to server through local network interface, error: "+err.Error())
		return nil, err
	}

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(60 * time.Second)
	return conn, nil
}

func (m *Mocker) handleClientSocket(mockerConn *MockerConn) {
	defer mockerConn.Close()

	logCode := Logger.GetCode(mockerConn.ConnId, "C")
	Logger.Log(logCode, "Listening on client socket...")
	resDataFiles := m.ReqDataResFiles.Elements()
	buffer := make([]byte, 64*1024)
	for {
		n, err := mockerConn.ClientConn.Read(buffer)
		if err != nil {
			Logger.LogError(logCode, "Error reading from client: "+err.Error())
			return
		}
		if n == cap(buffer) {
			Logger.LogWarn(logCode, "The length of data read from client is "+strconv.Itoa(n)+", which has reached the capacity of the client read buffer. "+
				"The request data may have been fragmented, which could result in the MockedReqDataResData not matching.")
		}
		request := buffer[:n]
		Logger.LogBytes(logCode, "Read new client data, length: "+strconv.Itoa(n), request, m.MockerConfig.PrintDetails)

		//handle ReqDataResFiles
		var fileUris []string
		for _, item := range resDataFiles {
			rr := item.(*ReqDataResFiles)
			if bytes.Equal(*rr.ReqData, request) {
				Logger.Log(logCode, "Response is found in ReqDataResFiles!")
				fileUris = rr.FileUris
				break
			}
		}
		if fileUris != nil {
			for _, fileUri := range fileUris {
				err := writeFileToClient(fileUri, mockerConn, logCode)
				if err != nil {
					return
				}
			}
			continue
		}

		Logger.Log(logCode, "Response isn't found in ReqDataResFiles, try to send request to server")
		_, err = mockerConn.ServerConn.Write(request)
		if err != nil {
			Logger.LogError(logCode, "Error sending request to server: "+err.Error())
			return
		}
	}
}

func (m *Mocker) handleServerSocket(mockerConn *MockerConn) {
	defer mockerConn.Close()

	logCode := Logger.GetCode(mockerConn.ConnId, "S")
	Logger.Log(logCode, "Listening on server socket...")
	postElements := m.ResLenResFiles.Elements()
	buffer := make([]byte, 128*1024)
	for {
		n, err := mockerConn.ServerConn.Read(buffer)
		if err != nil {
			Logger.LogError(logCode, "Error reading from server: "+err.Error())
			return
		}
		if n == cap(buffer) {
			Logger.LogWarn(logCode, "The length of data read from server is "+strconv.Itoa(n)+", which has reached the capacity of the server read buffer. "+
				"The response data may have been fragmented, which could result in the ResLenResFiles not matching.")
		}
		response := buffer[:n]
		Logger.LogBytes(logCode, "Read new server data, length: "+strconv.Itoa(n), response, m.MockerConfig.PrintDetails)

		//handle ResLenResFiles
		var fileUris []string
		for _, item := range postElements {
			rr := item.(*ResLenResFiles)
			if rr.ResLen == len(response) {
				Logger.Log(logCode, "Response is matched in ResLenResFiles!")
				fileUris = rr.FileUris
				break
			}
		}
		if fileUris != nil {
			for _, fileUri := range fileUris {
				err := writeFileToClient(fileUri, mockerConn, logCode)
				if err != nil {
					return
				}
			}
			continue
		}

		_, err = mockerConn.ClientConn.Write(response)
		if err != nil {
			Logger.LogError(logCode, "Error sending response to client: "+err.Error())
			return
		}
	}
}

func writeFileToClient(fileUri string, mockerConn *MockerConn, logCode string) error {
	hexFile, err := file.NewHexFile(fileUri)
	defer hexFile.Close()
	if err != nil {
		Logger.LogError(logCode, "Error opening file: "+fileUri+", error: "+err.Error())
		return err
	}
	buf := make([]byte, 64*1024)
	for {
		isEOF := false
		n, err := hexFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				isEOF = true
			} else {
				Logger.LogError(logCode, "Error reading file: "+fileUri+", error: "+err.Error())
				return err
			}
		}
		_, err = mockerConn.ClientConn.Write(buf[:n])
		if err != nil {
			Logger.LogError(logCode, "Error sending file to client: "+err.Error())
			return err
		}
		if isEOF {
			return nil
		}
	}

}
