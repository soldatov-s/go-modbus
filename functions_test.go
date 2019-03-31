// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"testing"
)

type testModbusFunctionCodeStringpair struct {
	fc   ModbusFunctionCode
	name string
}

var testsModbusFunctionCodeString = []testModbusFunctionCodeStringpair{
	{FcReadCoilStatus, "ReadCoilStatus"},
	{FcReadDescreteInputs, "ReadDescreteInputs"},
	{FcReadHoldingRegisters, "ReadHoldingRegisters"},
	{FcReadInputRegisters, "ReadInputRegisters"},
	{FcForceSingleCoil, "ForceSingleCoil"},
	{FcPresetSingleRegister, "PresetSingleRegister"},
	{FcForceMultipleCoils, "ForceMultipleCoils"},
	{FcPresetMultipleRegisters, "PresetMultipleRegisters"},
}

func TestModbusFunctionCode_String(t *testing.T) {
	for _, pair := range testsModbusFunctionCodeString {
		if pair.fc.String() != pair.name {
			t.Error("Expected ", pair.name, "got ", pair.fc.String())
		}
	}
}

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
	test_cnt := byte(5)
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
