// main
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	mode = flag.Bool("mode", false, "client or server mode")
	// Common options
	protocol   = flag.String("protocol", "tcp", "type of protocol, tcp/udp")
	port       = flag.String("port", "502", "port number")
	host       = flag.String("host", "localhost", "hostname or host ip")
	mbprotocol = flag.String("mbprotocol", "ModbusRTUviaTCP", "type of modbus protocol: ModbusTCP or ModbusRTUviaTCP")
	// For server
	holding_reg_cnt = flag.Int("holding_reg_cnt", 100, "holding register counter")
	input_reg_cnt   = flag.Int("input_reg_cnt", 100, "input register counter")
	// For client
	fcode      = flag.String("fcode", "0x03", "function code")
	slave_addr = flag.String("slave", "1", "slave address")
	data       = flag.String("data", "0000000A", "data for send")
)

type ModbusApp struct {
	host       string
	protocol   string
	port       string
	mbprotocol ModbusTypeProtocol
}

func (app *ModbusApp) address() string {
	return app.host + ":" + app.port
}

func main() {
	var (
		app *ModbusApp
		err error
	)
	fmt.Println("Modbus app!")

	flag.Parse()

	app = &ModbusApp{
		host:       *host,
		protocol:   *protocol,
		port:       *port,
		mbprotocol: StringToModbusTypeProtocol(*mbprotocol)}

	if *mode {
		//app.ClientStart()
	} else {
		err = app.ServerStart()
	}
	if err != nil {
		os.Exit(1)
	}
}
