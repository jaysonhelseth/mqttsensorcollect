package models

import (
	"encoding/json"
	"log"
)

type TempMsg struct {
	Temp float64 `json:"temp"`
}

func GetAirTemp(serialData []byte) float64 {
	tempMsg := &TempMsg{}

	err := json.Unmarshal(serialData, tempMsg)
	if err != nil {
		log.Printf(err.Error())
	}
	return tempMsg.Temp
}
