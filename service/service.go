package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/Go-routine-4595/ude-alert/domain"
	"github.com/rs/zerolog"
)

const (
	maxEngineOilTemp = 250
	minEngineOilTemp = 120
	maxEngineOitPre  = 70
	minEngineOitPre  = 30
	maxTransOilTemp  = 220
	minTransOilPre   = 100
)

type Storer interface {
	WriteEquipmentAndData(ctx context.Context, e *domain.Equipment) (err error)
	LoadEquipment(ctx context.Context, v int) ([]domain.Equipment, error)
	InsertFuelLevel(ctx context.Context, fld *domain.FuelLevel, equipmentUUID string)
	InsertOilPressure(ctx context.Context, opd *domain.OilPressure, equipmentUUID string)
	InsertOilEngineTemperature(ctx context.Context, oetd *domain.OilEngineTemperature, equipmentUUID string)
	InsertTransmissionOilTemperature(ctx context.Context, totd *domain.TransmissionOilTemperature, equipmentUUID string)
}

type Service struct {
	store Storer
	eql   []domain.Equipment
	log   zerolog.Logger
}

func NewService(store Storer) *Service {

	return &Service{
		store: store.(Storer),
		log:   zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.DebugLevel).With().Timestamp().Logger(),
	}
}

func (s *Service) AddEquipment(e []byte) error {
	var (
		eq  *domain.Equipment
		ctx context.Context
		err error
	)

	eq = new(domain.Equipment)

	err = json.Unmarshal(e, eq)
	if err != nil {
		err = fmt.Errorf("AddEquipment json unmarshall error: [%w]", err)
		return err
	}
	ctx = context.Background()
	err = s.store.WriteEquipmentAndData(ctx, eq)

	return err
}

func (s *Service) LoadEquipment(v int) error {
	var (
		err error
		ctx context.Context
	)

	ctx = context.Background()
	s.eql, err = s.store.LoadEquipment(ctx, v)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (s *Service) UpdateEquipment(count int) error {

	var (
		fl  domain.FuelLevel
		op  domain.OilPressure
		ot  domain.OilEngineTemperature
		tot domain.TransmissionOilTemperature
		ctx context.Context
	)
	if count > len(s.eql) {
		count = len(s.eql)
	}

	for i := 0; i < count; i++ {
		s.eql[i].FuelLevelItem.FuelLevelDecimal = giveValue(s.eql[i].FuelLevelItem.FuelLevelDecimal, 100, 10, 2, true)
		s.eql[i].FuelLevelItem.Timestamp = time.Now()
		s.eql[i].OilPressureItem.OilPressureDecimal = giveValue(s.eql[i].OilPressureItem.OilPressureDecimal, maxEngineOitPre, minEngineOitPre, 2, false)
		s.eql[i].OilPressureItem.Timestamp = time.Now()
		s.eql[i].OilEngineTemperatureItem.OilEngineTemperatureDecimal = giveValue(s.eql[i].OilEngineTemperatureItem.OilEngineTemperatureDecimal, maxEngineOilTemp, minEngineOilTemp, 3, false)
		s.eql[i].OilEngineTemperatureItem.Timestamp = time.Now()
		s.eql[i].TransmissionOilTemperatureItem.TransmissionOilTemperatureDecimal = giveValue(s.eql[i].TransmissionOilTemperatureItem.TransmissionOilTemperatureDecimal, maxTransOilTemp, minTransOilPre, 3, false)
		s.eql[i].TransmissionOilTemperatureItem.Timestamp = time.Now()

		fl = s.eql[i].FuelLevelItem
		op = s.eql[i].OilPressureItem
		ot = s.eql[i].OilEngineTemperatureItem
		tot = s.eql[i].TransmissionOilTemperatureItem

		ctx = context.Background()

		s.store.InsertFuelLevel(ctx, &fl, s.eql[i].EquipmentID.String())
		s.store.InsertOilPressure(ctx, &op, s.eql[i].EquipmentID.String())
		s.store.InsertOilEngineTemperature(ctx, &ot, s.eql[i].EquipmentID.String())
		s.store.InsertTransmissionOilTemperature(ctx, &tot, s.eql[i].EquipmentID.String())

		s.log.Debug().Str(s.eql[i].EquipmentName, "EquipmentName").Float64("FuelLevel", fl.FuelLevelDecimal).Msg(("Fuel Level Decimal"))
		s.log.Debug().Str(s.eql[i].EquipmentName, "EquipmentName").Float64("OilPressure", op.OilPressureDecimal).Msg(("Oil Pressure Decimal"))
		s.log.Debug().Str(s.eql[i].EquipmentName, "EquipmentName").Float64("OilTemperature", ot.OilEngineTemperatureDecimal).Msg(("Oil Temperature Decimal"))
		s.log.Debug().Str(s.eql[i].EquipmentName, "EquipmentName").Float64("TranissionOilTemperature", tot.TransmissionOilTemperatureDecimal).Msg(("Tremission Oil Temperature Decimal"))

	}

	return nil
}

func giveValue(v float64, max float64, min float64, f int, onlydown bool) float64 {
	var (
		r int
		d int
	)

	r = rand.IntN(1 * f)
	d = rand.IntN(10)

	if onlydown == true && d%2 == 0 {
		return v
	}

	if d%2 == 0 {
		if v+float64(r) <= max {
			return v + float64(r)
		} else {
			return v
		}
	} else {
		if v-float64(r) > min {
			return v - float64(r)
		} else {
			return v
		}
	}
}
