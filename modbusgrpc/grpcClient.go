// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbusgrpc

import (
	"fmt"

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

func NewgRPCClient(port, host string) (*ModbusgRPCClient, error) {
	var err error
	cl := new(ModbusgRPCClient)
	cl.Port = port
	cl.Host = host
	cl.Conn, err = grpc.Dial(cl.String(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Can not connected: %v", err)
	}
	cl.ServiceClient = NewModbusServiceClient(cl.Conn)
	return cl, err
}

// Close client
func (cl *ModbusgRPCClient) Close() {
	// Close the connection when you're done with it.
	cl.Conn.Close()
}

func (cl *ModbusgRPCClient) ReadHoldingRegisters(addr, cnt int32) ([]int32, error) {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadHoldingRegisters(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("Can't read holding registers: %v", err)
	}
	return answer.Data, nil
}

func (cl *ModbusgRPCClient) ReadInputRegisters(addr, cnt int32) ([]int32, error) {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadInputRegisters(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("Can't read input registers: %v", err)
	}
	return answer.Data, nil
}

func (cl *ModbusgRPCClient) ReadCoilStatus(addr, cnt int32) ([]bool, error) {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadCoilStatus(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("Can't read coils: %v", err)
	}
	return answer.Data, nil

}

func (cl *ModbusgRPCClient) ReadDescreteInputs(addr, cnt int32) ([]bool, error) {
	request := &ModbusRequest{Addr: addr, Cnt: cnt}
	answer, err := cl.ServiceClient.ReadDescreteInputs(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("Can't read inputs: %v", err)
	}
	return answer.Data, nil

}

func (cl *ModbusgRPCClient) PresetMultipleRegisters(addr int32, data []int32) ([]int32, error) {
	request := &ModbusWriteRegistersRequest{Addr: addr, Data: data}
	answer, err := cl.ServiceClient.PresetMultipleRegisters(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("Can't write holding registers: %v", err)
	}
	return answer.Data, nil

}

func (cl *ModbusgRPCClient) ForceMultipleCoils(addr int32, data []bool) ([]bool, error) {
	request := &ModbusWriteBitsRequest{Addr: addr, Data: data}
	answer, err := cl.ServiceClient.ForceMultipleCoils(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("Can't write coils: %v", err)
	}
	return answer.Data, nil
}
