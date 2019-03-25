// Example of Modbus Master device (Client)
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/soldatov-s/go-modbus"
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
	var (
		err    error
		cl     *modbus.ModbusClient
		answer *modbus.ModbusPacket
	)
	fmt.Println("Modbus client app!")
	flag.Parse()

	cl, err = modbus.NewClient(*host, *port,
		modbus.StringToModbusTypeProtocol(*mbprotocol))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	request := &modbus.ModbusPacket{
		Data:   []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5},
		Length: 8,
		MTP:    cl.MTP}

	request.ModbusDump()

	answer, err = cl.SendRequest(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	answer.ModbusDump()

	cl.Close()
}
