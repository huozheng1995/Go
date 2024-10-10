### To use Mocker, follow the steps bellow:

- Edit 'config.json', see the "Common Config"
- Create request and response data files, see "Format for RequestFile and ResponseFiles"
- Run mocker.exe as administrator
- Run CData driver

### Common Config
- ServerIP: The server IP to connect
- ServerPort: The server port to connect
- MockerIP: Mocker will use this IP to create a local network interface and listen on it
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
  - ResponseDataLength: Mocker will send the raw request to server, then match the length of the response
  - ResponseFiles: Same as itself in MockDataGroup1
  ```
  For example:
  "MockDataGroup2": [
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


### Format for RequestFile and ResponseFiles
- The format is hexadecimal bytes, the same as CData driver log file, such as:
   ```
   12 01 01 13 00 00 02 00 16 03 03 01 06 01 00 01 
   02 03 03 64 DF 13 EF BB D7 EC 0F B6 83 36 8C F2 
   1B 46 BB 9C D0 2D 3C A5 9F 19 85 7D 77 CF E9 AF 
   2D D5 EF 00 00 74 C0 2C C0 2B C0 2F C0 30 C0 2E 
   00 9D 00 9C C0 2D 00 
   ```
- The following regex will help you extract hexadecimal bytes from the log:
  - Hexadecimal byte data lines in CData logs:
    ^([0-9A-F][0-9A-F] ){1, 16}.*\r\n 
  - Non-hex byte data lines in CData logs:
    ^(?!([0-9A-F][0-9A-F] ){1, 16}).*\r\n
