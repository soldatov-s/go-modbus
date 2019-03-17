// packet
package modbus

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
	if mp.mtp == ModbusRTUviaTCP {
		return mp.data[0]
	} else if mp.mtp == ModbusTCP {
		return mp.data[5]
	} else {
		return 0
	}
}

func (mp *ModbusPacket) GetFC() ModbusFunctionCode {
	if mp.mtp == ModbusRTUviaTCP {
		return ModbusFunctionCode(mp.data[1])
	} else if mp.mtp == ModbusTCP {
		return ModbusFunctionCode(mp.data[6])
	} else {
		return 0
	}
}

func (mp *ModbusPacket) HandlerRequest(md *ModbusData) (*ModbusPacket, error) {
	return mp.GetFC().Handler(mp, md)
}

func (mp *ModbusPacket) GetPrefix() []byte {
	if mp.mtp == ModbusRTUviaTCP {
		return mp.data[0:2]
	} else if mp.mtp == ModbusTCP {
		return mp.data[5:7]
	} else {
		return []byte{0x0}
	}
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

func (mp *ModbusPacket) ModbusDump() {
	fmt.Printf("\nDump Modbus Packet\n")

	fmt.Printf("Packet length: \t\t\t%d\n", mp.length)
	fmt.Printf("Slave addr: \t\t\t%x\n", mp.GetAddr())
	fmt.Printf("Code function: \t\t\t%s(0x%x)\n", mp.GetFC(), byte(mp.GetFC()))
	fmt.Println("Packet data:")
	fmt.Println(hex.Dump(mp.GetData()))
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, mp.GetCrc())
	fmt.Printf("Modbus CRC16: \t\t\t%x %x\n\n", bs[0], bs[1])
}
