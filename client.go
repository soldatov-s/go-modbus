// client
package modbus

/*func (mb *ModbusApp) read() ([]byte, int) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 260)
	// Read the incoming connection into the buffer.
	reqLen, err := mb.conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	fmt.Printf("Src->: \t\t\t\t%s\n", mb.conn.RemoteAddr().String())
	fmt.Printf("Dst<-: \t\t\t\t%s\n", mb.conn.LocalAddr().String())

	return buf, reqLen
}*/

/*func (mb *ModbusApp) write() {
	// Send a response back to person contacting us.
	mb.conn.Write([]byte("Message received."))
}

func (mb *ModbusApp) sendAnswer(answer []byte) {
	mb.conn.Write(answer)
}*/

/*func (mb *ModbusApp) client() {
	var err error

	fmt.Println("Client mode")

	mb.conn, err = net.Dial(mb.protocol, mb.address())
	if err != nil {
		fmt.Println("Error connect:", err.Error())
		os.Exit(1)
	}
	mb.read()
	mb.write()
	// Close the connection when you're done with it.
	mb.conn.Close()
}*/

/*func (mp *ModbusPacket) HexStrToData(str string) {
	data, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	mp.data = make([]byte, 0, len(data))
	mp.length = len(data)
	copy(data, mp.data)
}*/
