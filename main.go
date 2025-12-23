
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/genai"
)

func main() {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set.")
	}

	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	if youtubeAPIKey == "" {
		log.Fatal("YOUTUBE_API_KEY environment variable not set.")
	}

	ctx := context.Background()
	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: geminiAPIKey,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer genaiClient.Close()

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	model := genaiClient.GenerativeModel("gemini-1.5-pro-latest")

	http.HandleFunc("/api/tutor", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var reqBody struct {
			Prompt string `json:"prompt"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := model.GenerateContent(ctx, genai.Text(reqBody.Prompt))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/api/youtube", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Missing search query", http.StatusBadRequest)
			return
		}

		call := youtubeService.Search.List([]string{"snippet"}).Q(query).MaxResults(10)
		response, err := call.Do()
		if err != nil {
			http.Error(w, "Error calling YouTube API", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Serve the frontend files
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
