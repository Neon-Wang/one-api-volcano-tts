package volcanoTTS

// VolcanoTTSRequest represents the request format for Volcano TTS API v3
type VolcanoTTSRequest struct {
	User      VolcanoUser      `json:"user"`
	ReqParams VolcanoReqParams `json:"req_params"`
}

type VolcanoUser struct {
	UID string `json:"uid"`
}

type VolcanoReqParams struct {
	Text        string            `json:"text"`
	Speaker     string            `json:"speaker"`
	Model       string            `json:"model"`
	AudioParams VolcanoAudioParams `json:"audio_params"`
	Additions   VolcanoAdditions  `json:"additions,omitempty"`
}

type VolcanoAudioParams struct {
	Format       string `json:"format"`
	SampleRate   int    `json:"sample_rate"`
	Emotion      string `json:"emotion,omitempty"`
	EmotionScale int    `json:"emotion_scale,omitempty"`
}

type VolcanoAdditions struct {
	DisableMarkdownFilter  bool               `json:"disable_markdown_filter"`
	DisableEmojiFilter     bool               `json:"disable_emoji_filter"`
	EnableLanguageDetector bool               `json:"enable_language_detector"`
	CacheConfig            VolcanoCacheConfig `json:"cache_config,omitempty"`
}

type VolcanoCacheConfig struct {
	TextType int  `json:"text_type"`
	UseCache bool `json:"use_cache"`
}

// VoiceMapping maps OpenAI voice names to Volcano speaker IDs
var VoiceMapping = map[string]string{
	// OpenAI standard voices
	"alloy":   "zh_female_shuangkuaisisi_moon_bigtts",
	"echo":    "zh_male_chunhou_moon_bigtts",
	"fable":   "zh_female_qiaopi_moon_bigtts",
	"onyx":    "zh_male_wennuan_moon_bigtts",
	"nova":    "zh_female_shuangkuaisisi_moon_bigtts",
	"shimmer": "zh_female_tianmei_moon_bigtts",

	// Allow direct Volcano speaker IDs
	"zh_female_shuangkuaisisi_moon_bigtts": "zh_female_shuangkuaisisi_moon_bigtts",
	"zh_male_chunhou_moon_bigtts":          "zh_male_chunhou_moon_bigtts",
	"zh_female_qiaopi_moon_bigtts":         "zh_female_qiaopi_moon_bigtts",
	"zh_male_wennuan_moon_bigtts":          "zh_male_wennuan_moon_bigtts",
	"zh_female_tianmei_moon_bigtts":        "zh_female_tianmei_moon_bigtts",
}

// GetVolcanoSpeaker converts an OpenAI voice name to a Volcano speaker ID
func GetVolcanoSpeaker(voice string) string {
	if speaker, ok := VoiceMapping[voice]; ok {
		return speaker
	}
	// Default speaker
	return "zh_female_shuangkuaisisi_moon_bigtts"
}

// ConvertOpenAIToVolcano converts an OpenAI TTS request to Volcano format
func ConvertOpenAIToVolcano(input string, voice string, userID string) *VolcanoTTSRequest {
	speaker := GetVolcanoSpeaker(voice)

	return &VolcanoTTSRequest{
		User: VolcanoUser{
			UID: userID,
		},
		ReqParams: VolcanoReqParams{
			Text:    input,
			Speaker: speaker,
			Model:   "seed-tts-1.1",
			AudioParams: VolcanoAudioParams{
				Format:     "mp3",
				SampleRate: 24000,
			},
			Additions: VolcanoAdditions{
				DisableMarkdownFilter:  true,
				DisableEmojiFilter:     true,
				EnableLanguageDetector: true,
				CacheConfig: VolcanoCacheConfig{
					TextType: 1,
					UseCache: true,
				},
			},
		},
	}
}
