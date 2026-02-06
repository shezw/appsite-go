package dbs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Map represents a JSON object
type Map map[string]interface{}

func (m Map) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Map) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	
	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*m = Map(result)
	return err
}

// Slice represents a JSON array
type Slice []interface{}

func (s Slice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *Slice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result []interface{}
	err := json.Unmarshal(bytes, &result)
	*s = Slice(result)
	return err
}

// StringArray represents a JSON string array (e.g. tags)
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result []string
	err := json.Unmarshal(bytes, &result)
	*s = StringArray(result)
	return err
}
