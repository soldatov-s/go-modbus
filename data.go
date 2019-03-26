// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"errors"
	"fmt"
	"sync"
)

type ModbusDataType int

// Data types enumeration
const (
	Coils ModbusDataType = iota
	DiscreteInputs
	HoldingRegisters
	InputRegisters
)

// ModbusData implements data interface
type ModbusData struct {
	coils, discrete_inputs    []bool
	holding_reg, input_reg    []uint16
	mu_holding_regs, mu_coils *sync.Mutex
}

// Checks that requested data is not outside the present range
func (md *ModbusData) checkOutside(dataType ModbusDataType, addr, cnt uint16) (bool, error) {
	var (
		err error
		l   uint16
	)
	res := true
	switch dataType {
	case HoldingRegisters:
		l = len(md.holding_reg)
	case DiscreteInputs:
		l = len(md.discrete_inputs)
	case InputRegisters:
		l = len(md.input_reg)
	case Coils:
		l = len(md.coils)
	default:
		l = 0
	}

	if addr+cnt > uint16(l) {
		err_str := fmt.Sprintf("Requested data %d...%d outside the valid range 0...%d", addr, addr+cnt, l)
		err = errors.New(err_str)
		res = false
	}

	return res, err
}

// Initializate new instance of ModbusData
func (md *ModbusData) Init(coils_cnt, discrete_inputs_cnt, holding_reg_cnt, input_reg_cnt int) error {
	md.mu_holding_regs = &sync.Mutex{}
	md.mu_coils = &sync.Mutex{}
	md.coils = make([]bool, coils_cnt)
	md.discrete_inputs = make([]bool, discrete_inputs_cnt)
	md.holding_reg = make([]uint16, holding_reg_cnt)
	md.input_reg = make([]uint16, input_reg_cnt)

	return nil
}

// Preset Single Register
func (md *ModbusData) PresetSingleRegister(addr uint16, data uint16) error {
	cnt := uint16(1)
	_, err := md.checkOutside(HoldingRegisters, addr, cnt)
	md.mu_holding_regs.Lock()
	defer md.mu_holding_regs.Unlock()
	md.holding_reg[addr] = data
	return err
}

// Set Preset Multiple Registers
func (md *ModbusData) PresetMultipleRegisters(addr uint16, data []uint16) error {
	cnt := uint16(len(data))
	_, err := md.checkOutside(HoldingRegisters, addr, cnt)
	md.mu_holding_regs.Lock()
	defer md.mu_holding_regs.Unlock()
	copy(md.holding_reg[addr:addr+cnt], data)
	return err
}

// Read Holding Registers
func (md *ModbusData) ReadHoldingRegisters(addr, cnt uint16) ([]uint16, error) {
	_, err := md.checkOutside(HoldingRegisters, addr, cnt)
	if err != nil {
		return nil, err
	}
	return md.holding_reg[addr : addr+cnt], err
}

// Read Input Registers
func (md *ModbusData) ReadInputRegisters(addr, cnt uint16) ([]uint16, error) {
	_, err := md.checkOutside(InputRegisters, addr, cnt)
	if err != nil {
		return nil, err
	}
	return md.input_reg[addr : addr+cnt], err
}

// Read Coil Status
func (md *ModbusData) ReadCoilStatus(addr, cnt uint16) ([]bool, error) {
	_, err := md.checkOutside(Coils, addr, cnt)
	return md.coils[addr : addr+cnt], err
}

// Force Single Coil
func (md *ModbusData) ForceSingleCoil(addr uint16, data bool) error {
	cnt := uint16(1)
	_, err := md.checkOutside(Coils, addr, cnt)
	md.mu_coils.Lock()
	defer md.mu_coils.Unlock()
	md.coils[addr] = data
	return err
}

// Force Multiple Coils
func (md *ModbusData) ForceMultipleCoils(addr uint16, data []bool) error {
	cnt := uint16(len(data))
	_, err := md.checkOutside(Coils, addr, cnt)
	md.mu_coils.Lock()
	defer md.mu_coils.Unlock()
	copy(md.coils[addr:addr+cnt], data)
	return err
}

// Read Descrete Inputs
func (md *ModbusData) ReadDescreteInputs(addr, cnt uint16) ([]bool, error) {
	_, err := md.checkOutside(DiscreteInputs, addr, cnt)
	return md.coils[addr : addr+cnt], err
}
