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
  req := buildPacket(ModbusRTUviaTCP, 1, ModbusFunctionCode, 0, 0x5)
  answ := ReadHoldingRegistersHndl(req, md)
  for i, v := range test_data {
    res := binary.BigEndian.Uint16(mp.Data[2+i:4+i])
    if v != res {
      t.Error("Expected ", v, "got ", res)
    }
  }
}
