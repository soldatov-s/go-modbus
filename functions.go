// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
)

// Type of Modbus function
type ModbusFunctionCode byte

const (
	FcReadCoilStatus          ModbusFunctionCode = 0x01
	FcReadDescreteInputs      ModbusFunctionCode = 0x02
	FcReadHoldingRegisters    ModbusFunctionCode = 0x03
	FcReadInputRegisters      ModbusFunctionCode = 0x04
	FcForceSingleCoil         ModbusFunctionCode = 0x05
	FcPresetSingleRegister    ModbusFunctionCode = 0x06
	FcForceMultipleCoils      ModbusFunctionCode = 0x0F
	FcPresetMultipleRegisters ModbusFunctionCode = 0x10
)

// Get the name of this function
func (fc ModbusFunctionCode) String() string {
	switch fc {
	case FcReadCoilStatus:
		return "ReadCoilStatus"
	case FcReadDescreteInputs:
		return "ReadDescreteInputs"
	case FcReadHoldingRegisters:
		return "ReadHoldingRegisters"
	case FcReadInputRegisters:
		return "ReadInputRegisters"
	case FcForceSingleCoil:
		return "ForceSingleCoil"
	case FcPresetSingleRegister:
		return "PresetSingleRegister"
	case FcForceMultipleCoils:
		return "ForceMultipleCoils"
	case FcPresetMultipleRegisters:
		return "PresetMultipleRegisters"
	default:
		return "Unknown"
	}
}

// Countes how many bytes we need to store bool array
func boolCntToByteCnt(cnt uint16) uint16 {
	q, r := cnt/8, cnt%8
	if r > 0 {
		q++
	}
	return q
}

// Convert word array to byte array
func wordArrToByteArr(data []uint16) []byte {
	byte_data := make([]byte, len(data)*2)
	for i, value := range data {
		binary.BigEndian.PutUint16(byte_data[i*2:(i+1)*2], value)
	}
	return byte_data
}

// Convert byte array to word array
func byteArrToWordArr(data []byte) []uint16 {
	reg_data := make([]uint16, 0, len(data)/2)
	for i := uint16(0); i < uint16(len(data)/2); i++ {
		reg_data = append(reg_data, binary.BigEndian.Uint16(data[2*i:2*i+2]))
	}
	return reg_data
}

// Convert bool array to byte array
func boolArrToByteArr(data []bool) []byte {
	byte_data := make([]byte, boolCntToByteCnt(uint16(len(data))))
	shift := 0
	for i, value := range data {
		if value {
			shift = i % 8
			byte_data[i/8] |= byte(1 << uint(shift))
		}
	}
	return byte_data
}

// Convert byte array to bool array
func byteArrToBoolArr(data []byte, cnt byte) []bool {
	bool_data := make([]bool, cnt)
	j := byte(0)
	for _, value := range data {
		for k := uint(0); k < 8; k++ {
			bool_data[j] = bool(value&byte(1<<k) != 0)
			j++
			if j == cnt {
				break
			}

		}
		if j == cnt {
			break
		}
	}
	return bool_data
}
