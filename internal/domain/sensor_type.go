package domain

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Sensor_type struct {
	Id                string                      `json:"id"`
	Code              string                      `json:"code"`
	Topic             string                      `json:"topic"`
	Description       *string                     `json:"description,omitempty"`
	ReadingsTableName string                      `json:"readings_table_name"`
	ColumnSchema      map[string]ColumnSchemaType `json:"column_schema"`
	ValueMapping      json.RawMessage             `json:"value_mapping,omitempty"`
	PayloadFormat     string                      `json:"payload_format"`
	QosMqtt           *int16                      `json:"qos_mqtt,omitempty"`
}

type ColumnSchemaType struct {
	Column string `json:"column" validate:"required"`
	Type   string `json:"type" validate:"required,oneof=float int bool string"`
}

type Sensor_typeRequest struct {
	Code              string                      `json:"code" validate:"required"`
	Topic             string                      `json:"topic" validate:"required,max=50"`
	Description       *string                     `json:"description,omitempty"`
	ReadingsTableName string                      `json:"readings_table_name" validate:"required"`
	ColumnSchema      map[string]ColumnSchemaType `json:"column_schema" validate:"required,min=1"`
	ValueMapping      json.RawMessage             `json:"value_mapping,omitempty" validate:"omitempty,json"`
	PayloadFormat     string                      `json:"payload_format" validate:"required,oneof=json raw"`
	QosMqtt           *int16                      `json:"qos_mqtt,omitempty" validate:"omitempty,min=0,max=2"`
}

// Validation function for structure SensorTypeRequest
func SensorTypeReqValidation(sl validator.StructLevel) {
	s := sl.Current().Interface().(Sensor_typeRequest)

	for k := range s.ColumnSchema {
		if strings.TrimSpace(k) == "" {
			sl.ReportError(
				s.ColumnSchema,
				"ColumnSchema",
				"column_schema",
				"empty_key",
				"",
			)
			return
		}
	}

	if s.PayloadFormat != "raw" {
		return
	}

	// Payload raw -> column_schema deve avere $payload
	if len(s.ColumnSchema) == 0 {
		sl.ReportError(
			s.ColumnSchema,
			"ColumnSchema",
			"column_schema",
			"required",
			"",
		)
		return
	}
	cs, ok := s.ColumnSchema["$payload"]
	if !ok {
		sl.ReportError(
			s.ColumnSchema,
			"ColumnSchema",
			"column_schema",
			"payload_required",
			"",
		)
		return
	}
	// bool + raw -> value_mapping obbligatorio
	if cs.Type == "bool" && len(s.ValueMapping) == 0 {
		sl.ReportError(
			s.ValueMapping,
			"ValueMapping",
			"value_mapping",
			"required_for_bool_raw",
			"",
		)
	}
}
