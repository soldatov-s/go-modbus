// server
package main

import (
	"encoding/hex"
	"fmt"
	"net"
)

type ModbusData struct {
	holding_reg_cnt, input_reg_cnt int
	holding_reg, input_reg         []byte
}

func (md *ModbusData) init() error {
	md.holding_reg = make([]byte, md.holding_reg_cnt*2)
	md.input_reg = make([]byte, md.input_reg_cnt*2)

	return nil
}

func (md *ModbusData) SetHoldRegs(addr, cnt uint16, data []byte) error {
	copy(md.holding_reg[addr*2:(addr+cnt)*2], data)

	return nil
}

func (md *ModbusData) ReadHoldRegs(addr, cnt uint16) []byte {
	return md.holding_reg[addr*2 : (addr+cnt)*2]
}

func (app *ModbusApp) ServerStart() error {
	var (
		md *ModbusData = &ModbusData{
			holding_reg_cnt: *holding_reg_cnt,
			input_reg_cnt:   *input_reg_cnt}
		ln  net.Listener
		err error
	)

	fmt.Println("Server mode")
	fmt.Println("Listening at " + app.address())

	md.init()

	// Listen for incoming connections.
	ln, err = net.Listen(app.protocol, app.address())
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}

	defer ln.Close()

	for {
		// Listen for an incoming connection.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			conn.Close()
			continue
		}
		fmt.Printf(
			"Src->: \t\t\t\t%s\nDst<-: \t\t\t\t%s\n",
			conn.RemoteAddr().String(),
			conn.LocalAddr().String())
		// Handle connections in a new goroutine.
		go app.handleRequest(conn, md)
	}

}

// Handles incoming requests.
func (app *ModbusApp) handleRequest(conn net.Conn, md *ModbusData) error {
	var (
		id_packet int
		err       error
		mp        *ModbusPacket = &ModbusPacket{
			mtp: app.mbprotocol}
	)
	// Close the connection when you're done with it.
	defer conn.Close()
	mp.Init()

	// Read the incoming connection into the buffer.
	for {
		mp.length, err = conn.Read(mp.data)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		id_packet++
		fmt.Println("Packet ID:", id_packet)
		var (
			answer []byte
		)

		if mp.length > 0 {
			mp.ModbusDumper()
			switch mp.GetFC() {
			case ReadHoldingRegisters:
				answer = mp.ReadHoldRegs(md)
			case PresetMultipleRegisters:
				answer = mp.PresetMultipleRegs(md)
			default:

			}
			fmt.Println(hex.Dump(answer))
			conn.Write(answer)
		}
	}
	return err
}
