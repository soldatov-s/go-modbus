// server
package modbus

import (
	"fmt"
	"math"
	"net"
)

type ModbusServer struct {
	host, protocol, port string
	mbprotocol           ModbusTypeProtocol
	Data                 *ModbusData
	ln                   net.Listener
	done                 chan struct{}
}

func NewServer(host, protocol, port string,
	mbprotocol ModbusTypeProtocol,
	coils_cnt, discrete_inputs_cnt,
	holding_reg_cnt, input_reg_cnt int) *ModbusServer {

	srv := new(ModbusServer)
	srv.host = host
	srv.protocol = protocol
	srv.port = port
	srv.mbprotocol = mbprotocol
	srv.Data = new(ModbusData)
	srv.Data.Init(coils_cnt, discrete_inputs_cnt, holding_reg_cnt, input_reg_cnt)
	srv.done = make(chan struct{})

	return srv
}

func (srv *ModbusServer) String() string {
	return srv.host + ":" + srv.port
}

func (srv *ModbusServer) ServerStop() error {
	var err error

	fmt.Println("Server stop")
	close(srv.done)
	return err
}

func (srv *ModbusServer) ServerStart() error {
	var err error

	fmt.Println("Server start")
	fmt.Println("Listening at", srv.String())

	// Listen for incoming connections.
	srv.ln, err = net.Listen(srv.protocol, srv.String())
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}

	defer srv.ln.Close()

	for {
		select {
		case <-srv.done:
			fmt.Println("Stop listen incoming connection")
			return nil
		default:
			// Listen for an incoming connection.
			conn, err := srv.ln.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				conn.Close()
				continue
			}
			fmt.Printf(
				"Src->: \t\t\t\t%s\nDst<-: \t\t\t\t%s\n",
				conn.RemoteAddr(),
				conn.LocalAddr())
			// Handle connections in a new goroutine.
			go srv.handleRequest(conn, srv.done)
		}
	}
}

// Handles incoming requests.
func (srv *ModbusServer) handleRequest(conn net.Conn, done <-chan struct{}) error {
	var (
		id_packet int
		err       error
		request   *ModbusPacket = &ModbusPacket{
			mtp: srv.mbprotocol}
	)
	// Close the connection when you're done with it.
	defer conn.Close()
	request.Init()

	// Read the incoming connection into the buffer.
	for {
		select {
		case <-done:
			fmt.Println("Close connection", conn.RemoteAddr())
			conn.Close()
			return nil
		default:
			request.length, err = conn.Read(request.data)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				break
			}

			if request.length == 0 {
				continue
			}

			id_packet++
			if id_packet == math.MaxInt32 {
				id_packet = 0
			}
			fmt.Println("Packet ID:", id_packet)

			// fmt.Println(hex.Dump(request.data))
			request.ModbusDump()
			var answer *ModbusPacket
			answer, err = request.HandlerRequest(srv.Data)

			answer.ModbusDump()
			conn.Write(answer.data)
		}
	}
	return err
}
