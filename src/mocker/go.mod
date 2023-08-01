module github.com/edward/mocker

go 1.20

require myutil v0.0.0

require (
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20230704135630-469159ecf7d1 // indirect
)

replace myutil => ../myutil
