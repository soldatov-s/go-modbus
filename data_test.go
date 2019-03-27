// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
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
