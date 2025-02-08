package request

import "time"

type Waterrequest struct {
	RecordID           string    `json:"recordid"`
	FactoryID          string    `json:"factoryid"`
	RecordYear         int       `json:"recordyear"`
	RecordMonth        int       `json:"recordmonth"`
	TapWaterMeter      float64   `json:"tapWaterMeter"`
	RecycledWaterMeter float64   `json:"recycledWaterMeter"`
	UserID             string    `json:"userid"`
	UserDate           time.Time `json:"userdate"`
}
