// functions
package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
)

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

func (fc ModbusFunctionCode) Handler(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	switch fc {
	case ReadCoilStatus:
		return ReadCoilStatusHndl(mp, md)
	case ReadDescreteInputs:
		return ReadDescreteInputsHndl(mp, md)
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

func wordArrToByteArr(data []uint16) []byte {
	var byte_data []byte
	for i := uint16(0); i < uint16(len(data)); i++ {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, data[i])
		byte_data = append(byte_data, b...)
	}

	return byte_data
}

func ReadHoldingRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	answer.mtp = mp.mtp
	// Copy addr and code function
	answer.data = append(answer.data, mp.GetPrefix()...)
	// Answer length in byte
	answer.data = append(answer.data, byte(cnt*2))
	// Data for answer
	data, err := md.ReadHoldingRegisters(addr, cnt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	answer.data = append(answer.data, wordArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.data)

	answer.length = len(answer.data)

	return answer, err
}

func ReadInputRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	answer.mtp = mp.mtp
	// Copy addr and code function
	answer.data = append(answer.data, mp.GetPrefix()...)
	// Answer length in byte
	answer.data = append(answer.data, byte(cnt*2))
	// Data for answer
	data, err := md.ReadInputRegisters(addr, cnt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	answer.data = append(answer.data, wordArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.data)

	answer.length = len(answer.data)

	return answer, err
}

func byteArrToWordArr(data []byte) []uint16 {
	var reg_data []uint16
	for i := uint16(0); i < uint16(len(data)/2); i++ {
		reg_data = append(reg_data, binary.BigEndian.Uint16(data[2*i:2*i+2]))
	}
	return reg_data
}

func PresetMultipleRegistersHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	answer.mtp = mp.mtp
	// Set values in ModbusData
	err := md.PresetMultipleRegisters(addr, cnt, byteArrToWordArr(mp.data[7:7+cnt*2]))
	// Copy addr and code function
	answer.data = append(answer.data, mp.GetPrefix()...)
	// Copy body
	answer.data = append(answer.data, mp.data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer.data)

	answer.length = len(answer.data)

	return answer, err
}

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

func ReadCoilStatusHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	answer.mtp = mp.mtp
	// Copy addr and code function
	answer.data = append(answer.data, mp.GetPrefix()...)
	// Data for answer
	data, err := md.ReadCoilStatus(addr, cnt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Answer length in byte
	answer.data = append(answer.data, byte(len(data)))

	answer.data = append(answer.data, boolArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.data)

	answer.length = len(answer.data)

	return answer, err
}

func ReadDescreteInputsHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	answer.mtp = mp.mtp
	// Copy addr and code function
	answer.data = append(answer.data, mp.GetPrefix()...)
	// Data for answer
	data, err := md.ReadDescreteInputs(addr, cnt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Answer length in byte
	answer.data = append(answer.data, byte(len(data)))

	answer.data = append(answer.data, boolArrToByteArr(data)...)
	// Crc Answer
	AppendCrc16(&answer.data)

	answer.length = len(answer.data)

	return answer, err
}

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

func ForceMultipleCoilsHndl(mp *ModbusPacket, md *ModbusData) (*ModbusPacket, error) {
	answer := new(ModbusPacket)
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])
	cnt_byte := mp.data[6]

	answer.mtp = mp.mtp
	// Set values in ModbusData
	err := md.ForceMultipleCoils(addr, cnt, byteArrToBoolArr(mp.data[8:8+cnt_byte*2], cnt))
	// Copy addr and code function
	answer.data = append(answer.data, mp.GetPrefix()...)
	// Copy body
	answer.data = append(answer.data, mp.data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer.data)
	answer.length = len(answer.data)

	return answer, err
}
