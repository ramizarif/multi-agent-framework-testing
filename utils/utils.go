package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"sync"
	"time"
)

func GenerateID(prefix string) string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes))
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func RoundToDecimal(value float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(value*shift) / shift
}

func CalculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return RoundToDecimal((part/total)*100, 2)
}

func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	return RoundToDecimal((celsius*9/5)+32, 1)
}

func ConvertFahrenheitToCelsius(fahrenheit float64) float64 {
	return RoundToDecimal((fahrenheit-32)*5/9, 1)
}

func IsWithinRange(value, target, tolerance float64) bool {
	return math.Abs(value-target) <= tolerance
}

func LimitValue(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ParseSchedule(schedule string) (time.Duration, error) {
	scheduleMap := map[string]time.Duration{
		"every_minute":     time.Minute,
		"every_5_minutes":  5 * time.Minute,
		"every_15_minutes": 15 * time.Minute,
		"every_30_minutes": 30 * time.Minute,
		"hourly":           time.Hour,
		"daily":            24 * time.Hour,
		"weekly":           7 * 24 * time.Hour,
	}
	
	if duration, ok := scheduleMap[schedule]; ok {
		return duration, nil
	}
	
	return 0, fmt.Errorf("unknown schedule: %s", schedule)
}

func GetTimeOfDay(t time.Time) string {
	hour := t.Hour()
	
	switch {
	case hour >= 5 && hour < 12:
		return "morning"
	case hour >= 12 && hour < 17:
		return "afternoon"
	case hour >= 17 && hour < 21:
		return "evening"
	default:
		return "night"
	}
}

func IsBusinessHours(t time.Time) bool {
	hour := t.Hour()
	weekday := t.Weekday()
	
	return weekday >= time.Monday && weekday <= time.Friday && hour >= 9 && hour < 17
}

func CalculateEnergyRate(timeOfDay string, isWeekend bool) float64 {
	baseRate := 0.12
	
	if isWeekend {
		return baseRate * 0.8
	}
	
	switch timeOfDay {
	case "morning":
		return baseRate * 1.2
	case "afternoon":
		return baseRate * 1.5
	case "evening":
		return baseRate * 1.3
	case "night":
		return baseRate * 0.7
	default:
		return baseRate
	}
}

func ValidateDeviceProperties(deviceType string, properties map[string]interface{}) error {
	switch deviceType {
	case "light":
		if brightness, ok := properties["brightness"].(int); ok {
			if brightness < 0 || brightness > 100 {
				return fmt.Errorf("brightness must be between 0 and 100")
			}
		}
		
	case "thermostat":
		if temp, ok := properties["target_temp"].(float64); ok {
			if temp < 10 || temp > 35 {
				return fmt.Errorf("temperature must be between 10 and 35 degrees")
			}
		}
		
	case "sensor":
		if sensitivity, ok := properties["sensitivity"].(int); ok {
			if sensitivity < 1 || sensitivity > 10 {
				return fmt.Errorf("sensitivity must be between 1 and 10")
			}
		}
	}
	
	return nil
}

func GenerateRandomFloat(min, max float64) float64 {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	
	value := float64(bytes[0]) / 255.0
	return min + (max-min)*value
}

func GenerateRandomInt(min, max int) int {
	if min >= max {
		return min
	}
	
	bytes := make([]byte, 4)
	rand.Read(bytes)
	
	value := int(bytes[0])
	return min + value%(max-min)
}

func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func RemoveFromSlice(slice []string, str string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != str {
			result = append(result, s)
		}
	}
	return result
}

func MergeStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func SafeMapAccess(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371
	
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

type CircuitBreaker struct {
	maxFailures int
	resetTime   time.Duration
	failures    int
	lastFailure time.Time
	state       string
}

func NewCircuitBreaker(maxFailures int, resetTime time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		resetTime:   resetTime,
		state:       "closed",
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	if cb.state == "open" {
		if time.Since(cb.lastFailure) > cb.resetTime {
			cb.state = "half-open"
			cb.failures = 0
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}
	
	err := fn()
	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		
		if cb.failures >= cb.maxFailures {
			cb.state = "open"
		}
		
		return err
	}
	
	if cb.state == "half-open" {
		cb.state = "closed"
	}
	
	cb.failures = 0
	return nil
}

type RateLimiter struct {
	rate     int
	interval time.Duration
	tokens   int
	lastTime time.Time
	mu       sync.Mutex
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		interval: interval,
		tokens:   rate,
		lastTime: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rl.lastTime)
	
	tokensToAdd := int(elapsed / rl.interval * time.Duration(rl.rate))
	if tokensToAdd > 0 {
		rl.tokens = min(rl.tokens+tokensToAdd, rl.rate)
		rl.lastTime = now
	}
	
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}