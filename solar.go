package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
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
	DateTime                time.Time
	Status                  string
	PvVolts                 float64
	PvAmps                  float64
	PvWatts                 float64
	AcWatts                 float64
	AcHerz                  float64
	AcVolts                 float64
	AcAmps                  float64
	TotalKiloWattsHourToday float64
	TotalKiloWattsHour      float64
	RuntimeSeconds          int
	TemperatureCelcius      float64
}

func main() {
	modbusClient := newModbusClient()

	registers := registersFromModbusClient(modbusClient)

	solarData := solarDataFromRegisters(registers)

	storeSolarData(solarData)

	json, err := json.MarshalIndent(solarData, "", "    ")

	if err != nil {
		println(err.Error())
	}

	fmt.Printf("results: %s\n", json)
}

func newModbusClient() modbus.Client {
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

	return modbus.NewClient(handler)
}

func registersFromModbusClient(client modbus.Client) Registers {
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

	return registers
}

func solarDataFromRegisters(registers Registers) SolarData {
	status := map[int16]string{
		0: "Waiting",
		1: "Normal",
		2: "Fault",
	}

	return SolarData{
		time.Now(),
		status[registers.Status],
		float64(registers.PvDeciVolts) / 10,
		float64(registers.PvDeciAmps) / 10,
		float64(registers.PvDeciWatts) / 10,
		float64(registers.AcDeciWatts) / 10,
		float64(registers.AcCentiHerz) / 100,
		float64(registers.AcDeciVolts) / 10,
		float64(registers.AcDeciAmps) / 10,
		float64(registers.TotalHectaWattsHourToday) / 10,
		float64(registers.TotalHectaWattsHour) / 10,
		int(registers.RuntimeSeconds),
		float64(registers.TemperatureDeciCelcius) / 10,
	}
}

func storeSolarData(solarData SolarData) {
	json, _ := json.Marshal(solarData)
	filename := fmt.Sprintf("/data/solar/%s.json", solarData.DateTime.Format("2006-01"))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println(err.Error())
	}

	if _, err := file.Write([]byte(fmt.Sprintf("%s\n", json))); err != nil {
		println(err.Error())
	}

	if err := file.Close(); err != nil {
		println(err.Error())
	}
}
