package io

import (
	"MqttSensorCollect/models"
	"bufio"
	"encoding/json"
	"github.com/tarm/serial"
	"log"
	"strings"
)

var (
	deviceConfig serial.Config
	device       *serial.Port
)

func setup(devicePath string) {
	deviceConfig = serial.Config{Name: devicePath, Baud: 9600}

	var err error
	device, err = serial.OpenPort(&deviceConfig)

	if err != nil {
		log.Fatal("Failed to setup listening device.")
	}
}

func read() ([]byte, error) {
	reader := bufio.NewReader(device)
	bytes, err := reader.ReadBytes('\x0a')

	if err != nil {
		return errorTempMsg(), err
	}

	return bytes, nil
}

func errorTempMsg() []byte {
	dict := map[string]float32{
		"temp": 0.0,
	}
	msg, _ := json.Marshal(dict)
	return msg
}

func ReadFromSerial(devicePath string, airTemp *models.Temp) {
	setup(devicePath)

	for {
		serialData, err := read()

		var msg string
		if err == nil {
			msg = strings.TrimSpace(string(serialData))
		} else {
			msg = string(serialData)
		}

		log.Println(msg)
		temp := models.GetAirTemp(serialData)
		airTemp.Write(temp)
	}
}
