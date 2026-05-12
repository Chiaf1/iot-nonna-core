package repository

import (
	"context"
	"encoding/json"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
)

// Get all sensors_types associated to a device
func (r *Repository) GetDevicesSensors(deviceId string) ([]domain.SensorType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	rows, err := r.DbPool.Query(ctx, `
	SELECT
	st.id, st.code, st.topic, st.description,
	st.readings_table_name, st.column_schema, st.value_mapping,
	st.payload_format, st.qos_mqtt
	FROM sensor_type st
	JOIN sensors_devices sd ON sd.sensor_id = st.id
	WHERE sd.device_id = $1
	ORDER BY st.code
	`, deviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []domain.SensorType
	for rows.Next() {
		var st domain.SensorType
		var colSchemaBytes json.RawMessage

		err := rows.Scan(
			&st.Id, &st.Code, &st.Topic, &st.Description,
			&st.ReadingsTableName, &colSchemaBytes, &st.ValueMapping,
			&st.PayloadFormat, &st.QosMqtt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal column schema
		if err := json.Unmarshal(colSchemaBytes, &st.ColumnSchema); err != nil {
			return nil, err
		}

		sensors = append(sensors, st)
	}
	return sensors, rows.Err()
}

// Associat new sensor to device
func (r *Repository) AddSensorToDevice(deviceId, sensorId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_write)
	defer cancel()

	_, err := r.DbPool.Exec(ctx, `
		INSERT INTO sensors_devices (device_id, sensor_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, deviceId, sensorId)
	return err
}

// Remove associated sensor from device
func (r *Repository) RemoveSensorFromDevice(deviceId, sensorId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_write)
	defer cancel()

	result, err := r.DbPool.Exec(ctx, `
		DELETE FROM sensors_devices
		WHERE device_id = $1 AND sensor_id = $2
	`, deviceId, sensorId)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}
