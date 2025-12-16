package connect

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/jonathanvineet/linkedin-automation/internal/logger"
	"github.com/jonathanvineet/linkedin-automation/internal/stealth"
)

// ConnectionRequest handles sending connection requests
type ConnectionRequest struct {
	page   *rod.Page
	typing *stealth.TypingSimulator
}

// NewConnectionRequest creates a new connection handler
func NewConnectionRequest(page *rod.Page, wpm int, errorRate float64) *ConnectionRequest {
	return &ConnectionRequest{
		page:   page,
		typing: stealth.NewTypingSimulator(wpm, errorRate),
	}
}

// SendRequest sends a connection request to a profile
func (cr *ConnectionRequest) SendRequest(profileURL string, personalizedNote string) error {
	logger.Log.WithField("profile", profileURL).Info("Navigating to profile")
	
	// Navigate to profile
	err := cr.page.Navigate(profileURL)
	if err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}
	
	// Wait for page load
	time.Sleep(3 * time.Second)
	err = cr.page.WaitLoad()
	if err != nil {
		return fmt.Errorf("failed to load profile: %w", err)
	}
	
	// Simulate reading the profile (human behavior)
	cr.simulateProfileReading()
	
	// Find Connect button
	logger.Log.Info("Looking for Connect button")
	connectButton, err := cr.findConnectButton()
	if err != nil {
		return fmt.Errorf("could not find Connect button: %w", err)
	}
	
	// Human-like delay before clicking
	time.Sleep(time.Duration(800+time.Now().UnixNano()%700) * time.Millisecond)
	
	// Click Connect button
	logger.Log.Info("Clicking Connect button")
	err = connectButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click Connect button: %w", err)
	}
	
	// Wait for dialog to appear
	time.Sleep(2 * time.Second)
	
	// Check if note option is available and add note if provided
	if personalizedNote != "" {
		err = cr.addPersonalizedNote(personalizedNote)
		if err != nil {
			logger.Log.WithError(err).Warn("Could not add personalized note, sending without note")
		}
	}
	
	// Find and click Send button
	logger.Log.Info("Clicking Send button")
	sendButton, err := cr.page.Timeout(5 * time.Second).Element("button[aria-label='Send now']")
	if err != nil {
		// Try alternative selector
		sendButton, err = cr.page.Element("button[aria-label='Send']")
		if err != nil {
			return errors.New("could not find Send button")
		}
	}
	
	// Final delay before sending
	time.Sleep(time.Duration(500+time.Now().UnixNano()%500) * time.Millisecond)
	
	err = sendButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click Send button: %w", err)
	}
	
	// Wait for confirmation
	time.Sleep(2 * time.Second)
	
	logger.Log.Info("Connection request sent successfully")
	return nil
}

// findConnectButton locates the Connect button with various strategies
func (cr *ConnectionRequest) findConnectButton() (*rod.Element, error) {
	// Try multiple selectors as LinkedIn's UI can vary
	selectors := []string{
		"button[aria-label*='Connect']",
		"button[aria-label='Connect']",
		"button.pvs-profile-actions__action",
		"button:contains('Connect')",
	}
	
	for _, selector := range selectors {
		button, err := cr.page.Timeout(5 * time.Second).Element(selector)
		if err == nil {
			return button, nil
		}
	}
	
	return nil, errors.New("Connect button not found with any selector")
}

// addPersonalizedNote adds a custom note to the connection request
func (cr *ConnectionRequest) addPersonalizedNote(note string) error {
	// Look for "Add a note" button
	addNoteButton, err := cr.page.Timeout(3 * time.Second).Element("button[aria-label='Add a note']")
	if err != nil {
		return errors.New("could not find Add a note button")
	}
	
	// Click Add a note
	err = addNoteButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	
	// Wait for text area to appear
	time.Sleep(1 * time.Second)
	
	// Find the note text area
	noteField, err := cr.page.Timeout(3 * time.Second).Element("textarea[name='message']")
	if err != nil {
		return errors.New("could not find note text field")
	}
	
	// Click on text area
	err = noteField.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	
	// Type the note with human-like behavior
	time.Sleep(500 * time.Millisecond)
	
	// Truncate note if too long (LinkedIn limit is 300 characters)
	if len(note) > 300 {
		note = note[:297] + "..."
	}
	
	// Type character by character with delays
	events := cr.typing.TypeString(note)
	for _, event := range events {
		time.Sleep(event.Delay)
		
		if event.IsPause || event.Character == "" {
			continue
		}
		
		if event.IsBackspace {
			continue // Skip backspaces in notes for simplicity
		}
		
		err = noteField.Input(event.Character)
		if err != nil {
			return err
		}
	}
	
	logger.Log.Info("Personalized note added successfully")
	return nil
}

// simulateProfileReading simulates spending time reading the profile
func (cr *ConnectionRequest) simulateProfileReading() {
	// Scroll down to view experience section
	cr.page.Mouse.Scroll(0, 300, 10)
	time.Sleep(time.Duration(1000+time.Now().UnixNano()%1000) * time.Millisecond)
	
	// Scroll down more
	cr.page.Mouse.Scroll(0, 400, 10)
	time.Sleep(time.Duration(800+time.Now().UnixNano()%700) * time.Millisecond)
	
	// Scroll back up slightly (re-reading)
	cr.page.Mouse.Scroll(0, -150, 10)
	time.Sleep(time.Duration(600+time.Now().UnixNano()%400) * time.Millisecond)
}

// IsAlreadyConnected checks if already connected to this profile
func (cr *ConnectionRequest) IsAlreadyConnected() bool {
	// Look for "Message" button which indicates existing connection
	_, err := cr.page.Timeout(2 * time.Second).Element("button[aria-label*='Message']")
	return err == nil
}

// GetProfileName extracts the profile name
func (cr *ConnectionRequest) GetProfileName() (string, error) {
	nameElement, err := cr.page.Timeout(5 * time.Second).Element("h1")
	if err != nil {
		return "", err
	}
	
	name, err := nameElement.Text()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(name), nil
}

// HasPendingRequest checks if a connection request is already pending
func (cr *ConnectionRequest) HasPendingRequest() bool {
	// Look for "Pending" button
	_, err := cr.page.Timeout(2 * time.Second).Element("button[aria-label*='Pending']")
	return err == nil
}
