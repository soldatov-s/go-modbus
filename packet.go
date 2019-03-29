// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// ModbusPacket implements packet interface
type ModbusPacket struct {
	Data         []byte             // Packet Data
	Length       int                // Length of Data
	TypeProtocol ModbusTypeProtocol // Type Modbus Protocol
}

// Init ModbusPacket
func (mp *ModbusPacket) Init() {
	mp.Data = make([]byte, mp.TypeProtocol.MaxSize())
}

// Get device address field from packet
func (mp *ModbusPacket) GetAddr() byte {
	return mp.Data[mp.TypeProtocol.Offset()]
}

// Get function code field from packet
func (mp *ModbusPacket) GetFC() ModbusFunctionCode {
	return ModbusFunctionCode(mp.Data[1+mp.TypeProtocol.Offset()])
}

// Handler request by function code
func (mp *ModbusPacket) HandlerRequest(md *ModbusData) (*ModbusPacket, error) {
	return mp.GetFC().Handler(mp, md)
}

// Get prefix (device address & function code) from packet
func (mp *ModbusPacket) GetPrefix() []byte {
	if mp.TypeProtocol == ModbusRTUviaTCP || mp.TypeProtocol == ModbusTCP {
		return mp.Data[mp.TypeProtocol.Offset():mp.TypeProtocol.Offset()+2]
	}
	return []byte{0x0}
}

// Get body Modbus request from packet
func (mp *ModbusPacket) GetData() []byte {
	if mp.Length == 0 {
		return nil
	}
	return mp.Data[2 : mp.Length-2]
}

// Get CRC field from packet
func (mp *ModbusPacket) GetCrc() uint16 {
	if mp.Length == 0 || mp.TypeProtocol == ModbusTCP {
		return 0
	}
	return binary.BigEndian.Uint16(mp.Data[mp.Length-2 : mp.Length])
}

// Recalculate and check CRC of packet
func (mp *ModbusPacket) Crc16Check() bool {
	if mp.Length == 0 || mp.GetCrc() == 0 {
		return false
	}
	return Crc16Check(mp.Data[:mp.Length-2], mp.GetCrc())
}

// Print Modbus Packet dump
func (mp *ModbusPacket) ModbusDump(msg ...string) {
	fmt.Printf("\n%s\n", msg)
	fmt.Printf("Packet length: \t\t\t%d\n", mp.Length)
	fmt.Printf("Slave addr: \t\t\t%x\n", mp.GetAddr())
	fmt.Printf("Code function: \t\t\t%s(0x%x)\n", mp.GetFC(), byte(mp.GetFC()))
	fmt.Println("Packet data:")
	fmt.Println(hex.Dump(mp.GetData()))
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, mp.GetCrc())
	fmt.Printf("Modbus CRC16: \t\t\t%x %x\n\n", bs[0], bs[1])
}
