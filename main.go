
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

// TutorRequest represents the request body for the /api/tutor endpoint.
	type TutorRequest struct {
		Prompt string `json:"prompt"`
	}

// APIResponse is a generic struct for sending JSON responses.
	type APIResponse struct {
		Data  interface{} `json:"data,omitempty"`
		Error string      `json:"error,omitempty"`
	}

// newAPIResponse creates a new APIResponse.
	func newAPIResponse(data interface{}, err string) *APIResponse {
		return &APIResponse{Data: data, Error: err}
	}

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
			log.Fatalf("Failed to create GenAI client: %v", err)
		}
		defer genaiClient.Close()

		youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))
		if err != nil {
			log.Fatalf("Error creating YouTube service: %v", err)
		}

		model := genaiClient.GenerativeModel("gemini-1.5-pro-latest")

		http.HandleFunc("/api/tutor", handleTutorRequest(model, ctx))
		http.HandleFunc("/api/youtube", handleYoutubeSearch(youtubeService))

		// Serve the frontend files.
		fs := http.FileServer(http.Dir("."))
		http.Handle("/", fs)

		log.Println("Server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}

	func handleTutorRequest(model *genai.GenerativeModel, ctx context.Context) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				writeJSONError(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
			}

			var reqBody TutorRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				writeJSONError(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			resp, err := model.GenerateContent(ctx, genai.Text(reqBody.Prompt))
			if err != nil {
				writeJSONError(w, "Failed to generate content", http.StatusInternalServerError)
				return
			}

			writeJSON(w, resp, http.StatusOK)
		}
	}

	func handleYoutubeSearch(youtubeService *youtube.Service) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("q")
			if query == "" {
				writeJSONError(w, "Missing search query", http.StatusBadRequest)
				return
			}

			call := youtubeService.Search.List([]string{"snippet"}).Q(query).MaxResults(10)
			response, err := call.Do()
			if err != nil {
				log.Printf("Error calling YouTube API: %v", err)
				writeJSONError(w, "Error calling YouTube API", http.StatusInternalServerError)
				return
			}

			writeJSON(w, response, http.StatusOK)
		}
	}

	func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(newAPIResponse(data, "")); err != nil {
			log.Printf("Failed to write JSON response: %v", err)
		}
	}

	func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(newAPIResponse(nil, message)); err != nil {
			log.Printf("Failed to write JSON error response: %v", err)
		}
	}
