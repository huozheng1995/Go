### [WireGuard config](https://www.procustodibus.com/blog/2021/01/wireguard-endpoints-and-ip-addresses/)

```C
[Interface]
PrivateKey = AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEE=
Address = 10.0.0.1/32
ListenPort = 51821

[Peer]
PublicKey = fE/wdxzl0klVp/IR8UcaoGUMjqaWi3jAd7KzHKFS6Ds=
AllowedIPs = 192.168.200.0/24
Endpoint = 203.0.113.2:51822
