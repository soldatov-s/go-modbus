// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

// Package modbus provides a server for MODBUS RTU over TCP.

package modbus

import (
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
	MTP              ModbusTypeProtocol // Type of Modbus protocol: TCP or RTU over TCP
	ln               net.Listener       // Listener
	done             chan struct{}      // Chan for sending "done" command
	exited           chan struct{}      // Chan for sending to main app signal that server is fully stopped
	wg               sync.WaitGroup     // WaitGroup for waiting end all connection
}

// NewServer function initializate new instance of ModbusServer
func NewServer(host, port string, mbprotocol ModbusTypeProtocol, md *ModbusData) *ModbusServer {
	srv := new(ModbusServer)
	srv.Host = host
	srv.Port = port
	srv.MTP = mbprotocol
	srv.Data = md
	srv.done = make(chan struct{})
	srv.exited = make(chan struct{})

	return srv
}

// Stop function close listener and wait closing all connection
func (srv *ModbusServer) Stop() error {
	var err error

	log.Println("Shutting down server...")
	if srv.ln != nil {
		srv.ln.Close()
	}
	close(srv.done)
	srv.wg.Wait()
	log.Println("Server is stopped")
	return err
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
	return err
}

// Handles incoming requests.
func (srv *ModbusServer) handleRequest(conn net.Conn) error {
	var (
		id_packet int
		err       error
		request   *ModbusPacket = &ModbusPacket{
			TypeProtocol: srv.MTP}
	)
	// Close the connection when you're done with it.
	defer conn.Close()
	srv.wg.Add(1)
	request.Init()

	log.Printf(
		"Src->: %s Dst<-: %s\n",
		conn.RemoteAddr(),
		conn.LocalAddr())

	// Read the incoming connection into the buffer.
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		select {
		case <-srv.done:
			log.Println("Close connection", conn.RemoteAddr())
			srv.wg.Done()
			return nil
		default:
			request.Length, err = conn.Read(request.Data)
			if err != nil {
				log.Println("Error reading:", err.Error())
				break
			}

			if request.Length == 0 {
				continue
			}

			if id_packet++; id_packet == math.MaxInt32 {
				id_packet = 0
			}
			log.Printf("Src->: %s, Packet ID:%d\n", conn.RemoteAddr(), id_packet)

			// fmt.Println(hex.Dump(request.data))
			log.Println("****Request Dump****")
			request.ModbusDump()
			var answer *ModbusPacket
			answer, err = request.HandlerRequest(srv.Data)
			log.Println("****Answer Dump****")
			answer.ModbusDump()
			conn.Write(answer.Data)
		}
	}
	return err
}
