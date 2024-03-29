// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package http

import "encoding/json"
import "fmt"

type ExecutionItemResponse struct {
	// Seconds corresponds to the JSON schema field "seconds".
	Seconds int `json:"seconds" yaml:"seconds" mapstructure:"seconds"`

	// Zones corresponds to the JSON schema field "zones".
	Zones []string `json:"zones" yaml:"zones" mapstructure:"zones"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExecutionItemResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["seconds"]; !ok || v == nil {
		return fmt.Errorf("field seconds in ExecutionItemResponse: required")
	}
	if v, ok := raw["zones"]; !ok || v == nil {
		return fmt.Errorf("field zones in ExecutionItemResponse: required")
	}
	type Plain ExecutionItemResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ExecutionItemResponse(plain)
	return nil
}

type ProgramItemResponse struct {
	// Executions corresponds to the JSON schema field "executions".
	Executions []ExecutionItemResponse `json:"executions" yaml:"executions" mapstructure:"executions"`

	// Hour corresponds to the JSON schema field "hour".
	Hour string `json:"hour" yaml:"hour" mapstructure:"hour"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProgramItemResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["executions"]; !ok || v == nil {
		return fmt.Errorf("field executions in ProgramItemResponse: required")
	}
	if v, ok := raw["hour"]; !ok || v == nil {
		return fmt.Errorf("field hour in ProgramItemResponse: required")
	}
	type Plain ProgramItemResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ProgramItemResponse(plain)
	return nil
}

type TemperatureItemResponse struct {
	// Programs corresponds to the JSON schema field "programs".
	Programs []ProgramItemResponse `json:"programs" yaml:"programs" mapstructure:"programs"`

	// Temperature corresponds to the JSON schema field "temperature".
	Temperature float64 `json:"temperature" yaml:"temperature" mapstructure:"temperature"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TemperatureItemResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["programs"]; !ok || v == nil {
		return fmt.Errorf("field programs in TemperatureItemResponse: required")
	}
	if v, ok := raw["temperature"]; !ok || v == nil {
		return fmt.Errorf("field temperature in TemperatureItemResponse: required")
	}
	type Plain TemperatureItemResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = TemperatureItemResponse(plain)
	return nil
}

type WeeklyItemResponse struct {
	// Programs corresponds to the JSON schema field "programs".
	Programs []ProgramItemResponse `json:"programs" yaml:"programs" mapstructure:"programs"`

	// WeekDay corresponds to the JSON schema field "week_day".
	WeekDay string `json:"week_day" yaml:"week_day" mapstructure:"week_day"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *WeeklyItemResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["programs"]; !ok || v == nil {
		return fmt.Errorf("field programs in WeeklyItemResponse: required")
	}
	if v, ok := raw["week_day"]; !ok || v == nil {
		return fmt.Errorf("field week_day in WeeklyItemResponse: required")
	}
	type Plain WeeklyItemResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = WeeklyItemResponse(plain)
	return nil
}

// This schema defines the programs response
type ProgramsResponseJson struct {
	// Daily corresponds to the JSON schema field "daily".
	Daily []ProgramItemResponse `json:"daily" yaml:"daily" mapstructure:"daily"`

	// Even corresponds to the JSON schema field "even".
	Even []ProgramItemResponse `json:"even" yaml:"even" mapstructure:"even"`

	// Odd corresponds to the JSON schema field "odd".
	Odd []ProgramItemResponse `json:"odd" yaml:"odd" mapstructure:"odd"`

	// Temperature corresponds to the JSON schema field "temperature".
	Temperature []TemperatureItemResponse `json:"temperature" yaml:"temperature" mapstructure:"temperature"`

	// Weekly corresponds to the JSON schema field "weekly".
	Weekly []WeeklyItemResponse `json:"weekly" yaml:"weekly" mapstructure:"weekly"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProgramsResponseJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["daily"]; !ok || v == nil {
		return fmt.Errorf("field daily in ProgramsResponseJson: required")
	}
	if v, ok := raw["even"]; !ok || v == nil {
		return fmt.Errorf("field even in ProgramsResponseJson: required")
	}
	if v, ok := raw["odd"]; !ok || v == nil {
		return fmt.Errorf("field odd in ProgramsResponseJson: required")
	}
	if v, ok := raw["temperature"]; !ok || v == nil {
		return fmt.Errorf("field temperature in ProgramsResponseJson: required")
	}
	if v, ok := raw["weekly"]; !ok || v == nil {
		return fmt.Errorf("field weekly in ProgramsResponseJson: required")
	}
	type Plain ProgramsResponseJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ProgramsResponseJson(plain)
	return nil
}
