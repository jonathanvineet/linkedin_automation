package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/input"
	"github.com/jonathanvineet/linkedin-automation/internal/logger"
	"github.com/jonathanvineet/linkedin-automation/internal/stealth"
)

// LoginHandler manages LinkedIn authentication
type LoginHandler struct {
	email    string
	password string
	page     *rod.Page
	typing   *stealth.TypingSimulator
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(email, password string, page *rod.Page, wpm int, errorRate float64) *LoginHandler {
	return &LoginHandler{
		email:    email,
		password: password,
		page:     page,
		typing:   stealth.NewTypingSimulator(wpm, errorRate),
	}
}

// Login performs LinkedIn login with human-like behavior
func (lh *LoginHandler) Login() error {
	logger.Log.Info("Starting LinkedIn login process")
	
	// Navigate to LinkedIn login page
	err := lh.page.Navigate("https://www.linkedin.com/login")
	if err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}
	
	// Wait for page to load
	time.Sleep(2 * time.Second)
	err = lh.page.WaitLoad()
	if err != nil {
		return fmt.Errorf("failed to load login page: %w", err)
	}

	logger.Log.Info("Login page loaded successfully")
	
	// Check if already logged in
	if lh.IsLoggedIn() {
		logger.Log.Info("Already logged in, skipping login process")
		return nil
	}
	
	// Find username field
	usernameField, err := lh.page.Timeout(10 * time.Second).Element("#username")
	if err != nil {
		return errors.New("could not find username field")
	}
	
	// Human-like delay before typing
	time.Sleep(time.Duration(800+time.Now().UnixNano()%700) * time.Millisecond)
	
	// Type email with human-like behavior
	logger.Log.Info("Entering email address")
	err = lh.typeWithBehavior(usernameField, lh.email)
	if err != nil {
		return fmt.Errorf("failed to enter email: %w", err)
	}
	
	// Delay between fields
	time.Sleep(time.Duration(500+time.Now().UnixNano()%500) * time.Millisecond)
	
	// Find password field
	passwordField, err := lh.page.Element("#password")
	if err != nil {
		return errors.New("could not find password field")
	}
	
	// Type password
	logger.Log.Info("Entering password")
	err = lh.typeWithBehavior(passwordField, lh.password)
	if err != nil {
		return fmt.Errorf("failed to enter password: %w", err)
	}
	
	// Delay before clicking submit
	time.Sleep(time.Duration(800+time.Now().UnixNano()%400) * time.Millisecond)
	
	// Find and click login button
	logger.Log.Info("Clicking login button")
	loginButton, err := lh.page.Element("button[type='submit']")
	if err != nil {
		return errors.New("could not find login button")
	}
	
	err = loginButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click login button: %w", err)
	}
	
	// Wait for navigation
	logger.Log.Info("Waiting for login to complete")
	time.Sleep(5 * time.Second)
	
	// Check for 2FA or security challenge
	if lh.hasSecurityChallenge() {
		logger.Log.Warn("Security challenge detected - manual intervention required")
		return errors.New("security challenge detected: please complete manually")
	}
	
	// Verify login success
	if !lh.IsLoggedIn() {
		logger.Log.Error("Login failed - credentials may be incorrect")
		return errors.New("login failed: please check credentials")
	}
	
	logger.Log.Info("Login successful")
	return nil
}

// typeWithBehavior types text with human-like behavior
func (lh *LoginHandler) typeWithBehavior(element *rod.Element, text string) error {
	// Click on the field first
	err := element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	
	// Small delay after click
	time.Sleep(100 * time.Millisecond)
	
	// Type each character with realistic delays
	events := lh.typing.TypeString(text)
	
	for _, event := range events {
		time.Sleep(event.Delay)
		
		if event.IsPause {
			continue // Just a pause, no typing
		}
		
		if event.IsBackspace {
			// Simulate backspace
			err = lh.page.Keyboard.Press(input.Backspace)
			if err != nil {
				return err
			}
			continue
		}
		
		// Type character
		if event.Character != "" {
			err = element.Input(event.Character)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

// IsLoggedIn checks if user is currently logged in
func (lh *LoginHandler) IsLoggedIn() bool {
	// Check for elements that only appear when logged in
	_, err := lh.page.Timeout(3 * time.Second).Element("nav")
	return err == nil
}

// hasSecurityChallenge checks for 2FA or security prompts
func (lh *LoginHandler) hasSecurityChallenge() bool {
	// Check for common security challenge indicators
	selectors := []string{
		"#input__email_verification_pin",
		"#input__phone_verification_pin",
		"[data-test-id='security-challenge']",
	}
	
	for _, selector := range selectors {
		_, err := lh.page.Timeout(2 * time.Second).Element(selector)
		if err == nil {
			return true
		}
	}
	
	return false
}

// GetSessionCookies retrieves current session cookies
func (lh *LoginHandler) GetSessionCookies() ([]*proto.NetworkCookie, error) {
	return lh.page.Cookies([]string{"https://www.linkedin.com"})
}

// LoadSessionCookies loads previously saved cookies
func (lh *LoginHandler) LoadSessionCookies(cookies []*proto.NetworkCookie) error {
	if len(cookies) == 0 {
		return nil
	}
	
	// Convert NetworkCookie to NetworkCookieParam
	cookieParams := make([]*proto.NetworkCookieParam, len(cookies))
	for i, cookie := range cookies {
		cookieParams[i] = &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
			SameSite: cookie.SameSite,
		}
	}
	return lh.page.SetCookies(cookieParams)
}

// Logout performs logout
func (lh *LoginHandler) Logout() error {
	logger.Log.Info("Logging out of LinkedIn")
	
	// Navigate to logout URL
	err := lh.page.Navigate("https://www.linkedin.com/m/logout")
	if err != nil {
		return err
	}
	
	time.Sleep(2 * time.Second)
	logger.Log.Info("Logout complete")
	return nil
}
