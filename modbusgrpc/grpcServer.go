// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbusgrpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"

	. "github.com/soldatov-s/go-modbus"
)

// ModbusService is service for gRPC
type ModbusService struct {
	ModbusBaseServer
	ln net.Listener // Listener
}

func int16ArrToInt32Arr(data []int16) []int32{
	data_int32 := make([]int32, 0, len(data))
	for _, a := range answer {
		data_int32 = append(data_int32, int32(a))
	}
	return data_int32
}

// ReadHoldingRegisters handel request to gRPC server
func (s *ModbusService) ReadHoldingRegisters(ctx context.Context, req *ModbusRequest) (*RegisterResponse, error) {
	// Read holding registers from ModbusData
	answer, err := s.Data.ReadHoldingRegisters(uint16(req.Addr), uint16(req.Cnt))
	if err != nil {
		return nil, err
	}
	return &RegisterResponse{Data: int16ArrToInt32Arr(answer)}, nil
}

// ReadInputRegisters handel request to gRPC server
func (s *ModbusService) ReadInputRegisters(ctx context.Context, req *ModbusRequest) (*RegisterResponse, error) {
	// Read input registers from ModbusData
	answer, err := s.Data.ReadInputRegisters(uint16(req.Addr), uint16(req.Cnt))
	if err != nil {
		return nil, err
	}
	return &RegisterResponse{Data: int16ArrToInt32Arr(answer)}, nil
}

// ReadCoilStatus handel request to gRPC server
func (s *ModbusService) ReadCoilStatus(ctx context.Context, req *ModbusRequest) (*BitResponse, error) {
	// Read coils from ModbusData
	answer, err := s.Data.ReadCoilStatus(uint16(req.Addr), uint16(req.Cnt))
	if err != nil {
		return nil, err
	}
	return &BitResponse{Data: answer}, nil
}

// ReadDescreteInputs handel request to gRPC server
func (s *ModbusService) ReadDescreteInputs(ctx context.Context, req *ModbusRequest) (*BitResponse, error) {
	// Read inputs from ModbusData
	answer, err := s.Data.ReadDescreteInputs(uint16(req.Addr), uint16(req.Cnt))
	if err != nil {
		return nil, err
	}
	return &BitResponse{Data: answer}, nil
}

func int32ArrToInt16Arr(data []int32) []int16{
	data_int16 := make([]int16, 0, len(data))
	for _, a := range answer {
		data_int16 = append(data_int16, int16(a))
	}
	return data_int16
}

// PresetMultipleRegisters handel request to gRPC server
func (s *ModbusService) PresetMultipleRegisters(ctx context.Context, req *ModbusWriteRegistersRequest) (*RegisterResponse, error) {
	// Write holding registers to ModbusData
	err := s.Data.PresetMultipleRegisters(uint16(req.Addr), int32ArrToInt16Arr(req.Data)...)
	if err != nil {
		return nil, err
	}
	// Read holding registers
	return s.ReadHoldingRegisters(ctx, &ModbusRequest{Addr: req.Addr, Cnt: len(req.Data)}
}

// ForceMultipleCoils handel request to gRPC server
func (s *ModbusService) ForceMultipleCoils(ctx context.Context, req *ModbusWriteBitsRequest) (*BitResponse, error) {
	// Write coils to ModbusData
	err := s.Data.ForceMultipleCoils(uint16(req.Addr), req.Data...)
	if err != nil {
		return nil, err
	}
	// Read coils
	return s.ReadCoilStatus(ctx, &ModbusRequest{Addr: req.Addr, Cnt: len(req.Data)}
}

func NewgRPCService(host, port string, md *ModbusData) *ModbusService {
	srv := new(ModbusService)
	srv.Data = md
	srv.Host = host
	srv.Port = port

	return srv
}

// Start gRPC server
func (srv *ModbusService) Start() {
	var err error
	// Start gRPC server
	srv.ln, err = net.Listen("tcp", srv.String())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	RegisterModbusServiceServer(s, srv)

	// Register answer service at gRPC server
	reflection.Register(s)
	if err := s.Serve(srv.ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
