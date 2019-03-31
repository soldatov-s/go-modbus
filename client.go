// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"fmt"
	"log"
	"math"
	"net"
)

// ModbusClient implements client interface
type ModbusClient struct {
	ModbusBaseClient
	DevID         byte
	TypeProtocol  ModbusTypeProtocol // Type Modbus Protocol
	Conn          net.Conn           // Connection
	TranscationId uint16             // for ModbusTCP
}

// NewClient function initializate new instance of ModbusClient
func NewClient(port, host string, mbprotocol ModbusTypeProtocol, devID byte) (*ModbusClient, error) {
	var err error
	mc := &ModbusClient{TypeProtocol: mbprotocol, DevID: devID}
	mc.Host = host
	mc.Port = port
	mc.Conn, err = net.Dial("tcp", mc.String())
	return mc, err
}

// Read Answer from Slave device (Server)
func (mc *ModbusClient) ReadAnswer() (*ModbusPacket, error) {
	var err error
	answer := &ModbusPacket{isAnswer: true}
	answer.Init(mc.TypeProtocol)

	log.Printf(
		"Src->: %s Dst<-: %s\n",
		mc.Conn.RemoteAddr(),
		mc.Conn.LocalAddr())

	// Read the incoming connection into the buffer.
	answer.Length, err = mc.Conn.Read(answer.PDU)
	if err != nil {
		log.Println("Error reading:", err.Error())
	}

	return answer, err
}

// Send Request to Slave device (Server) and return Answer from it
func (mc *ModbusClient) SendRequest(mp *ModbusPacket) (*ModbusPacket, error) {
	var err error

	log.Println("Send request to", mc)
	_, err = mc.Conn.Write(mp.PDU[:mp.GetPDULength()])
	if err != nil {
		log.Println("Error connect:", err.Error())
		return nil, err
	}

	return mc.ReadAnswer()
}

// Send Request ReadHoldingRegisters
func (mc *ModbusClient) ReadHoldingRegisters(addr, cnt uint16) ([]uint16, error) {
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcReadHoldingRegisters, addr, cnt)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return nil, err
	}
	_, data := answer.GetData()
	return byteArrToWordArr(data), nil
}

// Send Request ReadInputRegisters
func (mc *ModbusClient) ReadInputRegisters(addr, cnt uint16) ([]uint16, error) {
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcReadInputRegisters, addr, cnt)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return nil, err
	}
	_, data := answer.GetData()
	return byteArrToWordArr(data), nil
}

// Send Request PresetSingleRegister
func (mc *ModbusClient) PresetSingleRegister(addr, value uint16) error {
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcReadInputRegisters, addr, value)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return err
	}
	if byte(answer.GetFunctionCode())&0x80 > 0 {
		return fmt.Errorf("Error code %x", answer.GetErrorCode())
	}
	return nil
}

// Send Request ReadCoilStatus
func (mc *ModbusClient) ReadCoilStatus(addr, cnt uint16) ([]bool, error) {
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcReadCoilStatus, addr, cnt)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return nil, err
	}
	_, data := answer.GetData()
	return byteArrToBoolArr(data, byte(cnt)), nil

}

// Send Request ReadCoilStatus
func (mc *ModbusClient) ReadDescreteInputs(addr, cnt uint16) ([]bool, error) {
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcReadDescreteInputs, addr, cnt)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return nil, err
	}
	_, data := answer.GetData()
	return byteArrToBoolArr(data, byte(cnt)), nil
}

// Send Request ForceSingleCoil
func (mc *ModbusClient) ForceSingleCoil(addr uint16, value bool) error {
	v := uint16(0)
	if value {
		v = 0xFF00
	}
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcForceSingleCoil, addr, v)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return err
	}
	if byte(answer.GetFunctionCode())&0x80 > 0 {
		return fmt.Errorf("Error code %x", answer.GetErrorCode())
	}
	return nil

}

// Send Request PresetMultipleRegisters
func (mc *ModbusClient) PresetMultipleRegisters(addr, cnt uint16, data ...uint16) error {
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcPresetMultipleRegisters, addr, cnt, wordArrToByteArr(data)...)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return err
	}
	if byte(answer.GetFunctionCode())&0x80 > 0 {
		return fmt.Errorf("Error code %x", answer.GetErrorCode())
	}
	return nil
}

// Send Request ForceMultipleCoils
func (mc *ModbusClient) ForceMultipleCoils(addr, cnt uint16, data ...bool) error {
	request := buildRequest(mc.GetTransactionId(), mc.TypeProtocol, mc.DevID, FcForceMultipleCoils, addr, cnt, boolArrToByteArr(data)...)
	answer, err := mc.SendRequest(request)
	if err != nil {
		return err
	}
	if byte(answer.GetFunctionCode())&0x80 > 0 {
		return fmt.Errorf("Error code %x", answer.GetErrorCode())
	}
	return nil
}

// Close client
func (mc *ModbusClient) Close() {
	// Close the connection when you're done with it.
	mc.Conn.Close()
}

// Get Transaction ID
func (mc *ModbusClient) GetTransactionId() uint16 {
	mc.TranscationId++
	if mc.TranscationId == math.MaxUint16 {
		mc.TranscationId = 0
	}
	return mc.TranscationId
}
