// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"sync"
	"testing"
)

func TestModbusServer_ReadHoldingRegisters(t *testing.T) {
	test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
	md := &ModbusData{holding_reg: make([]uint16, 10), mu_holding_regs: &sync.Mutex{}}
	md.PresetMultipleRegisters(0, test_data...)
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcReadHoldingRegisters, 0, 0x5)
	answ, _ := srv.ReadHoldingRegisters(req)
	_, answ_data := answ.GetData()
	for i, v := range test_data {
		res := binary.BigEndian.Uint16(answ_data[2*i : 2+i*2])
		if v != res {
			t.Error("Expected ", v, "got ", res)
		}
	}
}

func TestModbusServer_PresetMultipleRegisters(t *testing.T) {
	test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
	md := &ModbusData{holding_reg: make([]uint16, 10), mu_holding_regs: &sync.Mutex{}}
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcPresetMultipleRegisters, 0, 0x5, wordArrToByteArr(test_data)...)
	srv.PresetMultipleRegisters(req)
	for i, v := range test_data {
		if md.holding_reg[i] != v {
			t.Error("Expected ", v, "got ", md.holding_reg[i])
		}
	}
}

func TestModbusServer_ReadInputRegisters(t *testing.T) {
	test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
	md := &ModbusData{input_reg: make([]uint16, 10)}
	md.PresetMultipleInputsRegisters(0, test_data...)
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcReadInputRegisters, 0, 0x5)
	answ, _ := srv.ReadInputRegisters(req)
	_, answ_data := answ.GetData()
	for i, v := range test_data {
		res := binary.BigEndian.Uint16(answ_data[2*i : 2+i*2])
		if v != res {
			t.Error("Expected ", v, "got ", res)
		}
	}
}

func TestModbusServer_ReadCoilStatus(t *testing.T) {
	test_data := []bool{true, true, false, false, true}
	md := &ModbusData{coils: make([]bool, 10), mu_coils: &sync.Mutex{}}
	md.ForceMultipleCoils(0, test_data...)
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcReadCoilStatus, 0, 0x5)
	answ, _ := srv.ReadCoilStatus(req)
	_, answ_data := answ.GetData()
	bool_arr := byteArrToBoolArr(answ_data, byte(len(test_data)))
	for i, v := range test_data {
		if bool_arr[i] != v {
			t.Error("Expected ", v, "got ", bool_arr[i])
		}
	}
}

func TestModbusServer_TestForceMultipleCoils(t *testing.T) {
	test_data := []bool{true, true, false, false, true}
	md := &ModbusData{coils: make([]bool, 10), mu_coils: &sync.Mutex{}}
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcForceMultipleCoils, 0, 0x5, boolArrToByteArr(test_data)...)
	srv.ForceMultipleCoils(req)
	for i, v := range test_data {
		if md.coils[i] != v {
			t.Error("Expected ", v, "got ", md.coils[i])
		}
	}
}

func TestModbusServer_ReadDescreteInputs(t *testing.T) {
	test_data := []bool{true, true, false, false, true}
	md := &ModbusData{discrete_inputs: make([]bool, 10)}
	md.ForceMultipleDescreteInputs(0, test_data...)
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcReadDescreteInputs, 0, 0x5)
	answ, _ := srv.ReadDescreteInputs(req)
	_, answ_data := answ.GetData()
	bool_arr := byteArrToBoolArr(answ_data, byte(len(test_data)))
	for i, v := range test_data {
		if bool_arr[i] != v {
			t.Error("Expected ", v, "got ", bool_arr[i])
		}
	}
}

func TestModbusServer_TestPresetSingleRegister(t *testing.T) {
	test_data := uint16(0x01)
	test_addr := uint16(0x0)
	md := &ModbusData{holding_reg: make([]uint16, 10), mu_holding_regs: &sync.Mutex{}}
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcPresetSingleRegister, test_addr, test_data)
	srv.PresetSingleRegister(req)
	if md.holding_reg[test_addr] != test_data {
		t.Error("Expected ", test_data, "got ", md.holding_reg[test_addr])
	}
}

func TestModbusServer_TestForceSingleCoil(t *testing.T) {
	test_data := uint16(0xFF00)
	test_addr := uint16(0x0)
	md := &ModbusData{coils: make([]bool, 10), mu_coils: &sync.Mutex{}}
	srv := &ModbusServer{}
	srv.Data = md
	req := buildRequest(0, ModbusRTUviaTCP, 1, FcForceSingleCoil, test_addr, test_data)
	srv.ForceSingleCoil(req)
	if md.coils[test_addr] != true {
		t.Error("Expected ", true, "got ", md.coils[test_addr])
	}
}
