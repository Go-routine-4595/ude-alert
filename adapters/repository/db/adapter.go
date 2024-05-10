package db

import (
	"fmt"
	"github.com/Go-routine-4595/ude-alert/domain"
	uuid "github.com/satori/go.uuid"
)

func adaptFuleLevel(fl FuelLevel, reset bool) (domain.FuelLevel, error) {
	var (
		dfl domain.FuelLevel
		err error
		uid uuid.UUID
		tmp float64
	)

	uid, err = uuid.FromString(fl.EquipmentUUID)
	if err != nil {
		return dfl, fmt.Errorf("error creating equipement UUID from DB: [%w]", err)
	}

	if reset {
		tmp = 50
	} else {
		tmp = fl.FuelLevelDecimal
	}

	dfl = domain.FuelLevel{
		FuelLevelID:      uid,
		Timestamp:        fl.Timestamp,
		FuelLevelDecimal: tmp,
	}
	return dfl, nil
}

func adaptOilPressure(x OilPressure, reset bool) (domain.OilPressure, error) {
	var (
		i   domain.OilPressure
		err error
		uid uuid.UUID
		tmp float64
	)

	uid, err = uuid.FromString(x.EquipmentUUID)
	if err != nil {
		return i, fmt.Errorf("error creating equipement UUID from DB: [%w]", err)
	}

	if reset {
		tmp = 45
	} else {
		tmp = x.OilPressureDecimal
	}

	i = domain.OilPressure{
		OilPressureID:      uid,
		Timestamp:          x.Timestamp,
		OilPressureDecimal: tmp,
	}
	return i, nil
}

func adaptOilEngineTemperature(x OilEngineTemperature, reset bool) (domain.OilEngineTemperature, error) {
	var (
		i   domain.OilEngineTemperature
		err error
		uid uuid.UUID
		tmp float64
	)

	uid, err = uuid.FromString(x.EquipmentUUID)
	if err != nil {
		return i, fmt.Errorf("error creating equipement UUID from DB: [%w]", err)
	}

	if reset {
		tmp = 195
	} else {
		tmp = x.OilEngineTemperatureDecimal
	}

	i = domain.OilEngineTemperature{
		OilEngineTempID:             uid,
		Timestamp:                   x.Timestamp,
		OilEngineTemperatureDecimal: tmp,
	}
	return i, nil
}

func adaptTransmissionOilTemperature(x TransmissionOilTemperature, reset bool) (domain.TransmissionOilTemperature, error) {
	var (
		i   domain.TransmissionOilTemperature
		err error
		uid uuid.UUID
		tmp float64
	)

	uid, err = uuid.FromString(x.EquipmentUUID)
	if err != nil {
		return i, fmt.Errorf("error creating equipement UUID from DB: [%w]", err)
	}

	if reset {
		tmp = 150
	} else {
		tmp = x.TransmissionOilTemperatureDecimal
	}

	i = domain.TransmissionOilTemperature{
		TransmissionOilTempID:             uid,
		Timestamp:                         x.Timestamp,
		TransmissionOilTemperatureDecimal: tmp,
	}
	return i, nil
}
