package stealth

import (
	"fmt"
	"math/rand"
	"time"
)

// Scheduler manages when actions should occur
type Scheduler struct {
	businessHoursOnly  bool
	businessHoursStart int // 0-23
	businessHoursEnd   int // 0-23
	actionsToday       int
	lastActionTime     time.Time
	cooldownMinutes    int
}

// NewScheduler creates an action scheduler
func NewScheduler(businessHoursOnly bool, startHour, endHour, cooldownMin int) *Scheduler {
	return &Scheduler{
		businessHoursOnly:  businessHoursOnly,
		businessHoursStart: startHour,
		businessHoursEnd:   endHour,
		actionsToday:       0,
		lastActionTime:     time.Time{},
		cooldownMinutes:    cooldownMin,
	}
}

// CanPerformAction checks if action can be performed now
func (s *Scheduler) CanPerformAction() (bool, string) {
	// Check business hours
	if s.businessHoursOnly {
		if !IsWithinBusinessHours(s.businessHoursStart, s.businessHoursEnd) {
			nextStart := GetNextBusinessHourStart(s.businessHoursStart)
			return false, fmt.Sprintf("Outside business hours. Next window in %v", nextStart.Round(time.Minute))
		}
	}
	
	// Check cooldown
	if !s.lastActionTime.IsZero() {
		timeSinceLastAction := time.Since(s.lastActionTime)
		cooldownDuration := time.Duration(s.cooldownMinutes) * time.Minute
		
		if timeSinceLastAction < cooldownDuration {
			remaining := cooldownDuration - timeSinceLastAction
			return false, fmt.Sprintf("Cooldown active. Wait %v", remaining.Round(time.Second))
		}
	}
	
	return true, ""
}

// RecordAction records that an action was performed
func (s *Scheduler) RecordAction() {
	s.lastActionTime = time.Now()
	s.actionsToday++
}

// GetNextActionTime returns when next action can occur
func (s *Scheduler) GetNextActionTime() time.Time {
	if s.lastActionTime.IsZero() {
		return time.Now()
	}
	
	return s.lastActionTime.Add(time.Duration(s.cooldownMinutes) * time.Minute)
}

// GetCooldownRemaining returns remaining cooldown duration
func (s *Scheduler) GetCooldownRemaining() time.Duration {
	if s.lastActionTime.IsZero() {
		return 0
	}
	
	nextActionTime := s.GetNextActionTime()
	remaining := time.Until(nextActionTime)
	
	if remaining < 0 {
		return 0
	}
	
	return remaining
}

// GetRandomWorkingHourOffset returns offset for realistic timing
func (s *Scheduler) GetRandomWorkingHourOffset() time.Duration {
	// Spread actions throughout the day, not all at start
	now := time.Now()
	hour := now.Hour()
	
	// If early morning, maybe delay a bit
	if hour < 10 {
		return time.Duration(rand.Intn(60)) * time.Minute
	}
	
	// Lunchtime - longer delays
	if hour >= 12 && hour < 14 {
		return time.Duration(30+rand.Intn(60)) * time.Minute
	}
	
	// Normal delay
	return time.Duration(rand.Intn(30)) * time.Minute
}

// ShouldTakeBreak determines if a break should be taken
func (s *Scheduler) ShouldTakeBreak(actionsBeforeBreak int) bool {
	return s.actionsToday > 0 && s.actionsToday%actionsBeforeBreak == 0
}

// GetBreakDuration returns how long to break
func (s *Scheduler) GetBreakDuration() time.Duration {
	// Breaks range from 5-15 minutes
	return time.Duration(5+rand.Intn(10)) * time.Minute
}

// ResetDailyCounter resets the daily action counter
func (s *Scheduler) ResetDailyCounter() {
	s.actionsToday = 0
}

// GetActionsToday returns count of actions performed today
func (s *Scheduler) GetActionsToday() int {
	return s.actionsToday
}

// EstimateCompletionTime estimates when batch of actions will complete
func (s *Scheduler) EstimateCompletionTime(actionCount int, avgActionDuration time.Duration) time.Time {
	totalDuration := time.Duration(actionCount) * (avgActionDuration + time.Duration(s.cooldownMinutes)*time.Minute)
	return time.Now().Add(totalDuration)
}

// GetOptimalStartTime returns best time to start automation
func (s *Scheduler) GetOptimalStartTime() time.Time {
	now := time.Now()
	
	if !s.businessHoursOnly {
		return now
	}
	
	// If within business hours, start now
	if IsWithinBusinessHours(s.businessHoursStart, s.businessHoursEnd) {
		return now
	}
	
	// Calculate next business hour
	hour := now.Hour()
	if hour >= s.businessHoursEnd {
		// Tomorrow morning
		tomorrow := now.Add(24 * time.Hour)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 
			s.businessHoursStart, 0, 0, 0, tomorrow.Location())
	}
	
	// Later today
	return time.Date(now.Year(), now.Month(), now.Day(), 
		s.businessHoursStart, 0, 0, 0, now.Location())
}
