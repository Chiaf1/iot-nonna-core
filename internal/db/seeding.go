package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Run seeding queries with development data, one room, one device, status + dht sensors
func RunSeeding(pool *pgxpool.Pool, timeout_query time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout_query)
	defer cancel()
	// check if seeding already done
	var count int
	pool.QueryRow(ctx, "SELECT COUNT(*) FROM device_type").Scan(&count)
	if count > 0 {
		log.Println("Seed already applied, skipping")
		return nil
	}

	// Start seeding with development data
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := []string{
		`INSERT INTO rooms (name) VALUES ('test') ON CONFLICT (name) DO NOTHING`,

		`INSERT INTO device_type (code,topic,description) VALUES ('esp32-wroom-2dev','esp32','Devkit esp 32 wroom 2') ON CONFLICT (code) DO NOTHING`,

		`INSERT INTO sensor_type (code,topic,description,readings_table_name,column_schema,payload_format,qos_mqtt)
		VALUES ('dht22', 'dht', 'Generic temperature and humidity sensor', 'dht_readings', 
		'{
			"temperature": {
				"column": "temperature",
				"type": "float"
			},
			"humidity": {
				"column": "humidity",
				"type": "float"
			}
		}',
		'json',
		0)  ON CONFLICT (code) DO NOTHING`,

		`INSERT INTO sensor_type (code,topic,description,readings_table_name,column_schema,value_mapping,payload_format,qos_mqtt)
		VALUES ('Connection status', 'status', 'Not realy a sensor but more of an expansion of the device with the mqtt lastwill enabled',
		'device_status', 
		'{
			"$payload": {
				"column": "status",
				"type": "bool"
			}
		}',
		'{
			"online": true,
			"offline": false
		}',
		'raw',
		1)  ON CONFLICT (code) DO NOTHING`,

		`INSERT INTO devices (code,device_type_id,room_id) VALUES (
		'test-1',
		(SELECT id FROM device_type WHERE topic ILIKE 'esp32'),
		(SELECT id FROM rooms WHERE name ILIKE 'test')
		) ON CONFLICT (code) DO NOTHING`,

		`INSERT INTO sensors_devices (device_id,sensor_id) VALUES (
			(SELECT id FROM devices WHERE code ILIKE 'test-1'),
			(SELECT id FROM sensor_type WHERE topic ILIKE 'dht')
		),
		(
			(SELECT id FROM devices WHERE code ILIKE 'test-1'),
			(SELECT id FROM sensor_type WHERE topic ILIKE 'status')
		)  ON CONFLICT (device_id,sensor_id) DO NOTHING`,

	}

	for _, q := range queries {
		if _, err := tx.Exec(ctx, q); err != nil {
			log.Printf("Errore nella query: %s", q)
			return err
		}
	}

	log.Println("Seeding applied succesfully")
	return tx.Commit(ctx)
}
