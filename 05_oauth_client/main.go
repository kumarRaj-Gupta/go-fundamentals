package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
)

// --- CONFIGURATION CONSTANTS (Hardcoded for demo) ---

// In a real app, these come from environment variables or a config file.
const (
	clientID     = "YOUR_CLIENT_ID"
	clientSecret = "YOUR_CLIENT_SECRET"
	redirectURL  = "http://localhost:8080/callback"
	tokenPath    = "token.json" // File to store the access and refresh tokens
)

// Global OAuth Config
var config = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  redirectURL,
	Scopes:       []string{"https://www.googleapis.com/auth/drive.readonly"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
}

// ----------------------------------------------------

// saveToken persists the OAuth2 token to a file as JSON. (Replaces Python's pickle)
func saveToken(token *oauth2.Token) error {
	log.Printf("Attempting to save token to %s...", tokenPath)

	// 1. --- MARSHAL THE TOKEN TO JSON ---
	b, err := json.Marshal(token)
	if err != nil {
		// Using %w for error wrapping: allows caller to inspect the underlying error.
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// 2. --- WRITE JSON BYTES TO FILE ---
	// Use 0600 permission (owner read/write) for security.
	if err := os.WriteFile(tokenPath, b, 0600); err != nil {
		// Reusing the err variable in the if statement above.
		return fmt.Errorf("failed to write token file: %w", err)
	}

	log.Printf("Token successfully saved.")
	return nil
}

// loadToken loads the OAuth2 token from a JSON file.
func loadToken() (*oauth2.Token, error) {
	// Task 2: Implement this function in the next step.
	b, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}
	token := new(oauth2.Token) // or &oauth2.Token{}

	err = json.Unmarshal(b, token)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling token from %v ERROR: %w", tokenPath, err)
	}
	return token, nil
}

// main simulates the token lifecycle
func main() {
	fmt.Println("--- Project 5: OAuth2 Token Management ---")

	_ = context.Background() // Keep context here for future steps

	// ------------------------------------------------------------------------------------------------
	// NOTE: We cannot perform the full web-based OAuth flow here, so we manually create a mock token.
	// ------------------------------------------------------------------------------------------------

	mockToken := &oauth2.Token{
		AccessToken:  "mock-access-token-12345",
		TokenType:    "Bearer",
		RefreshToken: "mock-refresh-token-ABCDE",
		Expiry:       time.Now().Add(time.Hour).Round(0), // Valid for 1 hour
	}

	// TASK 1 EXECUTION: Save the mock token
	if err := saveToken(mockToken); err != nil {
		log.Fatalf("Fatal error during token save: %v", err)
	}

	// --- TASK 2 (NEXT): Load the token ---
	token, err := loadToken()
	if err != nil {
		log.Fatalf("Fatal error during token load: %v", err)
	}
	log.Printf("Loaded Token Refresh Token: %s", token.RefreshToken)
}
