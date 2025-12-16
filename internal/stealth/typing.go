package stealth

import (
	"math/rand"
	"strings"
	"time"
)

// TypingSimulator simulates human typing behavior
type TypingSimulator struct {
	wpm          int
	errorRate    float64
	correctChars []string
}

// NewTypingSimulator creates a typing simulator
func NewTypingSimulator(wpm int, errorRate float64) *TypingSimulator {
	return &TypingSimulator{
		wpm:          wpm,
		errorRate:    errorRate,
		correctChars: make([]string, 0),
	}
}

// TypeCharacter calculates delay for a single character
func (t *TypingSimulator) TypeCharacter() time.Duration {
	// Average: 60000ms / (WPM * 5 chars per word)
	avgDelay := 60000.0 / float64(t.wpm*5)
	
	// Add variance: ±40%
	variance := avgDelay * 0.4
	delay := avgDelay + (rand.Float64()*2-1)*variance
	
	if delay < 10 {
		delay = 10
	}
	
	return time.Duration(delay) * time.Millisecond
}

// ShouldMakeTypo determines if typo should occur
func (t *TypingSimulator) ShouldMakeTypo() bool {
	return rand.Float64()*100 < t.errorRate
}

// GetTypoCharacter returns a nearby keyboard character
func (t *TypingSimulator) GetTypoCharacter(intended rune) rune {
	// Keyboard proximity map (simplified QWERTY layout)
	proximityMap := map[rune][]rune{
		'a': {'s', 'q', 'z', 'w'},
		'b': {'v', 'g', 'h', 'n'},
		'c': {'x', 'd', 'f', 'v'},
		'd': {'s', 'e', 'r', 'f', 'c', 'x'},
		'e': {'w', 'r', 'd', 's'},
		'f': {'d', 'r', 't', 'g', 'v', 'c'},
		'g': {'f', 't', 'y', 'h', 'b', 'v'},
		'h': {'g', 'y', 'u', 'j', 'n', 'b'},
		'i': {'u', 'o', 'k', 'j'},
		'j': {'h', 'u', 'i', 'k', 'm', 'n'},
		'k': {'j', 'i', 'o', 'l', 'm'},
		'l': {'k', 'o', 'p'},
		'm': {'n', 'j', 'k'},
		'n': {'b', 'h', 'j', 'm'},
		'o': {'i', 'p', 'l', 'k'},
		'p': {'o', 'l'},
		'q': {'w', 'a'},
		'r': {'e', 't', 'f', 'd'},
		's': {'a', 'w', 'e', 'd', 'x', 'z'},
		't': {'r', 'y', 'g', 'f'},
		'u': {'y', 'i', 'j', 'h'},
		'v': {'c', 'f', 'g', 'b'},
		'w': {'q', 'e', 's', 'a'},
		'x': {'z', 's', 'd', 'c'},
		'y': {'t', 'u', 'h', 'g'},
		'z': {'a', 's', 'x'},
	}
	
	intended = toLower(intended)
	nearby, exists := proximityMap[intended]
	if !exists || len(nearby) == 0 {
		return intended
	}
	
	return nearby[rand.Intn(len(nearby))]
}

// TypeString simulates typing a full string with potential errors
func (t *TypingSimulator) TypeString(text string) []TypingEvent {
	events := make([]TypingEvent, 0)
	
	for i, char := range text {
		delay := t.TypeCharacter()
		
		// Should we make a typo?
		if t.ShouldMakeTypo() && char != ' ' {
			// Type wrong character
			typo := t.GetTypoCharacter(char)
			events = append(events, TypingEvent{
				Character: string(typo),
				Delay:     delay,
				IsTypo:    true,
			})
			
			// Pause (realize mistake)
			events = append(events, TypingEvent{
				Character: "",
				Delay:     time.Duration(100+rand.Intn(300)) * time.Millisecond,
				IsPause:   true,
			})
			
			// Backspace
			events = append(events, TypingEvent{
				Character: "\b",
				Delay:     time.Duration(50+rand.Intn(100)) * time.Millisecond,
				IsBackspace: true,
			})
			
			// Type correct character
			events = append(events, TypingEvent{
				Character: string(char),
				Delay:     delay,
				IsTypo:    false,
			})
		} else {
			// Normal typing
			events = append(events, TypingEvent{
				Character: string(char),
				Delay:     delay,
				IsTypo:    false,
			})
		}
		
		// Occasional longer pause (thinking)
		if i > 0 && i%20 == 0 && rand.Float64() < 0.3 {
			events = append(events, TypingEvent{
				Character: "",
				Delay:     time.Duration(300+rand.Intn(700)) * time.Millisecond,
				IsPause:   true,
			})
		}
	}
	
	return events
}

// TypingEvent represents a single typing action
type TypingEvent struct {
	Character   string
	Delay       time.Duration
	IsTypo      bool
	IsBackspace bool
	IsPause     bool
}

// toLower converts rune to lowercase
func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}

// GetWordDelay returns longer delay between words
func (t *TypingSimulator) GetWordDelay() time.Duration {
	// Slightly longer pause between words
	return t.TypeCharacter() + time.Duration(rand.Intn(100))*time.Millisecond
}

// SimulateBackspacing simulates deleting characters
func (t *TypingSimulator) SimulateBackspacing(count int) []TypingEvent {
	events := make([]TypingEvent, 0)
	
	for i := 0; i < count; i++ {
		events = append(events, TypingEvent{
			Character:   "\b",
			Delay:       time.Duration(30+rand.Intn(70)) * time.Millisecond,
			IsBackspace: true,
		})
	}
	
	return events
}

// CleanText removes special characters that might cause issues
func CleanText(text string) string {
	return strings.TrimSpace(text)
}
