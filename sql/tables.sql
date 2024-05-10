CREATE TABLE Equipment (
    equipment_id SERIAL PRIMARY KEY,
    equipment_uuid VARCHAR(36) NOT NULL,
    equipment_name VARCHAR(100) NOT NULL,
    equipment_type VARCHAR(100) NOT NULL,
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    production_year INT,
    location VARCHAR(100)
);

CREATE TABLE FuelLevel (
    fuel_level_id SERIAL PRIMARY KEY,
    equipment_uuid VARCHAR(36) NOT NULL,
    equipment_id INT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    fuel_level DECIMAL,
    FOREIGN KEY (equipment_id) REFERENCES Equipment(equipment_id)
);

CREATE TABLE OilPressure (
    oil_pressure_id SERIAL PRIMARY KEY,
    equipment_uuid VARCHAR(36) NOT NULL,
    equipment_id INT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    oil_pressure DECIMAL,
    FOREIGN KEY (equipment_id) REFERENCES Equipment(equipment_id)
);

CREATE TABLE OilEngineTemperature (
    oil_engine_temp_id SERIAL PRIMARY KEY,
    equipment_uuid VARCHAR(36) NOT NULL,
    equipment_id INT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    oil_engine_temperature DECIMAL,
    FOREIGN KEY (equipment_id) REFERENCES Equipment(equipment_id)
);

CREATE TABLE TransmissionOilTemperature (
    transmission_oil_temp_id SERIAL PRIMARY KEY,
    equipment_uuid VARCHAR(36) NOT NULL,
    equipment_id INT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    transmission_oil_temperature DECIMAL,
    FOREIGN KEY (equipment_id) REFERENCES Equipment(equipment_id)
);
