package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	// calendarOAuthConfig is the OAuth2 configuration for Google Calendar.
	calendarOAuthConfig *oauth2.Config
	// calendarState is a token to protect against CSRF attacks.
	calendarState = "random-string"
)

// initCalendarAuth initializes the Google Calendar OAuth2 configuration.
func initCalendarAuth() error {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	calendarOAuthConfig = config
	return nil
}

// getCalendarService creates a new Google Calendar service with a client that is
// authenticated with a valid token.
func getCalendarService(ctx context.Context) (*calendar.Service, error) {
	tok, err := tokenFromFile("token.json")
	if err != nil {
		return nil, err
	}

	client := calendarOAuthConfig.Client(ctx, tok)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}
	return srv, nil
}

// handleCalendarOAuth initiates the OAuth2 flow for Google Calendar.
func handleCalendarOAuth(w http.ResponseWriter, r *http.Request) {
	url := calendarOAuthConfig.AuthCodeURL(calendarState, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleCalendarCallback handles the OAuth2 callback from Google.
func handleCalendarCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != calendarState {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	tok, err := calendarOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("unable to retrieve token from web: %v", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	saveToken("token.json", tok)
	http.Redirect(w, r, "/?google_auth_success=true", http.StatusTemporaryRedirect)
}


// tokenFromFile retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// CreateCalendarEvent creates a new event in the user's primary calendar.
func CreateCalendarEvent(srv *calendar.Service, summary, description string, start, end time.Time) (*calendar.Event, error) {
	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
				DateTime: start.Format(time.RFC3339),
				TimeZone: "America/Los_Angeles", // You might want to make this configurable
			},
			End: &calendar.EventDateTime{
				DateTime: end.Format(time.RFC3339),
				TimeZone: "America/Los_Angeles", // You might want to make this configurable
			},
		}

	return srv.Events.Insert("primary", event).Do()
}
