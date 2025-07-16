package services

import (
	"fmt"
	"math/rand"
	"time"

	"multi-agent-framework-testing/models"
	"multi-agent-framework-testing/storage"
)

type DeviceService struct {
	store *storage.MemoryStore
}

func NewDeviceService(store *storage.MemoryStore) *DeviceService {
	service := &DeviceService{
		store: store,
	}
	
	service.initializeDefaultDevices()
	go service.simulateDeviceUpdates()
	
	return service
}

func (d *DeviceService) initializeDefaultDevices() {
	defaultDevices := []*models.Device{
		{
			ID:       "light_001",
			Name:     "Living Room Light",
			Type:     models.DeviceTypeLight,
			Status:   models.DeviceStatusOnline,
			Location: "Living Room",
			Properties: map[string]interface{}{
				"brightness": 80,
				"color":      "warm_white",
				"power":      true,
			},
		},
		{
			ID:       "thermostat_001",
			Name:     "Main Thermostat",
			Type:     models.DeviceTypeThermostat,
			Status:   models.DeviceStatusOnline,
			Location: "Hallway",
			Properties: map[string]interface{}{
				"temperature":     22.5,
				"target_temp":     21.0,
				"mode":            "auto",
				"heating":         false,
				"cooling":         false,
			},
		},
		{
			ID:       "camera_001",
			Name:     "Front Door Camera",
			Type:     models.DeviceTypeCamera,
			Status:   models.DeviceStatusOnline,
			Location: "Front Door",
			Properties: map[string]interface{}{
				"recording":      true,
				"motion_detect":  true,
				"night_vision":   true,
				"resolution":     "1080p",
			},
		},
		{
			ID:       "sensor_001",
			Name:     "Motion Sensor",
			Type:     models.DeviceTypeSensor,
			Status:   models.DeviceStatusOnline,
			Location: "Living Room",
			Properties: map[string]interface{}{
				"motion_detected": false,
				"sensitivity":     7,
				"battery_level":   85,
			},
		},
		{
			ID:       "lock_001",
			Name:     "Front Door Lock",
			Type:     models.DeviceTypeLock,
			Status:   models.DeviceStatusOnline,
			Location: "Front Door",
			Properties: map[string]interface{}{
				"locked":        true,
				"auto_lock":     true,
				"battery_level": 92,
			},
		},
	}
	
	for _, device := range defaultDevices {
		d.store.AddDevice(device)
	}
}

func (d *DeviceService) simulateDeviceUpdates() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			d.updateDeviceStates()
		}
	}
}

func (d *DeviceService) updateDeviceStates() {
	devices := d.store.ListDevices()
	
	for _, device := range devices {
		if rand.Float64() < 0.1 {
			d.simulateDeviceChange(device)
		}
		
		if device.Type == models.DeviceTypeThermostat {
			d.updateThermostat(device)
		}
		
		if device.Type == models.DeviceTypeSensor {
			d.updateSensor(device)
		}
	}
}

func (d *DeviceService) simulateDeviceChange(device *models.Device) {
	updates := make(map[string]interface{})
	
	switch device.Type {
	case models.DeviceTypeLight:
		if rand.Float64() < 0.5 {
			updates["brightness"] = rand.Intn(100) + 1
		}
		if rand.Float64() < 0.3 {
			updates["power"] = rand.Float64() < 0.7
		}
		
	case models.DeviceTypeCamera:
		if rand.Float64() < 0.2 {
			updates["motion_detect"] = rand.Float64() < 0.8
		}
		if rand.Float64() < 0.1 {
			updates["recording"] = rand.Float64() < 0.9
		}
		
	case models.DeviceTypeLock:
		if rand.Float64() < 0.1 {
			updates["locked"] = rand.Float64() < 0.9
		}
		if rand.Float64() < 0.05 {
			updates["battery_level"] = rand.Intn(100) + 1
		}
	}
	
	if len(updates) > 0 {
		d.store.UpdateDevice(device.ID, updates)
	}
}

func (d *DeviceService) updateThermostat(device *models.Device) {
	weather := d.store.GetWeather()
	if weather.Temperature == 0 {
		return
	}
	
	currentTemp := device.Properties["temperature"].(float64)
	targetTemp := device.Properties["target_temp"].(float64)
	
	tempDiff := targetTemp - currentTemp
	
	updates := make(map[string]interface{})
	
	if tempDiff > 1.0 {
		updates["heating"] = true
		updates["cooling"] = false
		updates["temperature"] = currentTemp + 0.5
	} else if tempDiff < -1.0 {
		updates["heating"] = false
		updates["cooling"] = true
		updates["temperature"] = currentTemp - 0.5
	} else {
		updates["heating"] = false
		updates["cooling"] = false
	}
	
	weatherInfluence := (weather.Temperature - currentTemp) * 0.1
	if weatherInfluence != 0 {
		updates["temperature"] = currentTemp + weatherInfluence
	}
	
	if len(updates) > 0 {
		d.store.UpdateDevice(device.ID, updates)
	}
}

func (d *DeviceService) updateSensor(device *models.Device) {
	if rand.Float64() < 0.05 {
		motionDetected := rand.Float64() < 0.3
		
		updates := map[string]interface{}{
			"motion_detected": motionDetected,
		}
		
		if motionDetected {
			security := d.store.GetSecurity()
			if security.State == models.SecurityStateArmed {
				security.State = models.SecurityStateTriggered
				security.LastTriggered = time.Now()
				security.TriggeredBy = device.ID
				d.store.UpdateSecurity(security)
			}
		}
		
		d.store.UpdateDevice(device.ID, updates)
	}
}

func (d *DeviceService) AddDevice(device *models.Device) error {
	if device.ID == "" {
		device.ID = fmt.Sprintf("%s_%d", device.Type, time.Now().UnixNano())
	}
	
	if device.Properties == nil {
		device.Properties = make(map[string]interface{})
	}
	
	switch device.Type {
	case models.DeviceTypeLight:
		if device.Properties["brightness"] == nil {
			device.Properties["brightness"] = 50
		}
		if device.Properties["power"] == nil {
			device.Properties["power"] = true
		}
		
	case models.DeviceTypeThermostat:
		if device.Properties["temperature"] == nil {
			device.Properties["temperature"] = 20.0
		}
		if device.Properties["target_temp"] == nil {
			device.Properties["target_temp"] = 21.0
		}
		if device.Properties["mode"] == nil {
			device.Properties["mode"] = "auto"
		}
		
	case models.DeviceTypeCamera:
		if device.Properties["recording"] == nil {
			device.Properties["recording"] = true
		}
		if device.Properties["motion_detect"] == nil {
			device.Properties["motion_detect"] = true
		}
		
	case models.DeviceTypeSensor:
		if device.Properties["motion_detected"] == nil {
			device.Properties["motion_detected"] = false
		}
		if device.Properties["sensitivity"] == nil {
			device.Properties["sensitivity"] = 5
		}
		if device.Properties["battery_level"] == nil {
			device.Properties["battery_level"] = 100
		}
		
	case models.DeviceTypeLock:
		if device.Properties["locked"] == nil {
			device.Properties["locked"] = true
		}
		if device.Properties["battery_level"] == nil {
			device.Properties["battery_level"] = 100
		}
	}
	
	if device.Status == "" {
		device.Status = models.DeviceStatusOnline
	}
	
	return d.store.AddDevice(device)
}

func (d *DeviceService) GetDevice(id string) (*models.Device, error) {
	return d.store.GetDevice(id)
}

func (d *DeviceService) UpdateDevice(id string, updates map[string]interface{}) error {
	return d.store.UpdateDevice(id, updates)
}

func (d *DeviceService) ListDevices() []*models.Device {
	return d.store.ListDevices()
}

func (d *DeviceService) DeleteDevice(id string) error {
	return d.store.DeleteDevice(id)
}

func (d *DeviceService) GetDevicesByType(deviceType models.DeviceType) []*models.Device {
	allDevices := d.store.ListDevices()
	var filteredDevices []*models.Device
	
	for _, device := range allDevices {
		if device.Type == deviceType {
			filteredDevices = append(filteredDevices, device)
		}
	}
	
	return filteredDevices
}

func (d *DeviceService) GetDevicesByLocation(location string) []*models.Device {
	allDevices := d.store.ListDevices()
	var filteredDevices []*models.Device
	
	for _, device := range allDevices {
		if device.Location == location {
			filteredDevices = append(filteredDevices, device)
		}
	}
	
	return filteredDevices
}

func (d *DeviceService) GetDevicesByStatus(status models.DeviceStatus) []*models.Device {
	allDevices := d.store.ListDevices()
	var filteredDevices []*models.Device
	
	for _, device := range allDevices {
		if device.Status == status {
			filteredDevices = append(filteredDevices, device)
		}
	}
	
	return filteredDevices
}

func (d *DeviceService) CalculateEnergyUsage() []models.EnergyUsage {
	devices := d.store.ListDevices()
	var usageData []models.EnergyUsage
	
	for _, device := range devices {
		var usage float64
		
		switch device.Type {
		case models.DeviceTypeLight:
			if power, ok := device.Properties["power"].(bool); ok && power {
				if brightness, ok := device.Properties["brightness"].(int); ok {
					usage = float64(brightness) * 0.001
				} else {
					usage = 0.05
				}
			}
			
		case models.DeviceTypeThermostat:
			if heating, ok := device.Properties["heating"].(bool); ok && heating {
				usage = 2.5
			}
			if cooling, ok := device.Properties["cooling"].(bool); ok && cooling {
				usage = 3.0
			}
			if usage == 0 {
				usage = 0.1
			}
			
		case models.DeviceTypeCamera:
			if recording, ok := device.Properties["recording"].(bool); ok && recording {
				usage = 0.8
			} else {
				usage = 0.2
			}
			
		case models.DeviceTypeSensor:
			usage = 0.01
			
		case models.DeviceTypeLock:
			usage = 0.005
			
		default:
			usage = 0.1
		}
		
		if device.Status == models.DeviceStatusOnline {
			usageData = append(usageData, models.EnergyUsage{
				DeviceID:   device.ID,
				DeviceName: device.Name,
				Usage:      usage,
				Cost:       usage * 0.12,
				Timestamp:  time.Now(),
			})
		}
	}
	
	return usageData
}