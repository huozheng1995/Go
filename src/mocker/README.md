### To use Mocker, follow the steps bellow:
#### Tunnel mode off
- Edit 'config.json', see "Common Config"
- Create request and response data files, see "Format for RequestFile and ResponseFiles"
- Modify CData driver's connection string. Server property should be "127.0.0.1", Port property should be the same as "MockerPort"
- Start mocker.exe
- Execute CData driver

#### Tunnel mode on
- Edit 'config.json', see the "Common Config" and "Tunnel Mode Config" below
- Create request and response data files, see "Format for RequestFile and ResponseFiles"
- Start mocker.exe
- Mocker will prompt you to set the correct IP address for WireGuard, see "Waiting for network interface to be 
  created... IP: xxx.xxx.xxx.xxx"
- Start WireGuard, then add empty tunnel (press Ctrl+N), input a name and add the IP address, then click "Activate"
- ```
   [Interface]
   PrivateKey = XXXXXXXXXXXXXXXXXXXXXXXXXXXX
   Address = xxx.xxx.xxx.xxx/32
   ```
- execute CData driver



### Common Config
 - ServerIP: The server ip to connect
 - ServerPort: The server port to connect
 - MockerPort: Mocker will listen on this port
 - PrintDetails: Whether to print the byte data
 - MockDataLocation: The root folder to save the mocked request and response data
 - MockDataGroup1: 
   - RequestFile: Mocker will match the exact data of the RequestFile
   - ResponseFiles: If matched, return the data of the ResponseFiles. Large response file can be split.
    ```
    For example:
    "MockDataGroup1": [
      {
        "RequestFile": "/req1.log",
        "ResponseFiles": [
          "/res1.log",
        ]
      },
      {
        "RequestFile": "/req2.log",
        "ResponseFiles": [
          "/res2-part1.log",
          "/res2-part2.log",
        ]
      }
    ],
    ```
 - MockDataGroup2: 
   - ResponseDataLength: Mocker will send the real request to server, then match the length of the response
   - ResponseFiles: Same as itself in MockDataGroup1
    ```
    For example:
    "MockDataGroup1": [
      {
        "ResponseDataLength": 123,
        "ResponseFiles": [
          "/res1.log",
        ]
      },
      {
        "ResponseDataLength": 123456,
        "ResponseFiles": [
          "/res2-part1.log",
          "/res2-part2.log",
        ]
      }
    ],
    ```



### Tunnel Mode Config
It is used to handle http-based drivers such as databricks and snowflake. In these drivers, the server property 
cannot be changed, otherwise the server will not recognize it.

#### Working Principle
It will create a virtual network interface to route the request to Mocker, so you don't need to change 
   the connection string in the driver.

#### To enable Tunnel Mode, you should configure these
 - TunnelMode: True
 - CreateNetworkInterfaceManual: Only support True for now



### Format for RequestFile and ResponseFiles
 - The format is hexadecimal bytes, the same as CData driver log file, such as:
    ```
    12 01 01 13 00 00 02 00 16 03 03 01 06 01 00 01 
    02 03 03 64 DF 13 EF BB D7 EC 0F B6 83 36 8C F2 
    1B 46 BB 9C D0 2D 3C A5 9F 19 85 7D 77 CF E9 AF 
    2D D5 EF 00 00 74 C0 2C C0 2B C0 2F C0 30 C0 2E 
    00 9D 00 9C C0 2D 00 
    ```



### [WireGuard config](https://www.procustodibus.com/blog/2021/01/wireguard-endpoints-and-ip-addresses/)
#### Basic (sufficient in most cases)

```
[Interface]
PrivateKey = aHZ0FJBG3kdNhiEhLiTjBwHwYpzdIHd9dTStJoQClGk=
Address = xxx.xxx.xxx.xxx/32
```

#### Advanced (for example: use the network proxy)

```
[Interface]
PrivateKey = aHZ0FJBG3kdNhiEhLiTjBwHwYpzdIHd9dTStJoQClGk=
Address = xxx.xxx.xxx.xxx/32

[Peer]
PublicKey = fE/wdxzl0klVp/IR8UcaoGUMjqaWi3jAd7KzHKFS6Ds=
AllowedIPs = xxx.xxx.xxx.xxx/32
Endpoint = 127.0.0.1:10809
```




