package domain

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Service interface {
	LoadEquipment(v int) error
	UpdateEquipment(count int) error
}

// Equipment represents the 'Equipment' table
type Equipment struct {
	EquipmentID                    uuid.UUID `json:"equipment_id"`
	EquipmentName                  string    `json:"equipment_name"`
	EquipmentType                  string    `json:"equipment_type"`
	Manufacturer                   string    `json:"manufacturer"`
	Model                          string    `json:"model"`
	ProductionYear                 int       `json:"production_year"`
	Location                       string    `json:"location"`
	FuelLevelItem                  FuelLevel
	OilPressureItem                OilPressure
	OilEngineTemperatureItem       OilEngineTemperature
	TransmissionOilTemperatureItem TransmissionOilTemperature
}

// FuelLevel represents the 'FuelLevel' table
type FuelLevel struct {
	FuelLevelID      uuid.UUID
	Timestamp        time.Time
	FuelLevelDecimal float64
}

// OilPressure represents the 'OilPressure' table
type OilPressure struct {
	OilPressureID      uuid.UUID
	Timestamp          time.Time
	OilPressureDecimal float64
}

// OilEngineTemperature represents the 'OilEngineTemperature' table
type OilEngineTemperature struct {
	OilEngineTempID             uuid.UUID
	Timestamp                   time.Time
	OilEngineTemperatureDecimal float64
}

// TransmissionOilTemperature represents the 'TransmissionOilTemperature' table
type TransmissionOilTemperature struct {
	TransmissionOilTempID             uuid.UUID
	Timestamp                         time.Time
	TransmissionOilTemperatureDecimal float64
}
