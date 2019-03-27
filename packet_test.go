// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"testing"
	"fmt"
)

type testCrc16Checkpair struct {
	TypeProtocol ModbusTypeProtocol
	data         []byte
	res          bool
}

var testsCrc16Check = []testCrc16Checkpair{
	{ModbusRTUviaTCP, []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0xA, 0xCD, 0xC5}, true},
}

func TestModbusPacket_Crc16Check(t *testing.T) {
	mp := new(ModbusPacket)

	for _, pair := range testsCrc16Check {
		mp.TypeProtocol = pair.TypeProtocol
		mp.Init()
		mp.Data = pair.data
		mp.Length = len(pair.data)
		fmt.Println(mp.GetCrc())
		res := mp.Crc16Check()
		if res != pair.res {
			t.Error(
				"For", pair.data,
				"expected", pair.res,
				"got", res,
			)

		}
	}
}
