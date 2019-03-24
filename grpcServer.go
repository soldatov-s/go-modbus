// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
)

// ModbusService is service for gRPC
type ModbusService struct {
	ModbusBase
	ln net.Listener // Listener
}

// Return string with host ip/name and port
func (s *ModbusService) String() string {
	return s.Host + ":" + s.Port
}

// ReadHoldingRegisters handel request to gRPC server
func (s *ModbusService) ReadHoldingRegisters(ctx context.Context, req *ModbusRequest) (*RegisterResponse, error) {

	// Read holding registers from ModbusData
	answer, err := s.Data.ReadHoldingRegisters(uint16(req.Addr), uint16(req.Cnt))
	if err != nil {
		return nil, err
	}

	answer_int32 := make([]int32, 0, len(answer))
	for _, a := range answer {
		answer_int32 = append(answer_int32, int32(a))
	}

	return &RegisterResponse{Data: answer_int32}, nil
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
