// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

// Type of Modbus Protocol:
// - ModbusTCP
// - ModbusRTUviaTCP
type ModbusTypeProtocol int

const (
	ModbusTCP       ModbusTypeProtocol = 0
	ModbusRTUviaTCP ModbusTypeProtocol = 1
	ModbusUnknown                      = -1
)

const (
	ModbusRTUviaTCPMaxSize int = 256
	ModbusTCPMaxSize       int = 260
)

// Get MaxSize of packet for this protocol
func (p ModbusTypeProtocol) MaxSize() int {
	names := []int{
		ModbusTCPMaxSize,
		ModbusRTUviaTCPMaxSize}

	if p < ModbusTCP || p > ModbusRTUviaTCP {
		return 0
	}

	return names[p]
}

// PDU offest
func (p ModbusTypeProtocol) Offset() int {
	if p == ModbusTCP {
		return 6
	}
	return 0
}

// Get the name of this protocol
func (p ModbusTypeProtocol) String() string {
	names := []string{
		"ModbusTCP",
		"ModbusRTUviaTCP"}

	if p < ModbusTCP || p > ModbusRTUviaTCP {
		return "Unknown"
	}

	return names[p]
}

// Convert name protocol to type
func StringToModbusTypeProtocol(name string) ModbusTypeProtocol {
	switch name {
	case "ModbusTCP":
		return ModbusTCP
	case "ModbusRTUviaTCP":
		return ModbusRTUviaTCP
	default:
		return ModbusUnknown
	}
}
