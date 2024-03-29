// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package http

import "encoding/json"
import "fmt"

// This schema defines the request to create a zone
type CreateZoneRequestJson struct {
	// Id corresponds to the JSON schema field "id".
	Id string `json:"id" yaml:"id" mapstructure:"id"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Relays corresponds to the JSON schema field "relays".
	Relays []int `json:"relays" yaml:"relays" mapstructure:"relays"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateZoneRequestJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["id"]; !ok || v == nil {
		return fmt.Errorf("field id in CreateZoneRequestJson: required")
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in CreateZoneRequestJson: required")
	}
	if v, ok := raw["relays"]; !ok || v == nil {
		return fmt.Errorf("field relays in CreateZoneRequestJson: required")
	}
	type Plain CreateZoneRequestJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CreateZoneRequestJson(plain)
	return nil
}
