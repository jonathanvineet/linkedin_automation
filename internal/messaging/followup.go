package messaging

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

// FollowUp handles sending follow-up messages
type FollowUp struct {
	page   *rod.Page
	typing *stealth.TypingSimulator
}

// NewFollowUp creates a new follow-up message handler
func NewFollowUp(page *rod.Page, wpm int, errorRate float64) *FollowUp {
	return &FollowUp{
		page:   page,
		typing: stealth.NewTypingSimulator(wpm, errorRate),
	}
}

// SendMessage sends a message to a connection
func (fu *FollowUp) SendMessage(profileURL string, message string) error {
	logger.Log.WithField("profile", profileURL).Info("Sending message to connection")
	
	// Navigate to profile
	err := fu.page.Navigate(profileURL)
	if err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}
	
	// Wait for page load
	time.Sleep(3 * time.Second)
	err = fu.page.WaitLoad()
	if err != nil {
		return fmt.Errorf("failed to load profile: %w", err)
	}
	
	// Find Message button
	logger.Log.Info("Looking for Message button")
	messageButton, err := fu.findMessageButton()
	if err != nil {
		return fmt.Errorf("could not find Message button: %w", err)
	}
	
	// Human-like delay before clicking
	time.Sleep(time.Duration(700+time.Now().UnixNano()%600) * time.Millisecond)
	
	// Click Message button
	logger.Log.Info("Clicking Message button")
	err = messageButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click Message button: %w", err)
	}
	
	// Wait for message box to appear
	time.Sleep(2 * time.Second)
	
	// Find message text field
	messageField, err := fu.page.Timeout(5 * time.Second).Element("div[role='textbox']")
	if err != nil {
		return errors.New("could not find message text field")
	}
	
	// Click on message field
	err = messageField.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	
	// Delay before typing
	time.Sleep(800 * time.Millisecond)
	
	// Type the message with human-like behavior
	logger.Log.Info("Typing message")
	err = fu.typeMessage(messageField, message)
	if err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}
	
	// Find and click Send button
	logger.Log.Info("Clicking Send button")
	sendButton, err := fu.findSendButton()
	if err != nil {
		return fmt.Errorf("could not find Send button: %w", err)
	}
	
	// Final delay before sending
	time.Sleep(time.Duration(500+time.Now().UnixNano()%500) * time.Millisecond)
	
	err = sendButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click Send button: %w", err)
	}
	
	// Wait for confirmation
	time.Sleep(2 * time.Second)
	
	logger.Log.Info("Message sent successfully")
	return nil
}

// findMessageButton locates the Message button
func (fu *FollowUp) findMessageButton() (*rod.Element, error) {
	selectors := []string{
		"button[aria-label*='Message']",
		"a[href*='/messaging/thread']",
		"button:contains('Message')",
	}
	
	for _, selector := range selectors {
		button, err := fu.page.Timeout(5 * time.Second).Element(selector)
		if err == nil {
			return button, nil
		}
	}
	
	return nil, errors.New("Message button not found")
}

// findSendButton locates the Send button in message dialog
func (fu *FollowUp) findSendButton() (*rod.Element, error) {
	selectors := []string{
		"button[type='submit']",
		"button[aria-label='Send']",
		"button.msg-form__send-button",
	}
	
	for _, selector := range selectors {
		button, err := fu.page.Timeout(5 * time.Second).Element(selector)
		if err == nil {
			// Check if button is enabled
			disabled, _ := button.Attribute("disabled")
			if disabled == nil {
				return button, nil
			}
		}
	}
	
	return nil, errors.New("Send button not found or disabled")
}

// typeMessage types a message with human-like behavior
func (fu *FollowUp) typeMessage(element *rod.Element, message string) error {
	// Type message character by character
	events := fu.typing.TypeString(message)
	
	for _, event := range events {
		time.Sleep(event.Delay)
		
		if event.IsPause || event.Character == "" {
			continue
		}
		
		if event.IsBackspace {
			// For message boxes, we'll skip backspace simulation for simplicity
			continue
		}
		
		err := element.Input(event.Character)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// PersonalizeMessage replaces template variables with actual values
func PersonalizeMessage(template string, firstName string, company string, title string) string {
	message := template
	message = strings.ReplaceAll(message, "{firstName}", firstName)
	message = strings.ReplaceAll(message, "{company}", company)
	message = strings.ReplaceAll(message, "{title}", title)
	message = strings.ReplaceAll(message, "{industry}", "your industry") // Placeholder
	return message
}

// GetConversationHistory retrieves recent messages (read-only)
func (fu *FollowUp) GetConversationHistory() ([]string, error) {
	// Find message container
	messages, err := fu.page.Elements(".msg-s-message-list__event")
	if err != nil {
		return nil, err
	}
	
	history := make([]string, 0)
	for _, msg := range messages {
		text, err := msg.Text()
		if err == nil && text != "" {
			history = append(history, strings.TrimSpace(text))
		}
	}
	
	return history, nil
}

// HasExistingConversation checks if there's already a conversation
func (fu *FollowUp) HasExistingConversation() bool {
	_, err := fu.page.Timeout(2 * time.Second).Elements(".msg-s-message-list__event")
	return err == nil
}
