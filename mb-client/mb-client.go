// test-client
package main

import (
	"flag"
	"fmt"
)

var (
	protocol   = flag.String("protocol", "tcp", "type of protocol, tcp/udp")
	port       = flag.String("port", "502", "port number")
	host       = flag.String("host", "localhost", "hostname or host ip")
	mbprotocol = flag.String("mbprotocol", "ModbusRTUviaTCP", "type of modbus protocol: ModbusTCP or ModbusRTUviaTCP")

	fcode      = flag.String("fcode", "0x03", "function code")
	slave_addr = flag.String("slave", "1", "slave address")
	data       = flag.String("data", "0000000A", "data for send")
)

func main() {
	fmt.Println("Hello World!")
}
