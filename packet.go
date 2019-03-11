// packet
package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

const (
	ModbusTCPMaxSize int = 260
)

type ModbusPacket struct {
	data   []byte
	length int
	mtp    ModbusTypeProtocol
}

func (mp *ModbusPacket) Init() {
	mp.data = make([]byte, 260)
	mp.length = 260
}

func (mp *ModbusPacket) addr() byte {
	return mp.data[0]
}

func (mp *ModbusPacket) fcode() ModbusFunctionCode {
	return ModbusFunctionCode(mp.data[1])
}

func (mp *ModbusPacket) crc() uint16 {
	if mp.mtp == ModbusRTUviaTCP {
		return binary.BigEndian.Uint16(mp.data[mp.length-2 : mp.length])
	} else {
		return 0
	}
}

func (mp *ModbusPacket) Crc16Check() bool {
	res := false
	if mp.mtp == ModbusRTUviaTCP {
		res = Crc16Check(mp.data[:mp.length-2], mp.crc())
	} else if mp.mtp == ModbusTCP {
		res = true
	}

	if res {
		fmt.Println("CRC16 is OK")
	} else {
		fmt.Println("CRC16 is FAIL")
	}

	return res
}

func (mp *ModbusPacket) ErrorHandler() {

}

func (mp *ModbusPacket) HexStrToData(str string) {
	data, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	mp.data = make([]byte, 0, len(data))
	mp.length = len(data)
	copy(data, mp.data)
}

func (mp *ModbusPacket) ReadHoldRegs(md *ModbusData) []byte {
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	length := cnt*2 + 5
	answer := make([]byte, length)
	// Copy addr and code function
	copy(answer[0:2], mp.data[0:2])
	// Answer length in byte
	answer[2] = byte(length - 5)
	// Data for answer
	copy(answer[3:length-2], md.holding_reg[addr*2:(addr+cnt)*2])
	// Crc Answer
	AppendCrc16(answer)
	return answer
}

func (mp *ModbusPacket) PresetMultipleRegs(md *ModbusData) []byte {
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	// Set values in ModbusData
	md.SetHoldRegs(addr, cnt, mp.data[7:7+cnt*2])
	length := uint16(mp.length) - cnt*2 - 1
	answer := make([]byte, length)
	// Copy addr and code function
	copy(answer[0:2], mp.data[0:7])
	// Crc Answer
	AppendCrc16(answer)
	return answer

}

func (mp *ModbusPacket) ModbusDumper() {
	fmt.Printf("\nDump Modbus Packet\n")

	fmt.Printf("Packet length: \t\t\t%d\n", mp.length)
	fmt.Printf("Slave addr: \t\t\t%x\n", mp.addr())
	fmt.Printf("Code function: \t\t\t%s(0x%x)\n", mp.fcode(), byte(mp.fcode()))
	fmt.Printf("Packet data: \t\t\t%s", hex.Dump((mp.data[2 : mp.length-2])))
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, mp.crc())
	fmt.Printf("Modbus CRC16: \t\t\t%s\n", hex.Dump(bs))
}
