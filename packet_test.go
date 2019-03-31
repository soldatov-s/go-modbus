// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"testing"
)

type testCrc16Checkpair struct {
	TypeProtocol ModbusTypeProtocol
	data         []byte
	res          bool
}

var testsCrc16Check = []testCrc16Checkpair{
	{ModbusRTUviaTCP, []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}, true},
}

func TestModbusPacket_Crc16Check(t *testing.T) {
	mp := new(ModbusPacket)

	for _, pair := range testsCrc16Check {
		mp.Init(pair.TypeProtocol)
		copy(mp.PDU, pair.data)
		mp.Length = len(pair.data)
		res := mp.IsCrcGood()
		if res != pair.res {
			t.Error(
				"For", pair.data,
				"expected", pair.res,
				"got", res,
			)

		}
	}
}

func TestModbusPacket_GetFunctionCode(t *testing.T) {
	test_data := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}
	mp := &ModbusPacket{}
	mp.Init(ModbusRTUviaTCP)
	copy(mp.PDU, test_data)
	mp.Length = len(test_data)
	res := mp.GetFunctionCode()
	if res != ModbusFunctionCode(0x3) {
		t.Error("Expected", byte(0x3), "got", res)
	}
	test_data = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x3, 0x0, 0x0, 0x0, 0xA}
	mp.Init(ModbusTCP)
	copy(mp.PDU, test_data)
	mp.Length = len(test_data)
	res = mp.GetFunctionCode()
	if res != ModbusFunctionCode(0x3) {
		t.Error("Expected", byte(0x3), "got", res)
	}
}

func TestModbusPacket_GetFC(t *testing.T) {
	test_data := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}
	mp := &ModbusPacket{}
	mp.Init(ModbusRTUviaTCP)
	copy(mp.PDU, test_data)
	mp.Length = len(test_data)
	res := mp.GetFunctionCode()
	if res != FcReadHoldingRegisters {
		t.Error("Expected", FcReadHoldingRegisters, "got", res)
	}
	test_data = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x3, 0x0, 0x0, 0x0, 0xA}
	mp = &ModbusPacket{}
	mp.Init(ModbusTCP)
	copy(mp.PDU, test_data)
	mp.Length = len(test_data)
	res = mp.GetFunctionCode()
	if res != FcReadHoldingRegisters {
		t.Error("Expected", FcReadHoldingRegisters, "got", res)
	}
}

func TestModbusPacket_GetCrc(t *testing.T) {
	test_data := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}
	mp := &ModbusPacket{}
	mp.Init(ModbusRTUviaTCP)
	copy(mp.PDU, test_data)
	mp.Length = len(test_data)
	res := mp.GetCrc()
	if res != uint16(0xCDC5) {
		t.Error("Expected", uint16(0xCDC5), "got", res)
	}
	test_data = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x3, 0x0, 0x0, 0x0, 0xA}
	mp = &ModbusPacket{}
	mp.Init(ModbusTCP)
	copy(mp.PDU, test_data)
	mp.Length = len(test_data)
	res = mp.GetCrc()
	if res != uint16(0x0) {
		t.Error("Expected", uint16(0x0), "got", res)
	}
}

type testbuildPacketpair struct {
	isReq        bool
	TypeProtocol ModbusTypeProtocol
	dev_id       byte
	fc           ModbusFunctionCode
	par1, par2   uint16
	PDU          []byte
}

var testsbuildPacket = []testbuildPacketpair{
	{true, ModbusRTUviaTCP, 1, FcReadHoldingRegisters, 0x0, 0xA, []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xC5, 0xCD}},
}

func TestModbusPacket_buildPDU(t *testing.T) {
	for _, pair := range testsbuildPacket {
		mp := &ModbusPacket{}
		mp.Init(pair.TypeProtocol)
		mp.buildPDU(pair.dev_id, pair.fc, pair.par1, pair.par2)
		for i, v := range mp.PDU[:mp.GetPDULength()] {
			if v != pair.PDU[i] {
				t.Error(
					"For", pair,
					"expected", pair.PDU[i],
					"got", v,
				)
			}
		}
	}
}
