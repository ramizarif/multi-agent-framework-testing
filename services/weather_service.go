package services

import (
	"math"
	"math/rand"
	"time"

	"multi-agent-framework-testing/models"
	"multi-agent-framework-testing/storage"
)

type WeatherService struct {
	store   *storage.MemoryStore
	current *models.WeatherData
}

func NewWeatherService() *WeatherService {
	service := &WeatherService{
		current: &models.WeatherData{
			Temperature: 20.0,
			Humidity:    60.0,
			Pressure:    1013.25,
			Condition:   "clear",
			WindSpeed:   5.0,
			WindDir:     "N",
			Timestamp:   time.Now(),
		},
	}
	
	return service
}

func (w *WeatherService) SetStore(store *storage.MemoryStore) {
	w.store = store
	w.store.UpdateWeather(w.current)
	go w.startWeatherUpdates()
}

func (w *WeatherService) startWeatherUpdates() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			w.updateWeather()
		}
	}
}

func (w *WeatherService) updateWeather() {
	now := time.Now()
	hour := now.Hour()
	
	baseTemp := w.calculateBaseTemperature(hour)
	
	tempVariation := (rand.Float64() - 0.5) * 4.0
	w.current.Temperature = baseTemp + tempVariation
	
	w.current.Humidity = w.calculateHumidity()
	w.current.Pressure = w.calculatePressure()
	w.current.Condition = w.calculateCondition()
	w.current.WindSpeed = w.calculateWindSpeed()
	w.current.WindDir = w.calculateWindDirection()
	w.current.Timestamp = now
	
	if w.store != nil {
		w.store.UpdateWeather(w.current)
	}
}

func (w *WeatherService) calculateBaseTemperature(hour int) float64 {
	dayOfYear := time.Now().YearDay()
	seasonalVariation := 10.0 * math.Sin(2.0*math.Pi*float64(dayOfYear)/365.0)
	
	hourlyVariation := 5.0 * math.Sin(2.0*math.Pi*float64(hour-6)/24.0)
	
	baseTemp := 18.0 + seasonalVariation + hourlyVariation
	
	return baseTemp
}

func (w *WeatherService) calculateHumidity() float64 {
	currentHumidity := w.current.Humidity
	
	variation := (rand.Float64() - 0.5) * 10.0
	newHumidity := currentHumidity + variation
	
	if newHumidity < 30.0 {
		newHumidity = 30.0
	} else if newHumidity > 90.0 {
		newHumidity = 90.0
	}
	
	return newHumidity
}

func (w *WeatherService) calculatePressure() float64 {
	currentPressure := w.current.Pressure
	
	variation := (rand.Float64() - 0.5) * 20.0
	newPressure := currentPressure + variation
	
	if newPressure < 980.0 {
		newPressure = 980.0
	} else if newPressure > 1040.0 {
		newPressure = 1040.0
	}
	
	return newPressure
}

func (w *WeatherService) calculateCondition() string {
	pressure := w.current.Pressure
	humidity := w.current.Humidity
	
	if pressure < 1000.0 && humidity > 80.0 {
		conditions := []string{"rainy", "stormy", "cloudy"}
		return conditions[rand.Intn(len(conditions))]
	} else if pressure < 1010.0 {
		conditions := []string{"cloudy", "partly_cloudy", "overcast"}
		return conditions[rand.Intn(len(conditions))]
	} else if humidity < 40.0 {
		return "clear"
	} else {
		conditions := []string{"clear", "partly_cloudy", "sunny"}
		return conditions[rand.Intn(len(conditions))]
	}
}

func (w *WeatherService) calculateWindSpeed() float64 {
	currentSpeed := w.current.WindSpeed
	
	variation := (rand.Float64() - 0.5) * 5.0
	newSpeed := currentSpeed + variation
	
	if newSpeed < 0.0 {
		newSpeed = 0.0
	} else if newSpeed > 30.0 {
		newSpeed = 30.0
	}
	
	return newSpeed
}

func (w *WeatherService) calculateWindDirection() string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	
	currentIndex := w.findDirectionIndex(w.current.WindDir)
	
	change := rand.Intn(3) - 1
	newIndex := (currentIndex + change + len(directions)) % len(directions)
	
	return directions[newIndex]
}

func (w *WeatherService) findDirectionIndex(direction string) int {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	for i, dir := range directions {
		if dir == direction {
			return i
		}
	}
	return 0
}

func (w *WeatherService) GetCurrentWeather() *models.WeatherData {
	return w.current
}

func (w *WeatherService) SetWeather(weather *models.WeatherData) {
	w.current = weather
	w.current.Timestamp = time.Now()
	
	if w.store != nil {
		w.store.UpdateWeather(w.current)
	}
}

func (w *WeatherService) SimulateWeatherScenario(scenario string) {
	switch scenario {
	case "storm":
		w.current.Temperature = 15.0 + rand.Float64()*5.0
		w.current.Humidity = 85.0 + rand.Float64()*10.0
		w.current.Pressure = 980.0 + rand.Float64()*20.0
		w.current.Condition = "stormy"
		w.current.WindSpeed = 20.0 + rand.Float64()*10.0
		w.current.WindDir = "SW"
		
	case "heatwave":
		w.current.Temperature = 35.0 + rand.Float64()*10.0
		w.current.Humidity = 20.0 + rand.Float64()*15.0
		w.current.Pressure = 1020.0 + rand.Float64()*10.0
		w.current.Condition = "clear"
		w.current.WindSpeed = 2.0 + rand.Float64()*5.0
		w.current.WindDir = "S"
		
	case "cold_snap":
		w.current.Temperature = -5.0 + rand.Float64()*10.0
		w.current.Humidity = 40.0 + rand.Float64()*20.0
		w.current.Pressure = 1030.0 + rand.Float64()*10.0
		w.current.Condition = "clear"
		w.current.WindSpeed = 15.0 + rand.Float64()*10.0
		w.current.WindDir = "N"
		
	case "rain":
		w.current.Temperature = 12.0 + rand.Float64()*8.0
		w.current.Humidity = 80.0 + rand.Float64()*15.0
		w.current.Pressure = 995.0 + rand.Float64()*15.0
		w.current.Condition = "rainy"
		w.current.WindSpeed = 8.0 + rand.Float64()*7.0
		w.current.WindDir = "W"
		
	case "fog":
		w.current.Temperature = 8.0 + rand.Float64()*5.0
		w.current.Humidity = 95.0 + rand.Float64()*5.0
		w.current.Pressure = 1015.0 + rand.Float64()*5.0
		w.current.Condition = "foggy"
		w.current.WindSpeed = 1.0 + rand.Float64()*2.0
		w.current.WindDir = "Calm"
		
	default:
		w.current.Temperature = 20.0 + rand.Float64()*10.0
		w.current.Humidity = 50.0 + rand.Float64()*20.0
		w.current.Pressure = 1013.0 + rand.Float64()*10.0
		w.current.Condition = "clear"
		w.current.WindSpeed = 5.0 + rand.Float64()*5.0
		w.current.WindDir = "N"
	}
	
	w.current.Timestamp = time.Now()
	
	if w.store != nil {
		w.store.UpdateWeather(w.current)
	}
}

func (w *WeatherService) GetForecast(days int) []models.WeatherData {
	forecast := make([]models.WeatherData, days)
	
	for i := 0; i < days; i++ {
		futureDate := time.Now().AddDate(0, 0, i+1)
		
		baseTemp := w.calculateBaseTemperature(12)
		tempVariation := (rand.Float64() - 0.5) * 8.0
		
		forecast[i] = models.WeatherData{
			Temperature: baseTemp + tempVariation,
			Humidity:    40.0 + rand.Float64()*40.0,
			Pressure:    1000.0 + rand.Float64()*30.0,
			Condition:   w.calculateCondition(),
			WindSpeed:   rand.Float64() * 20.0,
			WindDir:     w.calculateWindDirection(),
			Timestamp:   futureDate,
		}
	}
	
	return forecast
}

func (w *WeatherService) GetWeatherHistory(hours int) []models.WeatherData {
	history := make([]models.WeatherData, hours)
	
	for i := 0; i < hours; i++ {
		pastTime := time.Now().Add(-time.Duration(hours-i) * time.Hour)
		hour := pastTime.Hour()
		
		baseTemp := w.calculateBaseTemperature(hour)
		tempVariation := (rand.Float64() - 0.5) * 6.0
		
		history[i] = models.WeatherData{
			Temperature: baseTemp + tempVariation,
			Humidity:    45.0 + rand.Float64()*35.0,
			Pressure:    1005.0 + rand.Float64()*25.0,
			Condition:   w.calculateCondition(),
			WindSpeed:   rand.Float64() * 15.0,
			WindDir:     w.calculateWindDirection(),
			Timestamp:   pastTime,
		}
	}
	
	return history
}

func (w *WeatherService) IsExtremeWeather() bool {
	return w.current.Temperature < 0.0 || w.current.Temperature > 35.0 ||
		w.current.WindSpeed > 20.0 || w.current.Condition == "stormy" ||
		w.current.Pressure < 990.0
}

func (w *WeatherService) GetWeatherAlert() *string {
	if w.current.Temperature > 35.0 {
		alert := "Extreme heat warning - Temperature above 35°C"
		return &alert
	}
	
	if w.current.Temperature < 0.0 {
		alert := "Freezing temperature alert - Temperature below 0°C"
		return &alert
	}
	
	if w.current.WindSpeed > 20.0 {
		alert := "High wind warning - Wind speed above 20 km/h"
		return &alert
	}
	
	if w.current.Condition == "stormy" {
		alert := "Storm warning - Severe weather conditions"
		return &alert
	}
	
	if w.current.Pressure < 990.0 {
		alert := "Low pressure system - Potential weather instability"
		return &alert
	}
	
	return nil
}