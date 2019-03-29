// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"sync"
	"testing"
)

type testcheckOutsidepair struct {
	dataType  ModbusDataType
	addr, cnt uint16
	res       bool
}

var testscheckOutside = []testcheckOutsidepair{
	{HoldingRegisters, 0, 2, true},
	{DiscreteInputs, 0, 3, true},
	{InputRegisters, 0, 4, false},
}

func TestModbusData_checkOutside(t *testing.T) {
	md := &ModbusData{
		coils:           []bool{false, false, false},
		discrete_inputs: []bool{false, false, false},
		holding_reg:     []uint16{0, 0, 0},
		input_reg:       []uint16{0, 0, 0}}
	for _, pair := range testscheckOutside {
		res, _ := md.isNotOutside(pair.dataType, pair.addr, pair.cnt)
		if res != pair.res {
			t.Error(
				"For ModbusDataType", pair.dataType, "addr=", pair.addr, "cnt=", pair.cnt,
				"expected", pair.res,
				"got", res,
			)
		}
	}
}

func TestModbusData_PresetMultipleRegisters(t *testing.T) {
	test_addr := uint16(5)
	test_data := []uint16{10, 20, 30}
	md := &ModbusData{holding_reg: make([]uint16, 10), mu_holding_regs: &sync.Mutex{}}
	md.PresetMultipleRegisters(test_addr, test_data...)
	for i, v := range test_data {
		if v != md.holding_reg[test_addr+uint16(i)] {
			t.Error("Expected", test_data[i], "got", md.holding_reg[test_addr+uint16(i)])
		}
	}
}

func TestModbusData_ForceMultipleCoils(t *testing.T) {
	test_addr := uint16(5)
	test_data := []bool{true, false, true}
	md := &ModbusData{coils: make([]bool, 10), mu_coils: &sync.Mutex{}}
	md.ForceMultipleCoils(test_addr, test_data...)
	for i, v := range test_data {
		if v != md.coils[test_addr+uint16(i)] {
			t.Error("Expected", test_data[i], "got", md.holding_reg[test_addr+uint16(i)])
		}
	}
}

func TestModbusData_ReadHoldingRegisters(t *testing.T) {
	test_addr := uint16(0)
	test_data := []uint16{10, 20, 30}
	test_cnt := uint16(3)
	md := new(ModbusData)
	md.holding_reg = test_data
	res_data, _ := md.ReadHoldingRegisters(test_addr, test_cnt)
	for i, v := range test_data {
		if v != res_data[i] {
			t.Error("Expected", test_data[i], "got", res_data[i])
		}
	}
}

func TestModbusData_ReadCoilStatus(t *testing.T) {
	test_addr := uint16(0)
	test_data := []bool{true, false, true}
	test_cnt := uint16(3)
	md := new(ModbusData)
	md.coils = test_data
	res_data, _ := md.ReadCoilStatus(test_addr, test_cnt)
	for i, v := range test_data {
		if v != res_data[i] {
			t.Error("Expected", test_data[i], "got", res_data[i])
		}
	}
}
