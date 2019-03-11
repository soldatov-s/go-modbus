// protocols
package main

type ModbusTypeProtocol int

const (
	ModbusTCP       ModbusTypeProtocol = 0
	ModbusRTUviaTCP ModbusTypeProtocol = 1
	Unknown                            = -1
)

func (p ModbusTypeProtocol) String() string {
	names := []string{
		"ModbusTCP",
		"ModbusRTUviaTCP"}

	if p < ModbusTCP || p > ModbusRTUviaTCP {
		return "Unknown"
	}

	return names[p]
}

func StringToModbusTypeProtocol(name string) ModbusTypeProtocol {
	if name == "ModbusTCP" {
		return ModbusTCP
	} else if name == "ModbusRTUviaTCP" {
		return ModbusRTUviaTCP
	}

	return Unknown
}
