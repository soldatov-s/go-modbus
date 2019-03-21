// protocols
package modbus

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

func (p ModbusTypeProtocol) MaxSize() int {
	names := []int{
		ModbusTCPMaxSize,
		ModbusRTUviaTCPMaxSize}

	if p < ModbusTCP || p > ModbusRTUviaTCP {
		return 0
	}

	return names[p]
}

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
	switch name {
	case "ModbusTCP":
		return ModbusTCP
	case "ModbusRTUviaTCP":
		return ModbusRTUviaTCP
	default:
		return ModbusUnknown
	}
}
