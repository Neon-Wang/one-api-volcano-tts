package volcanoTTS

// VolcanoTTSChunk represents a single chunk in the streaming response
type VolcanoTTSChunk struct {
	Code    int               `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Data    string            `json:"data,omitempty"` // Base64 encoded audio
	Usage   *VolcanoTTSUsage  `json:"usage,omitempty"`
	LogID   string            `json:"logid,omitempty"`
}

type VolcanoTTSUsage struct {
	TextWords int `json:"text_words"`
}

// IsSuccess checks if the response indicates success
func (c *VolcanoTTSChunk) IsSuccess() bool {
	return c.Code == 20000000 || c.Code == 0
}

// HasAudio checks if the chunk contains audio data
func (c *VolcanoTTSChunk) HasAudio() bool {
	return c.Data != "" && len(c.Data) > 0
}

// GetUsageTokens returns the usage in tokens (text_words)
// This is used for quota deduction
func (c *VolcanoTTSChunk) GetUsageTokens() int {
	if c.Usage != nil {
		return c.Usage.TextWords
	}
	return 0
}
