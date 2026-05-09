-- 1. Elimino gli indici (opzionale, cadono con le tabelle, ma è buona pratica)
DROP INDEX IF EXISTS idx_dht_device_ts;
DROP INDEX IF EXISTS idx_device_status_device_ts;

-- 2. Elimino le View
DROP VIEW IF EXISTS mqtt_topic_list_metadata;
DROP VIEW IF EXISTS mqtt_topic_list;

-- 3. Elimino le tabelle che dipendono da altre (le foglie della gerarchia)
DROP TABLE IF EXISTS device_status;
DROP TABLE IF EXISTS dht_readings;
DROP TABLE IF EXISTS sensors_devices;
DROP TABLE IF EXISTS devices;

-- 4. Elimino le tabelle principali
DROP TABLE IF EXISTS sensor_type;
DROP TABLE IF EXISTS device_type;
DROP TABLE IF EXISTS rooms;

-- 5. Elimino l'estensione (attenzione: solo se non ti serve per altro)
-- DROP EXTENSION IF EXISTS pgcrypto;