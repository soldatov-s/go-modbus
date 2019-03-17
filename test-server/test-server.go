// main
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/soldatov-s/go-modbus/modbus"
)

var (
	protocol   = flag.String("protocol", "tcp", "type of protocol, tcp/udp")
	port       = flag.String("port", "502", "port number")
	host       = flag.String("host", "localhost", "hostname or host ip")
	mbprotocol = flag.String("mbprotocol", "ModbusRTUviaTCP", "type of modbus protocol: ModbusTCP or ModbusRTUviaTCP")

	coils_cnt           = flag.Int("coils_cnt", 100, "coils counter")
	discrete_inputs_cnt = flag.Int("discrete_inputs_cnt", 100, "discrete inputs counter")
	holding_reg_cnt     = flag.Int("holding_reg_cnt", 100, "holding register counter")
	input_reg_cnt       = flag.Int("input_reg_cnt", 100, "input register counter")
)

func getFireSignalsChannel() chan os.Signal {

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP,  // "terminal is disconnected"
	)
	return c

}

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

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(s *modbus.ModbusServer) {
		<-c
		fmt.Println("\nExit program")
		s.ServerStop()
		os.Exit(1)
	}(srv)

	srv.Data.PresetMultipleRegisters(0, 5, []uint16{0x01, 0x02, 0x03, 0x04, 0x05})
	err = srv.ServerStart()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
