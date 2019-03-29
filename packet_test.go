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
		mp.TypeProtocol = pair.TypeProtocol
		mp.Init()
		mp.Data = pair.data
		mp.Length = len(pair.data)
		res := mp.Crc16Check()
		if res != pair.res {
			t.Error(
				"For", pair.data,
				"expected", pair.res,
				"got", res,
			)

		}
	}
}

func TestModbusPacket_GetAddr(t *testing.T) {
	test_data := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}
	mp := &ModbusPacket{TypeProtocol: ModbusRTUviaTCP,
		Data:   test_data,
		Length: len(test_data)}
	res := mp.GetAddr()
	if res != byte(0x1) {
		t.Error("Expected", byte(0x1), "got", res)
	}
	test_data = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x3, 0x0, 0x0, 0x0, 0xA}
	mp = &ModbusPacket{TypeProtocol: ModbusTCP,
		Data:   test_data,
		Length: len(test_data)}
	res = mp.GetAddr()
	if res != byte(0x1) {
		t.Error("Expected", byte(0x1), "got", res)
	}
}

func TestModbusPacket_GetFC(t *testing.T) {
	test_data := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}
	mp := &ModbusPacket{TypeProtocol: ModbusRTUviaTCP,
		Data:   test_data,
		Length: len(test_data)}
	res := mp.GetFC()
	if res != ReadHoldingRegisters {
		t.Error("Expected", ReadHoldingRegisters, "got", res)
	}
	test_data = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x3, 0x0, 0x0, 0x0, 0xA}
	mp = &ModbusPacket{TypeProtocol: ModbusTCP,
		Data:   test_data,
		Length: len(test_data)}
	res = mp.GetFC()
	if res != ReadHoldingRegisters {
		t.Error("Expected", ReadHoldingRegisters, "got", res)
	}
}

func TestModbusPacket_GetCrc(t *testing.T) {
	test_data := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}
	mp := &ModbusPacket{TypeProtocol: ModbusRTUviaTCP,
		Data:   test_data,
		Length: len(test_data)}
	res := mp.GetCrc()
	if res != uint16(0xCDC5) {
		t.Error("Expected", uint16(0xCDC5), "got", res)
	}
	test_data = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x3, 0x0, 0x0, 0x0, 0xA}
	mp = &ModbusPacket{TypeProtocol: ModbusTCP,
		Data:   test_data,
		Length: len(test_data)}
	res = mp.GetCrc()
	if res != uint16(0x0) {
		t.Error("Expected", uint16(0x0), "got", res)
	}
}

func TestModbusPacket_GetData(t *testing.T) {
	test_data := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}
	mp := &ModbusPacket{TypeProtocol: ModbusRTUviaTCP,
		Data:   test_data,
		Length: len(test_data)}
	res := mp.GetData(0, len(test_data))
	for i, v := range test_data {
		if res[i] != v {
			t.Error("Expected", v, "got", res[i])
		}
	}
	test_data = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x3, 0x0, 0x0, 0x0, 0xA}
	mp = &ModbusPacket{TypeProtocol: ModbusTCP,
		Data:   test_data,
		Length: len(test_data)}
	res = mp.GetData(0, len(test_data)-mp.TypeProtocol.Offset())
	for i, v := range test_data[mp.TypeProtocol.Offset():] {
		if res[i] != v {
			t.Error("Expected", v, "got", res[i])
		}
	}
}
