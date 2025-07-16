package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"multi-agent-framework-testing/models"
)

type MemoryStore struct {
	devices      map[string]*models.Device
	weather      *models.WeatherData
	security     *models.SecuritySystem
	tasks        map[string]*models.ScheduledTask
	energyUsage  []models.EnergyUsage
	systemEvents []models.SystemEvent
	mu           sync.RWMutex
	startTime    time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		devices:      make(map[string]*models.Device),
		weather:      &models.WeatherData{},
		security:     &models.SecuritySystem{State: models.SecurityStateDisarmed},
		tasks:        make(map[string]*models.ScheduledTask),
		energyUsage:  make([]models.EnergyUsage, 0),
		systemEvents: make([]models.SystemEvent, 0),
		startTime:    time.Now(),
	}
}

func (s *MemoryStore) AddDevice(device *models.Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.devices[device.ID]; exists {
		return fmt.Errorf("device with ID %s already exists", device.ID)
	}
	
	device.CreatedAt = time.Now()
	device.LastUpdated = time.Now()
	s.devices[device.ID] = device
	
	s.addSystemEvent("device_added", "storage", fmt.Sprintf("Device %s added", device.Name), map[string]interface{}{
		"device_id": device.ID,
		"device_type": device.Type,
	})
	
	return nil
}

func (s *MemoryStore) GetDevice(id string) (*models.Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	device, exists := s.devices[id]
	if !exists {
		return nil, fmt.Errorf("device with ID %s not found", id)
	}
	
	return device, nil
}

func (s *MemoryStore) UpdateDevice(id string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	device, exists := s.devices[id]
	if !exists {
		return fmt.Errorf("device with ID %s not found", id)
	}
	
	if device.Properties == nil {
		device.Properties = make(map[string]interface{})
	}
	
	for key, value := range updates {
		switch key {
		case "name":
			if name, ok := value.(string); ok {
				device.Name = name
			}
		case "status":
			if status, ok := value.(string); ok {
				device.Status = models.DeviceStatus(status)
			}
		case "location":
			if location, ok := value.(string); ok {
				device.Location = location
			}
		default:
			device.Properties[key] = value
		}
	}
	
	device.LastUpdated = time.Now()
	
	s.addSystemEvent("device_updated", "storage", fmt.Sprintf("Device %s updated", device.Name), map[string]interface{}{
		"device_id": device.ID,
		"updates": updates,
	})
	
	return nil
}

func (s *MemoryStore) ListDevices() []*models.Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	devices := make([]*models.Device, 0, len(s.devices))
	for _, device := range s.devices {
		devices = append(devices, device)
	}
	
	return devices
}

func (s *MemoryStore) DeleteDevice(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	device, exists := s.devices[id]
	if !exists {
		return fmt.Errorf("device with ID %s not found", id)
	}
	
	delete(s.devices, id)
	
	s.addSystemEvent("device_deleted", "storage", fmt.Sprintf("Device %s deleted", device.Name), map[string]interface{}{
		"device_id": device.ID,
	})
	
	return nil
}

func (s *MemoryStore) UpdateWeather(weather *models.WeatherData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.weather = weather
	
	s.addSystemEvent("weather_updated", "storage", "Weather data updated", map[string]interface{}{
		"temperature": weather.Temperature,
		"condition": weather.Condition,
	})
}

func (s *MemoryStore) GetWeather() *models.WeatherData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.weather
}

func (s *MemoryStore) UpdateSecurity(security *models.SecuritySystem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.security = security
	
	s.addSystemEvent("security_updated", "storage", fmt.Sprintf("Security state changed to %s", security.State), map[string]interface{}{
		"state": security.State,
		"sensors": security.ActiveSensors,
	})
}

func (s *MemoryStore) GetSecurity() *models.SecuritySystem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.security
}

func (s *MemoryStore) AddTask(task *models.ScheduledTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}
	
	task.CreatedAt = time.Now()
	s.tasks[task.ID] = task
	
	s.addSystemEvent("task_added", "storage", fmt.Sprintf("Scheduled task %s added", task.Name), map[string]interface{}{
		"task_id": task.ID,
		"device_id": task.DeviceID,
		"schedule": task.Schedule,
	})
	
	return nil
}

func (s *MemoryStore) GetTask(id string) (*models.ScheduledTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}
	
	return task, nil
}

func (s *MemoryStore) ListTasks() []*models.ScheduledTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	tasks := make([]*models.ScheduledTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	
	return tasks
}

func (s *MemoryStore) UpdateTask(id string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}
	
	for key, value := range updates {
		switch key {
		case "enabled":
			if enabled, ok := value.(bool); ok {
				task.Enabled = enabled
			}
		case "next_run":
			if nextRun, ok := value.(time.Time); ok {
				task.NextRun = nextRun
			}
		case "last_run":
			if lastRun, ok := value.(time.Time); ok {
				task.LastRun = lastRun
			}
		}
	}
	
	return nil
}

func (s *MemoryStore) AddEnergyUsage(usage models.EnergyUsage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.energyUsage = append(s.energyUsage, usage)
	
	if len(s.energyUsage) > 1000 {
		s.energyUsage = s.energyUsage[len(s.energyUsage)-1000:]
	}
}

func (s *MemoryStore) GetEnergyUsage(limit int) []models.EnergyUsage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if limit <= 0 || limit > len(s.energyUsage) {
		return s.energyUsage
	}
	
	return s.energyUsage[len(s.energyUsage)-limit:]
}

func (s *MemoryStore) AddSystemEvent(event models.SystemEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.systemEvents = append(s.systemEvents, event)
	
	if len(s.systemEvents) > 500 {
		s.systemEvents = s.systemEvents[len(s.systemEvents)-500:]
	}
}

func (s *MemoryStore) GetSystemEvents(limit int) []models.SystemEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if limit <= 0 || limit > len(s.systemEvents) {
		return s.systemEvents
	}
	
	return s.systemEvents[len(s.systemEvents)-limit:]
}

func (s *MemoryStore) GetSystemState() *models.SystemState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	devices := make([]models.Device, 0, len(s.devices))
	for _, device := range s.devices {
		devices = append(devices, *device)
	}
	
	tasks := make([]models.ScheduledTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, *task)
	}
	
	return &models.SystemState{
		Devices:      devices,
		Weather:      *s.weather,
		Security:     *s.security,
		Tasks:        tasks,
		EnergyUsage:  s.energyUsage,
		SystemEvents: s.systemEvents,
		Uptime:       time.Since(s.startTime),
		Timestamp:    time.Now(),
	}
}

func (s *MemoryStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.devices = make(map[string]*models.Device)
	s.weather = &models.WeatherData{}
	s.security = &models.SecuritySystem{State: models.SecurityStateDisarmed}
	s.tasks = make(map[string]*models.ScheduledTask)
	s.energyUsage = make([]models.EnergyUsage, 0)
	s.systemEvents = make([]models.SystemEvent, 0)
	s.startTime = time.Now()
	
	s.addSystemEvent("system_reset", "storage", "System state reset", map[string]interface{}{})
}

func (s *MemoryStore) SaveToFile(filename string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	state := s.GetSystemState()
	
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}

func (s *MemoryStore) addSystemEvent(eventType, source, message string, data map[string]interface{}) {
	event := models.SystemEvent{
		ID:        fmt.Sprintf("event_%d", time.Now().UnixNano()),
		Type:      eventType,
		Source:    source,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		Severity:  "info",
	}
	
	s.systemEvents = append(s.systemEvents, event)
	
	if len(s.systemEvents) > 500 {
		s.systemEvents = s.systemEvents[len(s.systemEvents)-500:]
	}
}