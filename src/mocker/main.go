package main

import (
	"bufio"
	"fmt"
	"net"
)

// RequestResponse stores the raw request data and response data for a request
type RequestResponse struct {
	Request  []byte
	Response []byte
}

// Mocker stores the map of RequestResponse objects and the TCP listener
type Mocker struct {
	Requests map[uint64]RequestResponse
	Listener net.Listener
}

// NewMocker creates a new Mocker with an empty request map
func NewMocker() (*Mocker, error) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	return &Mocker{
		Requests: make(map[uint64]RequestResponse),
		Listener: ln,
	}, nil
}

// GenerateRequestResponse generates a new RequestResponse object and stores it in the request map
func (m *Mocker) GenerateRequestResponse(request []byte, response []byte) {
	hash := hash(request)

	m.Requests[hash] = RequestResponse{
		Request:  request,
		Response: response,
	}
}

// Start starts the TCP server and listens for incoming connections
func (m *Mocker) Start() {
	for {
		conn, err := m.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go m.handleConnection(conn)
	}
}

// handleConnection handles an incoming TCP connection by reading the request and sending the response
func (m *Mocker) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		request, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}

		hash := hash(request)

		if rr, ok := m.Requests[hash]; ok {
			conn.Write(rr.Response)
		} else {
			fmt.Println("Error: no response found for request:", string(request))
			conn.Write([]byte("Error: no response found"))
		}
	}
}

// hash generates a hash code for a byte slice using the fnv64 algorithm
func hash(b []byte) uint64 {
	var h uint64 = 14695981039346656037
	var fnv64prime uint64 = 1099511628211

	for _, c := range b {
		h = h*fnv64prime + uint64(c)
	}

	return h
}

func main() {
	m, err := NewMocker()
	if err != nil {
		fmt.Println("Error creating Mocker:", err)
		return
	}

	m.GenerateRequestResponse([]byte("abc\n"), []byte("def\n"))

	fmt.Println("Mocker started on port 8080")
	m.Start()
}
