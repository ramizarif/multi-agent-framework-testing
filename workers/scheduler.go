package workers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"multi-agent-framework-testing/models"
	"multi-agent-framework-testing/services"
	"multi-agent-framework-testing/storage"
)

type Scheduler struct {
	store          *storage.MemoryStore
	deviceService  *services.DeviceService
	weatherService *services.WeatherService
	running        bool
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

func NewScheduler(store *storage.MemoryStore, deviceService *services.DeviceService, weatherService *services.WeatherService) *Scheduler {
	scheduler := &Scheduler{
		store:          store,
		deviceService:  deviceService,
		weatherService: weatherService,
		running:        false,
		stopChan:       make(chan struct{}),
	}
	
	weatherService.SetStore(store)
	
	return scheduler
}

func (s *Scheduler) Start() {
	if s.running {
		return
	}
	
	s.running = true
	log.Println("Scheduler started")
	
	s.wg.Add(4)
	go s.taskRunner()
	go s.energyMonitor()
	go s.securityMonitor()
	go s.systemHealthMonitor()
	
	s.wg.Wait()
}

func (s *Scheduler) Stop() {
	if !s.running {
		return
	}
	
	s.running = false
	close(s.stopChan)
	s.wg.Wait()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) taskRunner() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.processTasks()
		case <-s.stopChan:
			return
		}
	}
}

func (s *Scheduler) processTasks() {
	tasks := s.store.ListTasks()
	now := time.Now()
	
	for _, task := range tasks {
		if task.Enabled && (task.NextRun.IsZero() || task.NextRun.Before(now)) {
			s.executeTask(task)
		}
	}
}

func (s *Scheduler) executeTask(task *models.ScheduledTask) {
	log.Printf("Executing task: %s", task.Name)
	
	device, err := s.deviceService.GetDevice(task.DeviceID)
	if err != nil {
		log.Printf("Task %s failed: device not found", task.Name)
		return
	}
	
	switch task.Action {
	case "turn_on":
		s.deviceService.UpdateDevice(device.ID, map[string]interface{}{
			"power": true,
		})
		
	case "turn_off":
		s.deviceService.UpdateDevice(device.ID, map[string]interface{}{
			"power": false,
		})
		
	case "set_temperature":
		if temp, ok := task.Parameters["temperature"].(float64); ok {
			s.deviceService.UpdateDevice(device.ID, map[string]interface{}{
				"target_temp": temp,
			})
		}
		
	case "set_brightness":
		if brightness, ok := task.Parameters["brightness"].(float64); ok {
			s.deviceService.UpdateDevice(device.ID, map[string]interface{}{
				"brightness": int(brightness),
			})
		}
		
	case "lock":
		s.deviceService.UpdateDevice(device.ID, map[string]interface{}{
			"locked": true,
		})
		
	case "unlock":
		s.deviceService.UpdateDevice(device.ID, map[string]interface{}{
			"locked": false,
		})
		
	case "arm_security":
		security := s.store.GetSecurity()
		security.State = models.SecurityStateArmed
		security.LastArmed = time.Now()
		s.store.UpdateSecurity(security)
		
	case "disarm_security":
		security := s.store.GetSecurity()
		security.State = models.SecurityStateDisarmed
		s.store.UpdateSecurity(security)
		
	default:
		log.Printf("Unknown task action: %s", task.Action)
		return
	}
	
	nextRun := s.calculateNextRun(task.Schedule)
	s.store.UpdateTask(task.ID, map[string]interface{}{
		"last_run": time.Now(),
		"next_run": nextRun,
	})
	
	s.store.AddSystemEvent(models.SystemEvent{
		ID:        fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Type:      "task_executed",
		Source:    "scheduler",
		Message:   fmt.Sprintf("Task %s executed successfully", task.Name),
		Data: map[string]interface{}{
			"task_id":   task.ID,
			"device_id": task.DeviceID,
			"action":    task.Action,
		},
		Timestamp: time.Now(),
		Severity:  "info",
	})
}

func (s *Scheduler) calculateNextRun(schedule string) time.Time {
	switch schedule {
	case "daily":
		return time.Now().AddDate(0, 0, 1)
	case "hourly":
		return time.Now().Add(time.Hour)
	case "every_30_minutes":
		return time.Now().Add(30 * time.Minute)
	case "every_15_minutes":
		return time.Now().Add(15 * time.Minute)
	case "weekly":
		return time.Now().AddDate(0, 0, 7)
	default:
		return time.Now().Add(time.Hour)
	}
}

func (s *Scheduler) energyMonitor() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.collectEnergyUsage()
		case <-s.stopChan:
			return
		}
	}
}

func (s *Scheduler) collectEnergyUsage() {
	usageData := s.deviceService.CalculateEnergyUsage()
	
	for _, usage := range usageData {
		s.store.AddEnergyUsage(usage)
	}
	
	totalUsage := 0.0
	for _, usage := range usageData {
		totalUsage += usage.Usage
	}
	
	if totalUsage > 10.0 {
		s.store.AddSystemEvent(models.SystemEvent{
			ID:        fmt.Sprintf("energy_alert_%d", time.Now().UnixNano()),
			Type:      "energy_alert",
			Source:    "scheduler",
			Message:   fmt.Sprintf("High energy usage detected: %.2f kWh", totalUsage),
			Data: map[string]interface{}{
				"total_usage": totalUsage,
				"threshold":   10.0,
			},
			Timestamp: time.Now(),
			Severity:  "warning",
		})
	}
}

func (s *Scheduler) securityMonitor() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.checkSecurityStatus()
		case <-s.stopChan:
			return
		}
	}
}

func (s *Scheduler) checkSecurityStatus() {
	security := s.store.GetSecurity()
	
	if security.State == models.SecurityStateTriggered {
		timeSinceTriggered := time.Since(security.LastTriggered)
		
		if timeSinceTriggered > 5*time.Minute {
			security.State = models.SecurityStateArmed
			s.store.UpdateSecurity(security)
			
			s.store.AddSystemEvent(models.SystemEvent{
				ID:        fmt.Sprintf("security_reset_%d", time.Now().UnixNano()),
				Type:      "security_reset",
				Source:    "scheduler",
				Message:   "Security system automatically reset after timeout",
				Data: map[string]interface{}{
					"previous_state": "triggered",
					"timeout_minutes": 5,
				},
				Timestamp: time.Now(),
				Severity:  "info",
			})
		}
	}
	
	devices := s.deviceService.GetDevicesByType(models.DeviceTypeSensor)
	for _, device := range devices {
		if device.Status == models.DeviceStatusOffline {
			s.store.AddSystemEvent(models.SystemEvent{
				ID:        fmt.Sprintf("sensor_offline_%d", time.Now().UnixNano()),
				Type:      "sensor_offline",
				Source:    "scheduler",
				Message:   fmt.Sprintf("Security sensor %s is offline", device.Name),
				Data: map[string]interface{}{
					"device_id": device.ID,
					"location":  device.Location,
				},
				Timestamp: time.Now(),
				Severity:  "warning",
			})
		}
	}
}

func (s *Scheduler) systemHealthMonitor() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.checkSystemHealth()
		case <-s.stopChan:
			return
		}
	}
}

func (s *Scheduler) checkSystemHealth() {
	devices := s.deviceService.ListDevices()
	
	offlineCount := 0
	for _, device := range devices {
		if device.Status == models.DeviceStatusOffline {
			offlineCount++
		}
	}
	
	if offlineCount > len(devices)/2 {
		s.store.AddSystemEvent(models.SystemEvent{
			ID:        fmt.Sprintf("system_health_%d", time.Now().UnixNano()),
			Type:      "system_health_warning",
			Source:    "scheduler",
			Message:   fmt.Sprintf("High number of offline devices: %d/%d", offlineCount, len(devices)),
			Data: map[string]interface{}{
				"offline_count": offlineCount,
				"total_devices": len(devices),
			},
			Timestamp: time.Now(),
			Severity:  "warning",
		})
	}
	
	weather := s.weatherService.GetCurrentWeather()
	if alert := s.weatherService.GetWeatherAlert(); alert != nil {
		s.store.AddSystemEvent(models.SystemEvent{
			ID:        fmt.Sprintf("weather_alert_%d", time.Now().UnixNano()),
			Type:      "weather_alert",
			Source:    "scheduler",
			Message:   *alert,
			Data: map[string]interface{}{
				"temperature": weather.Temperature,
				"condition":   weather.Condition,
				"wind_speed":  weather.WindSpeed,
				"pressure":    weather.Pressure,
			},
			Timestamp: time.Now(),
			Severity:  "warning",
		})
	}
}

func (s *Scheduler) AddTask(task *models.ScheduledTask) error {
	if task.ID == "" {
		task.ID = fmt.Sprintf("task_%d", time.Now().UnixNano())
	}
	
	if task.NextRun.IsZero() {
		task.NextRun = s.calculateNextRun(task.Schedule)
	}
	
	task.Enabled = true
	task.CreatedAt = time.Now()
	
	return s.store.AddTask(task)
}

func (s *Scheduler) GetTask(id string) (*models.ScheduledTask, error) {
	return s.store.GetTask(id)
}

func (s *Scheduler) ListTasks() []*models.ScheduledTask {
	return s.store.ListTasks()
}

func (s *Scheduler) UpdateTask(id string, updates map[string]interface{}) error {
	return s.store.UpdateTask(id, updates)
}

func (s *Scheduler) TriggerTask(id string) error {
	task, err := s.store.GetTask(id)
	if err != nil {
		return err
	}
	
	s.executeTask(task)
	return nil
}

func (s *Scheduler) CreateAutomationRule(name string, trigger map[string]interface{}, actions []models.TestAction) error {
	task := &models.ScheduledTask{
		ID:         fmt.Sprintf("automation_%d", time.Now().UnixNano()),
		Name:       name,
		DeviceID:   "",
		Action:     "automation_rule",
		Parameters: trigger,
		Schedule:   "trigger_based",
		Enabled:    true,
		CreatedAt:  time.Now(),
	}
	
	return s.store.AddTask(task)
}

func (s *Scheduler) TriggerAutomationScenario(scenario string) error {
	switch scenario {
	case "morning_routine":
		return s.executeMorningRoutine()
	case "evening_routine":
		return s.executeEveningRoutine()
	case "away_mode":
		return s.executeAwayMode()
	case "sleep_mode":
		return s.executeSleepMode()
	case "security_breach":
		return s.executeSecurityBreach()
	default:
		return fmt.Errorf("unknown automation scenario: %s", scenario)
	}
}

func (s *Scheduler) executeMorningRoutine() error {
	lights := s.deviceService.GetDevicesByType(models.DeviceTypeLight)
	for _, light := range lights {
		s.deviceService.UpdateDevice(light.ID, map[string]interface{}{
			"power":      true,
			"brightness": 80,
		})
	}
	
	thermostats := s.deviceService.GetDevicesByType(models.DeviceTypeThermostat)
	for _, thermostat := range thermostats {
		s.deviceService.UpdateDevice(thermostat.ID, map[string]interface{}{
			"target_temp": 22.0,
		})
	}
	
	security := s.store.GetSecurity()
	security.State = models.SecurityStateDisarmed
	s.store.UpdateSecurity(security)
	
	return nil
}

func (s *Scheduler) executeEveningRoutine() error {
	lights := s.deviceService.GetDevicesByType(models.DeviceTypeLight)
	for _, light := range lights {
		s.deviceService.UpdateDevice(light.ID, map[string]interface{}{
			"power":      true,
			"brightness": 40,
		})
	}
	
	thermostats := s.deviceService.GetDevicesByType(models.DeviceTypeThermostat)
	for _, thermostat := range thermostats {
		s.deviceService.UpdateDevice(thermostat.ID, map[string]interface{}{
			"target_temp": 20.0,
		})
	}
	
	return nil
}

func (s *Scheduler) executeAwayMode() error {
	lights := s.deviceService.GetDevicesByType(models.DeviceTypeLight)
	for _, light := range lights {
		s.deviceService.UpdateDevice(light.ID, map[string]interface{}{
			"power": false,
		})
	}
	
	thermostats := s.deviceService.GetDevicesByType(models.DeviceTypeThermostat)
	for _, thermostat := range thermostats {
		s.deviceService.UpdateDevice(thermostat.ID, map[string]interface{}{
			"target_temp": 18.0,
		})
	}
	
	security := s.store.GetSecurity()
	security.State = models.SecurityStateArmed
	security.LastArmed = time.Now()
	s.store.UpdateSecurity(security)
	
	return nil
}

func (s *Scheduler) executeSleepMode() error {
	lights := s.deviceService.GetDevicesByType(models.DeviceTypeLight)
	for _, light := range lights {
		s.deviceService.UpdateDevice(light.ID, map[string]interface{}{
			"power": false,
		})
	}
	
	locks := s.deviceService.GetDevicesByType(models.DeviceTypeLock)
	for _, lock := range locks {
		s.deviceService.UpdateDevice(lock.ID, map[string]interface{}{
			"locked": true,
		})
	}
	
	return nil
}

func (s *Scheduler) executeSecurityBreach() error {
	security := s.store.GetSecurity()
	security.State = models.SecurityStateTriggered
	security.LastTriggered = time.Now()
	security.TriggeredBy = "simulation"
	s.store.UpdateSecurity(security)
	
	lights := s.deviceService.GetDevicesByType(models.DeviceTypeLight)
	for _, light := range lights {
		s.deviceService.UpdateDevice(light.ID, map[string]interface{}{
			"power":      true,
			"brightness": 100,
		})
	}
	
	return nil
}