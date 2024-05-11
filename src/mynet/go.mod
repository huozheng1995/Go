module github.com/edward/mynet

go 1.20

require (
	golang.org/x/net v0.25.0
	golang.org/x/sys v0.20.0
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20231211153847-12269c276173
	myutil v0.0.0
)

replace myutil => ../myutil
