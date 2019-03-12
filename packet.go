// packet
package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ModbusPacket struct {
	data   []byte
	length int
	mtp    ModbusTypeProtocol
}

func (mp *ModbusPacket) Init() {
	mp.data = make([]byte, mp.mtp.MaxSize())
}

func (mp *ModbusPacket) GetAddr() byte {
	return mp.data[0]
}

func (mp *ModbusPacket) GetFC() ModbusFunctionCode {
	return ModbusFunctionCode(mp.data[1])
}

func (mp *ModbusPacket) GetPrefix() []byte {
	return mp.data[0:2]
}

func (mp *ModbusPacket) GetData() []byte {
	if mp.length == 0 {
		return nil
	}

	return mp.data[2 : mp.length-2]
}

func (mp *ModbusPacket) GetCrc() uint16 {
	if mp.length == 0 {
		return 0
	}

	if mp.mtp == ModbusRTUviaTCP {
		return binary.BigEndian.Uint16(mp.data[mp.length-2 : mp.length])
	} else {
		return 0
	}
}

func (mp *ModbusPacket) Crc16Check() bool {
	if mp.length == 0 {
		return false
	}

	res := true
	if mp.mtp == ModbusRTUviaTCP {
		res = Crc16Check(mp.data[:mp.length-2], mp.GetCrc())
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

/*func (mp *ModbusPacket) HexStrToData(str string) {
	data, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	mp.data = make([]byte, 0, len(data))
	mp.length = len(data)
	copy(data, mp.data)
}*/

func (mp *ModbusPacket) ReadHoldRegs(md *ModbusData) []byte {
	var answer []byte
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	// Copy addr and code function
	answer = append(answer, mp.GetPrefix()...)
	// Answer length in byte
	answer = append(answer, byte(cnt*2))
	// Data for answer
	answer = append(answer, md.ReadHoldRegs(addr, cnt)...)
	// Crc Answer
	AppendCrc16(&answer)
	return answer
}

func (mp *ModbusPacket) PresetMultipleRegs(md *ModbusData) []byte {
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	// Set values in ModbusData
	md.SetHoldRegs(addr, cnt, mp.data[7:7+cnt*2])
	// Copy addr and code function
	answer = append(answer, mp.GetPrefix()...)
	//
	answer = append(answer, mp.data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer)
	return answer
}

func (mp *ModbusPacket) ModbusDumper() {
	fmt.Printf("\nDump Modbus Packet\n")

	fmt.Printf("Packet length: \t\t\t%d\n", mp.length)
	fmt.Printf("Slave addr: \t\t\t%x\n", mp.GetAddr())
	fmt.Printf("Code function: \t\t\t%s(0x%x)\n", mp.GetFC(), mp.GetFC())
	fmt.Printf("Packet data: \t\t\t%s", hex.Dump(mp.GetData()))
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, mp.GetCrc())
	fmt.Printf("Modbus CRC16: \t\t\t%s\n", hex.Dump(bs))
}
