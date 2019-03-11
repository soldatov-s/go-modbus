// functions
package main

type ModbusFunctionCode byte

const (
	ReadCoilStatus          ModbusFunctionCode = 0x01
	ReadDescreteInputs      ModbusFunctionCode = 0x02
	ReadHoldingRegisters    ModbusFunctionCode = 0x03
	ReadInputRegisters      ModbusFunctionCode = 0x04
	ForceSingleCoil         ModbusFunctionCode = 0x05
	PresetSingleRegister    ModbusFunctionCode = 0x06
	ForceMultipleCoils      ModbusFunctionCode = 0x0F
	PresetMultipleRegisters ModbusFunctionCode = 0x10
)

func (p ModbusFunctionCode) String() string {
	names := []string{
		"Unknown",
		"ReadCoilStatus",
		"ReadDescreteInputs",
		"ReadHoldingRegisters",
		"ReadInputRegisters",
		"ForceSingleCoil",
		"PresetSingleRegister",
		"Unknown",
		"Unknown",
		"Unknown",
		"Unknown",
		"Unknown",
		"Unknown",
		"Unknown",
		"Unknown",
		"ForceMultipleCoils",
		"PresetMultipleRegisters"}

	if p < ReadCoilStatus || p > PresetMultipleRegisters {
		return "Unknown"
	}

	return names[p]
}
