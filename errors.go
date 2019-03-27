// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"errors"
)

// Modbus Error code
type ModbusErrors int

const (
	ErrCantHandel ModbusErrors = 1
	ErrOutside    ModbusErrors = 2
	ErrBadVal     ModbusErrors = 3
)

func (e ModbusErrors) Error() error {
	switch e {
	case ErrCantHandel:
		return errors.New("Can't handel request")
	case ErrOutside:
		return errors.New("Requested outside valid range")
	case ErrBadVal:
		return errors.New("Bad value in request")
	default:
		return errors.New("Unknown Error")
	}
}
