// main
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/soldatov-s/go-modbus"
)

var (
	protocol   = flag.String("protocol", "tcp", "type of protocol, tcp/udp")
	port       = flag.String("port", "502", "port number")
	host       = flag.String("host", "localhost", "hostname or host ip")
	mbprotocol = flag.String("mbprotocol", "ModbusRTUviaTCP", "type of modbus protocol: ModbusTCP or ModbusRTUviaTCP")

	coils_cnt           = flag.Int("coils_cnt", 9999, "coils counter")
	discrete_inputs_cnt = flag.Int("discrete_inputs_cnt", 9999, "discrete inputs counter")
	holding_reg_cnt     = flag.Int("holding_reg_cnt", 9999, "holding register counter")
	input_reg_cnt       = flag.Int("input_reg_cnt", 9999, "input register counter")
)

func main() {
	var err error

	fmt.Println("Modbus server app!")
	flag.Parse()

	srv := modbus.NewServer(
		*host,
		*protocol,
		*port,
		modbus.StringToModbusTypeProtocol(*mbprotocol),
		*coils_cnt, *discrete_inputs_cnt, *holding_reg_cnt, *input_reg_cnt)

	// Exit handler
	exit := make(chan struct{})
	closeSignal := make(chan os.Signal)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeSignal
		srv.Stop()
		fmt.Println("Exit program")
		close(exit)
	}()

	srv.Data.PresetMultipleRegisters(0, 5, []uint16{0x01, 0x02, 0x03, 0x04, 0x05})
	err = srv.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Exit app if chan is closed
	<-exit
}
