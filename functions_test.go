// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"fmt"
	"sync"
	"testing"
)

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
	fmt.Println(req)
	answ, _ := ReadDescreteInputsHndl(req, md)
	cnt_byte := uint16(answ.Data[2])
	bool_arr := byteArrToBoolArr(answ.Data[3:3+cnt_byte], uint16(len(test_data)))
	for i, v := range test_data {
		if bool_arr[i] != v {
			t.Error("Expected ", v, "got ", bool_arr[i])
		}
	}
}
