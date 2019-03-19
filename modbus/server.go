// server
package modbus

import (
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
)

type ModbusServer struct {
	host, protocol, port string
	mbprotocol           ModbusTypeProtocol
	Data                 *ModbusData
	ln                   net.Listener
	done                 chan struct{}
	exited               chan struct{}
	wg                   sync.WaitGroup
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
	srv.exited = make(chan struct{})

	return srv
}

func (srv *ModbusServer) String() string {
	return srv.host + ":" + srv.port
}

func (srv *ModbusServer) Stop() error {
	var err error

	fmt.Println("Shutting down server...")
	srv.ln.Close()
	close(srv.done)
	srv.wg.Wait()
	fmt.Println("Server is stopped")
	return err
}

func (srv *ModbusServer) Start() error {
	var err error

	fmt.Println("Server startup...")
	fmt.Println("Listening at", srv.String())

	// Listen for incoming connections.
	srv.ln, err = net.Listen(srv.protocol, srv.String())
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	srv.wg.Add(1)
	go func() {
		defer func() {
			fmt.Println("Stop listen incoming connection")
			srv.ln.Close()
			srv.wg.Done()
		}()
		for {
			select {
			case <-srv.done:
				return
			default:
				// Listen for an incoming connection.
				conn, err := srv.ln.Accept()
				if err != nil {
					if strings.Contains(err.Error(), "use of closed network connection") {
						return
					}
					fmt.Println("Error accepting: ", err.Error())
					continue
				}
				// Handle connections in a new goroutine.
				go srv.handleRequest(conn)
			}
		}
	}()

	fmt.Println("Server is started")
	return err
}

// Handles incoming requests.
func (srv *ModbusServer) handleRequest(conn net.Conn) error {
	var (
		id_packet int
		err       error
		request   *ModbusPacket = &ModbusPacket{
			mtp: srv.mbprotocol}
	)
	// Close the connection when you're done with it.
	defer conn.Close()
	srv.wg.Add(1)
	request.Init()

	fmt.Printf(
		"Src->: \t\t\t\t%s\nDst<-: \t\t\t\t%s\n",
		conn.RemoteAddr(),
		conn.LocalAddr())

	// Read the incoming connection into the buffer.
	for {
		select {
		case <-srv.done:
			fmt.Println("Close connection", conn.RemoteAddr())
			srv.wg.Done()
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
			fmt.Printf("Src->: \t\t\t\t%s, Packet ID:%d\n", conn.RemoteAddr(), id_packet)

			// fmt.Println(hex.Dump(request.data))
			//request.ModbusDump()
			var answer *ModbusPacket
			answer, err = request.HandlerRequest(srv.Data)

			//answer.ModbusDump()
			conn.Write(answer.data)
		}
	}
	return err
}
