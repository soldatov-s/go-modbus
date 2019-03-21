// rest
package modbus

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var md ModbusData

type ModbusReq struct {
	addr uint16 `json:"addr"`
	cnt  uint16 `json:"cnt"`
}

type ModbusRegAnswer struct {
	data []uint16 `json:"data"`
}

type ModbusBoolAnswer struct {
	data []bool `json:"data"`
}

type ModbusWriteRegReq struct {
	addr int      `json:"addr"`
	data []uint16 `json:"data"`
}

type ModbusWriteBoolReq struct {
	addr uint16 `json:"addr"`
	data []bool `json:"data"`
}

func readHoldReg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		req    ModbusReq
		answer ModbusRegAnswer
		err    error
	)
	_ = json.NewDecoder(r.Body).Decode(&req)
	answer.data, err = md.ReadHoldingRegisters(req.addr, req.cnt)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(answer)
}

func CreateRest() {
	r := mux.NewRouter()

	/*	r.HandleFunc("/coils", readModbusData).Methods("GET")
		r.HandleFunc("/coils", readModbusData).Methods("POST")
		r.HandleFunc("/d_in", writeModbusData).Methods("GET")
		r.HandleFunc("/d_in", writeModbusData).Methods("POST")*/
	r.HandleFunc("/hold_reg", readHoldReg).Methods("GET")
	//	r.HandleFunc("/hold_reg", writeHoldReg).Methods("POST")
	/*	r.HandleFunc("/in_reg", writeModbusData).Methods("GET")
		r.HandleFunc("/in_reg", writeModbusData).Methods("POST")*/
	http.ListenAndServe(":8000", r)
}
