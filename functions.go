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

// Build ModbusPacket
func buildPacket(TypeProtocol ModbusTypeProtocol, dev_id byte, fc ModbusFunctionCode, 
		 par1, par2 uint16, data ...byte) *ModbusPacket {
	isError := bool(byte(mp.GetFC() & byte(0x80) != 0)
	mp := new(ModbusPacket)
	mp.TypeProtocol = TypeProtocol
	mp.Length = 4 // dev_id, fc, crc

	// ReadRegs/ReadIns - 2 + 2; Answer - 1 + data_len
	// WriteRegs/WriteOuts - 2 + 2 + 1 + data_len; Answer - 2 + 2
	// WriteReg/WriteOut - 2 + 2; Answer - 2 + 2
	switch fc {
	case ReadCoilStatus, ReadDescreteInputs, ReadHoldingRegistersHndl:
		if data == nill {
			mp.Length += 4
		} else {
			mp.Length += 1
		}
	case ForceMultipleCoils, PresetMultipleRegisters:
		if data != nill {
			mp.Length += 5
		} else {
			mp.Length += 4
		}
	case PresetSingleRegister, ForceSingleCoil:
		mp.Length += 4
	default:
		if isError {
			mp.Length += 1
		}
		mp.Length = 0
	}
	if mp.Length == 0 {
		return nil
	}
	mp.Length += len(data)
	mp.Data = make([]byte, 2, mp.Length)
	mp.Data[0] = dev_id
	mp.Data[1] = byte(fc)
	binary.BigEndian.PutUint16(mp.Data[2:3], par1)
	if !isError {
		binary.BigEndian.PutUint16(mp.Data[4:6], par2)
		if data != nill {
			mp.Data = append(mp.Data, data...)
		}
	}
	AppendCrc16(&mp.Data)
}

// Build answer for request
func buildAnswer(req *ModbusPacket, fc ModbusFunctionCode, data ...byte) *ModbusPacket {
	par1 := binary.BigEndian.Uint16(mp.Data[2:3])
	par2 := binary.BigEndian.Uint16(mp.Data[4:6])
	return buildPacket(req.TypeProtocol, req.GetAddr(), byte(mp.GetFC()), par1, par2, data...)
}

// Error handler, builds ModbusPacket for error answer
func buildErrAnswer(mp *ModbusPacket, errCode byte) *ModbusPacket {
	return buildPacket(req.TypeProtocol, req.GetAddr(), byte(mp.GetFC()) | byte(0x80), errCode, 0)
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
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ReadHoldingRegisters, wordArrToByteArr(data)...), nil
}

// Read Inputs registers
func ReadInputRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	// Try get data for answer
	data, err := md.ReadInputRegisters(addr, cnt)
	if err != nil {
		log.Println(err)
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ReadInputRegisters, wordArrToByteArr(data)...), nil
}

// Preset Single Register
func PresetSingleRegisterHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])

	// Set values in ModbusData
	err := md.PresetSingleRegister(addr, binary.BigEndian.Uint16(mp.Data[4:6]))
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, PresetSingleRegister), nil
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
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, PresetMultipleRegisters), nil
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
		return buildErrAnswer(mp, 2), err
	}
	// Init Answer Data
	q, r := cnt/8, cnt%8
	if r > 0 {
		q++
	}
	return buildAnswer(mp, ReadCoilStatus, boolArrToByteArr(data)...), nil
}

// Read Descrete Inputs
func ReadDescreteInputsHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.Data[2:4])
	cnt := binary.BigEndian.Uint16(mp.Data[4:6])

	// Data for answer
	data, err := md.ReadDescreteInputs(addr, cnt)
	if err != nil {
		log.Println(err)
		return buildErrAnswer(mp, 2), err
	}
	// Init Answer Data
	q, r := cnt/8, cnt%8
	if r > 0 {
		q++
	}
	return buildAnswer(mp, ReadDescreteInputs, boolArrToByteArr(data)...), nil
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
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ForceSingleCoil), nil
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
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ForceMultipleCoils), nill
}
