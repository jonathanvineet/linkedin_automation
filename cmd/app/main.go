package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	
	"github.com/jonathanvineet/linkedin-automation/internal/auth"
	"github.com/jonathanvineet/linkedin-automation/internal/behavior"
	"github.com/jonathanvineet/linkedin-automation/internal/browser"
	"github.com/jonathanvineet/linkedin-automation/internal/connect"
	"github.com/jonathanvineet/linkedin-automation/internal/logger"
	"github.com/jonathanvineet/linkedin-automation/internal/messaging"
	"github.com/jonathanvineet/linkedin-automation/internal/search"
	"github.com/jonathanvineet/linkedin-automation/internal/state"
	"github.com/jonathanvineet/linkedin-automation/internal/stealth"
)

// App holds application state
type App struct {
	store     *state.Store
	session   *browser.Session
	persona   *behavior.Persona
	scheduler *stealth.Scheduler
	isRunning bool
}

var app *App

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logLevel := getEnv("LOG_LEVEL", "info")
	if err := logger.Initialize(logLevel, "./logs/app.log"); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	logger.Log.Info("🚀 LinkedIn Automation PoC starting...")
	logger.Log.Warn("⚠️  EDUCATIONAL USE ONLY - DO NOT USE IN PRODUCTION")

	// Initialize database
	dbPath := getEnv("DB_PATH", "./data/automation.db")
	store, err := state.NewStore(dbPath)
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to initialize database")
	}
	defer store.Close()

	// Initialize app
	app = &App{
		store:     store,
		isRunning: false,
	}

	// Setup API routes
	router := mux.NewRouter()
	
	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/status", handleStatus).Methods("GET")
	apiRouter.HandleFunc("/start", handleStart).Methods("POST")
	apiRouter.HandleFunc("/stop", handleStop).Methods("POST")
	apiRouter.HandleFunc("/stats", handleStats).Methods("GET")
	apiRouter.HandleFunc("/activity", handleActivity).Methods("GET")
	apiRouter.HandleFunc("/persona", handlePersona).Methods("POST")
	apiRouter.HandleFunc("/search", handleSearch).Methods("POST")
	apiRouter.HandleFunc("/connect", handleConnect).Methods("POST")
	apiRouter.HandleFunc("/message", handleMessage).Methods("POST")

	// CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start server
	port := getEnv("API_PORT", "8090")
	addr := ":" + port
	
	logger.Log.WithField("port", port).Info("API server starting")
	log.Fatal(http.ListenAndServe(addr, corsHandler.Handler(router)))
}

// API Handlers

func handleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"running":    app.isRunning,
		"logged_in":  app.session != nil,
		"persona":    "none",
		"stealth":    true,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	
	if app.persona != nil {
		status["persona"] = app.persona.Name
	}
	
	jsonResponse(w, status)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if app.isRunning {
		jsonError(w, "Automation already running", http.StatusBadRequest)
		return
	}

	// Start browser session
	config := browser.Config{
		Headless:       getEnvBool("HEADLESS_MODE", false),
		UserAgent:      getEnv("USER_AGENT", browser.GetRandomUserAgent()),
		ViewportWidth:  1920,
		ViewportHeight: 1080,
		Timeout:        30 * time.Second,
	}

	session, err := browser.NewSession(config)
	if err != nil {
		jsonError(w, fmt.Sprintf("Failed to start browser: %v", err), http.StatusInternalServerError)
		return
	}

	app.session = session
	
	// Initialize persona
	personaType := getEnv("DEFAULT_PERSONA", "recruiter")
	app.persona = behavior.GetDefaultPersona(behavior.PersonaType(personaType))
	
	// Initialize scheduler
	app.scheduler = stealth.NewScheduler(
		getEnvBool("BUSINESS_HOURS_ONLY", true),
		9,  // start hour
		17, // end hour
		2,  // cooldown minutes
	)

	// Perform login
	email := getEnv("LINKEDIN_EMAIL", "")
	password := getEnv("LINKEDIN_PASSWORD", "")
	
	if email == "" || password == "" {
		jsonError(w, "LinkedIn credentials not configured", http.StatusBadRequest)
		return
	}

	loginHandler := auth.NewLoginHandler(email, password, session.GetPage(), app.persona.TypingSpeedWPM, app.persona.ErrorRate)
	
	logger.Log.Info("Attempting LinkedIn login...")
	if err := loginHandler.Login(); err != nil {
		jsonError(w, fmt.Sprintf("Login failed: %v", err), http.StatusUnauthorized)
		session.Close()
		app.session = nil
		return
	}

	app.isRunning = true
	app.store.LogActivity("Session started", "success", "Browser session initialized and logged in")

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "Automation started successfully",
		"persona": app.persona.Name,
	})
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	if !app.isRunning {
		jsonError(w, "Automation not running", http.StatusBadRequest)
		return
	}

	if app.session != nil {
		app.session.Close()
		app.session = nil
	}

	app.isRunning = false
	app.store.LogActivity("Session stopped", "info", "Automation stopped by user")

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "Automation stopped successfully",
	})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	connectionsSent, _ := app.store.GetConnectionsSentToday()
	messagesSent, _ := app.store.GetMessagesSentToday()
	
	var cooldownRemaining int64 = 0
	if app.scheduler != nil {
		remaining := app.scheduler.GetCooldownRemaining()
		cooldownRemaining = int64(remaining.Seconds())
	}

	stats := map[string]interface{}{
		"connections_sent": connectionsSent,
		"messages_sent":    messagesSent,
		"cooldown_seconds": cooldownRemaining,
		"daily_limit": map[string]int{
			"connections": getEnvInt("DAILY_CONNECTION_LIMIT", 20),
			"messages":    getEnvInt("DAILY_MESSAGE_LIMIT", 10),
		},
	}

	jsonResponse(w, stats)
}

func handleActivity(w http.ResponseWriter, r *http.Request) {
	logs, err := app.store.GetRecentActivityLogs(20)
	if err != nil {
		jsonError(w, "Failed to fetch activity logs", http.StatusInternalServerError)
		return
	}

	// Format logs for frontend
	formattedLogs := make([]map[string]interface{}, 0)
	for _, log := range logs {
		formattedLogs = append(formattedLogs, map[string]interface{}{
			"id":        strconv.FormatInt(log.ID, 10),
			"timestamp": log.Timestamp.Format("15:04:05"),
			"action":    log.Action,
			"type":      log.Type,
			"details":   log.Details,
		})
	}

	jsonResponse(w, formattedLogs)
}

func handlePersona(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PersonaType string `json:"persona_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	app.persona = behavior.GetDefaultPersona(behavior.PersonaType(req.PersonaType))
	app.store.LogActivity("Persona changed", "info", fmt.Sprintf("Switched to %s persona", app.persona.Name))

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"persona": app.persona.Name,
	})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if !app.isRunning || app.session == nil {
		jsonError(w, "Session not started", http.StatusBadRequest)
		return
	}

	var req struct {
		Keywords string `json:"keywords"`
		Location string `json:"location"`
		Company  string `json:"company"`
		JobTitle string `json:"job_title"`
		MaxResults int `json:"max_results"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MaxResults == 0 {
		req.MaxResults = 10
	}

	searcher := search.NewPeopleSearch(app.session.GetPage())
	criteria := search.SearchCriteria{
		Keywords: req.Keywords,
		Location: req.Location,
		Company:  req.Company,
		JobTitle: req.JobTitle,
	}

	results, err := searcher.Search(criteria, req.MaxResults)
	if err != nil {
		jsonError(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	app.store.LogActivity("Search executed", "info", fmt.Sprintf("Found %d profiles", len(results)))

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"results": results,
		"count":   len(results),
	})
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	if !app.isRunning || app.session == nil {
		jsonError(w, "Session not started", http.StatusBadRequest)
		return
	}

	var req struct {
		ProfileURL string `json:"profile_url"`
		Note       string `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check rate limits
	sent, _ := app.store.GetConnectionsSentToday()
	limit := getEnvInt("DAILY_CONNECTION_LIMIT", 20)
	if sent >= limit {
		jsonError(w, "Daily connection limit reached", http.StatusTooManyRequests)
		return
	}

	// Send connection request
	connector := connect.NewConnectionRequest(app.session.GetPage(), app.persona.TypingSpeedWPM, app.persona.ErrorRate)
	
	if err := connector.SendRequest(req.ProfileURL, req.Note); err != nil {
		jsonError(w, fmt.Sprintf("Connection request failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Save to database
	name, _ := connector.GetProfileName()
	app.store.SaveConnectionRequest(req.ProfileURL, name, req.Note)
	app.store.LogActivity("Connection sent", "success", fmt.Sprintf("Sent request to %s", name))

	if app.scheduler != nil {
		app.scheduler.RecordAction()
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "Connection request sent",
	})
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	if !app.isRunning || app.session == nil {
		jsonError(w, "Session not started", http.StatusBadRequest)
		return
	}

	var req struct {
		ProfileURL string `json:"profile_url"`
		Message    string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check rate limits
	sent, _ := app.store.GetMessagesSentToday()
	limit := getEnvInt("DAILY_MESSAGE_LIMIT", 10)
	if sent >= limit {
		jsonError(w, "Daily message limit reached", http.StatusTooManyRequests)
		return
	}

	// Send message
	messenger := messaging.NewFollowUp(app.session.GetPage(), app.persona.TypingSpeedWPM, app.persona.ErrorRate)
	
	if err := messenger.SendMessage(req.ProfileURL, req.Message); err != nil {
		jsonError(w, fmt.Sprintf("Message failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Save to database
	app.store.SaveMessage(req.ProfileURL, req.Message)
	app.store.LogActivity("Message sent", "success", "Follow-up message delivered")

	if app.scheduler != nil {
		app.scheduler.RecordAction()
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "Message sent successfully",
	})
}

// Helper functions

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
