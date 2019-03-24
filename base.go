// base
package modbus

// ModbusBase implements base server interface
type ModbusBase struct {
	Host string      // Host Name/IP
	Port string      // Server port
	Data *ModbusData // Modbus Data
}
