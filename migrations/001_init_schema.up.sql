-- extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- tables
CREATE TABLE rooms(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT UNIQUE NOT NULL
);
CREATE TABLE device_type(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code TEXT UNIQUE NOT NULL, --tipologia scheda lettura
	topic VARCHAR(50) NOT NULL, --base topic mqtt
	description TEXT
);
CREATE TABLE sensor_type(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code TEXT UNIQUE NOT NULL, --tipologia sensore, es codice DHT22
	topic VARCHAR(50) UNIQUE NOT NULL, --topic mqtt del sensore
	description TEXT,
	readings_table_name TEXT NOT NULL, --nome tabella salvataggio dati
	column_schema JSONB NOT NULL, --struttura json per decode topic to colum
	value_mapping JSONB,
	payload_format TEXT NOT NULL DEFAULT 'json',
	qos_mqtt SMALLINT
);
CREATE TABLE devices(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code VARCHAR(50) UNIQUE NOT NULL, --nome dispositivo, parte topic mqtt
	device_type_id UUID REFERENCES device_type(id) NOT NULL, 
	room_id UUID REFERENCES rooms(id) ON DELETE SET NULL,
	created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE sensors_devices(
	device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
	sensor_id UUID REFERENCES sensor_type(id) ON DELETE RESTRICT,
	PRIMARY KEY (device_id, sensor_id)
);
CREATE TABLE dht_readings(
	id BIGSERIAL PRIMARY KEY,
	device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
	sensor_id UUID REFERENCES sensor_type(id) ON DELETE RESTRICT,
	timestamp TIMESTAMPTZ DEFAULT now(),
	temperature REAL,
	humidity REAL
);
CREATE TABLE device_status(
	id BIGSERIAL PRIMARY KEY,
	device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
	sensor_id UUID REFERENCES sensor_type(id) ON DELETE RESTRICT,
	timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
	status BOOL NOT NULL
);
-- views
CREATE VIEW mqtt_topic_list AS
SELECT dt.topic || '/' || d.code || '/' || s.topic AS topics FROM device_type AS dt
INNER JOIN devices AS d ON d.device_type_id = dt.id
INNER JOIN sensors_devices AS sd ON d.id = sd.device_id
INNER JOIN sensor_type AS s ON sd.sensor_id = s.id;

CREATE VIEW mqtt_topic_list_metadata AS
SELECT (((dt.topic::text || '/'::text) || d.code::text) || '/'::text) || s.topic::text AS topics,
s.readings_table_name,
d.id AS device_id,
s.id AS sensor_id,
s.column_schema,
s.value_mapping,
s.payload_format,
s.qos_mqtt 
FROM device_type dt
JOIN devices d ON d.device_type_id = dt.id
JOIN sensors_devices sd ON d.id = sd.device_id
JOIN sensor_type s ON sd.sensor_id = s.id;

-- indexes
CREATE INDEX idx_device_status_device_ts
ON device_status(device_id, timestamp DESC);

CREATE INDEX idx_dht_device_ts
ON dht_readings(device_id, timestamp DESC); 