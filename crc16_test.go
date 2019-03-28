// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"fmt"
	"testing"
)

type testCrcPair struct {
	data []byte
	crc  uint16
}

var testsCrc = []testCrcPair{
	{[]byte{0x1, 0x2, 0x3, 0x4, 0x5}, 0xBB2A},
	{[]byte{0x10, 0x20, 0x30, 0x40, 0x50}, 0xF0DF},
	{nil, 0},
}

func TestCrc16(t *testing.T) {
	for _, pair := range testsCrc {
		crc := Crc16(pair.data)
		if crc != pair.crc {
			t.Error(
				"For", pair.data,
				"expected", pair.crc,
				"got", crc,
			)
		}

	}
}

func BenchmarkCrc16(b *testing.B) {
    for i := 0; i < b.N; i++ {
        crc := Crc16([]byte{0x1, 0x2, 0x3, 0x4, 0x5})
    }
}

func ExampleCrc16() {
	fmt.Printf("0x%X", Crc16([]byte{0x1, 0x2, 0x3, 0x4, 0x5}))
	// Output: 0xBB2A
}

func TestCrc16Check(t *testing.T) {
	data := []byte{0x1, 0x2, 0x3, 0x4, 0x5}
	r := Crc16Check(data, 0xBB2A)
	if !r {
		t.Error("Expected true, got ", r)
	}
}
