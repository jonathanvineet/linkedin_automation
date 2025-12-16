package stealth

import (
	"math/rand"
	"time"
)

// TimingJitter adds variability to action timing
type TimingJitter struct {
	baselineDelay time.Duration
}

// NewTimingJitter creates a new timing controller
func NewTimingJitter(baseline time.Duration) *TimingJitter {
	return &TimingJitter{
		baselineDelay: baseline,
	}
}

// GetJitteredDelay returns delay with random jitter
func (tj *TimingJitter) GetJitteredDelay(variance float64) time.Duration {
	// variance is a percentage (e.g., 0.3 for ±30%)
	baseMs := float64(tj.baselineDelay.Milliseconds())
	jitter := baseMs * variance * (rand.Float64()*2 - 1)
	
	totalMs := baseMs + jitter
	if totalMs < 0 {
		totalMs = baseMs
	}
	
	return time.Duration(totalMs) * time.Millisecond
}

// GetDelayBetweenActions returns context-aware delay
func GetDelayBetweenActions(minSeconds, maxSeconds int, contextFactor float64) time.Duration {
	// Base delay
	baseSeconds := minSeconds + rand.Intn(maxSeconds-minSeconds+1)
	
	// Apply context factor (fatigue, complexity, etc.)
	adjustedSeconds := float64(baseSeconds) * contextFactor
	
	return time.Duration(adjustedSeconds*1000) * time.Millisecond
}

// GetPageLoadWaitTime returns realistic page load wait
func GetPageLoadWaitTime() time.Duration {
	// Humans wait 1-4 seconds for page content
	return time.Duration(1000+rand.Intn(3000)) * time.Millisecond
}

// GetScrollPauseTime returns pause duration during scrolling
func GetScrollPauseTime() time.Duration {
	// Pause 200-800ms while scrolling
	return time.Duration(200+rand.Intn(600)) * time.Millisecond
}

// GetReadingTime estimates time to read text
func GetReadingTime(wordCount int, wpm int) time.Duration {
	// Average reading speed: 200-250 WPM
	// Add variance
	readingSeconds := float64(wordCount) / float64(wpm) * 60.0
	
	// Add ±20% variance
	variance := readingSeconds * 0.2 * (rand.Float64()*2 - 1)
	
	totalSeconds := readingSeconds + variance
	if totalSeconds < 0.5 {
		totalSeconds = 0.5
	}
	
	return time.Duration(totalSeconds*1000) * time.Millisecond
}

// GetButtonClickDelay returns delay after button visible before clicking
func GetButtonClickDelay() time.Duration {
	// React time: 200-600ms after button appears
	return time.Duration(200+rand.Intn(400)) * time.Millisecond
}

// GetFormFieldDelay returns delay between form fields
func GetFormFieldDelay() time.Duration {
	// Pause 300-1000ms between form fields
	return time.Duration(300+rand.Intn(700)) * time.Millisecond
}

// SimulateNetworkLatency adds realistic wait for network operations
func SimulateNetworkLatency() time.Duration {
	// Simulate 100-500ms network delay
	return time.Duration(100+rand.Intn(400)) * time.Millisecond
}

// GetRandomMicroDelay returns very short random delay
func GetRandomMicroDelay() time.Duration {
	// 10-100ms micro delays
	return time.Duration(10+rand.Intn(90)) * time.Millisecond
}

// IsWithinBusinessHours checks if current time is business hours
func IsWithinBusinessHours(startHour, endHour int) bool {
	now := time.Now()
	hour := now.Hour()
	
	// Check day of week (Monday = 1, Sunday = 0)
	weekday := now.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	return hour >= startHour && hour < endHour
}

// GetNextBusinessHourStart returns time until next business hour
func GetNextBusinessHourStart(startHour int) time.Duration {
	now := time.Now()
	currentHour := now.Hour()
	
	if currentHour < startHour {
		// Wait until start hour today
		hoursToWait := startHour - currentHour
		return time.Duration(hoursToWait) * time.Hour
	}
	
	// Wait until start hour tomorrow
	hoursToWait := (24 - currentHour) + startHour
	return time.Duration(hoursToWait) * time.Hour
}
