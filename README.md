[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/soldatov-s/go-modbus)
[![Build Status](https://travis-ci.org/soldatov-s/go-modbus.svg?branch=master)](https://travis-ci.org/soldatov-s/go-modbus)
[![Coverage Status](http://codecov.io/github/soldatov-s/go-modbus/coverage.svg?branch=master)](http://codecov.io/github/soldatov-s/go-modbus?branch=master)
[![codebeat badge](https://codebeat.co/badges/b671ecf0-3e82-4e48-b220-e369d0ced46c)](https://codebeat.co/projects/github-com-soldatov-s-go-modbus-master)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
# go-modbus

## About
Modbus protocol framework

## Supported
 * Modbus RTU over TCP
 * Modbus Slave mode (Modbus Server)
 * Rest server for read/write Modbus Data
 * gRPC service
 * Dump Modbus packets
 * Function:  
 - Read Coil Status (0x1)
 - Read Discrete Inputs (0x2)
 - Read Holding Registers (0x3)
 - Read Input Registers (0x4)
 - Force Single Coil (0x5)
 - Preset Single Register (0x6)
 - Force Multiple Coils (0xF)
 - Preset Multiple Registers (0x10)

## Installation
```sh
go get github.com/soldatov-s/go-modbus
```
Next, build and run the examples:

 * [mb-server.go](mb-server/mb-server.go) for an Modbus RTU over TCP server example

## Rest Server
 - /coils - Coils (GET and PUT)
 - /d_in - Discrete Inputs (only GET)
 - /hold_reg - Holding Registers (GET and PUT)
 - /in_reg - Input Registers (only GET)

Example Read Holding Registers:
```sh
curl -X GET -i 'http://localhost:8000/hold_reg?addr=0&cnt=5'
```

Example Write Holding Registers:
```sh
curl -X POST -i http://localhost:8000/hold_reg --data '{
"addr": 0,
"data": [77, 11]
}'
```

## More Documentation

More documentation about Modbus is available on the
- [Wiki](https://en.wikipedia.org/wiki/Modbus).
- [Modbus Technical Specifications](http://www.modbus.org/specs.php)
