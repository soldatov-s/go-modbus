// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

// ModbusBase implements base server interface
type ModbusBase struct {
	Host string      // Host Name/IP
	Port string      // Server port
	Data *ModbusData // Modbus Data
}
