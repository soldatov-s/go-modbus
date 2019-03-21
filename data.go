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

const (
	Coils ModbusDataType = iota
	DiscreteInputs
	HoldingRegisters
	InputRegisters
)

type ModbusData struct {
	coils, discrete_inputs    []bool
	holding_reg, input_reg    []uint16
	mu_holding_regs, mu_coils *sync.Mutex
}

func (md *ModbusData) checkOutside(dataType ModbusDataType, addr, cnt uint16) error {
	var (
		err error
		l   uint16
	)
	switch dataType {
	case HoldingRegisters:
		l = uint16(len(md.holding_reg))
		if addr > l || addr+cnt > l {
			err_str := fmt.Sprintf("Requested data %d...%d outside the valid range 0...%d", addr, addr+cnt, l)
			err = errors.New(err_str)
		}
	default:
		l = 0
	}

	return err
}

func (md *ModbusData) Init(coils_cnt, discrete_inputs_cnt, holding_reg_cnt, input_reg_cnt int) error {
	md.mu_holding_regs = &sync.Mutex{}
	md.mu_coils = &sync.Mutex{}
	md.coils = make([]bool, coils_cnt)
	md.discrete_inputs = make([]bool, discrete_inputs_cnt)
	md.holding_reg = make([]uint16, holding_reg_cnt)
	md.input_reg = make([]uint16, input_reg_cnt)

	return nil
}

func (md *ModbusData) PresetMultipleRegisters(addr, cnt uint16, data []uint16) error {
	err := md.checkOutside(HoldingRegisters, addr, cnt)
	md.mu_holding_regs.Lock()
	defer md.mu_holding_regs.Unlock()
	copy(md.holding_reg[addr:addr+cnt], data)
	return err
}

func (md *ModbusData) ReadHoldingRegisters(addr, cnt uint16) ([]uint16, error) {
	err := md.checkOutside(HoldingRegisters, addr, cnt)
	if err != nil {
		return nil, err
	}
	return md.holding_reg[addr : addr+cnt], err
}

func (md *ModbusData) ReadInputRegisters(addr, cnt uint16) ([]uint16, error) {
	err := md.checkOutside(InputRegisters, addr, cnt)
	if err != nil {
		return nil, err
	}
	return md.input_reg[addr : addr+cnt], err
}

func (md *ModbusData) ReadCoilStatus(addr, cnt uint16) ([]bool, error) {

	err := md.checkOutside(Coils, addr, cnt)
	return md.coils[addr : addr+cnt], err
}

func (md *ModbusData) ForceMultipleCoils(addr, cnt uint16, data []bool) error {
	md.mu_coils.Lock()
	defer md.mu_coils.Unlock()
	err := md.checkOutside(Coils, addr, cnt)
	copy(md.coils[addr:addr+cnt], data)
	return err
}

func (md *ModbusData) ReadDescreteInputs(addr, cnt uint16) ([]bool, error) {
	err := md.checkOutside(DiscreteInputs, addr, cnt)
	return md.coils[addr : addr+cnt], err
}
