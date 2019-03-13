// data
package main

import (
	"encoding/binary"
	"sync"
)

type ModbusData struct {
	holding_reg_cnt, input_reg_cnt int
	holding_reg, input_reg         []byte
	mu                             *sync.Mutex
}

func (md *ModbusData) init() error {
	md.mu = &sync.Mutex{}
	md.holding_reg = make([]byte, md.holding_reg_cnt*2)
	md.input_reg = make([]byte, md.input_reg_cnt*2)

	return nil
}

func (md *ModbusData) PresetMultipleRegs(mp *ModbusPacket) []byte {
	md.mu.Lock()
	defer md.mu.Unlock()

	var answer []byte
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	// Set values in ModbusData
	copy(md.holding_reg[addr*2:(addr+cnt)*2], mp.data[7:7+cnt*2])
	// Copy addr and code function
	answer = append(answer, mp.GetPrefix()...)
	//
	answer = append(answer, mp.data[2:6]...)
	// Crc Answer
	AppendCrc16(&answer)

	md.mu.Unlock()

	return answer
}

func (md *ModbusData) ReadHoldRegs(mp *ModbusPacket) []byte {
	var answer []byte
	addr := binary.BigEndian.Uint16(mp.data[2:4])
	cnt := binary.BigEndian.Uint16(mp.data[4:6])

	// Copy addr and code function
	answer = append(answer, mp.GetPrefix()...)
	// Answer length in byte
	answer = append(answer, byte(cnt*2))
	// Data for answer
	answer = append(answer, md.holding_reg[addr*2:(addr+cnt)*2]...)
	// Crc Answer
	AppendCrc16(&answer)
	return answer
}
