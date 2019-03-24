// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"errors"
	"log"
)

// Type of Modbus function
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

// Get the name of this function
func (fc ModbusFunctionCode) String() string {
	switch fc {
	case ReadCoilStatus:
		return "ReadCoilStatus"
	case ReadDescreteInputs:
		return "ReadDescreteInputs"
	case ReadHoldingRegisters:
		return "ReadHoldingRegisters"
	case ReadInputRegisters:
		return "ReadInputRegisters"
	case ForceSingleCoil:
		return "ForceSingleCoil"
	case PresetSingleRegister:
		return "PresetSingleRegister"
	case ForceMultipleCoils:
		return "ForceMultipleCoils"
	case PresetMultipleRegisters:
		return "PresetMultipleRegisters"
	default:
		return "Unknown"
	}
}

// Handler of the function by it code
func (fc ModbusFunctionCode) Handler(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	switch fc {
	case ReadCoilStatus:
		return ReadCoilStatusHndl(mp, md)
	case ReadDescreteInputs:
		return ReadDescreteInputsHndl(mp, md)
	case ForceSingleCoil:
		return ForceSingleCoilHndl(mp, md)
	case PresetSingleRegister:
		return PresetSingleRegisterHndl(mp, md)
	case ReadHoldingRegisters:
		return ReadHoldingRegistersHndl(mp, md)
	case ReadInputRegisters:
		return ReadInputRegistersHndl(mp, md)
	case ForceMultipleCoils:
		return ForceMultipleCoilsHndl(mp, md)
	case PresetMultipleRegisters:
		return PresetMultipleRegistersHndl(mp, md)
	default:
		return nil, errors.New("Unknown function code")
	}
}

// Convert word array to byte array
func wordArrToByteArr(data []uint16) []byte {
	var byte_data []byte
	for i := uint16(0); i < uint16(len(data)); i++ {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, data[i])
		byte_data = append(byte_data, b...)
	}

	return byte_data
}

// Read Holding registers
func ReadHoldingRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	answer.MTP = mp.MTP
	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Answer length in byte
	answer.Data = append(answer.Data, byte(cnt*2))
	// Data for answer
	data, err := md.ReadHoldingRegisters(addr, cnt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	answer.Data = append(answer.Data, wordArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.Data)

	answer.Length = len(answer.Data)

	return answer, err
}

// Read Inputs registers
func ReadInputRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	answer.MTP = mp.MTP
	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Answer length in byte
	answer.Data = append(answer.Data, byte(cnt*2))
	// Data for answer
	data, err := md.ReadInputRegisters(addr, cnt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	answer.Data = append(answer.Data, wordArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.Data)

	answer.Length = len(answer.Data)

	return answer, err
}

// Preset Single Register
func PresetSingleRegisterHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])

	answer.MTP = mp.MTP
	// Set values in ModbusData
	err := md.PresetSingleRegister(addr, binary.BigEndian.Uint16(mp.Data[4:6]))
	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Copy body
	answer.Data = append(answer.Data, mp.Data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer.Data)

	answer.Length = len(answer.Data)

	return answer, err
}

// Convert byte array to word array
func byteArrToWordArr(data []byte) []uint16 {
	var reg_data []uint16
	for i := uint16(0); i < uint16(len(data)/2); i++ {
		reg_data = append(reg_data, binary.BigEndian.Uint16(data[2*i:2*i+2]))
	}
	return reg_data
}

// Preset Multiple Holding Registers
func PresetMultipleRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	answer.MTP = mp.MTP
	// Set values in ModbusData
	err := md.PresetMultipleRegisters(addr, byteArrToWordArr(mp.Data[7:7+cnt*2]))
	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Copy body
	answer.Data = append(answer.Data, mp.Data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer.Data)

	answer.Length = len(answer.Data)

	return answer, err
}

// Convert bool array to byte array
func boolArrToByteArr(data []bool) []byte {
	var (
		byte_data []byte
		b, j      byte
	)
	for i := uint16(0); i < uint16(len(data)); i++ {
		if data[i] {
			b = b | 1<<j
		}
		j++

		if j == 7 {
			byte_data = append(byte_data, b)
			j = 0
			b = 0
		}
	}
	if j < 7 {
		byte_data = append(byte_data, b)
	}

	return byte_data
}

// Read Coil Status
func ReadCoilStatusHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	answer.MTP = mp.MTP
	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Data for answer
	data, err := md.ReadCoilStatus(addr, cnt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Answer length in byte
	answer.Data = append(answer.Data, byte(len(data)))

	answer.Data = append(answer.Data, boolArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.Data)

	answer.Length = len(answer.Data)

	return answer, err
}

// Read Descrete Inputs
func ReadDescreteInputsHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	answer.MTP = mp.MTP
	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Data for answer
	data, err := md.ReadDescreteInputs(addr, cnt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Answer length in byte
	answer.Data = append(answer.Data, byte(len(data)))

	answer.Data = append(answer.Data, boolArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.Data)

	answer.Length = len(answer.Data)

	return answer, err
}

// Force Single Coil
func ForceSingleCoilHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	var err error
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])

	answer.MTP = mp.MTP
	// Set values in ModbusData
	b := binary.BigEndian.Uint16(mp.Data[4:6])
	if b == 0xFF00 {
		err = md.ForceSingleCoil(addr, true)
	} else {
		err = md.ForceSingleCoil(addr, false)
	}

	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Copy body
	answer.Data = append(answer.Data, mp.Data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer.Data)
	answer.Length = len(answer.Data)

	return answer, err
}

// Convert byte array to bool array
func byteArrToBoolArr(data []byte, cnt uint16) []bool {
	var (
		bool_data []bool
		j         uint16
	)
	for i := 0; i < len(data); i++ {
		for k := uint(0); k < 8; k++ {
			if (data[i] & byte(1<<k)) != 0 {
				bool_data = append(bool_data, true)
			} else {
				bool_data = append(bool_data, false)
			}

			j++
		}
		if j == cnt {
			break
		}
	}
	return bool_data
}

// Force Multiple Coils
func ForceMultipleCoilsHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])
	cnt_byte := mp.Data[6]

	answer.MTP = mp.MTP
	// Set values in ModbusData
	err := md.ForceMultipleCoils(addr, byteArrToBoolArr(mp.Data[8:8+cnt_byte*2], cnt))
	// Copy addr and code function
	answer.Data = append(answer.Data, mp.GetPrefix()...)
	// Copy body
	answer.Data = append(answer.Data, mp.Data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer.Data)
	answer.Length = len(answer.Data)

	return answer, err
}
