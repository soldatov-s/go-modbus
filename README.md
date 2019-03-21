[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/soldatov-s/go-modbus)
# go-modbus

## About
Modbus protocol framework

## Supported
 * Modbus RTU over TCP
 * Modbus Slave as Server
 * Dump Modbus packets
 * Function 0x1, 0x2, 0x3, 0x4, 0xF, 0x10

## Instalation
```sh
go get github.com/gorilla/mux
go get github.com/soldatov-s/go-modbus
```
Next, build and run the examples:

 * [mb-server.go](mb-server/mb-server.go) for an Modbus RTU over TCP server example

