package handler

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/openai-go"
	"github.com/charmbracelet/openai-go/option"
	"github.com/gin-gonic/gin"
)

type TranscribeHandler struct {
	client *openai.Client
	model  openai.AudioModel
}

func NewTranscribeHandler(apiKey string, model string) *TranscribeHandler {
	if model == "" {
		model = openai.AudioModelWhisper1
	}
	c := openai.NewClient(option.WithAPIKey(apiKey))
	return &TranscribeHandler{
		client: &c,
		model:  openai.AudioModel(model),
	}
}

func (h *TranscribeHandler) HandleTranscribe(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file in multipart form"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open uploaded file: %v", err)})
		return
	}
	defer f.Close()

	resp, err := h.client.Audio.Transcriptions.New(c.Request.Context(), openai.AudioTranscriptionNewParams{
		File:  f,
		Model: h.model,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Transcription failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"text": resp.Text,
	})
}
