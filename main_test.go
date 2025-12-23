
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// TestTutorHandler tests the /api/tutor endpoint.
func TestTutorHandler(t *testing.T) {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		t.Skip("GEMINI_API_KEY not set, skipping tutor handler test")
	}

	ctx := context.Background()
	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		t.Fatalf("Failed to create GenAI client: %v", err)
	}

	model := genaiClient.GenerativeModel("gemini-pro")

	h := handleTutorRequest(model, ctx)

	// Create a mock request
	reqBody := TutorRequest{Prompt: "testing prompt"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/tutor", strings.NewReader(string(body)))
	w := httptest.NewRecorder()

	h(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// TestYoutubeSearchHandler tests the /api/youtube endpoint.
func TestYoutubeSearchHandler(t *testing.T) {
	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	if youtubeAPIKey == "" {
		t.Skip("YOUTUBE_API_KEY not set, skipping YouTube search handler test")
	}

	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))
	if err != nil {
		t.Fatalf("Error creating YouTube service: %v", err)
	}

	h := handleYoutubeSearch(youtubeService)

	// Create a mock request with a query
	req := httptest.NewRequest(http.MethodGet, "/api/youtube?q=testing", nil)
	w := httptest.NewRecorder()

	h(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
