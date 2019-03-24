// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/json"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	. "github.com/soldatov-s/go-modbus"
)

// ModbusClient implements client interface
type ModbusgRPCClient struct {
	ModbusBaseClient
	Conn *grpc.ClientConn // Connection
}

func NewgRPCClient(port, host string) *ModbusgRPCClient {
	var err error
	cl := new(ModbusgRPCClient)
	cl.Port = port
	cl.Host = host

	cl.Conn, err = grpc.Dial(cl.String(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can not connected: %v", err)
	}

	return cl
}

func (cl *ModbusgRPCClient) Request(data []byte) {
	var (
		err     error
		request *ModbusRequest
	)

	client := NewModbusServiceClient(cl.Conn)
	json.Unmarshal(data, &request)

	if err != nil {
		log.Fatalf("Can't unmarshal: %v", err)
	}

	r, err := client.ReadHoldingRegisters(context.Background(), request)
	if err != nil {
		log.Fatalf("Can't read holding registers: %v", err)
	}

	log.Println(r.Data)
}
