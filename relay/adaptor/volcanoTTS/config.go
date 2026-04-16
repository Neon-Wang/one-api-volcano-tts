package volcanoTTS

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// VolcanoConfig holds the credentials for Volcano TTS API
type VolcanoConfig struct {
	AppID      string `json:"app_id"`
	AccessKey  string `json:"access_key"`
	ResourceID string `json:"resource_id"`
}

// ParseVolcanoConfig parses the API key field which contains JSON-encoded credentials
func ParseVolcanoConfig(apiKey string) (*VolcanoConfig, error) {
	var config VolcanoConfig
	if err := json.Unmarshal([]byte(apiKey), &config); err != nil {
		return nil, fmt.Errorf("invalid Volcano TTS config: %w", err)
	}

	if config.AppID == "" || config.AccessKey == "" {
		return nil, fmt.Errorf("missing required fields: app_id and access_key")
	}

	// Default resource ID for TTS service
	if config.ResourceID == "" {
		config.ResourceID = "volc.service_type.10029"
	}

	return &config, nil
}

// generateRequestID creates a unique request ID for tracing
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
