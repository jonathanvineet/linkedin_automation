package stealth

import (
	"math"
	"math/rand"
	"time"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// MouseMovement handles human-like mouse cursor movement
type MouseMovement struct {
	currentPos Point
}

// NewMouseMovement creates a new mouse movement controller
func NewMouseMovement() *MouseMovement {
	return &MouseMovement{
		currentPos: Point{X: 0, Y: 0},
	}
}

// GenerateBezierCurve creates a natural curved path between two points
func (m *MouseMovement) GenerateBezierCurve(start, end Point, steps int) []Point {
	// Add randomness to control points for natural curves
	midX := (start.X + end.X) / 2
	midY := (start.Y + end.Y) / 2
	
	// Random offset for control points
	offsetX := (rand.Float64()*2 - 1) * math.Abs(end.X-start.X) * 0.3
	offsetY := (rand.Float64()*2 - 1) * math.Abs(end.Y-start.Y) * 0.3
	
	// Quadratic Bezier curve control points
	cp1 := Point{
		X: midX + offsetX,
		Y: midY + offsetY,
	}
	
	points := make([]Point, steps)
	
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		
		// Quadratic Bezier formula: B(t) = (1-t)²P₀ + 2(1-t)tP₁ + t²P₂
		oneMinusT := 1 - t
		
		points[i] = Point{
			X: oneMinusT*oneMinusT*start.X + 2*oneMinusT*t*cp1.X + t*t*end.X,
			Y: oneMinusT*oneMinusT*start.Y + 2*oneMinusT*t*cp1.Y + t*t*end.Y,
		}
	}
	
	// Add overshoot and correction at the end (human tendency)
	if rand.Float64() < 0.4 {
		overshoot := 5 + rand.Intn(10)
		direction := 1.0
		if rand.Float64() < 0.5 {
			direction = -1.0
		}
		
		points = append(points, Point{
			X: end.X + direction*float64(overshoot),
			Y: end.Y + direction*float64(overshoot),
		})
		points = append(points, end) // Correct back
	}
	
	return points
}

// GetMovementDelay calculates delay between mouse movements
func (m *MouseMovement) GetMovementDelay() time.Duration {
	// Human mouse moves at varying speeds: 1-3ms between coordinate updates
	delay := 1 + rand.Intn(3)
	return time.Duration(delay) * time.Millisecond
}

// AddMicroCorrections simulates small mouse adjustments
func (m *MouseMovement) AddMicroCorrections(point Point, precision int) Point {
	// Lower precision = more deviation
	maxDeviation := float64(100-precision) / 10.0
	
	return Point{
		X: point.X + (rand.Float64()*2-1)*maxDeviation,
		Y: point.Y + (rand.Float64()*2-1)*maxDeviation,
	}
}

// SimulateMouseWander generates idle mouse movement
func (m *MouseMovement) SimulateMouseWander(currentPos Point, wanderRadius float64, steps int) []Point {
	points := make([]Point, steps)
	
	for i := 0; i < steps; i++ {
		angle := rand.Float64() * 2 * math.Pi
		distance := rand.Float64() * wanderRadius
		
		points[i] = Point{
			X: currentPos.X + distance*math.Cos(angle),
			Y: currentPos.Y + distance*math.Sin(angle),
		}
		
		currentPos = points[i]
	}
	
	return points
}

// GetHoverPosition calculates natural hover position near element
func (m *MouseMovement) GetHoverPosition(elementCenter Point, elementWidth, elementHeight float64) Point {
	// Humans don't hover exactly at center
	offsetX := (rand.Float64()*2 - 1) * elementWidth * 0.3
	offsetY := (rand.Float64()*2 - 1) * elementHeight * 0.3
	
	return Point{
		X: elementCenter.X + offsetX,
		Y: elementCenter.Y + offsetY,
	}
}

// Distance calculates Euclidean distance between two points
func Distance(p1, p2 Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// SetCurrentPosition updates the mouse's current position
func (m *MouseMovement) SetCurrentPosition(pos Point) {
	m.currentPos = pos
}

// GetCurrentPosition returns the mouse's current position
func (m *MouseMovement) GetCurrentPosition() Point {
	return m.currentPos
}
