package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"multi-agent-framework-testing/models"
	"multi-agent-framework-testing/services"
	"multi-agent-framework-testing/storage"
	"multi-agent-framework-testing/workers"
)

type Handler struct {
	store          *storage.MemoryStore
	deviceService  *services.DeviceService
	weatherService *services.WeatherService
	scheduler      *workers.Scheduler
	upgrader       *websocket.Upgrader
	wsClients      map[*websocket.Conn]bool
	wsClientsMu    sync.RWMutex
	rateLimiter    map[string]time.Time
	rateLimiterMu  sync.Mutex
}

func NewHandler(store *storage.MemoryStore, deviceService *services.DeviceService, 
	weatherService *services.WeatherService, scheduler *workers.Scheduler, 
	upgrader *websocket.Upgrader) *Handler {
	
	handler := &Handler{
		store:          store,
		deviceService:  deviceService,
		weatherService: weatherService,
		scheduler:      scheduler,
		upgrader:       upgrader,
		wsClients:      make(map[*websocket.Conn]bool),
		rateLimiter:    make(map[string]time.Time),
	}
	
	go handler.broadcastUpdates()
	
	return handler
}

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/health") || strings.HasPrefix(r.URL.Path, "/ws") {
			next.ServeHTTP(w, r)
			return
		}
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			h.respondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization header")
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "smarthome-secret-token" {
			h.respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		
		h.rateLimiterMu.Lock()
		lastRequest, exists := h.rateLimiter[clientIP]
		
		if exists && time.Since(lastRequest) < time.Millisecond*10 {
			h.rateLimiterMu.Unlock()
			h.respondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}
		
		h.rateLimiter[clientIP] = time.Now()
		
		if len(h.rateLimiter) > 10000 {
			for ip, timestamp := range h.rateLimiter {
				if time.Since(timestamp) > time.Minute {
					delete(h.rateLimiter, ip)
				}
			}
		}
		
		h.rateLimiterMu.Unlock()
		
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	state := h.store.GetSystemState()
	
	devices := make([]*models.Device, len(state.Devices))
	for i := range state.Devices {
		devices[i] = &state.Devices[i]
	}
	
	health := map[string]interface{}{
		"status":         "healthy",
		"uptime":         state.Uptime.String(),
		"total_devices":  len(state.Devices),
		"online_devices": h.countOnlineDevices(devices),
		"weather_status": state.Weather.Condition,
		"security_state": state.Security.State,
		"timestamp":      time.Now(),
	}
	
	h.respondWithJSON(w, http.StatusOK, health)
}

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices := h.deviceService.ListDevices()
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    devices,
	})
}

func (h *Handler) AddDevice(w http.ResponseWriter, r *http.Request) {
	var device models.Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if device.Name == "" || device.Type == "" {
		h.respondWithError(w, http.StatusBadRequest, "Device name and type are required")
		return
	}
	
	if err := h.deviceService.AddDevice(&device); err != nil {
		h.respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	
	h.broadcastMessage("device_added", device)
	
	h.respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    device,
		Message: "Device added successfully",
	})
}

func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["id"]
	
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if err := h.deviceService.UpdateDevice(deviceID, updates); err != nil {
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	
	device, _ := h.deviceService.GetDevice(deviceID)
	h.broadcastMessage("device_updated", device)
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    device,
		Message: "Device updated successfully",
	})
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	weather := h.weatherService.GetCurrentWeather()
	
	response := map[string]interface{}{
		"current":  weather,
		"forecast": h.weatherService.GetForecast(5),
		"alert":    h.weatherService.GetWeatherAlert(),
	}
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func (h *Handler) GetEnergyUsage(w http.ResponseWriter, r *http.Request) {
	usage := h.store.GetEnergyUsage(100)
	currentUsage := h.deviceService.CalculateEnergyUsage()
	
	totalUsage := 0.0
	totalCost := 0.0
	for _, u := range currentUsage {
		totalUsage += u.Usage
		totalCost += u.Cost
	}
	
	response := map[string]interface{}{
		"current_usage":     currentUsage,
		"historical_usage":  usage,
		"total_usage_kwh":   totalUsage,
		"total_cost_usd":    totalCost,
		"usage_by_device":   h.aggregateUsageByDevice(usage),
	}
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func (h *Handler) ArmSecurity(w http.ResponseWriter, r *http.Request) {
	security := h.store.GetSecurity()
	security.State = models.SecurityStateArmed
	security.LastArmed = time.Now()
	
	sensors := h.deviceService.GetDevicesByType(models.DeviceTypeSensor)
	security.ActiveSensors = make([]string, 0)
	for _, sensor := range sensors {
		if sensor.Status == models.DeviceStatusOnline {
			security.ActiveSensors = append(security.ActiveSensors, sensor.ID)
		}
	}
	
	h.store.UpdateSecurity(security)
	h.broadcastMessage("security_armed", security)
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    security,
		Message: "Security system armed",
	})
}

func (h *Handler) DisarmSecurity(w http.ResponseWriter, r *http.Request) {
	security := h.store.GetSecurity()
	security.State = models.SecurityStateDisarmed
	security.ActiveSensors = []string{}
	
	h.store.UpdateSecurity(security)
	h.broadcastMessage("security_disarmed", security)
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    security,
		Message: "Security system disarmed",
	})
}

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	devices := h.deviceService.ListDevices()
	energyUsage := h.store.GetEnergyUsage(500)
	events := h.store.GetSystemEvents(100)
	weather := h.weatherService.GetCurrentWeather()
	tasks := h.scheduler.ListTasks()
	
	analytics := models.AnalyticsData{
		TotalDevices:     len(devices),
		OnlineDevices:    h.countOnlineDevices(devices),
		OfflineDevices:   h.countOfflineDevices(devices),
		TotalEnergyUsage: h.calculateTotalEnergy(energyUsage),
		TotalEnergyCost:  h.calculateTotalCost(energyUsage),
		DeviceUsage:      h.aggregateUsageByDevice(energyUsage),
		SecurityEvents:   h.countSecurityEvents(events),
		ScheduledTasks:   len(tasks),
		WeatherSummary:   *weather,
		GeneratedAt:      time.Now(),
	}
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    analytics,
	})
}

func (h *Handler) CreateScheduledTask(w http.ResponseWriter, r *http.Request) {
	var task models.ScheduledTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if task.Name == "" || task.DeviceID == "" || task.Action == "" {
		h.respondWithError(w, http.StatusBadRequest, "Task name, device ID, and action are required")
		return
	}
	
	if err := h.scheduler.AddTask(&task); err != nil {
		h.respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	
	h.respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    task,
		Message: "Scheduled task created successfully",
	})
}

func (h *Handler) DebugState(w http.ResponseWriter, r *http.Request) {
	state := h.store.GetSystemState()
	h.respondWithJSON(w, http.StatusOK, state)
}

func (h *Handler) ResetSystem(w http.ResponseWriter, r *http.Request) {
	h.store.Reset()
	h.deviceService = services.NewDeviceService(h.store)
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "System reset successfully",
	})
}

func (h *Handler) TriggerScenario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scenario := vars["scenario"]
	
	switch scenario {
	case "weather":
		weatherScenarios := []string{"storm", "heatwave", "cold_snap", "rain", "fog"}
		for _, ws := range weatherScenarios {
			if scenario == ws || strings.Contains(scenario, ws) {
				h.weatherService.SimulateWeatherScenario(ws)
				break
			}
		}
		
	case "device_failure":
		devices := h.deviceService.ListDevices()
		if len(devices) > 0 {
			device := devices[0]
			h.deviceService.UpdateDevice(device.ID, map[string]interface{}{
				"status": models.DeviceStatusOffline,
			})
		}
		
	case "power_surge":
		usageData := h.deviceService.CalculateEnergyUsage()
		for _, usage := range usageData {
			usage.Usage *= 5
			usage.Cost *= 5
			h.store.AddEnergyUsage(usage)
		}
		
	case "morning_routine", "evening_routine", "away_mode", "sleep_mode", "security_breach":
		if err := h.scheduler.TriggerAutomationScenario(scenario); err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		
	default:
		h.respondWithError(w, http.StatusBadRequest, "Unknown scenario")
		return
	}
	
	h.respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Scenario '%s' triggered successfully", scenario),
	})
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	
	h.wsClientsMu.Lock()
	h.wsClients[conn] = true
	h.wsClientsMu.Unlock()
	
	defer func() {
		h.wsClientsMu.Lock()
		delete(h.wsClients, conn)
		h.wsClientsMu.Unlock()
	}()
	
	state := h.store.GetSystemState()
	if err := conn.WriteJSON(models.WebSocketMessage{
		Type:      "initial_state",
		Data:      state,
		Timestamp: time.Now(),
	}); err != nil {
		log.Printf("WebSocket write error: %v", err)
		return
	}
	
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func (h *Handler) broadcastMessage(messageType string, data interface{}) {
	h.wsClientsMu.RLock()
	defer h.wsClientsMu.RUnlock()
	
	message := models.WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}
	
	for client := range h.wsClients {
		if err := client.WriteJSON(message); err != nil {
			log.Printf("WebSocket broadcast error: %v", err)
			client.Close()
		}
	}
}

func (h *Handler) broadcastUpdates() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			h.broadcastMessage("state_update", map[string]interface{}{
				"devices": h.deviceService.ListDevices(),
				"weather": h.weatherService.GetCurrentWeather(),
				"security": h.store.GetSecurity(),
			})
		}
	}
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *Handler) respondWithError(w http.ResponseWriter, status int, message string) {
	h.respondWithJSON(w, status, models.APIResponse{
		Success: false,
		Error:   message,
	})
}

func (h *Handler) countOnlineDevices(devices []*models.Device) int {
	count := 0
	for _, device := range devices {
		if device.Status == models.DeviceStatusOnline {
			count++
		}
	}
	return count
}

func (h *Handler) countOfflineDevices(devices []*models.Device) int {
	count := 0
	for _, device := range devices {
		if device.Status == models.DeviceStatusOffline {
			count++
		}
	}
	return count
}

func (h *Handler) calculateTotalEnergy(usage []models.EnergyUsage) float64 {
	total := 0.0
	for _, u := range usage {
		total += u.Usage
	}
	return total
}

func (h *Handler) calculateTotalCost(usage []models.EnergyUsage) float64 {
	total := 0.0
	for _, u := range usage {
		total += u.Cost
	}
	return total
}

func (h *Handler) aggregateUsageByDevice(usage []models.EnergyUsage) map[string]float64 {
	aggregate := make(map[string]float64)
	for _, u := range usage {
		aggregate[u.DeviceName] += u.Usage
	}
	return aggregate
}

func (h *Handler) countSecurityEvents(events []models.SystemEvent) int {
	count := 0
	for _, event := range events {
		if strings.Contains(event.Type, "security") {
			count++
		}
	}
	return count
}