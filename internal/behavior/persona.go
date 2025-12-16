package behavior

import (
	"math"
	"math/rand"
	"time"
)

// Persona defines human behavior characteristics
type Persona struct {
	Name                string  `yaml:"name"`
	TypingSpeedWPM      int     `yaml:"typing_speed_wpm"`
	MousePrecision      int     `yaml:"mouse_precision"`
	ErrorRate           float64 `yaml:"error_rate"`
	AttentionSpanSec    int     `yaml:"attention_span_seconds"`
	BreakFrequency      int     `yaml:"break_frequency_actions"`
	ScrollImpatience    string  `yaml:"scroll_impatience"` // low, medium, high
	actionCount         int
	lastBreakTime       time.Time
}

// PersonaType represents different user types
type PersonaType string

const (
	RecruiterPersona PersonaType = "recruiter"
	FounderPersona   PersonaType = "founder"
	SalesPersona     PersonaType = "sales"
)

// GetDefaultPersona returns a pre-configured persona
func GetDefaultPersona(personaType PersonaType) *Persona {
	personas := map[PersonaType]*Persona{
		RecruiterPersona: {
			Name:             "Professional Recruiter",
			TypingSpeedWPM:   65,
			MousePrecision:   87,
			ErrorRate:        3.5,
			AttentionSpanSec: 180,
			BreakFrequency:   5,
			ScrollImpatience: "medium",
			lastBreakTime:    time.Now(),
		},
		FounderPersona: {
			Name:             "Startup Founder",
			TypingSpeedWPM:   85,
			MousePrecision:   75,
			ErrorRate:        5.0,
			AttentionSpanSec: 90,
			BreakFrequency:   8,
			ScrollImpatience: "high",
			lastBreakTime:    time.Now(),
		},
		SalesPersona: {
			Name:             "Sales Professional",
			TypingSpeedWPM:   72,
			MousePrecision:   82,
			ErrorRate:        4.0,
			AttentionSpanSec: 150,
			BreakFrequency:   6,
			ScrollImpatience: "low",
			lastBreakTime:    time.Now(),
		},
	}

	p, exists := personas[personaType]
	if !exists {
		return personas[RecruiterPersona]
	}
	return p
}

// GetTypingDelay calculates delay between keystrokes based on WPM
func (p *Persona) GetTypingDelay() time.Duration {
	// Average WPM to milliseconds per character
	// WPM = (characters / 5) / minutes
	// So characters per minute = WPM * 5
	// Milliseconds per character = 60000 / (WPM * 5)
	avgDelayMs := 60000.0 / float64(p.TypingSpeedWPM*5)
	
	// Add variance ±30%
	variance := avgDelayMs * 0.3
	delay := avgDelayMs + (rand.Float64()*2-1)*variance
	
	return time.Duration(math.Max(delay, 10)) * time.Millisecond
}

// ShouldMakeTypo determines if a typing error should occur
func (p *Persona) ShouldMakeTypo() bool {
	return rand.Float64()*100 < p.ErrorRate
}

// GetMouseDeviation returns pixel deviation based on precision
func (p *Persona) GetMouseDeviation() int {
	// Lower precision = more deviation
	maxDeviation := 100 - p.MousePrecision
	if maxDeviation < 1 {
		return 0
	}
	return rand.Intn(maxDeviation)
}

// ShouldTakeBreak determines if user should pause
func (p *Persona) ShouldTakeBreak() bool {
	p.actionCount++
	timeSinceBreak := time.Since(p.lastBreakTime)
	
	// Break based on action count or time
	if p.actionCount >= p.BreakFrequency {
		p.actionCount = 0
		p.lastBreakTime = time.Now()
		return true
	}
	
	// Random breaks if attention span exceeded
	if timeSinceBreak.Seconds() > float64(p.AttentionSpanSec) {
		if rand.Float64() < 0.4 { // 40% chance
			p.lastBreakTime = time.Now()
			return true
		}
	}
	
	return false
}

// GetBreakDuration returns how long to pause
func (p *Persona) GetBreakDuration() time.Duration {
	// Short breaks: 5-15 seconds
	// Long breaks: 30-90 seconds
	if rand.Float64() < 0.7 {
		return time.Duration(5+rand.Intn(10)) * time.Second
	}
	return time.Duration(30+rand.Intn(60)) * time.Second
}

// GetScrollSpeed returns scroll behavior based on impatience
func (p *Persona) GetScrollSpeed() string {
	return p.ScrollImpatience
}

// String returns persona description
func (p *Persona) String() string {
	return p.Name
}
