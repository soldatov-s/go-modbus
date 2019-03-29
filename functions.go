// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"errors"
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

// Get FC parameters length
func (fc ModbusFunctionCode) Length(isReq bool) int {
	isError := bool(byte(fc)&byte(0x80) != 0)
	if isError {
		return 1
	}
	// ReadRegs/ReadIns - 2 + 2; Answer - 1 + data_len
	// WriteRegs/WriteOuts - 2 + 2 + 1 + data_len; Answer - 2 + 2
	// WriteReg/WriteOut - 2 + 2; Answer - 2 + 2
	switch fc {
	case ReadCoilStatus, ReadDescreteInputs, ReadHoldingRegisters, ReadInputRegisters:
		if isReq {
			return 4
		}
		return 1
	case ForceMultipleCoils, PresetMultipleRegisters:
		if !isReq {
			return 5
		}
		return 4
	case PresetSingleRegister, ForceSingleCoil:
		return 4
	default:
		return 0
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
		return buildErrAnswer(mp, 1), errors.New("Unknown function code")
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

// Build ModbusPacket
func buildPacket(isReq bool,
	TypeProtocol ModbusTypeProtocol,
	dev_id byte,
	fc ModbusFunctionCode,
	par1, par2 uint16,
	data ...byte) *ModbusPacket {
	isError := bool(byte(fc)&byte(0x80) != 0)
	mp := &ModbusPacket{TypeProtocol: TypeProtocol}
	mp.Length += fc.Length(isReq)
	if mp.Length == 0 {
		return nil
	}
	mp.Length += 4 // dev_id, fc, crc
	// Data length
	if (fc == ReadCoilStatus || fc == ReadDescreteInputs) && data == nil {
		mp.Length += int(boolCntToByteCnt(par2))
	} else {
		mp.Length += len(data)
	}
	mp.Length += mp.TypeProtocol.Offset()
	mp.Data = make([]byte, 2+mp.TypeProtocol.Offset(), mp.Length)
	mp.Data[mp.TypeProtocol.Offset() + 0] = dev_id
	mp.Data[mp.TypeProtocol.Offset() + 1] = byte(fc)
	tmp := make([]byte, 2)
	if !isError {
		if !isReq {
			mp.Data = append(mp.Data, []byte{byte(len(data))}...)
		} else {
			binary.BigEndian.PutUint16(tmp, par1)
			mp.Data = append(mp.Data, tmp...)
			binary.BigEndian.PutUint16(tmp, par2)
			mp.Data = append(mp.Data, tmp...)
			if data != nil {
				mp.Data = append(mp.Data, []byte{byte(len(data))}...)
			}
		}
		mp.Data = append(mp.Data, data...)
	}
	AppendCrc16(&mp.Data)

	return mp
}

// Build answer for request
func buildAnswer(req *ModbusPacket, fc ModbusFunctionCode, data ...byte) *ModbusPacket {
	par1 := binary.BigEndian.Uint16(req.GetData(2,4))
	par2 := binary.BigEndian.Uint16(req.GetData(4,6))
	return buildPacket(false, req.TypeProtocol, req.GetAddr(), req.GetFC(), par1, par2, data...)
}

// Error handler, builds ModbusPacket for error answer
func buildErrAnswer(req *ModbusPacket, errCode uint16) *ModbusPacket {
	return buildPacket(false, req.TypeProtocol, req.GetAddr(), req.GetFC()|ModbusFunctionCode(0x80), errCode, 0)
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
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))
	cnt := binary.BigEndian.Uint16(mp.GetData(4,6))

	// Try get data for answer
	data, err := md.ReadHoldingRegisters(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ReadHoldingRegisters, wordArrToByteArr(data)...), nil
}

// Read Inputs registers
func ReadInputRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))
	cnt := binary.BigEndian.Uint16(mp.GetData(4,6))

	// Try get data for answer
	data, err := md.ReadInputRegisters(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ReadInputRegisters, wordArrToByteArr(data)...), nil
}

// Preset Single Register
func PresetSingleRegisterHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))

	// Set values in ModbusData
	err := md.PresetSingleRegister(addr, binary.BigEndian.Uint16(mp.GetData(4,6)))
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, PresetSingleRegister), nil
}

// Convert byte array to word array
func byteArrToWordArr(data []byte) []uint16 {
	reg_data := make([]uint16, 0, len(data)/2)
	for i := uint16(0); i < uint16(len(data)/2); i++ {
		reg_data = append(reg_data, binary.BigEndian.Uint16(data[2*i:2*i+2]))
	}
	return reg_data
}

// Preset Multiple Holding Registers
func PresetMultipleRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))
	cnt := binary.BigEndian.Uint16(mp.GetData(4,6))
	// Set values in ModbusData
	err := md.PresetMultipleRegisters(addr, byteArrToWordArr(mp.GetData(7,7+cnt*2))...)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, PresetMultipleRegisters), nil
}

// Convert bool array to byte array
func boolArrToByteArr(data []bool) []byte {
	byte_data := make([]byte, boolCntToByteCnt(uint16(len(data))))
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
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))
	cnt := binary.BigEndian.Uint16(mp.GetData(4,6))

	// Data for answer
	data, err := md.ReadCoilStatus(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ReadCoilStatus, boolArrToByteArr(data)...), nil
}

// Read Descrete Inputs
func ReadDescreteInputsHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))
	cnt := binary.BigEndian.Uint16(mp.GetData(4,6))

	// Data for answer
	data, err := md.ReadDescreteInputs(addr, cnt)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ReadDescreteInputs, boolArrToByteArr(data)...), nil
}

// Force Single Coil
func ForceSingleCoilHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))

	// Set values in ModbusData
	b := binary.BigEndian.Uint16(mp.GetData(4,6))
	if b == 0xFF00 {
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
	bool_data := make([]bool, 0, cnt)
	j := uint16(0)
	for _, value := range data {
		for k := uint(0); k < 8; k++ {
			bool_data = append(bool_data, bool(value&byte(1<<k) != 0))
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
	addr := binary.BigEndian.Uint16(mp.GetData(2,4))
	cnt := binary.BigEndian.Uint16(mp.GetData(4,6))
	cnt_byte := mp.Data[6]

	// Set values in ModbusData
	err := md.ForceMultipleCoils(addr, byteArrToBoolArr(mp.GetData(7,7+cnt_byte), cnt)...)
	if err != nil {
		return buildErrAnswer(mp, 2), err
	}
	return buildAnswer(mp, ForceMultipleCoils), nil
}
