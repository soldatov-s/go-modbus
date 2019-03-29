// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"testing"
)

type testTypeProtocolpair struct {
	str string
	res ModbusTypeProtocol
}

var testsTypeProtocol = []testTypeProtocolpair{
	{"ModbusRTUviaTCP", ModbusRTUviaTCP},
	{"ModbusTCP", ModbusTCP},
}

func TestStringToModbusTypeProtocol(t *testing.T) {
	for _, pair := range testsTypeProtocol {
		res := StringToModbusTypeProtocol(pair.str)
		if res != pair.res {
			t.Error(
				"For", pair.str,
				"expected", pair.res,
				"got", res,
			)

		}
	}
}

type testModbusTypeProtocolpair struct {
	TypeProtocol   ModbusTypeProtocol
	name string
	maxSize int
	offset int
}

var testsModbusTypeProtocol = []testModbusTypeProtocolpair{
	{ModbusTCP, "ModbusTCP", 260, 6},
	{ModbusRTUviaTCP, "ModbusRTUviaTCP", 256, 0},
}

func TestModbusTypeProtocol_String(t *testing.T) {
	for _, pair := range testsModbusTypeProtocol {
		if pair.TypeProtocol.String() != pair.name {
			t.Error("Expected ", pair.name, "got ", pair.TypeProtocol.String())
		}
	}
}

func TestModbusTypeProtocol_MaxSize(t *testing.T) {
	for _, pair := range testsModbusTypeProtocol {
		if pair.TypeProtocol.MaxSize() != pair.maxSize {
			t.Error("Expected ", pair.maxSize, "got ", pair.TypeProtocol.MaxSize())
		}
	}
}

func TestModbusTypeProtocol_Offset(t *testing.T) {
	for _, pair := range testsModbusTypeProtocol {
		if pair.TypeProtocol.Offset() != pair.offset {
			t.Error("Expected ", pair.offset, "got ", pair.TypeProtocol.Offset())
		}
	}
}
