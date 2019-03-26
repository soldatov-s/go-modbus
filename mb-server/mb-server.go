// Example of Modbus Slave device (Server)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/soldatov-s/go-modbus"
	"github.com/soldatov-s/go-modbus/modbusrest"
)

var (
	port                = flag.String("port", "502", "port number")
	host                = flag.String("host", "localhost", "hostname or host ip")
	rest_port           = flag.String("rest_port", "8000", "port number")
	rest_host           = flag.String("rest_host", "localhost", "hostname or host ip")
	mbprotocol          = flag.String("mbprotocol", "ModbusRTUviaTCP", "type of modbus protocol: ModbusTCP or ModbusRTUviaTCP")
	coils_cnt           = flag.Int("coils_cnt", 65535, "coils counter")
	discrete_inputs_cnt = flag.Int("discrete_inputs_cnt", 65535, "discrete inputs counter")
	holding_reg_cnt     = flag.Int("holding_reg_cnt", 65535, "holding register counter")
	input_reg_cnt       = flag.Int("input_reg_cnt", 65535, "input register counter")
)

func main() {
	var err error

	log.Println("Modbus server app!")
	flag.Parse()

	md := new(modbus.ModbusData)
	md.Init(*coils_cnt, *discrete_inputs_cnt, *holding_reg_cnt, *input_reg_cnt)
	md.PresetMultipleRegisters(0, []uint16{0x01, 0x02, 0x03, 0x04, 0x05})

	srv := modbus.NewServer(*host, *port,
		modbus.StringToModbusTypeProtocol(*mbprotocol), md)

	rest := modbusrest.NewRest(*rest_host, *rest_port, md)

	servers := [2]modbus.IModbusBaseServer{srv, rest}
	// Exit handler
	exit := make(chan struct{})
	closeSignal := make(chan os.Signal)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeSignal
		for _, s := range servers {
			err = s.Stop()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		fmt.Println("Exit program")
		close(exit)
	}()

	for _, s := range servers {
		err = s.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Exit app if chan is closed
	<-exit
}
