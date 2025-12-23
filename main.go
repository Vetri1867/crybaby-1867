package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// TutorRequest represents the request body for the /api/tutor endpoint.
	type TutorRequest struct {
		Prompt string `json:"prompt"`
	}

// CalendarEventRequest represents the request body for the /api/calendar/event endpoint.
	type CalendarEventRequest struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
		Start       string `json:"start"`
		End         string `json:"end"`
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
		// Load environment variables from .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using environment variables from OS")
		}

		if err := initCalendarAuth(); err != nil {
			log.Fatalf("Failed to initialize calendar auth: %v", err)
		}


		geminiAPIKey := os.Getenv("GEMINI_API_KEY")
		if geminiAPIKey == "" {
			log.Fatal("GEMINI_API_KEY environment variable not set.")
		}

		youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
		if youtubeAPIKey == "" {
			log.Fatal("YOUTUBE_API_KEY environment variable not set.")
		}

		ctx := context.Background()
		genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
		if err != nil {
			log.Fatalf("Failed to create GenAI client: %v", err)
		}

		youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))
		if err != nil {
			log.Fatalf("Error creating YouTube service: %v", err)
		}

		model := genaiClient.GenerativeModel("gemini-pro")

		http.HandleFunc("/api/tutor", handleTutorRequest(model, ctx))
		http.HandleFunc("/api/youtube", handleYoutubeSearch(youtubeService))

		// Add handlers for Google Calendar authentication
		http.HandleFunc("/login/google", handleCalendarOAuth)
		http.HandleFunc("/callback/google", handleCalendarCallback)

		// Add handler for creating a calendar event
		http.HandleFunc("/api/calendar/event", handleCreateCalendarEvent)


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

	func handleCreateCalendarEvent(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONError(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
		}

		var reqBody CalendarEventRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				writeJSONError(w, "Invalid request body", http.StatusBadRequest)
				return
		}

		// Basic validation
		if reqBody.Summary == "" || reqBody.Start == "" || reqBody.End == "" {
			writeJSONError(w, "Missing required fields", http.StatusBadRequest)
				return
		}


		layout := time.RFC3339
		startTime, err := time.Parse(layout, reqBody.Start)
		if err != nil {
			writeJSONError(w, "Invalid start time format", http.StatusBadRequest)
				return
		}

		endTime, err := time.Parse(layout, reqBody.End)
		if err != nil {
			writeJSONError(w, "Invalid end time format", http.StatusBadRequest)
				return
		}


		srv, err := getCalendarService(r.Context())
		if err != nil {
			// This could mean the token is missing or expired.
			// The frontend will need to handle this by prompting the user to log in.
			writeJSONError(w, "Not authenticated with Google Calendar. Please log in.", http.StatusUnauthorized)
			return
		}

		event, err := CreateCalendarEvent(srv, reqBody.Summary, reqBody.Description, startTime, endTime)
		if err != nil {
			log.Printf("Error creating calendar event: %v", err)
			writeJSONError(w, "Failed to create calendar event", http.StatusInternalServerError)
			return
		}

		writeJSON(w, event, http.StatusCreated)
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
