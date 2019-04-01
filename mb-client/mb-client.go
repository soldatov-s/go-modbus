// Example of Modbus Master device (Client)
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/soldatov-s/go-modbus"
	"github.com/soldatov-s/go-modbus/modbusgrpc"
)

var (
	protocol   = flag.String("protocol", "tcp", "type of protocol, tcp/udp")
	port       = flag.String("port", "502", "port number")
	grpc_port  = flag.String("grpc_port", "9000", "port number")
	host       = flag.String("host", "localhost", "hostname or host ip")
	mbprotocol = flag.String("mbprotocol", "ModbusRTUviaTCP", "type of modbus protocol: ModbusTCP or ModbusRTUviaTCP")
)

func main() {
	flag.Parse()
	fmt.Println("Modbus client app!")

	cl, err := modbus.NewClient(*port, *host,
		modbus.StringToModbusTypeProtocol(*mbprotocol), 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hold_regs, err := cl.ReadHoldingRegisters(0, 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Result ", hold_regs)

	coils, err := cl.ReadCoilStatus(0, 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Result ", coils)

	cl.Close()

	// gRPC test
	fmt.Println("gRPC client app!")
	clRPC, err := modbusgrpc.NewgRPCClient(*grpc_port, *host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hold_regs_int32, err := clRPC.ReadHoldingRegisters(0, 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Result ", hold_regs_int32)

	coils, err = clRPC.ReadCoilStatus(0, 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Result ", coils)
	clRPC.Close()
}
