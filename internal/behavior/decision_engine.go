package behavior

import (
	"math/rand"
	"time"
)

// DecisionEngine makes context-aware decisions for human-like behavior
type DecisionEngine struct {
	persona       *Persona
	actionsToday  int
	pageComplexity int // 1-10 scale
	timeOfDay     int // hour 0-23
}

// NewDecisionEngine creates a new decision engine
func NewDecisionEngine(persona *Persona) *DecisionEngine {
	return &DecisionEngine{
		persona:        persona,
		actionsToday:   0,
		pageComplexity: 5,
		timeOfDay:      time.Now().Hour(),
	}
}

// SetPageComplexity updates current page complexity (affects think time)
func (d *DecisionEngine) SetPageComplexity(complexity int) {
	if complexity < 1 {
		complexity = 1
	}
	if complexity > 10 {
		complexity = 10
	}
	d.pageComplexity = complexity
}

// GetThinkTime calculates realistic "think time" before action
func (d *DecisionEngine) GetThinkTime() time.Duration {
	// Base think time: 1-3 seconds
	baseSeconds := 1.0 + rand.Float64()*2.0
	
	// Adjust for page complexity
	complexityFactor := 1.0 + (float64(d.pageComplexity) / 10.0)
	
	// Adjust for persona attention span
	attentionFactor := 1.0
	if d.persona.AttentionSpanSec < 120 {
		attentionFactor = 0.8 // More impatient
	} else if d.persona.AttentionSpanSec > 180 {
		attentionFactor = 1.2 // More patient
	}
	
	// Fatigue factor (more actions = slower thinking)
	fatigueFactor := 1.0 + (float64(d.actionsToday) / 100.0)
	
	thinkSeconds := baseSeconds * complexityFactor * attentionFactor * fatigueFactor
	
	// Add random jitter ±20%
	jitter := thinkSeconds * 0.2 * (rand.Float64()*2 - 1)
	
	return time.Duration((thinkSeconds+jitter)*1000) * time.Millisecond
}

// ShouldHoverFirst decides whether to hover over element before clicking
func (d *DecisionEngine) ShouldHoverFirst() bool {
	// 60% of the time, hover first
	return rand.Float64() < 0.6
}

// GetHoverDuration returns how long to hover
func (d *DecisionEngine) GetHoverDuration() time.Duration {
	// Hover for 200-800ms
	return time.Duration(200+rand.Intn(600)) * time.Millisecond
}

// ShouldHesitate decides if there should be a micro-pause
func (d *DecisionEngine) ShouldHesitate() bool {
	// 30% chance of hesitation
	return rand.Float64() < 0.3
}

// GetHesitationDuration returns micro-pause duration
func (d *DecisionEngine) GetHesitationDuration() time.Duration {
	// Hesitate for 300-1000ms
	return time.Duration(300+rand.Intn(700)) * time.Millisecond
}

// ShouldScrollBeforeAction decides if user should scroll before interacting
func (d *DecisionEngine) ShouldScrollBeforeAction() bool {
	// Higher for impatient personas
	chance := 0.4
	if d.persona.ScrollImpatience == "high" {
		chance = 0.6
	} else if d.persona.ScrollImpatience == "low" {
		chance = 0.2
	}
	return rand.Float64() < chance
}

// ShouldReReadContent decides if user re-reads before proceeding
func (d *DecisionEngine) ShouldReReadContent() bool {
	// More likely for careful personas
	chance := 0.25
	if d.persona.MousePrecision > 85 {
		chance = 0.4
	}
	return rand.Float64() < chance
}

// GetActionDelay returns delay between major actions
func (d *DecisionEngine) GetActionDelay(minSeconds, maxSeconds int) time.Duration {
	// Base delay in range
	delaySeconds := minSeconds + rand.Intn(maxSeconds-minSeconds+1)
	
	// Adjust for time of day (slower during "lunch" or "tired" hours)
	hour := time.Now().Hour()
	if hour == 12 || hour == 13 || hour > 16 {
		delaySeconds = int(float64(delaySeconds) * 1.3)
	}
	
	// Adjust for persona impatience
	if d.persona.ScrollImpatience == "high" {
		delaySeconds = int(float64(delaySeconds) * 0.8)
	}
	
	return time.Duration(delaySeconds) * time.Second
}

// IncrementActionCount tracks daily actions
func (d *DecisionEngine) IncrementActionCount() {
	d.actionsToday++
}

// IsWithinBusinessHours checks if current time is appropriate for actions
func (d *DecisionEngine) IsWithinBusinessHours(startHour, endHour int) bool {
	hour := time.Now().Hour()
	return hour >= startHour && hour < endHour
}

// ShouldRandomlyIdle decides if user should idle/wander
func (d *DecisionEngine) ShouldRandomlyIdle() bool {
	// 15% chance to simulate distraction
	return rand.Float64() < 0.15
}

// GetIdleDuration returns idle/distraction duration
func (d *DecisionEngine) GetIdleDuration() time.Duration {
	// Idle for 2-8 seconds
	return time.Duration(2+rand.Intn(6)) * time.Second
}
