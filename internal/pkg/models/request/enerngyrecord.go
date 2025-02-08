package request

import "time"

type Energyrequest struct {
	RecordID             string    `json:"recordid"`
	FactoryID            string    `json:"FactoryID"`
	RecordYear           int       `json:"RecordYear"`
	RecordMonth          int       `json:"RecordMonth"`
	GridElectricityMeter float64   `json:"GridElectricityMeter"`
	SolarEnergyMeter     float64   `json:"SolarEnergyMeter"`
	UserID               string    `json:"UserID"`
	UserDate             time.Time `json:"userdate"`
}
