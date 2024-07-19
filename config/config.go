package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// Config defines the struct of a configuration in general.
type Config struct {
	Data map[string]interface{}
	Raw  []byte
}

func newConfig() *Config {
	result := new(Config)
	result.Data = make(map[string]interface{})
	return result
}

// LoadConfigFile loads config information from a JSON file.
func LoadConfigFile(filename string) (*Config, error) {
	result := newConfig()
	err := result.parse(filename)
	if err != nil {
		log.Printf("error loading config file %s: %s", filename, err)
	}
	return result, err
}

func (c *Config) SaveConfigFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	err = json.NewEncoder(file).Encode(c.Data)
	if err != nil {
		return err
	}
	defer func() {
		file.Sync()
		file.Close()
	}()
	return err
}

// LoadConfigString loads config information from a JSON string.
func LoadConfigString(s string) (*Config, error) {
	result := newConfig()
	decoder := json.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()
	err := decoder.Decode(&result.Data)
	if err != nil {
		log.Fatalf("error parsing config string %s: %s", s, err)
	}
	return result, err
}

func (c *Config) parse(fileName string) error {
	jsonFileBytes, err := os.ReadFile(fileName)
	c.Raw = jsonFileBytes
	if err == nil {
		decoder := json.NewDecoder(strings.NewReader(string(jsonFileBytes)))
		decoder.UseNumber()
		err = decoder.Decode(&c.Data)
	}
	return err
}

// GetString returns a string for the config key.
func (c *Config) GetString(key string) string {
	x, present := c.Data[key]
	if !present {
		return ""
	}
	if result, isString := x.(string); isString {
		return result
	}
	return ""
}

// GetString returns a string for the config key.
func (c *Config) SetString(key, val string) {
	c.Data[key] = val
}

// GetFloat returns a float value for the config key.
func (c *Config) GetFloat(key string) float64 {
	x, present := c.Data[key]
	if !present {
		return -1
	}
	if result, isNumber := x.(json.Number); isNumber {
		number, err := result.Float64()
		if err != nil {
			return 0
		}
		return number
	}
	return 0
}

// returns a bool value for the config key with default val when not present
func (c *Config) GetBoolWithDefault(key string, defval bool) bool {
	_, present := c.Data[key]
	if !present {
		return defval
	}
	return c.GetBool(key)
}

// GetBool returns a bool value for the config key.
func (c *Config) GetBool(key string) bool {
	x, present := c.Data[key]
	if !present {
		return false
	}
	if result, isBool := x.(bool); isBool {
		return result
	}
	if result, isString := x.(string); isString {
		if result == "true" {
			return true
		}
	}
	return false
}

// GetInt returns a int value for the config key.
func (c *Config) GetInt(key string) int {
	return int(c.GetInt64(key))
}

// GetInt64 returns a int64 value for the config key.
func (c *Config) GetInt64(key string) int64 {
	x, present := c.Data[key]
	if !present {
		return 0
	}
	if result, isNumber := x.(json.Number); isNumber {
		number, err := result.Int64()
		if err != nil {
			return 0
		}
		return number
	}
	return 0
}

func (c *Config) HasKey(key string) bool {
	_, present := c.Data[key]
	return present
}

// GetInt64WithDefault returns a int64 value for the config key.
func (c *Config) GetInt64WithDefault(key string, defaultVal int64) int64 {
	if val := c.GetInt64(key); val == 0 {
		return defaultVal
	} else {
		return val
	}
}

// GetSlice returns an array for the config key.
func (c *Config) GetSlice(key string) []interface{} {
	result, present := c.Data[key]
	if !present {
		return []interface{}(nil)
	}
	return result.([]interface{})
}

func (c *Config) GetStringSlice(key string) []string {
	s := c.GetSlice(key)
	result := make([]string, 0, len(s))
	for _, item := range s {
		result = append(result, item.(string))
	}
	return result
}

// Check and get a string for the config key.
func (c *Config) CheckAndGetString(key string) (string, bool) {
	x, present := c.Data[key]
	if !present {
		return "", false
	}
	if result, isString := x.(string); isString {
		return result, true
	}
	return "", false
}

// GetBool returns a bool value for the config key.
func (c *Config) CheckAndGetBool(key string) (bool, bool) {
	x, present := c.Data[key]
	if !present {
		return false, false
	}
	if result, isBool := x.(bool); isBool {
		return result, true
	}
	// Take string value "true" and "false" as well.
	if result, isString := x.(string); isString {
		if result == "true" {
			return true, true
		}
		if result == "false" {
			return false, true
		}
	}
	return false, false
}
