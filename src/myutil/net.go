package myutil

type Network struct {
	Name        string `json:"Name"`
	IPv4Address string `json:"IPv4Address"`
	SubnetMask  int    `json:"SubnetMask"`
}
