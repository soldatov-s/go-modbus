// server
package main

import (
	"fmt"
	"math"
	"net"
)

func (app *ModbusApp) ServerStart() error {
	var (
		md *ModbusData = &ModbusData{
			holding_reg_cnt: *holding_reg_cnt,
			input_reg_cnt:   *input_reg_cnt}
		ln   net.Listener
		err  error
		done = make(chan struct{})
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

	defer func() {
		ln.Close()
		close(done)
	}()

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
		go app.handleRequest(conn, md, done)
	}
}

// Handles incoming requests.
func (app *ModbusApp) handleRequest(conn net.Conn, md *ModbusData, done <-chan struct{}) error {
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
		select {
		case <-done:
			conn.Close()
			return nil
		default:
			mp.length, err = conn.Read(mp.data)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				break
			}
			id_packet++
			if id_packet == math.MaxInt32 {
				id_packet = 0
			}
			fmt.Println("Packet ID:", id_packet)
			var (
				answer []byte
			)

			if mp.length > 0 {
				mp.ModbusDumper()
				switch mp.GetFC() {
				case ReadHoldingRegisters:
					answer = md.ReadHoldRegs(mp)
				case PresetMultipleRegisters:
					answer = md.PresetMultipleRegs(mp)
				default:

				}
				// fmt.Println(hex.Dump(answer))
				conn.Write(answer)
			}
		}
	}
	return err
}
