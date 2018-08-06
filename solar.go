package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goburrow/modbus"
)

type Registers struct {
	Status                   int16 // Register 0
	_                        int16
	_                        int16
	PvDeciVolts              int16 // Register 3
	PvDeciAmps               int16 // Register 4
	_                        int16
	PvDeciWatts              int16 // Register 6
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	AcDeciWatts              int16 // Register 12
	AcCentiHerz              int16 // Register 13
	AcDeciVolts              int16 // Register 14
	AcDeciAmps               int16 // Register 15
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	_                        int16
	TotalHectaWattsHourToday int16 // Register 27
	_                        int16
	TotalHectaWattsHour      int16 // Register 29
	RuntimeSeconds           int32 // Register 30 and 31
	TemperatureDeciCelcius   int16 // Register 32
}

type SolarData struct {
	Status                  string
	PvVolts                 float32
	PvAmps                  float32
	PvWatts                 float32
	AcWatts                 float32
	AcHerz                  float32
	AcVolts                 float32
	AcAmps                  float32
	TotalKiloWattsHourToday float32
	TotalKiloWattsHour      float32
	RuntimeSeconds          int32
	TemperatureCelcius      float32
}

func main() {
	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 1 * time.Second

	err := handler.Connect()
	handler.Close()

	if err != nil {
		println(err.Error())
	}

	client := modbus.NewClient(handler)

	results, err := client.ReadInputRegisters(0, 33)

	if err != nil {
		println(err.Error())
	}

	reader := bytes.NewReader(results)

	registers := Registers{}

	err = binary.Read(reader, binary.BigEndian, &registers)

	if err != nil {
		println(err.Error())
	}

	solarData := solarDataFromRegisters(registers)

	json, _ := json.MarshalIndent(solarData, "", "    ")

	fmt.Printf("results: %s\n", json)
}

func solarDataFromRegisters(registers Registers) SolarData {
	status := map[int16]string{
		0: "Waiting",
		1: "Normal",
		2: "Fault",
	}

	return SolarData{
		status[registers.Status],
		float32(registers.PvDeciVolts) / 10,
		float32(registers.PvDeciAmps) / 10,
		float32(registers.PvDeciWatts) / 10,
		float32(registers.AcDeciWatts) / 10,
		float32(registers.AcCentiHerz) / 100,
		float32(registers.AcDeciVolts) / 10,
		float32(registers.AcDeciAmps) / 10,
		float32(registers.TotalHectaWattsHourToday) / 10,
		float32(registers.TotalHectaWattsHour) / 10,
		registers.RuntimeSeconds,
		float32(registers.TemperatureDeciCelcius) / 10,
	}
}
