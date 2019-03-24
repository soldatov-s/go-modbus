// Copyright 2019 Sergey Soldatov. All rights reserved.
// This software may be modified and distributed under the terms
// of the Apache license. See the LICENSE file for details.

package modbusrest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	. "github.com/soldatov-s/go-modbus"
)

// ModbusServer implements server interface
type ModbusRest struct {
	ModbusBaseServer                // Anonim ModbusBase implementation
	Router           *http.ServeMux // HTTP request multiplexer
	Server           *http.Server   // HTTP server
}

// Rest answer for Holding/Input registers request
type ModbusRegAnswer struct {
	Data []uint16 `json:"data"` // Read values
}

// Rest answer for Coils/DigitInputs
type ModbusBoolAnswer struct {
	Data []bool `json:"data"` // Read values
}

// Rest request to write Holding/Input registers
type ModbusWriteRegReq struct {
	Addr uint16   `json:"addr"` // Addres first element
	Data []uint16 `json:"data"` // Values for writing
}

// Rest request to write Coils/DigitInputs
type ModbusWriteBoolReq struct {
	Addr uint16 `json:"addr"` // Addres first element
	Data []bool `json:"data"` // Values for writing
}

// Build a response to an unknown request
func errAnswer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s; Method: %s; URL: %s", "GO AWAY", r.Method, r.URL.Path)
}

// Handler for GET/PUT request Coils
func (rest *ModbusRest) hndlCoils(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		answer ModbusBoolAnswer
		err    error
		addr   int
		cnt    int
	)

	switch r.Method {
	case "POST":
		var req ModbusWriteBoolReq
		_ = json.NewDecoder(r.Body).Decode(&req)
		rest.Data.ForceMultipleCoils(req.Addr, req.Data)

		answer.Data, err = rest.Data.ReadCoilStatus(req.Addr, uint16(len(req.Data)))
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(answer.Data)
	case "GET":
		query := r.URL.Query()
		addr, err = strconv.Atoi(query.Get("addr"))
		cnt, err = strconv.Atoi(query.Get("cnt"))

		answer.Data, err = rest.Data.ReadCoilStatus(uint16(addr), uint16(cnt))
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(answer.Data)

	default:
		errAnswer(w, r)
	}
}

// Handler for GET request DigitInputs
func (rest *ModbusRest) hndlDigitInputs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		answer ModbusBoolAnswer
		err    error
		addr   int
		cnt    int
	)

	switch r.Method {
	case "GET":
		query := r.URL.Query()
		addr, err = strconv.Atoi(query.Get("addr"))
		cnt, err = strconv.Atoi(query.Get("cnt"))

		answer.Data, err = rest.Data.ReadDescreteInputs(uint16(addr), uint16(cnt))
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(answer.Data)

	default:
		errAnswer(w, r)
	}

}

// Handler for GET/PUT request Holding Registers
func (rest *ModbusRest) hndlHoldReg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		answer ModbusRegAnswer
		err    error
		addr   int
		cnt    int
	)

	switch r.Method {
	case "POST":
		var req ModbusWriteRegReq
		_ = json.NewDecoder(r.Body).Decode(&req)
		rest.Data.PresetMultipleRegisters(req.Addr, req.Data)

		answer.Data, err = rest.Data.ReadHoldingRegisters(req.Addr, uint16(len(req.Data)))
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(answer.Data)
	case "GET":
		query := r.URL.Query()
		addr, err = strconv.Atoi(query.Get("addr"))
		cnt, err = strconv.Atoi(query.Get("cnt"))

		answer.Data, err = rest.Data.ReadHoldingRegisters(uint16(addr), uint16(cnt))
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(answer.Data)

	default:
		errAnswer(w, r)
	}
}

// Handler for GE request Inputs Registers
func (rest *ModbusRest) hndlInputReg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		answer ModbusRegAnswer
		err    error
		addr   int
		cnt    int
	)

	switch r.Method {
	case "GET":
		query := r.URL.Query()
		addr, err = strconv.Atoi(query.Get("addr"))
		cnt, err = strconv.Atoi(query.Get("cnt"))

		answer.Data, err = rest.Data.ReadInputRegisters(uint16(addr), uint16(cnt))
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(answer.Data)

	default:
		errAnswer(w, r)
	}
}

// Create new Rest-server for Modbus Data
func NewRest(host, port string, md *ModbusData) *ModbusRest {
	rest := new(ModbusRest)
	rest.Host = host
	rest.Port = port
	rest.Data = md
	rest.Router = http.NewServeMux()

	rest.Router.HandleFunc("/coils", rest.hndlCoils)
	rest.Router.HandleFunc("/d_in", rest.hndlDigitInputs)
	rest.Router.HandleFunc("/hold_reg", rest.hndlHoldReg)
	rest.Router.HandleFunc("/in_reg", rest.hndlInputReg)

	return rest
}

// Start Rest-server for Modbus Data
func (rest *ModbusRest) Start() error {
	var err error
	log.Println("REST-Server startup...")
	log.Println("Listening at", rest.String())

	rest.Server = &http.Server{Addr: rest.String(), Handler: rest.Router}
	go func() {
		if err := rest.Server.ListenAndServe(); err != nil {
			log.Fatalf("listenAndServe failed: %v", err)
		}
	}()
	log.Println("REST-server started")

	return err
}

func (rest *ModbusRest) Stop() {
	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := rest.Server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
