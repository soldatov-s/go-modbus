// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"log"
	"net"
)

// ModbusClient implements client interface
type ModbusClient struct {
	ModbusBaseClient
	DevID byte
	MTP   ModbusTypeProtocol // Type Modbus Protocol
	Conn  net.Conn           // Connection
}

// NewClient function initializate new instance of ModbusClient
func NewClient(port, host string, mbprotocol ModbusTypeProtocol) (*ModbusClient, error) {
	var err error
	mc := new(ModbusClient)
	mc.Host = host
	mc.Port = port

	mc.Conn, err = net.Dial("tcp", mc.String())

	return mc, err
}

// Read Answer from Slave device (Server)
func (mc *ModbusClient) ReadAnswer() (*ModbusPacket, error) {
	var err error
	answer := new(ModbusPacket)
	answer.TypeProtocol = mc.MTP
	answer.Init()

	log.Printf(
		"Src->: %s\ Dst<-: %s\n",
		mc.Conn.RemoteAddr(),
		mc.Conn.LocalAddr())

	// Read the incoming connection into the buffer.
	_, err = mc.Conn.Read(answer.Data)
	if err != nil {
		log.Println("Error reading:", err.Error())
	}

	return answer, err
}

// Send Request to Slave device (Server) and return Answer from it
func (mc *ModbusClient) SendRequest(mp *ModbusPacket) (*ModbusPacket, error) {
	var err error

	log.Println("Send request to", mc)
	_, err = mc.Conn.Write(mp.Data)
	if err != nil {
		log.Println("Error connect:", err.Error())
		return nil, err
	}

	return mc.ReadAnswer()
}

// Send Request ReadHoldingRegisters
func (mc *ModbusClient) ReadHoldingRegisters(addr, cnt uint16) (*ModbusPacket, error) {
	request := buildPacket(mc.MTP, mc.DevID, ReadHoldingRegisters, addr, cnt)
	return mc.SendRequest(request)
}

// Send Request ReadInputRegisters
func (mc *ModbusClient) ReadInputRegisters(addr, cnt uint16) (*ModbusPacket, error) {
	request := buildPacket(mc.MTP, mc.DevID, ReadInputRegisters, addr, cnt)
	return mc.SendRequest(request)
}

// Send Request PresetSingleRegister
func (mc *ModbusClient) PresetSingleRegister(addr, value uint16) (*ModbusPacket, error) {
	request := buildPacket(mc.MTP, mc.DevID, ReadInputRegisters, addr, value)
	return mc.SendRequest(request)
}

// Send Request ReadCoilStatus
func (mc *ModbusClient) ReadCoilStatus(addr, cnt uint16) (*ModbusPacket, error) {
	request := buildPacket(mc.MTP, mc.DevID, ReadCoilStatus, addr, cnt)
	return mc.SendRequest(request)
}

// Send Request ReadCoilStatus
func (mc *ModbusClient) ReadDescreteInputs(addr, cnt uint16) (*ModbusPacket, error) {
	request := buildPacket(mc.MTP, mc.DevID, ReadDescreteInputs, addr, cnt)
	return mc.SendRequest(request)
}

// Send Request ForceSingleCoil
func (mc *ModbusClient) ForceSingleCoil(addr uint16, value bool) (*ModbusPacket, error) {
	v := uint16(0)
	if value {
		v = 0xFF00
	}
	request := buildPacket(mc.MTP, mc.DevID, ForceSingleCoil, addr, v)
	return mc.SendRequest(request)
}

// Send Request PresetMultipleRegisters
func (mc *ModbusClient) PresetMultipleRegisters(addr, cnt uint16, data ...uint16) (*ModbusPacket, error) {
	request := buildPacket(mc.MTP, mc.DevID, PresetMultipleRegisters, addr, cnt, wordArrToByteArr(data)...)
	return mc.SendRequest(request)
}

// Send Request ForceMultipleCoils
func (mc *ModbusClient) ForceMultipleCoils(addr, cnt uint16, data ...bool) (*ModbusPacket, error) {
	request := buildPacket(mc.MTP, mc.DevID, ForceMultipleCoils, addr, cnt, boolArrToByteArr(data)...)
	return mc.SendRequest(request)
}

func (mc *ModbusClient) Close() {
	// Close the connection when you're done with it.
	mc.Conn.Close()
}

/*func (mp *ModbusPacket) HexStrToData(str string) {
	data, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	mp.data = make([]byte, 0, len(data))
	mp.length = len(data)
	copy(data, mp.data)
}*/
