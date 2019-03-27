// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbusgrpc

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
	ServiceClient * ModbusServiceClient // Service
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
	
	cl.ServiceClient := NewModbusServiceClient(cl.Conn)

	return cl
}

func (cl *ModbusgRPCClient) ReadHoldingRegisters(addr, cnt int32) []int32{
	var err     error

	request := &ModbusRequest{addr: addr, cnt: cnt}

	answer, err := client.ReadHoldingRegisters(context.Background(), request)
	if err != nil {
		log.Fatalf("Can't read holding registers: %v", err)
	}

	log.Println(answer.Data)
	return answer.data
}

func (cl *ModbusgRPCClient) ReadInputRegisters(addr, cnt int32) []int32{
	var err     error

	request := &ModbusRequest{addr: addr, cnt: cnt}

	answer, err := client.ReadInputRegisters(context.Background(), request)
	if err != nil {
		log.Fatalf("Can't read input registers: %v", err)
	}

	log.Println(answer.Data)
	return answer.data
}

func (cl *ModbusgRPCClient) ReadCoilStatus(addr, cnt int32) []bool{
	var err     error

	request := &ModbusRequest{addr: addr, cnt: cnt}

	answer, err := client.ReadCoilStatus(context.Background(), request)
	if err != nil {
		log.Fatalf("Can't read coils: %v", err)
	}

	log.Println(answer.Data)
	return answer.data
}

func (cl *ModbusgRPCClient) ReadDescreteInputs(addr, cnt int32) []bool{
	var err     error

	request := &ModbusRequest{addr: addr, cnt: cnt}

	answer, err := client.ReadDescreteInputs(context.Background(), request)
	if err != nil {
		log.Fatalf("Can't read inputs: %v", err)
	}

	log.Println(answer.Data)
	return answer.data
}

func (cl *ModbusgRPCClient) PresetMultipleRegisters(addr int32, data []int32) []int32{
	var err     error

	request := &ModbusRequest{addr: addr, data: data}

	answer, err := client.PresetMultipleRegisters(context.Background(), request)
	if err != nil {
		log.Fatalf("Can't write holding registers: %v", err)
	}

	log.Println(answer.Data)
	return answer.data
}

func (cl *ModbusgRPCClient) ForceMultipleCoils(addr int32, data []bool) []bool{
	var err     error

	request := &ModbusRequest{addr: addr, data: data}

	answer, err := client.ForceMultipleCoils(context.Background(), request)
	if err != nil {
		log.Fatalf("Can't write coils: %v", err)
	}

	log.Println(answer.Data)
	return answer.data
}
