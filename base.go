// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"fmt"
)

// ModbusBase implements base server interface
type ModbusBaseServer struct {
	Host string      // Host Name/IP
	Port string      // Server port
	Data *ModbusData // Modbus Data
}

type IModbusBaseServer interface {
	Start() error
	Stop() error
}

// Return string with host ip/name and port
func (b *ModbusBaseServer) String() string {
	return fmt.Sprintf("%s:%s", b.Host, b.Port)
}

// ModbusBase implements base server interface
type ModbusBaseClient struct {
	Host string // Host Name/IP
	Port string // Server port
}

// Return string with host ip/name and port
func (b *ModbusBaseClient) String() string {
	return fmt.Sprintf("%s:%s", b.Host, b.Port)
}
