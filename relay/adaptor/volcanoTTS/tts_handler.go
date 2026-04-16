package volcanoTTS

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/client"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
)

// HandleTTSRequest processes TTS requests for Volcano TTS
// This function is called from the TTS controller when the channel type is VolcanoTTS
func HandleTTSRequest(c *gin.Context, config *VolcanoConfig, ttsRequest *openai.TextToSpeechRequest, userID string) error {
	// Convert OpenAI request to Volcano format
	volcanoReq := ConvertOpenAIToVolcano(ttsRequest.Input, ttsRequest.Voice, userID)

	// Serialize request body
	reqBody, err := json.Marshal(volcanoReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://openspeech.bytedance.com/api/v3/tts/unidirectional", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Api-App-Id", config.AppID)
	req.Header.Set("X-Api-Access-Key", config.AccessKey)
	req.Header.Set("X-Api-Resource-Id", config.ResourceID)
	req.Header.Set("X-Api-Request-Id", generateRequestID())
	req.Header.Set("X-Control-Require-Usage-Tokens-Return", "text_words")
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upstream error: %d - %s", resp.StatusCode, string(body))
	}

	// Stream response back to client
	// Volcano returns NDJSON with base64-encoded audio chunks
	// We need to decode and stream as raw audio

	c.Writer.Header().Set("Content-Type", "audio/mpeg")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.WriteHeader(http.StatusOK)

	var totalUsage int
	scanner := bufio.NewScanner(resp.Body)
	// Increase buffer size for large audio chunks
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var chunk VolcanoTTSChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		// Check for errors
		if chunk.Code != 0 && chunk.Code != 20000000 && chunk.Message != "" {
			return fmt.Errorf("volcano error: %s", chunk.Message)
		}

		// Write audio data
		if chunk.HasAudio() {
			// Decode base64 and write raw audio
			audioData, err := decodeBase64(chunk.Data)
			if err != nil {
				continue
			}
			c.Writer.Write(audioData)
			c.Writer.Flush()
		}

		// Track usage
		if chunk.Usage != nil {
			totalUsage = chunk.Usage.TextWords
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	// Set usage header for billing
	c.Writer.Header().Set("X-TTS-Usage-Tokens", strconv.Itoa(totalUsage))

	return nil
}

func decodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
