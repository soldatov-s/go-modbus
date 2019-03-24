// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"testing"
)

func TestCrc16(t *testing.T) {
	data := []byte{0x1, 0x2, 0x3, 0x4, 0x5}
	crc := Crc16(data)
	if crc != 0xBB2A {
		t.Error("Expected BB2A, got ", crc)
	}
}
