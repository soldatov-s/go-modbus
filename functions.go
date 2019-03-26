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
		return errorHndl(mp, 1), errors.New("Unknown function code")
	}
}

func prepareAnswer(req *ModbusPacket, data_len byte) *ModbusPacket {
	answer := new(ModbusPacket)
	answer.TypeProtocol = req.TypeProtocol
	answer.Data = make([]byte, 0, data_len)
	// Copy addr and code function
	answer.Data = append(answer.Data, req.GetPrefix()...)
	return answer
}

func endAnswer(mp *ModbusPacket) {
	// Crc Answer
	AppendCrc16(&mp.Data)
	mp.Length = len(mp.Data)
}

func buildAnswer(req *ModbusPacket, pref_len byte, data_len byte, data ...byte) *ModbusPacket {
	// Init Answer Data
	answer := prepareAnswer(req, pref_len+data_len)
	if data_len > 0 {
		answer.Data = append(answer.Data, data_len)
	}
	answer.Data = append(answer.Data, data...)
	// End answer
	endAnswer(answer)

	return answer
}

// Error handler, builds ModbusPacket for error answer
func errorHndl(mp *ModbusPacket, errCode byte) *ModbusPacket {
	answer := prepareAnswer(mp, 5)
	answer.Data[2] = byte(mp.GetFC()) | byte(0x80)
	answer.Data = append(answer.Data, errCode)
	endAnswer(answer)
	return answer
}

// Convert word array to byte array
func wordArrToByteArr(data []uint16) []byte {
	byte_data := make([]byte, len(data)*2)
	for i, value := range data {
		binary.BigEndian.PutUint16(byte_data[i*2:(i+1)*2], value)
	}
	return byte_data
}

// Read Holding registers
func ReadHoldingRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	// Try get data for answer
	data, err := md.ReadHoldingRegisters(addr, cnt)
	if err != nil {
		log.Println(err)
		return errorHndl(mp, 2), err
	}
	answer := buildAnswer(mp, 5, byte(cnt*2), wordArrToByteArr(data)...)
	return answer, err
}

// Read Inputs registers
func ReadInputRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	// Try get data for answer
	data, err := md.ReadInputRegisters(addr, cnt)
	if err != nil {
		log.Println(err)
		return errorHndl(mp, 2), err
	}
	answer := buildAnswer(mp, 5, byte(cnt*2), wordArrToByteArr(data)...)

	return answer, err
}

// Preset Single Register
func PresetSingleRegisterHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])

	// Set values in ModbusData
	err := md.PresetSingleRegister(addr, binary.BigEndian.Uint16(mp.Data[4:6]))
	if err != nil {
		return errorHndl(mp, 2), err
	}
	answer := buildAnswer(mp, 8, 0, mp.GetData()...)

	return answer, err
}

// Convert byte array to word array
func byteArrToWordArr(data []byte) []uint16 {
	reg_data := make([]uint16, len(data)/2)
	for i := uint16(0); i < uint16(len(data)/2); i++ {
		reg_data = append(reg_data, binary.BigEndian.Uint16(data[2*i:2*i+2]))
	}
	return reg_data
}

// Preset Multiple Holding Registers
func PresetMultipleRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	// Set values in ModbusData
	err := md.PresetMultipleRegisters(addr, byteArrToWordArr(mp.Data[7:7+cnt*2]))
	if err != nil {
		return errorHndl(mp, 2), err
	}
	answer := buildAnswer(mp, 8, 0, mp.GetData()...)

	return answer, err
}

// Convert bool array to byte array
func boolArrToByteArr(data []bool) []byte {
	q, r := len(data)/8, len(data)%8
	if r > 0 {
		q++
	}

	byte_data := make([]byte, q)
	shift := 0
	for i, value := range data {
		if value {
			shift = i % 8
		}
		byte_data[i/8] |= byte(1 << uint(shift))
	}
	return byte_data
}

// Read Coil Status
func ReadCoilStatusHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	// Data for answer
	data, err := md.ReadCoilStatus(addr, cnt)
	if err != nil {
		log.Println(err)
		return errorHndl(mp, 2), err
	}
	// Init Answer Data
	q, r := cnt/8, cnt%8
	if r > 0 {
		q++
	}
	answer := buildAnswer(mp, 5, byte(q), boolArrToByteArr(data)...)

	return answer, err
}

// Read Descrete Inputs
func ReadDescreteInputsHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	// Data for answer
	data, err := md.ReadDescreteInputs(addr, cnt)
	if err != nil {
		log.Println(err)
		return errorHndl(mp, 2), err
	}
	// Init Answer Data
	q, r := cnt/8, cnt%8
	if r > 0 {
		q++
	}
	answer := buildAnswer(mp, 5, byte(q), boolArrToByteArr(data)...)

	return answer, err
}

// Force Single Coil
func ForceSingleCoilHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])

	// Set values in ModbusData
	b := binary.BigEndian.Uint16(mp.Data[4:6])
	if b != 0xFF00 {
		b = 1
	}
	err := md.ForceSingleCoil(addr, bool((b&1) == 1))
	if err != nil {
		return errorHndl(mp, 2), err
	}
	answer := buildAnswer(mp, 8, 0, mp.Data[2:6]...)

	return answer, err
}

// Convert byte array to bool array
func byteArrToBoolArr(data []byte, cnt uint16) []bool {
	bool_data := make([]bool, cnt)
	j := uint16(0)
	for _, value := range data {
		for k := uint(0); k < 8; k++ {
			bool_data = append(bool_data, bool(value&byte(1<<k) == 1))
			if j == cnt {
				break
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
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])
	cnt_byte := mp.Data[6]

	// Set values in ModbusData
	err := md.ForceMultipleCoils(addr, byteArrToBoolArr(mp.Data[8:8+cnt_byte*2], cnt))
	if err != nil {
		return errorHndl(mp, 2), err
	}
	answer := buildAnswer(mp, 8, 0, mp.Data[2:6]...)

	return answer, err
}
