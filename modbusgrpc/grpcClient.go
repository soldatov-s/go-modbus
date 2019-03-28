// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbusgrpc

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	. "github.com/soldatov-s/go-modbus"
)

// ModbusClient implements client interface
type ModbusgRPCClient struct {
	ModbusBaseClient
	Conn          *grpc.ClientConn    // Connection
	ServiceClient ModbusServiceClient // Service
}

func NewgRPCClient(port, host string) *ModbusgRPCClient {
	var err error
	cl := &ModbusgRPCClient{Port: port, Host: host}
	cl.Conn, err = grpc.Dial(cl.String(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can not connected: %v", err)
	}
	cl.ServiceClient = NewModbusServiceClient(cl.Conn)
	return cl
}

func handelErr(err error, fmsg string) {
	if err != nil {
		log.Fatalf(fmsg, err)
	}
}

func (cl *ModbusgRPCClient) ReadHoldingRegisters(addr, cnt int32) []int32 {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadHoldingRegisters(context.Background(), request)
	handelErr(err, "Can't read holding registers: %v")
	return answer.Data
}

func (cl *ModbusgRPCClient) ReadInputRegisters(addr, cnt int32) []int32 {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadInputRegisters(context.Background(), request)
	handelErr(err, "Can't read input registers: %v")
	return answer.Data
}

func (cl *ModbusgRPCClient) ReadCoilStatus(addr, cnt int32) []bool {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadCoilStatus(context.Background(), request)
	handelErr(err, "Can't read coils: %v")
	return answer.Data
}

func (cl *ModbusgRPCClient) ReadDescreteInputs(addr, cnt int32) []bool {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadDescreteInputs(context.Background(), request)
	handelErr(err, "Can't read inputs: %v")
	return answer.Data
}

func (cl *ModbusgRPCClient) PresetMultipleRegisters(addr int32, data []int32) []int32 {
	request := &ModbusWriteRegistersRequest{Addr: addr, Data: data}
	answer, err := cl.ServiceClient.PresetMultipleRegisters(context.Background(), request)
	handelErr(err, "Can't write holding registers: %v")
	return answer.Data
}

func (cl *ModbusgRPCClient) ForceMultipleCoils(addr int32, data []bool) []bool {
	request := &ModbusWriteBitsRequest{Addr: addr, Data: data}
	answer, err := cl.ServiceClient.ForceMultipleCoils(context.Background(), request)
	handelErr(err, "Can't write coils: %v")
	return answer.Data
}
