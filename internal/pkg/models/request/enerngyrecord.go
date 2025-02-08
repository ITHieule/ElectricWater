package request

import "time"

type Energyrequest struct {
	RecordID             string    `json:"recordid"`
	FactoryID            string    `json:"factoryid"`
	RecordYear           int       `json:"recordyear"`
	RecordMonth          int       `json:"recordmonth"`
	GridElectricityMeter float64   `json:"grid_electricity_meter"`
	SolarEnergyMeter     float64   `json:"solar_energy_meter"`
	UserID               string    `json:"userid"`
	UserDate             time.Time `json:"userdate"`
}
