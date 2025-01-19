package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

// Config struct with fields for application configuration
type Config struct {
	BanDuration           time.Duration `json:"ban_duration"`            // ban duration
	UpdateInterval        time.Duration `json:"update_interval"`         // auto-update ban list interval
	XRayLogsFile          string        `json:"xray_logs_file"`          // XRay access.log file path
	XRayBlacklistOutbound string        `json:"xray_blacklist_outbound"` // XRay outbound tag for IPs be banned
	BannedDatabaseFile    string        `json:"banned_database_file"`    // banned database file
	DebugMode             bool          `json:"debug_mode"`              // debug mode, don't ban anybody, just log
}

//go:embed default_config.json
var defaultConfigFile embed.FS

// LoadConfig loads the configuration from a file or falls back to the embedded default
func LoadConfig(filePath string) (*Config, error) {
	var data []byte
	var err error

	if filePath != "" {
		// Try loading from the specified file
		data, err = os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to load config from file: %v, falling back to default", err)
		}
	}

	if data == nil {
		// Load the embedded default config
		data, err = defaultConfigFile.ReadFile("default_config.json")
		if err != nil {
			return nil, fmt.Errorf("failed to load embedded default config: %w", err)
		}
	}

	// Parse the YAML data. SYKA BLYAT wtf
	var jsonKeyValueMap = make(map[string]interface{})
	if err := json.Unmarshal(data, &jsonKeyValueMap); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	config := Config{}
	configValue := reflect.ValueOf(&config)
	configType := reflect.TypeOf(config)

	for i := 0; i < configType.NumField(); i++ {
		fieldType := configType.Field(i)
		fieldJsonTag := fieldType.Tag.Get("json")
		fieldValue := configValue.Elem().Field(i)
		if fieldJsonTag == "" {
			return nil, fmt.Errorf("struct %s does not have json tag for field %s", configType.String(), fieldType.Name)
		}

		fieldJsonValue, jsonKeyExists := jsonKeyValueMap[fieldJsonTag]
		if !jsonKeyExists {
			return nil, fmt.Errorf("config json does not have tag %s for field %s", fieldJsonTag, fieldType.Name)
		}

		if fieldValue.Kind() == reflect.TypeOf(fieldJsonValue).Kind() {
			fieldValue.Set(reflect.ValueOf(fieldJsonValue))
		} else {
			var stubDuration time.Duration
			if fieldValue.Kind() == reflect.TypeOf(stubDuration).Kind() {
				stubDuration, err := time.ParseDuration(fieldJsonValue.(string))
				if err != nil {
					return nil, fmt.Errorf("could not parse %s json tags for field %s", fieldJsonTag, fieldType.Name)
				}
				fieldValue.Set(reflect.ValueOf(stubDuration))
			}
		}
	}

	return &config, nil
}

// parse
