package volcanoTTS

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
)

type Adaptor struct{}

func (a *Adaptor) Init(meta *meta.Meta) {
	// No initialization needed
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	// Volcano TTS v3 unidirectional streaming API
	return "https://openspeech.bytedance.com/api/v3/tts/unidirectional", nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)

	// Volcano TTS uses custom headers for authentication
	// The APIKey field contains JSON-encoded credentials:
	// {"app_id": "xxx", "access_key": "yyy", "resource_id": "zzz"}
	config, err := ParseVolcanoConfig(meta.APIKey)
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-App-Id", config.AppID)
	req.Header.Set("X-Api-Access-Key", config.AccessKey)
	req.Header.Set("X-Api-Resource-Id", config.ResourceID)
	req.Header.Set("X-Api-Request-Id", generateRequestID())
	req.Header.Set("X-Control-Require-Usage-Tokens-Return", "text_words")
	req.Header.Set("Content-Type", "application/json")

	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	// This is called for chat/completion requests, not used for TTS
	return request, nil
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	return nil, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	// For TTS, we stream the response directly back to the client
	// The actual streaming is handled in the TTS controller
	return nil, nil
}

func (a *Adaptor) GetModelList() []string {
	return []string{
		"volcano-tts",
		"seed-tts-1.1",
	}
}

func (a *Adaptor) GetChannelName() string {
	return "Volcano TTS"
}
