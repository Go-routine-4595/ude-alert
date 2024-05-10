package db

import (
	"time"

	"github.com/uptrace/bun"
)

// Equipment represents the 'Equipment' table
type Equipment struct {
	bun.BaseModel  `bun:"table:equipment,alias:e"`
	EquipmentID    int64  `bun:"equipment_id,pk,autoincrement"`
	EquipmentUUID  string `bun:"equipment_uuid,notnull,unique"`
	EquipmentName  string `bun:"equipment_name,notnull"`
	EquipmentType  string `bun:"equipment_type,notnull"`
	Manufacturer   string `bun:"manufacturer"`
	Model          string `bun:"model"`
	ProductionYear int    `bun:"production_year"`
	Location       string `bun:"location"`
}

// FuelLevel represents the 'FuelLevel' table
type FuelLevel struct {
	bun.BaseModel    `bun:"table:fuellevel,alias:f"`
	FuelLevelID      int64     `bun:"fuel_level_id,pk,autoincrement"`
	EquipmentUUID    string    `bun:"equipment_uuid,notnull"`
	EquipmentID      int64     `bun:"equipment_id,notnull"`
	Timestamp        time.Time `bun:"timestamp,default:current_timestamp"`
	FuelLevelDecimal float64   `bun:"fuel_level"`
}

// OilPressure represents the 'OilPressure' table
type OilPressure struct {
	bun.BaseModel      `bun:"table:oilpressure,alias:o"`
	OilPressureID      int64     `bun:"oil_pressure_id,pk,autoincrement"`
	EquipmentUUID      string    `bun:"equipment_uuid,notnull"`
	EquipmentID        int64     `bun:"equipment_id,notnull"`
	Timestamp          time.Time `bun:"timestamp,default:current_timestamp"`
	OilPressureDecimal float64   `bun:"oil_pressure"`
}

// OilEngineTemperature represents the 'OilEngineTemperature' table
type OilEngineTemperature struct {
	bun.BaseModel               `bun:"table:oilenginetemperature,alias:t"`
	OilEngineTempID             int64     `bun:"oil_engine_temp_id,pk,autoincrement"`
	EquipmentUUID               string    `bun:"equipment_uuid,notnull"`
	EquipmentID                 int64     `bun:"equipment_id,notnull"`
	Timestamp                   time.Time `bun:"timestamp,default:current_timestamp"`
	OilEngineTemperatureDecimal float64   `bun:"oil_engine_temperature"`
}

// TransmissionOilTemperature represents the 'TransmissionOilTemperature' table
type TransmissionOilTemperature struct {
	bun.BaseModel                     `bun:"table:transmissionoiltemperature,alias:trans"`
	TransmissionOilTempID             int64     `bun:"transmission_oil_temp_id,pk,autoincrement"`
	EquipmentUUID                     string    `bun:"equipment_uuid,notnull"`
	EquipmentID                       int64     `bun:"equipment_id,notnull"`
	Timestamp                         time.Time `bun:"timestamp,default:current_timestamp"`
	TransmissionOilTemperatureDecimal float64   `bun:"transmission_oil_temperature"`
}

type EquipmentWithLatestData struct {
	Equipment                                            // Embedding Equipment struct to hold equipment data.
	LatestFuelLevel           FuelLevel                  `bun:"rel:has-one,join:equipment_id=equipment_id"` // Latest fuel level.
	LatestOilPressure         OilPressure                `bun:"rel:has-one,join:equipment_id=equipment_id"` // Latest oil pressure.
	LatestOilEngineTemp       OilEngineTemperature       `bun:"rel:has-one,join:equipment_id=equipment_id"` // Latest oil engine temperature.
	LatestTransmissionOilTemp TransmissionOilTemperature `bun:"rel:has-one,join:equipment_id=equipment_id"` // Latest transmission oil temp.
}
