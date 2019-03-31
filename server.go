// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

// Package modbus provides a server for MODBUS RTU over TCP.

package modbus

import (
	"errors"
	"log"
	"math"
	"net"
	"strings"
	"sync"
	"time"
)

// ModbusServer implements server interface
type ModbusServer struct {
	ModbusBaseServer                    // Anonim ModbusBase implementation
	TypeProtocol     ModbusTypeProtocol // Type of Modbus protocol: TCP or RTU over TCP
	ln               net.Listener       // Listener
	done             chan struct{}      // Chan for sending "done" command
	exited           chan struct{}      // Chan for sending to main app signal that server is fully stopped
	wg               sync.WaitGroup     // WaitGroup for waiting end all connection
}

// NewServer function initializate new instance of ModbusServer
func NewServer(host, port string, mbprotocol ModbusTypeProtocol, md *ModbusData) *ModbusServer {
	srv := &ModbusServer{
		TypeProtocol: mbprotocol,
		done:         make(chan struct{}),
		exited:       make(chan struct{})}
	srv.Port = port
	srv.Host = host
	srv.Data = md
	return srv
}

// Stop function close listener and wait closing all connection
func (srv *ModbusServer) Stop() error {
	log.Println("Shutting down server...")
	if srv.ln != nil {
		srv.ln.Close()
	}
	close(srv.done)
	srv.wg.Wait()
	log.Println("Server is stopped")
	return nil
}

// Start function begin listen incoming connection
func (srv *ModbusServer) Start() error {
	var err error

	log.Println("Server startup...")
	log.Println("Listening at", srv)

	// Listen for incoming connections.
	srv.ln, err = net.Listen("tcp", srv.String())
	if err != nil {
		log.Println("Error listening:", err.Error())
		return err
	}
	srv.wg.Add(1)
	go func() {
		defer func() {
			log.Println("Stop listen incoming connection")
			srv.ln.Close()
			srv.wg.Done()
		}()
		for {
			select {
			case <-srv.done:
				return
			default:
				// Listen for an incoming connection.
				conn, err := srv.ln.Accept()
				if err != nil {
					if strings.Contains(err.Error(), "use of closed network connection") {
						return
					}
					log.Println("Error accepting: ", err.Error())
					continue
				}
				// Handle connections in a new goroutine.
				go srv.handleRequest(conn)
			}
		}
	}()

	log.Println("Server is started")
	return nil
}

// Handles incoming requests.
func (srv *ModbusServer) handleRequest(conn net.Conn) error {
	var (
		id_packet int
		err       error
	)
	// Close the connection when you're done with it.
	defer func() {
		log.Println("Close connection", conn.RemoteAddr())
		srv.wg.Done()
		conn.Close()
	}()
	srv.wg.Add(1)
	request := &ModbusPacket{}
	request.Init(srv.TypeProtocol)

	log.Printf(
		"Src->: %s Dst<-: %s\n",
		conn.RemoteAddr(),
		conn.LocalAddr())

	// Read the incoming connection into the buffer.
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		select {
		case <-srv.done:
			return nil
		default:
			request.Length, err = conn.Read(request.PDU)
			if err != nil {
				log.Println("Error reading:", err.Error())
				return err
			}

			if request.Length == 0 {
				continue
			}

			if id_packet++; id_packet == math.MaxInt32 {
				id_packet = 0
			}
			log.Printf("Src->: %s, Packet ID:%d\n", conn.RemoteAddr(), id_packet)
			request.Dump("****Request Dump****")
			var answer *ModbusPacket
			answer, err = srv.RequestHadler(request)
			if err != nil {
				log.Println("Error handle request:", err.Error())
				break
			}
			answer.Dump("****Answer Dump****")
			conn.Write(answer.PDU[:answer.GetPDULength()])
		}
	}
	return err
}

func (srv *ModbusServer) RequestHadler(mp *ModbusPacket) (*ModbusPacket, error) {
	switch mp.GetFunctionCode() {
	case FcReadCoilStatus:
		return srv.ReadCoilStatus(mp)
	case FcReadDescreteInputs:
		return srv.ReadDescreteInputs(mp)
	case FcForceSingleCoil:
		return srv.ForceSingleCoil(mp)
	case FcPresetSingleRegister:
		return srv.PresetSingleRegister(mp)
	case FcReadHoldingRegisters:
		return srv.ReadHoldingRegisters(mp)
	case FcReadInputRegisters:
		return srv.ReadInputRegisters(mp)
	case FcForceMultipleCoils:
		return srv.ForceMultipleCoils(mp)
	case FcPresetMultipleRegisters:
		return srv.PresetMultipleRegisters(mp)
	default:
		return buildErrAnswer(mp, ErrCantHandel), errors.New("Unknown function code")
	}
}

// Read Holding registers
func (srv *ModbusServer) ReadHoldingRegisters(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, cnt := mp.GetFunctionParameters()
	// Try get data for answer
	data, err := srv.Data.ReadHoldingRegisters(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, wordArrToByteArr(data)...), nil
}

// Read Inputs registers
func (srv *ModbusServer) ReadInputRegisters(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, cnt := mp.GetFunctionParameters()
	// Try get data for answer
	data, err := srv.Data.ReadInputRegisters(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, wordArrToByteArr(data)...), nil
}

// Preset Single Register
func (srv *ModbusServer) PresetSingleRegister(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, value := mp.GetFunctionParameters()
	// Set values in ModbusData
	err := srv.Data.PresetSingleRegister(addr, value)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp), nil
}

// Preset Multiple Holding Registers
func (srv *ModbusServer) PresetMultipleRegisters(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, _ := mp.GetFunctionParameters()
	_, data := mp.GetData()
	// Set values in ModbusData
	err := srv.Data.PresetMultipleRegisters(addr, byteArrToWordArr(data)...)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp), nil
}

// Read Coil Status
func (srv *ModbusServer) ReadCoilStatus(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, cnt := mp.GetFunctionParameters()
	// Data for answer
	data, err := srv.Data.ReadCoilStatus(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, boolArrToByteArr(data)...), nil
}

// Read Descrete Inputs
func (srv *ModbusServer) ReadDescreteInputs(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, cnt := mp.GetFunctionParameters()
	// Data for answer
	data, err := srv.Data.ReadDescreteInputs(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, boolArrToByteArr(data)...), nil
}

// Force Single Coil
func (srv *ModbusServer) ForceSingleCoil(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, value := mp.GetFunctionParameters()
	if value == 0xFF00 {
		value = 1
	}
	// Set values in ModbusData
	err := srv.Data.ForceSingleCoil(addr, bool((value&1) == 1))
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp), nil
}

// Force Multiple Coils
func (srv *ModbusServer) ForceMultipleCoils(mp *ModbusPacket) (*ModbusPacket, error) {
	addr, cnt := mp.GetFunctionParameters()
	_, data := mp.GetData()
	// Set values in ModbusData)
	err := srv.Data.ForceMultipleCoils(addr, byteArrToBoolArr(data, byte(cnt))...)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp), nil
}
