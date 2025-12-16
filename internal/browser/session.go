package browser

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

// Session manages the browser automation session
type Session struct {
	browser    *rod.Browser
	page       *rod.Page
	userAgent  string
	headless   bool
	cookiePath string
}

// Config holds browser configuration
type Config struct {
	Headless      bool
	UserAgent     string
	ViewportWidth int
	ViewportHeight int
	Timeout       time.Duration
}

// NewSession creates a new browser session
func NewSession(config Config) (*Session, error) {
	// Launch browser
	l := launcher.New().
		Headless(config.Headless).
		Set("disable-blink-features", "AutomationControlled")
	
	url, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(url).MustConnect()
	
	// Create page with stealth mode
	page := stealth.MustPage(browser)
	
	// Set viewport
	if config.ViewportWidth > 0 && config.ViewportHeight > 0 {
		page.MustSetViewport(config.ViewportWidth, config.ViewportHeight, 1, false)
	}
	
	// Set user agent if provided
	if config.UserAgent != "" {
		page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: config.UserAgent,
		})
	}
	
	// Set default timeout
	if config.Timeout > 0 {
		page = page.Timeout(config.Timeout)
	}

	session := &Session{
		browser:   browser,
		page:      page,
		userAgent: config.UserAgent,
		headless:  config.Headless,
	}

	return session, nil
}

// Navigate navigates to a URL with random delay
func (s *Session) Navigate(url string) error {
	// Add small random delay before navigation
	time.Sleep(time.Duration(500+rand.Intn(1000)) * time.Millisecond)
	
	return s.page.Navigate(url)
}

// WaitLoad waits for page to load
func (s *Session) WaitLoad() error {
	return s.page.WaitLoad()
}

// GetPage returns the current page
func (s *Session) GetPage() *rod.Page {
	return s.page
}

// Screenshot takes a screenshot (for debugging)
func (s *Session) Screenshot(path string) error {
	data, err := s.page.Screenshot(false, nil)
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// GetCurrentURL returns the current page URL
func (s *Session) GetCurrentURL() string {
	info, err := s.page.Info()
	if err != nil {
		return ""
	}
	return info.URL
}

// Close closes the browser session
func (s *Session) Close() error {
	if s.page != nil {
		s.page.Close()
	}
	if s.browser != nil {
		return s.browser.Close()
	}
	return nil
}

// SaveCookies saves current session cookies
func (s *Session) SaveCookies(path string) error {
	cookies, err := s.page.Cookies([]string{})
	if err != nil {
		return err
	}
	
	// Implement cookie serialization here
	// For now, return nil
	_ = cookies
	_ = path
	return nil
}

// LoadCookies loads session cookies
func (s *Session) LoadCookies(path string) error {
	// Implement cookie deserialization here
	// For now, return nil
	_ = path
	return nil
}

// WaitForNavigation waits for page navigation to complete
func (s *Session) WaitForNavigation() error {
	wait := s.page.WaitNavigation(proto.PageLifecycleEventNameNetworkIdle)
	wait()
	return nil
}

// EvaluateJS executes JavaScript in the page context
func (s *Session) EvaluateJS(js string) (*proto.RuntimeRemoteObject, error) {
	return s.page.Eval(js)
}

// IsElementPresent checks if an element exists
func (s *Session) IsElementPresent(selector string) bool {
	_, err := s.page.Timeout(2 * time.Second).Element(selector)
	return err == nil
}

// RandomDelay adds a random human-like delay
func (s *Session) RandomDelay(minMs, maxMs int) {
	delay := minMs + rand.Intn(maxMs-minMs+1)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
