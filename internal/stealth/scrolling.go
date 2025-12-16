package stealth

import (
	"math/rand"
	"time"
)

// ScrollBehavior handles realistic scrolling patterns
type ScrollBehavior struct {
	impatience string // low, medium, high
}

// NewScrollBehavior creates a scroll controller
func NewScrollBehavior(impatience string) *ScrollBehavior {
	return &ScrollBehavior{
		impatience: impatience,
	}
}

// GetScrollAmount returns pixel amount to scroll
func (sb *ScrollBehavior) GetScrollAmount() int {
	var baseScroll int
	
	switch sb.impatience {
	case "high":
		baseScroll = 300 + rand.Intn(200) // 300-500px
	case "medium":
		baseScroll = 150 + rand.Intn(150) // 150-300px
	case "low":
		baseScroll = 50 + rand.Intn(100) // 50-150px
	default:
		baseScroll = 150 + rand.Intn(150)
	}
	
	return baseScroll
}

// GetScrollDelay returns delay between scroll actions
func (sb *ScrollBehavior) GetScrollDelay() time.Duration {
	var delayMs int
	
	switch sb.impatience {
	case "high":
		delayMs = 100 + rand.Intn(200) // Fast scrolling
	case "medium":
		delayMs = 300 + rand.Intn(300)
	case "low":
		delayMs = 500 + rand.Intn(500) // Slow, careful
	default:
		delayMs = 300 + rand.Intn(300)
	}
	
	return time.Duration(delayMs) * time.Millisecond
}

// ShouldPauseWhileScrolling determines if user pauses mid-scroll
func (sb *ScrollBehavior) ShouldPauseWhileScrolling() bool {
	chance := 0.3
	
	if sb.impatience == "low" {
		chance = 0.5 // More careful readers pause more
	} else if sb.impatience == "high" {
		chance = 0.1
	}
	
	return rand.Float64() < chance
}

// GetScrollPauseDuration returns pause duration
func (sb *ScrollBehavior) GetScrollPauseDuration() time.Duration {
	return time.Duration(500+rand.Intn(2000)) * time.Millisecond
}

// GenerateScrollPattern creates realistic scroll sequence
func (sb *ScrollBehavior) GenerateScrollPattern(targetScrollPosition int) []ScrollAction {
	actions := make([]ScrollAction, 0)
	currentPosition := 0
	
	for currentPosition < targetScrollPosition {
		scrollAmount := sb.GetScrollAmount()
		
		// Don't overshoot too much
		if currentPosition+scrollAmount > targetScrollPosition+100 {
			scrollAmount = targetScrollPosition - currentPosition
		}
		
		actions = append(actions, ScrollAction{
			Amount: scrollAmount,
			Delay:  sb.GetScrollDelay(),
			Pause:  sb.ShouldPauseWhileScrolling(),
		})
		
		currentPosition += scrollAmount
		
		// Add occasional pause to "read"
		if sb.ShouldPauseWhileScrolling() {
			actions = append(actions, ScrollAction{
				Amount: 0,
				Delay:  sb.GetScrollPauseDuration(),
				Pause:  true,
			})
		}
	}
	
	// Sometimes scroll back up slightly (humans re-read)
	if rand.Float64() < 0.2 {
		actions = append(actions, ScrollAction{
			Amount: -50 - rand.Intn(100),
			Delay:  sb.GetScrollDelay(),
			Pause:  false,
		})
	}
	
	return actions
}

// ScrollAction represents a single scroll operation
type ScrollAction struct {
	Amount int           // Pixels to scroll (negative = up)
	Delay  time.Duration // Wait before scroll
	Pause  bool          // Whether this is a reading pause
}

// GetWheelScrollDelta returns realistic mouse wheel delta
func GetWheelScrollDelta() int {
	// Mouse wheel typically scrolls 100-120 pixels per tick
	return 100 + rand.Intn(20)
}

// ShouldUseKeyboardScroll decides if keyboard should be used instead
func ShouldUseKeyboardScroll() bool {
	// 10% of users prefer keyboard (Page Down, Space)
	return rand.Float64() < 0.1
}

// GetSmoothnessEasing returns easing value for smooth scrolling
func GetSmoothnessEasing(step, totalSteps int) float64 {
	// Ease-in-out curve
	t := float64(step) / float64(totalSteps)
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - 2*(1-t)*(1-t)
}
