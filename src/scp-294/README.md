### Build
- go build -o bin/scp-294.exe cmd/main.go

### Windows service:
- kill task: taskkill /pid xxx /F
- create service: sc create scp-294 binPath= D:\workspace\learn-go\src\scp-294\scp-294.exe start= auto
- delete service: sc delete scp-294

### linux install go:
- mkdir ~/go && cd ~/go
- wget https://go.dev/dl/go1.18.4.linux-amd64.tar.gz
- tar -C /usr/local -zxvf  go1.18.4.linux-amd64.tar.gz
- vi ~/.bashrc
- export GOROOT=/usr/local/go
- export PATH=$PATH:$GOROOT/bin
- source ~/.bashrc
- go version
