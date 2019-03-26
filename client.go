// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbus

import (
	"log"
	"net"
)

// ModbusClient implements client interface
type ModbusClient struct {
	ModbusBaseClient
	MTP  ModbusTypeProtocol // Type Modbus Protocol
	Conn net.Conn           // Connection
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
		"Src->: \t\t\t\t%s\nDst<-: \t\t\t\t%s\n",
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
	var (
		answer *ModbusPacket
		err    error
	)

	log.Println("Send request to", mc)
	_, err = mc.Conn.Write(mp.Data)
	if err != nil {
		log.Println("Error connect:", err.Error())
		return nil, err
	}

	return  mc.ReadAnswer(mp)
}

// Send Request ReadHoldingRegisters
func (mc *ModbusClient) ReadHoldingRegisters(addr, cnt uint16) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	request.Data = make([]byte, 0, 8)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, ReadHoldingRegisters)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], cnt)
	endAnswer(request)

	return mc.SendRequest(request)
}

// Send Request ReadInputRegisters
func (mc *ModbusClient) ReadInputRegisters(addr, cnt uint16) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	request.Data = make([]byte, 0, 8)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, ReadInputRegisters)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], cnt)
	endAnswer(request)

	return mc.SendRequest(request)
}

// Send Request PresetSingleRegister
func (mc *ModbusClient) PresetSingleRegister(addr, value uint16) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	request.Data = make([]byte, 0, 8)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, PresetSingleRegister)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], value)
	endAnswer(request)

	return mc.SendRequest(request)
}

// Send Request ReadCoilStatus
func (mc *ModbusClient) ReadCoilStatus(addr, cnt uint16) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	request.Data = make([]byte, 0, 8)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, ReadCoilStatus)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], cnt)
	endAnswer(request)

	return mc.SendRequest(request)
}

// Send Request ReadCoilStatus
func (mc *ModbusClient) ReadDescreteInputs(addr, cnt uint16) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	request.Data = make([]byte, 0, 8)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, ReadDescreteInputs)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], cnt)
	endAnswer(answer)

	return mc.SendRequest(request)
}

// Send Request ForceSingleCoil
func (mc *ModbusClient) ForceSingleCoil(addr, value uint16) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	request.Data = make([]byte, 0, 8)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, ForceSingleCoil)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], cnt)
	endAnswer(answer)

	return mc.SendRequest(request)
}

// Send Request PresetMultipleRegisters
func (mc *ModbusClient) PresetMultipleRegisters(addr, value uint16, data []uint16...) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	request.Data = make([]byte, 0, 8 + len(data) * 2 + 1)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, PresetMultipleRegisters)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], cnt)
	request.Data = append(request.Data, byte(len(data) * 2 + 1))
	request.Data = append(request.Data, wordArrToByteArr(data)...)
	endAnswer(answer)

	return mc.SendRequest(request)
}

// Send Request ForceMultipleCoils
func (mc *ModbusClient) ForceMultipleCoils(addr, value uint16, data []bool...) (*ModbusPacket, error) {
	var (
		answer  *ModbusPacket
		err     error
		request *ModbusPacket
	)

	request.TypeProtocol = mc.MTP
	q, r := len(data)/8, len(data)%8
	if r > 0 {
		q++
	}
	request.Data = make([]byte, 0, 8 + q + 1)
	// Copy addr and code function
	request.Data = append(answer.Data, 0x1, ForceMultipleCoils)
	binary.BigEndian.PutUint16(request.Data[2:4], addr)
	binary.BigEndian.PutUint16(request.Data[4:6], cnt)
	request.Data = append(request.Data, byte(len(data) * 2 + 1))
	request.Data = append(request.Data, boolArrToByteArr(data)...)
	endAnswer(answer)

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
