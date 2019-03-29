// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"fmt"
	"strings"
	"testing"
)

func TestModbusBaseServer_String(t *testing.T) {
	test_str := "test:504"
	srv := &ModbusBaseServer{Host: "test", Port: "504"}
	s := fmt.Sprintf("%s", srv)
	r := strings.Compare(s, test_str)
	if r != 0 {
		t.Error("Expected ", test_str, "got ", r)
	}
}

func TestModbusBaseClient_String(t *testing.T) {
	test_str := "test:504"
	srv := &ModbusBaseClient{Host: "test", Port: "504"}
	s := fmt.Sprintf("%s", srv)
	r := strings.Compare(s, test_str)
	if r != 0 {
		t.Error("Expected ", test_str, "got ", r)
	}
}
