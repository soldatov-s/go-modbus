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
	PDU          []byte             // PDU
	aPDU         []byte             // "real" modbus PDU
	Length       int                // Length of Data
	TypeProtocol ModbusTypeProtocol // Type Modbus Protocol
	isAnswer     bool               // Is it answer packet?
}

// Init ModbusPacket
func (mp *ModbusPacket) Init(typeProtocol ModbusTypeProtocol) {
	mp.TypeProtocol = typeProtocol
	mp.PDU = make([]byte, typeProtocol.MaxSize())
	mp.aPDU = mp.PDU[typeProtocol.Offset():]
}

// Get PDU length
func (mp *ModbusPacket) GetPDULength() int {
	return mp.Length + mp.TypeProtocol.Offset()
}

// Get device address field from packet
func (mp *ModbusPacket) GetDevID() byte {
	return mp.aPDU[0]
}

// Set device address field in packet
func (mp *ModbusPacket) SetDevID(devId byte) {
	mp.aPDU[0] = devId
}

// Get function code field from packet
func (mp *ModbusPacket) GetFunctionCode() ModbusFunctionCode {
	return ModbusFunctionCode(mp.aPDU[1])
}

// Set function code field in packet
func (mp *ModbusPacket) SetFunctionCode(fc ModbusFunctionCode) {
	mp.aPDU[1] = byte(fc)
}

// Get Error Code
func (mp *ModbusPacket) GetErrorCode() ModbusErrors {
	return ModbusErrors(mp.aPDU[2])
}

// Get fucntion parameters from request packet
func (mp *ModbusPacket) GetFunctionParameters() (uint16, uint16) {
	if mp.isAnswer {
		return 0, 0
	}
	return binary.BigEndian.Uint16(mp.aPDU[2:4]), binary.BigEndian.Uint16(mp.aPDU[4:6])
}

// Set function parameters to packet
func (mp *ModbusPacket) SetFunctionParameters(par1, par2 uint16) {
	binary.BigEndian.PutUint16(mp.aPDU[2:4], par1)
	binary.BigEndian.PutUint16(mp.aPDU[4:6], par2)
}

// Get data len byte + data bytes from Modbus packet
func (mp *ModbusPacket) GetData() (byte, []byte) {
	if mp.isAnswer {
		cnt := mp.aPDU[2]
		return cnt, mp.aPDU[3 : cnt+3]
	} else {
		cnt := mp.aPDU[6]
		return cnt, mp.aPDU[7 : cnt+7]
	}
}

// Set data with data len byte to Modbus packet
func (mp *ModbusPacket) SetData(cnt byte, data []byte) {
	if mp.isAnswer {
		mp.aPDU[2] = cnt
		copy(mp.aPDU[3:cnt+3], data)
	} else {
		mp.aPDU[6] = cnt
		copy(mp.aPDU[7:cnt+7], data)
	}
}

// Get CRC field from packet
func (mp *ModbusPacket) GetCrc() uint16 {
	if mp.Length == 0 || mp.TypeProtocol == ModbusTCP {
		return 0
	}
	return binary.BigEndian.Uint16(mp.aPDU[mp.Length-2 : mp.Length])
}

// Caculcate and adds crc to packet
// This function change ModbusPacket length
func (mp *ModbusPacket) SetCrc() {
	if mp.TypeProtocol == ModbusTCP {
		return
	}
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, Crc16(mp.aPDU[:mp.Length]))
	copy(mp.aPDU[mp.Length:mp.Length+2], bs)
	mp.Length += 2
}

// Recalculate and check CRC of packet
func (mp *ModbusPacket) IsCrcGood() bool {
	if mp.Length == 0 || mp.GetCrc() == 0 {
		return false
	}
	return Crc16Check(mp.aPDU[:mp.Length-2], mp.GetCrc())
}

// Get Transaction Id for ModbusTCP
func (mp *ModbusPacket) GetTransactionId() uint16 {
	if mp.TypeProtocol == ModbusRTUviaTCP {
		return 0
	}
	return binary.BigEndian.Uint16(mp.PDU[0:3])
}

// Set Transaction Id for ModbusTCP
func (mp *ModbusPacket) SetTransactionId(id uint16) {
	if mp.TypeProtocol == ModbusRTUviaTCP {
		return
	}
	binary.BigEndian.PutUint16(mp.PDU[0:3], id)
}

// Print Modbus Packet dump
func (mp *ModbusPacket) Dump(msg ...string) {
	if msg != nil {
		fmt.Printf("\n%s\n", msg)
	}
	if mp.TypeProtocol == ModbusTCP {
		fmt.Printf("Packet Transaction Id: \t\t%d\n", mp.GetTransactionId())
	}
	fmt.Printf("Packet length: \t\t\t%d\n", mp.Length-mp.TypeProtocol.Offset())
	fmt.Printf("Slave addr: \t\t\t%x\n", mp.GetDevID())
	fmt.Printf("Code function: \t\t\t%s(0x%x)\n", mp.GetFunctionCode(), byte(mp.GetFunctionCode()))
	fmt.Println("Packet data:")
	fmt.Println(hex.Dump(mp.aPDU[0 : mp.Length-mp.TypeProtocol.Offset()]))
	if mp.TypeProtocol == ModbusRTUviaTCP {
		bs := make([]byte, 2)
		binary.LittleEndian.PutUint16(bs, mp.GetCrc())
		fmt.Printf("Modbus CRC16: \t\t\t%x %x\n\n", bs[0], bs[1])
	}
}

func (mp *ModbusPacket) initLength() {
	if mp.TypeProtocol == ModbusTCP {
		mp.Length = 6
	}
	if mp.TypeProtocol == ModbusRTUviaTCP {
		mp.Length = 0
	}
}

// Build ModbusPacket
func (mp *ModbusPacket) buildPDU(devid byte, fc ModbusFunctionCode, par1, par2 uint16, data ...byte) {
	mp.initLength()
	// Set Device ID
	mp.SetDevID(devid)
	mp.Length++
	// Set Function Code
	mp.SetFunctionCode(fc)
	mp.Length++
	// Set parameters
	if mp.isAnswer && (fc == FcForceSingleCoil || fc == FcPresetSingleRegister ||
		fc == FcPresetMultipleRegisters || fc == FcForceMultipleCoils) {
		mp.SetFunctionParameters(par1, par2)
		mp.Length += 4
	}
	if !mp.isAnswer {
		mp.SetFunctionParameters(par1, par2)
		mp.Length += 4
	}
	// Set data
	if data != nil {
		mp.SetData(byte(len(data)), data)
		mp.Length += len(data) + 1
	}
	// Set Crc
	if mp.TypeProtocol == ModbusRTUviaTCP {
		mp.SetCrc()
	}
	// Set Message Length
	if mp.TypeProtocol == ModbusTCP {
		binary.BigEndian.PutUint16(mp.PDU[4:6], uint16(mp.Length-mp.TypeProtocol.Offset()))
	}

}

// Build ModbusPacket for errors
func (mp *ModbusPacket) buildErrPDU(devid byte, fc ModbusFunctionCode, errCode ModbusErrors) {
	mp.initLength()
	// Set Device ID
	mp.SetDevID(devid)
	mp.Length++
	// Set Function Code
	mp.SetFunctionCode(ModbusFunctionCode(byte(fc) | 0x80))
	mp.Length++
	// Set Error code
	mp.aPDU[mp.Length+1] = byte(errCode)
	mp.Length++
	// Set Crc
	if mp.TypeProtocol == ModbusRTUviaTCP {
		mp.SetCrc()
	}
}

// Build error answer ModbusPacket for src packet
func buildErrAnswer(src *ModbusPacket, errCode ModbusErrors) *ModbusPacket {
	answer := &ModbusPacket{isAnswer: true}
	answer.Init(src.TypeProtocol)
	answer.buildErrPDU(src.GetDevID(), src.GetFunctionCode(), errCode)
	answer.SetTransactionId(src.GetTransactionId())
	return answer
}

// Build answer ModbusPacket for src packet
func buildAnswer(src *ModbusPacket, data ...byte) *ModbusPacket {
	answer := &ModbusPacket{isAnswer: true}
	answer.Init(src.TypeProtocol)
	par1, par2 := src.GetFunctionParameters()
	answer.buildPDU(src.GetDevID(), src.GetFunctionCode(), par1, par2, data...)
	answer.SetTransactionId(src.GetTransactionId())
	return answer
}

// Build request
func buildRequest(transactionId uint16, typeProtocol ModbusTypeProtocol, devid byte, fc ModbusFunctionCode,
	par1, par2 uint16, data ...byte) *ModbusPacket {
	request := &ModbusPacket{}
	request.Init(typeProtocol)
	request.buildPDU(devid, fc, par1, par2, data...)
	request.SetTransactionId(transactionId)
	return request
}
