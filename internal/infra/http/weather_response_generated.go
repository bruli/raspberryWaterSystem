// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package http

import "encoding/json"
import "fmt"

// This schema defines the weather response
type WeatherResponseJson struct {
	// Humidity corresponds to the JSON schema field "humidity".
	Humidity float64 `json:"humidity" yaml:"humidity" mapstructure:"humidity"`

	// IsRaining corresponds to the JSON schema field "is_raining".
	IsRaining bool `json:"is_raining" yaml:"is_raining" mapstructure:"is_raining"`

	// Temperature corresponds to the JSON schema field "temperature".
	Temperature float64 `json:"temperature" yaml:"temperature" mapstructure:"temperature"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *WeatherResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["humidity"]; !ok || v == nil {
		return fmt.Errorf("field humidity in WeatherResponseJson: required")
	}
	if v, ok := raw["is_raining"]; !ok || v == nil {
		return fmt.Errorf("field is_raining in WeatherResponseJson: required")
	}
	if v, ok := raw["temperature"]; !ok || v == nil {
		return fmt.Errorf("field temperature in WeatherResponseJson: required")
	}
	type Plain WeatherResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = WeatherResponseJson(plain)
	return nil
}
