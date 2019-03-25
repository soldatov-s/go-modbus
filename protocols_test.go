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
