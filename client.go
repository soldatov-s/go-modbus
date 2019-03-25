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

	answer, err = mc.ReadAnswer()
	return answer, err
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
