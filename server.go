// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

// Package modbus provides a server for MODBUS RTU over TCP.

package modbus

import (
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"time"
)

// ModbusServer implements server interface
type ModbusServer struct {
	ModbusBase                    // Anonim ModbusBase implementation
	MbProtocol ModbusTypeProtocol // Type of Modbus protocol: TCP or RTU over TCP
	ln         net.Listener       // Listener
	done       chan struct{}      // Chan for sending "done" command
	exited     chan struct{}      // Chan for sending to main app signal that server is fully stopped
	wg         sync.WaitGroup     // WaitGroup for waiting end all connection
}

// NewServer function initializate new instance of ModbusServer
func NewServer(host, port string, mbprotocol ModbusTypeProtocol, md *ModbusData) *ModbusServer {
	srv := new(ModbusServer)
	srv.Host = host
	srv.Port = port
	srv.MbProtocol = mbprotocol
	srv.done = make(chan struct{})
	srv.exited = make(chan struct{})

	return srv
}

// Return string with host ip/name and port
func (srv *ModbusServer) String() string {
	return srv.Host + ":" + srv.Port
}

// Stop function close listener and wait closing all connection
func (srv *ModbusServer) Stop() error {
	var err error

	fmt.Println("Shutting down server...")
	if srv.ln != nil {
		srv.ln.Close()
	}
	close(srv.done)
	srv.wg.Wait()
	fmt.Println("Server is stopped")
	return err
}

// Start function begin listen incoming connection
func (srv *ModbusServer) Start() error {
	var err error

	fmt.Println("Server startup...")
	fmt.Println("Listening at", srv.String())

	// Listen for incoming connections.
	srv.ln, err = net.Listen("tcp", srv.String())
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	srv.wg.Add(1)
	go func() {
		defer func() {
			fmt.Println("Stop listen incoming connection")
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
					fmt.Println("Error accepting: ", err.Error())
					continue
				}
				// Handle connections in a new goroutine.
				go srv.handleRequest(conn)
			}
		}
	}()

	fmt.Println("Server is started")
	return err
}

// Handles incoming requests.
func (srv *ModbusServer) handleRequest(conn net.Conn) error {
	var (
		id_packet int
		err       error
		request   *ModbusPacket = &ModbusPacket{
			MTP: srv.MbProtocol}
	)
	// Close the connection when you're done with it.
	defer conn.Close()
	srv.wg.Add(1)
	request.Init()

	fmt.Printf(
		"Src->: \t\t\t\t%s\nDst<-: \t\t\t\t%s\n",
		conn.RemoteAddr(),
		conn.LocalAddr())

	// Read the incoming connection into the buffer.
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		select {
		case <-srv.done:
			fmt.Println("Close connection", conn.RemoteAddr())
			srv.wg.Done()
			return nil
		default:
			request.Length, err = conn.Read(request.Data)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				break
			}

			if request.Length == 0 {
				continue
			}

			id_packet++
			if id_packet == math.MaxInt32 {
				id_packet = 0
			}
			fmt.Printf("Src->: \t\t\t\t%s, Packet ID:%d\n", conn.RemoteAddr(), id_packet)

			// fmt.Println(hex.Dump(request.data))
			//request.ModbusDump()
			var answer *ModbusPacket
			answer, err = request.HandlerRequest(srv.Data)

			//answer.ModbusDump()
			conn.Write(answer.Data)
		}
	}
	return err
}
