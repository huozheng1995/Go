## [WireGuard config](https://www.procustodibus.com/blog/2021/01/wireguard-endpoints-and-ip-addresses/)

### TCP layer protocol

```C
[Interface]
PrivateKey = aHZ0FJBG3kdNhiEhLiTjBwHwYpzdIHd9dTStJoQClGk=
Address = 172.16.85.139/32
```

### Http protocol

```C
[Interface]
PrivateKey = aHZ0FJBG3kdNhiEhLiTjBwHwYpzdIHd9dTStJoQClGk=
Address = xxx.xxx.xx.xxx/32

[Peer]
PublicKey = fE/wdxzl0klVp/IR8UcaoGUMjqaWi3jAd7KzHKFS6Ds=
AllowedIPs = xxx.xxx.xx.xxx/32
Endpoint = 127.0.0.1:10809
```




