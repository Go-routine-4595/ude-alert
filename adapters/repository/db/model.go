package db

import (
	"context"
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"sync"
	"time"

	"github.com/Go-routine-4595/ude-alert/domain"
	"github.com/uptrace/bun"
)

type Model struct {
	db    *bun.DB
	sqldb *sql.DB
	wg    *sync.WaitGroup

	Postgres
	loglevel bool
	errorLog *log.Logger
}

func (p *Model) WriteEquipmentAndData(ctx context.Context, e *domain.Equipment) error {

	var (
		fuelLevel           *FuelLevel
		oilPressure         *OilPressure
		oilEngineTemp       *OilEngineTemperature
		transmissionOilTemp *TransmissionOilTemperature
		equipment           *Equipment
		tmp                 Equipment
	)

	equipment = &Equipment{
		EquipmentName:  e.EquipmentName,
		EquipmentType:  e.EquipmentType,
		Manufacturer:   e.Manufacturer,
		Model:          e.Model,
		ProductionYear: e.ProductionYear,
		Location:       e.Location,
		EquipmentUUID:  e.EquipmentID.String(),
	}

	oilPressure = &OilPressure{
		Timestamp:          time.Now(),
		OilPressureDecimal: 45,
		EquipmentUUID:      uuid.NewV4().String(),
	}

	fuelLevel = &FuelLevel{
		Timestamp:        time.Now(),
		FuelLevelDecimal: 50,
		EquipmentUUID:    uuid.NewV4().String(),
	}

	oilEngineTemp = &OilEngineTemperature{
		Timestamp:                   time.Now(),
		OilEngineTemperatureDecimal: 195,
		EquipmentUUID:               uuid.NewV4().String(),
	}

	transmissionOilTemp = &TransmissionOilTemperature{
		Timestamp:                         time.Now(),
		TransmissionOilTemperatureDecimal: 200,
		EquipmentUUID:                     uuid.NewV4().String(),
	}

	// check if this equipment is not already existing
	err := p.db.NewSelect().Model(&tmp).Where("equipment_uuid = ?", equipment.EquipmentUUID).Scan(ctx)
	if err != nil {
		p.errorLog.Fatal("Error could not read equipment: ", err)
	}
	if tmp.EquipmentUUID == equipment.EquipmentUUID {
		return fmt.Errorf("Equipment %s already exists", equipment.EquipmentName)
	}
	// Start a transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.errorLog.Fatal("Error starting transaction: ", err)
	}

	// Insert the Equipment data
	_, err = tx.NewInsert().Model(equipment).Exec(ctx)
	if err != nil {
		tx.Rollback() // Rollback the transaction in case of error
		p.errorLog.Fatal("Error inserting equipment: ", err)
	}

	// Prepare associated data with the correct Equipment ID
	fuelLevel.EquipmentID = equipment.EquipmentID
	oilPressure.EquipmentID = equipment.EquipmentID
	oilEngineTemp.EquipmentID = equipment.EquipmentID
	transmissionOilTemp.EquipmentID = equipment.EquipmentID

	// Insert FuelLevel data
	_, err = tx.NewInsert().Model(fuelLevel).Exec(ctx)
	if err != nil {
		tx.Rollback()
		p.errorLog.Fatal("Error inserting fuel level: ", err)
	}

	// Insert OilPressure data
	_, err = tx.NewInsert().Model(oilPressure).Exec(ctx)
	if err != nil {
		tx.Rollback()
		p.errorLog.Fatal("Error inserting oil pressure: ", err)
	}

	// Insert OilEngineTemperature data
	_, err = tx.NewInsert().Model(oilEngineTemp).Exec(ctx)
	if err != nil {
		tx.Rollback()
		p.errorLog.Fatal("Error inserting oil engine temperature: ", err)
	}

	// Insert TransmissionOilTemperature data
	_, err = tx.NewInsert().Model(transmissionOilTemp).Exec(ctx)
	if err != nil {
		tx.Rollback()
		p.errorLog.Fatal("Error inserting transmission oil temperature: ", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		p.errorLog.Fatal("Error committing transaction: ", err)
	}

	p.errorLog.Println("Successfully inserted equipment and all associated data.")
	return nil
}

func (p *Model) InsertEquipment(ctx context.Context, ed *domain.Equipment) {
	var e *Equipment

	e = &Equipment{
		EquipmentUUID:  ed.EquipmentID.String(),
		EquipmentName:  ed.EquipmentName,
		EquipmentType:  ed.EquipmentType,
		Manufacturer:   ed.Manufacturer,
		Model:          ed.Model,
		Location:       ed.Location,
		ProductionYear: ed.ProductionYear,
	}
	_, err := p.db.NewInsert().Model(e).Exec(ctx)
	if err != nil {
		p.errorLog.Fatal("Error inserting equipment: ", err)
	}
}

func (p *Model) InsertFuelLevel(ctx context.Context, fld *domain.FuelLevel, equipmentUUID string) {
	var (
		fl         *FuelLevel
		equipments []Equipment
	)

	err := p.db.NewSelect().
		Model(&equipments).
		Limit(10).
		Where(" equipment_uuid = ?", equipmentUUID).
		Scan(ctx)
	if err != nil {
		p.errorLog.Fatal("Error fetching equipment: ", err)
	}

	fl = &FuelLevel{
		Timestamp:        fld.Timestamp,
		FuelLevelDecimal: fld.FuelLevelDecimal,
		EquipmentUUID:    fld.FuelLevelID.String(),
		EquipmentID:      equipments[0].EquipmentID,
	}
	_, err = p.db.NewInsert().Model(fl).Exec(ctx)
	if err != nil {
		p.errorLog.Fatal("Error inserting fuel level: ", err)
	}
}

func (p *Model) InsertOilPressure(ctx context.Context, opd *domain.OilPressure, equipmentUUID string) {
	var (
		op         *OilPressure
		equipments []Equipment
	)

	err := p.db.NewSelect().
		Model(&equipments).
		Limit(10).
		Where(" equipment_uuid = ?", equipmentUUID).
		Scan(ctx)
	if err != nil {
		p.errorLog.Fatal("Error fetching equipment: ", err)
	}

	op = &OilPressure{
		Timestamp:          opd.Timestamp,
		OilPressureDecimal: opd.OilPressureDecimal,
		EquipmentUUID:      opd.OilPressureID.String(),
		EquipmentID:        equipments[0].EquipmentID,
	}
	_, err = p.db.NewInsert().Model(op).Exec(ctx)
	if err != nil {
		p.errorLog.Fatal("Error inserting oil pressure: ", err)
	}
}

func (p *Model) InsertOilEngineTemperature(ctx context.Context, oetd *domain.OilEngineTemperature, equipmentUUID string) {
	var (
		oet        *OilEngineTemperature
		equipments []Equipment
	)

	err := p.db.NewSelect().
		Model(&equipments).
		Limit(10).
		Where(" equipment_uuid = ?", equipmentUUID).
		Scan(ctx)
	if err != nil {
		p.errorLog.Fatal("Error fetching equipment: ", err)
	}

	oet = &OilEngineTemperature{
		Timestamp:                   oetd.Timestamp,
		OilEngineTemperatureDecimal: oetd.OilEngineTemperatureDecimal,
		EquipmentUUID:               oetd.OilEngineTempID.String(),
		EquipmentID:                 equipments[0].EquipmentID,
	}
	_, err = p.db.NewInsert().Model(oet).Exec(ctx)
	if err != nil {
		p.errorLog.Fatal("Error inserting oil engine temperature: ", err)
	}
}

func (p *Model) InsertTransmissionOilTemperature(ctx context.Context, totd *domain.TransmissionOilTemperature, equipmentUUID string) {
	var (
		tot        *TransmissionOilTemperature
		equipments []Equipment
	)

	err := p.db.NewSelect().
		Model(&equipments).
		Limit(10).
		Where(" equipment_uuid = ?", equipmentUUID).
		Scan(ctx)
	if err != nil {
		p.errorLog.Fatal("Error fetching equipment: ", err)
	}

	tot = &TransmissionOilTemperature{
		Timestamp:                         totd.Timestamp,
		TransmissionOilTemperatureDecimal: totd.TransmissionOilTemperatureDecimal,
		EquipmentUUID:                     totd.TransmissionOilTempID.String(),
		EquipmentID:                       equipments[0].EquipmentID,
	}

	_, err = p.db.NewInsert().Model(tot).Exec(ctx)
	if err != nil {
		p.errorLog.Fatal("Error inserting transmission oil temperature: ", err)
	}
}

func (p *Model) ReadEquipment(ctx context.Context) ([]domain.Equipment, error) {
	var dequipments []domain.Equipment

	var equipments []EquipmentWithLatestData
	err := p.db.NewSelect().
		Model(&equipments).
		Relation("LatestFuelLevel", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Relation("LatestOilPressure", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Relation("LatestOilEngineTemp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Relation("LatestTransmissionOilTemp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Order("equipment_id ASC").
		Scan(ctx)
	if err != nil {
		log.Fatal("Error reading latest equipment data: ", err)
	}

	for _, e := range equipments {
		var (
			ed    domain.Equipment
			fd    domain.FuelLevel
			otd   domain.OilEngineTemperature
			opd   domain.OilPressure
			totd  domain.TransmissionOilTemperature
			luuid uuid.UUID
			err   error
		)

		luuid, err = uuid.FromString(e.EquipmentUUID)
		if err != nil {
			p.errorLog.Printf("Error while creating equipement UUID from DB: %s", err)
			lerror := fmt.Errorf("error creating equipement UUID from DB: [%w]", err)
			return dequipments, lerror
		}

		fd, err = adaptFuleLevel(e.LatestFuelLevel, false)
		if err != nil {
			return dequipments, err
		}
		otd, err = adaptOilEngineTemperature(e.LatestOilEngineTemp, false)
		if err != nil {
			return dequipments, err
		}
		opd, err = adaptOilPressure(e.LatestOilPressure, false)
		if err != nil {
			return dequipments, err
		}
		totd, err = adaptTransmissionOilTemperature(e.LatestTransmissionOilTemp, false)
		if err != nil {
			return dequipments, err
		}

		ed = domain.Equipment{
			EquipmentID:                    luuid,
			EquipmentName:                  e.EquipmentName,
			EquipmentType:                  e.EquipmentType,
			Manufacturer:                   e.Manufacturer,
			Model:                          e.Model,
			Location:                       e.Location,
			ProductionYear:                 e.ProductionYear,
			FuelLevelItem:                  fd,
			OilEngineTemperatureItem:       otd,
			OilPressureItem:                opd,
			TransmissionOilTemperatureItem: totd,
		}
		fmt.Printf("Equipment ID: %d, Name: %s, Type: %s, Manufacturer: %s, Model: %s, Production Year: %d, Location: %s\n",
			e.EquipmentID, e.EquipmentName, e.EquipmentType, e.Manufacturer, e.Model, e.ProductionYear, e.Location)
		dequipments = append(dequipments, ed)
	}
	return dequipments, nil
}

func (p *Model) LoadEquipment(ctx context.Context, v int) ([]domain.Equipment, error) {
	var (
		dequipments []domain.Equipment
		equipments  []Equipment
		err         error
	)

	err = p.db.NewSelect().
		Model(&equipments).
		Limit(v).
		Scan(ctx)
	if err != nil {
		log.Fatal("Error reading latest equipment data: ", err)
	}

	var extendedEquipments []EquipmentWithLatestData

	for _, e := range equipments {
		var (
			ne  EquipmentWithLatestData
			fl  FuelLevel
			op  OilPressure
			ot  OilEngineTemperature
			tot TransmissionOilTemperature
		)

		err = p.db.NewSelect().Model(&fl).Where("equipment_id = ?", e.EquipmentID).Scan(ctx)
		if err != nil {
			log.Fatal("Error reading latest equipment data: ", err)
		}
		err = p.db.NewSelect().Model(&op).Where("equipment_id = ?", e.EquipmentID).Scan(ctx)
		if err != nil {
			log.Fatal("Error reading latest equipment data: ", err)
		}
		err = p.db.NewSelect().Model(&ot).Where("equipment_id = ?", e.EquipmentID).Scan(ctx)
		if err != nil {
			log.Fatal("Error reading latest equipment data: ", err)
		}
		err = p.db.NewSelect().Model(&tot).Where("equipment_id = ?", e.EquipmentID).Scan(ctx)
		if err != nil {
			log.Fatal("Error reading latest equipment data: ", err)
		}

		ne.Equipment = e
		ne.LatestFuelLevel = fl
		ne.LatestOilPressure = op
		ne.LatestOilEngineTemp = ot
		ne.LatestTransmissionOilTemp = tot

		extendedEquipments = append(extendedEquipments, ne)
	}

	for _, e := range extendedEquipments {
		var (
			ed    domain.Equipment
			fd    domain.FuelLevel
			otd   domain.OilEngineTemperature
			opd   domain.OilPressure
			totd  domain.TransmissionOilTemperature
			luuid uuid.UUID
			err   error
		)

		luuid, err = uuid.FromString(e.EquipmentUUID)
		if err != nil {
			p.errorLog.Printf("Error while creating equipement UUID from DB: %s", err)
			lerror := fmt.Errorf("error creating equipement UUID from DB: [%w]", err)
			return dequipments, lerror
		}

		fd, err = adaptFuleLevel(e.LatestFuelLevel, true)
		if err != nil {
			return dequipments, err
		}
		otd, err = adaptOilEngineTemperature(e.LatestOilEngineTemp, true)
		if err != nil {
			return dequipments, err
		}
		opd, err = adaptOilPressure(e.LatestOilPressure, true)
		if err != nil {
			return dequipments, err
		}
		totd, err = adaptTransmissionOilTemperature(e.LatestTransmissionOilTemp, true)
		if err != nil {
			return dequipments, err
		}

		ed = domain.Equipment{
			EquipmentID:                    luuid,
			EquipmentName:                  e.EquipmentName,
			EquipmentType:                  e.EquipmentType,
			Manufacturer:                   e.Manufacturer,
			Model:                          e.Model,
			Location:                       e.Location,
			ProductionYear:                 e.ProductionYear,
			FuelLevelItem:                  fd,
			OilEngineTemperatureItem:       otd,
			OilPressureItem:                opd,
			TransmissionOilTemperatureItem: totd,
		}
		fmt.Printf("Equipment ID: %d, Name: %s, Type: %s, Manufacturer: %s, Model: %s, Production Year: %d, Location: %s\n",
			e.EquipmentID, e.EquipmentName, e.EquipmentType, e.Manufacturer, e.Model, e.ProductionYear, e.Location)
		dequipments = append(dequipments, ed)
	}
	return dequipments, nil
}

func ReadLatestEquipmentData(ctx context.Context, db *bun.DB) {
	var equipments []EquipmentWithLatestData
	err := db.NewSelect().
		Model(&equipments).
		Relation("LatestFuelLevel", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Relation("LatestOilPressure", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Relation("LatestOilEngineTemp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Relation("LatestTransmissionOilTemp", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("timestamp DESC").Limit(1)
		}).
		Order("equipment_id ASC").
		Scan(ctx)
	if err != nil {
		log.Fatal("Error reading latest equipment data: ", err)
	}

	for _, e := range equipments {
		fmt.Printf("Equipment ID: %d, Name: %s\n", e.EquipmentID, e.EquipmentName)
		fmt.Printf("Latest Fuel Level: %f, Time: %v\n", e.LatestFuelLevel.FuelLevelDecimal, e.LatestFuelLevel.Timestamp)
		fmt.Printf("Latest Oil Pressure: %f, Time: %v\n", e.LatestOilPressure.OilPressureDecimal, e.LatestOilPressure.Timestamp)
		fmt.Printf("Latest Oil Engine Temperature: %f, Time: %v\n", e.LatestOilEngineTemp.OilEngineTemperatureDecimal, e.LatestOilEngineTemp.Timestamp)
		fmt.Printf("Latest Transmission Oil Temperature: %f, Time: %v\n", e.LatestTransmissionOilTemp.TransmissionOilTemperatureDecimal, e.LatestTransmissionOilTemp.Timestamp)
	}
}
