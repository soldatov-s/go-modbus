// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"testing"
)

func TestReadHoldingRegistersHndl(){
  test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
  md := &ModbusData(md.holding_reg: make([]uint16, 10))
  md.PresetMultipleRegisters(0, test_data)
  req := buildPacket(ModbusRTUviaTCP, 1, ReadHoldingRegisters, 0, 0x5)
  answ := ReadHoldingRegistersHndl(req, md)
  for i, v := range test_data {
    res := binary.BigEndian.Uint16(mp.Data[3+i:5+i])
    if v != res {
      t.Error("Expected ", v, "got ", res)
    }
  }
}

func TestPresetMultipleRegisters(){
  test_data := []uint16{00x01, 0x02, 0x03, 0x04, 0x05}
  md := &ModbusData(md.holding_reg: make([]uint16, 10))
  req := buildPacket(ModbusRTUviaTCP, 1, PresetMultipleRegisters, 0, 0x5, wordArrToByteArr(test_data))
  answ := PresetMultipleRegisters(req, md)
  for i, v := range test_data {
    if md.holding_reg[i] != v {
      t.Error("Expected ", v, "got ", md.holding_reg[i])
    }
  }
}

func TestReadInputRegistersHndl(){
  test_data := []uint16{0x01, 0x02, 0x03, 0x04, 0x05}
  md := &ModbusData(md.input_reg: make([]uint16, 10))
  md.PresetMultipleInputsRegisters(0, test_data)
  req := buildPacket(ModbusRTUviaTCP, 1, ReadInputRegisters, 0, 0x5)
  answ := ReadInputRegistersHndl(req, md)
  for i, v := range test_data {
    res := binary.BigEndian.Uint16(mp.Data[3+i:5+i])
    if v != res {
      t.Error("Expected ", v, "got ", res)
    }
  }
}

func TestReadCoilStatusHndll(){
  test_data := []bool{true, true, false, false, true}
  md := &ModbusData(md.coils: make([]bool, 10))
  md.ForceMultipleCoils(0, test_data)
  req := buildPacket(ModbusRTUviaTCP, 1, ReadCoilStatus, 0, 0x5)
  answ := ReadCoilStatusHndl(req, md)
  cnt_byte := mp.Data[2]
  bool_arr := byteArrToBoolArr(mp.Data[3:3+cnt_byte])
  for i, v := range test_data {
    if bool_arr[i] != v {
      t.Error("Expected ", v, "got ", bool_arr[i])
    }
  }
}

func TestForceMultipleCoilsHndl(){
  test_data := []bool{true, true, false, false, true}
  md := &ModbusData(md.holding_reg: make([]uint16, 10))
  req := buildPacket(ModbusRTUviaTCP, 1, ForceMultipleCoils, 0, 0x5, boolArrToByteArr(test_data))
  answ := ForceMultipleCoilsHndl(req, md)
  for i, v := range test_data {
    if md.coils[i] != v {
      t.Error("Expected ", v, "got ", md.coils[i])
    }
  }
}

func TestReadDescreteInputsHndl(){
  test_data := []bool{true, true, false, false, true}
  md := &ModbusData(md.discrete_inputs: make([]bool, 10))
  md.ForceMultipleDescreteInputs(0, test_data)
  req := buildPacket(ModbusRTUviaTCP, 1, ReadDescreteInputs, 0, 0x5)
  answ := ReadDescreteInputsHndl(req, md)
  cnt_byte := mp.Data[2]
  bool_arr := byteArrToBoolArr(mp.Data[3:3+cnt_byte])
  for i, v := range test_data {
    if bool_arr[i] != v {
      t.Error("Expected ", v, "got ", bool_arr[i])
    }
  }
}
