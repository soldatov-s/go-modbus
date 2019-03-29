// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"sync"
	"testing"
)

func TestboolCntToByteCnt(t *testing.T) {
	test_data := uint16(9)
	res := boolCntToByteCnt(test_data)
	v := uint16(2)
	if v != res {
		t.Error("Expected ", v, "got ", res)
	}
}

func TestboolArrToByteArr(t *testing.T) {
	test_data := []bool{true, true, false, false, true}
	res := boolArrToByteArr(test_data)
	v := []byte{19}
	if v[0] != res[0] {
		t.Error("Expected ", v[0], "got ", res[0])
	}
}

func TestwordArrToByteArr(t *testing.T) {
	test_data := []uint16{0x01, 0x02, 0x03, 0x04}
	res := wordArrToByteArr(test_data)
	for i, v := range test_data {
		r := binary.BigEndian.Uint16(res[2*i : 2+2*i])
		if v != r {
			t.Error("Expected ", v, "got ", r)
		}
	}
}

func TestbyteArrToBoolArr(t *testing.T) {
	test_data := []byte{0x19}
	test_cnt := uint16(5)
	target := []bool{true, true, false, false, true}
	res := byteArrToBoolArr(test_data, test_cnt)
	for i, v := range target {
		if v != res[i] {
			t.Error("Expected ", v, "got ", res[i])
		}
	}
}

func TestbyteArrToWordArr(t *testing.T) {
	test_data := []byte{0x0, 0x01, 0x0, 0x02, 0x0, 0x03, 0x0, 0x04}
	res := byteArrToWordArr(test_data)
	for i, v := range res {
		r := binary.BigEndian.Uint16(test_data[2*i : 2+2*i])
		if v != r {
			t.Error("Expected ", r, "got ", v)
		}
	}
}

type testbuildPacketpair struct {
	isReq        bool
	TypeProtocol ModbusTypeProtocol
	dev_id 	     byte
	fc           ModbusFunctionCode
	par1, par2   uint16
	data         []byte
}

var testsbuildPacket = []testbuildPacketpair{
	{true, ModbusRTUviaTCP,  1, ReadHoldingRegisters, 0x0, 0xA, []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}},
}

func TestbuildPacket(t *testing.T) {
	for _, pair := range testsbuildPacket {
		mp := buildPacket(pair.isReq, pair.TypeProtocol, pair.dev_id, pair.fc, pair.par1, pair.par2)
		for i, v := range mp.Data {
			if v != pair.data[i] {
				t.Error(
					"For", pair,
					"expected", pair.data[i],
					"got", v,
				)
			}
		}
	}
}

type testFunctionCodeLengthpair struct {
	fc     ModbusFunctionCode
	isReq  bool
	length int
}

var testsFunctionCodeLength = []testFunctionCodeLengthpair{
	{ReadCoilStatus, true, 4},
}

func TestModbusFunctionCode_Length(t *testing.T) {
	for _, pair := range testsFunctionCodeLength {
		res := pair.fc.Length(pair.isReq)
		if res != pair.length {
			t.Error(
				"For", pair,
				"expected", pair.length,
				"got", res,
			)
		}
	}
}

func TestReadHoldingRegistersHndl(t *testing.T) {
	test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
	md := &ModbusData{holding_reg: make([]uint16, 10), mu_holding_regs: &sync.Mutex{}}
	md.PresetMultipleRegisters(0, test_data...)
	req := buildPacket(true, ModbusRTUviaTCP, 1, ReadHoldingRegisters, 0, 0x5)
	answ, _ := ReadHoldingRegistersHndl(req, md)
	for i, v := range test_data {
		res := binary.BigEndian.Uint16(answ.Data[3+2*i : 5+2*i])
		if v != res {
			t.Error("Expected ", v, "got ", res)
		}
	}
}

func TestPresetMultipleRegistersHndl(t *testing.T) {
	test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
	md := &ModbusData{holding_reg: make([]uint16, 10), mu_holding_regs: &sync.Mutex{}}
	req := buildPacket(true, ModbusRTUviaTCP, 1, PresetMultipleRegisters, 0, 0x5, wordArrToByteArr(test_data)...)
	PresetMultipleRegistersHndl(req, md)
	for i, v := range test_data {
		if md.holding_reg[i] != v {
			t.Error("Expected ", v, "got ", md.holding_reg[i])
		}
	}
}

func TestReadInputRegistersHndl(t *testing.T) {
	test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
	md := &ModbusData{input_reg: make([]uint16, 10)}
	md.PresetMultipleInputsRegisters(0, test_data...)
	req := buildPacket(true, ModbusRTUviaTCP, 1, ReadInputRegisters, 0, 0x5)
	answ, _ := ReadInputRegistersHndl(req, md)
	for i, v := range test_data {
		res := binary.BigEndian.Uint16(answ.Data[3+2*i : 5+2*i])
		if v != res {
			t.Error("Expected ", v, "got ", res)
		}
	}
}

func TestReadCoilStatusHndll(t *testing.T) {
	test_data := []bool{true, true, false, false, true}
	md := &ModbusData{coils: make([]bool, 10), mu_coils: &sync.Mutex{}}
	md.ForceMultipleCoils(0, test_data...)
	req := buildPacket(true, ModbusRTUviaTCP, 1, ReadCoilStatus, 0, 0x5)
	answ, _ := ReadCoilStatusHndl(req, md)
	cnt_byte := uint16(answ.Data[2])
	bool_arr := byteArrToBoolArr(answ.Data[3:3+cnt_byte], uint16(len(test_data)))
	for i, v := range test_data {
		if bool_arr[i] != v {
			t.Error("Expected ", v, "got ", bool_arr[i])
		}
	}
}

func TestForceMultipleCoilsHndl(t *testing.T) {
	test_data := []bool{true, true, false, false, true}
	md := &ModbusData{coils: make([]bool, 10), mu_coils: &sync.Mutex{}}
	req := buildPacket(true, ModbusRTUviaTCP, 1, ForceMultipleCoils, 0, 0x5, boolArrToByteArr(test_data)...)
	ForceMultipleCoilsHndl(req, md)
	for i, v := range test_data {
		if md.coils[i] != v {
			t.Error("Expected ", v, "got ", md.coils[i])
		}
	}
}

func TestReadDescreteInputsHndl(t *testing.T) {
	test_data := []bool{true, true, false, false, true}
	md := &ModbusData{discrete_inputs: make([]bool, 10)}
	md.ForceMultipleDescreteInputs(0, test_data...)
	req := buildPacket(true, ModbusRTUviaTCP, 1, ReadDescreteInputs, 0, 0x5)
	answ, _ := ReadDescreteInputsHndl(req, md)
	cnt_byte := uint16(answ.Data[2])
	bool_arr := byteArrToBoolArr(answ.Data[3:3+cnt_byte], uint16(len(test_data)))
	for i, v := range test_data {
		if bool_arr[i] != v {
			t.Error("Expected ", v, "got ", bool_arr[i])
		}
	}
}

func TestPresetSingleRegisterHndl(t *testing.T) {
	test_data := uint16(0x01)
	test_addr := uint16(0x0)
	md := &ModbusData{holding_reg: make([]uint16, 10), mu_holding_regs: &sync.Mutex{}}
	req := buildPacket(true, ModbusRTUviaTCP, 1, PresetSingleRegister, test_addr, test_data)
	PresetSingleRegisterHndl(req, md)
	if md.holding_reg[test_addr] != test_data {
		t.Error("Expected ", test_data, "got ", md.holding_reg[test_addr])
	}
}

func TestForceSingleCoilHndl(t *testing.T) {
	test_data := uint16(0x01)
	test_addr := uint16(0xFF00)
	md := &ModbusData{coils: make([]bool, 10), mu_coils: &sync.Mutex{}}
	req := buildPacket(true, ModbusRTUviaTCP, 1, ForceSingleCoil, test_addr, test_data)
	ForceSingleCoilHndl(req, md)
	if md.coils[test_addr] != true {
		t.Error("Expected ", true, "got ", md.coils[test_addr])
	}
}
