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

func (mp *ModbusPacket) ModbusDumper() {
	fmt.Printf("\nDump Modbus Packet\n")

	fmt.Printf("Packet length: \t\t\t%d\n", mp.length)
	fmt.Printf("Slave addr: \t\t\t%x\n", mp.GetAddr())
	fmt.Printf("Code function: \t\t\t%s(0x%x)\n", mp.GetFC(), byte(mp.GetFC()))
	fmt.Printf("Packet data: \t\t\t%s", hex.Dump(mp.GetData()))
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, mp.GetCrc())
	fmt.Printf("Modbus CRC16: \t\t\t%x %x\n\n", bs[0], bs[1])
}
