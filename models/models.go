package models

import (
	"time"
)

type DeviceType string

const (
	DeviceTypeLight      DeviceType = "light"
	DeviceTypeThermostat DeviceType = "thermostat"
	DeviceTypeCamera     DeviceType = "camera"
	DeviceTypeSensor     DeviceType = "sensor"
	DeviceTypeLock       DeviceType = "lock"
	DeviceTypeAlarm      DeviceType = "alarm"
)

type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
	DeviceStatusError   DeviceStatus = "error"
)

type Device struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        DeviceType             `json:"type"`
	Status      DeviceStatus           `json:"status"`
	Properties  map[string]interface{} `json:"properties"`
	Location    string                 `json:"location"`
	LastUpdated time.Time              `json:"last_updated"`
	CreatedAt   time.Time              `json:"created_at"`
}

type WeatherData struct {
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	Condition   string    `json:"condition"`
	WindSpeed   float64   `json:"wind_speed"`
	WindDir     string    `json:"wind_direction"`
	Timestamp   time.Time `json:"timestamp"`
}

type EnergyUsage struct {
	DeviceID    string    `json:"device_id"`
	DeviceName  string    `json:"device_name"`
	Usage       float64   `json:"usage_kwh"`
	Cost        float64   `json:"cost_usd"`
	Timestamp   time.Time `json:"timestamp"`
}

type SecurityState string

const (
	SecurityStateDisarmed SecurityState = "disarmed"
	SecurityStateArmed    SecurityState = "armed"
	SecurityStateTriggered SecurityState = "triggered"
)

type SecuritySystem struct {
	State          SecurityState `json:"state"`
	LastArmed      time.Time     `json:"last_armed"`
	LastTriggered  time.Time     `json:"last_triggered"`
	ActiveSensors  []string      `json:"active_sensors"`
	TriggeredBy    string        `json:"triggered_by"`
}

type ScheduledTask struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	DeviceID    string                 `json:"device_id"`
	Action      string                 `json:"action"`
	Parameters  map[string]interface{} `json:"parameters"`
	Schedule    string                 `json:"schedule"`
	Enabled     bool                   `json:"enabled"`
	NextRun     time.Time              `json:"next_run"`
	LastRun     time.Time              `json:"last_run"`
	CreatedAt   time.Time              `json:"created_at"`
}

type AnalyticsData struct {
	TotalDevices     int                        `json:"total_devices"`
	OnlineDevices    int                        `json:"online_devices"`
	OfflineDevices   int                        `json:"offline_devices"`
	TotalEnergyUsage float64                    `json:"total_energy_usage_kwh"`
	TotalEnergyCost  float64                    `json:"total_energy_cost_usd"`
	DeviceUsage      map[string]float64         `json:"device_usage"`
	SecurityEvents   int                        `json:"security_events"`
	ScheduledTasks   int                        `json:"scheduled_tasks"`
	WeatherSummary   WeatherData                `json:"weather_summary"`
	GeneratedAt      time.Time                  `json:"generated_at"`
}

type SystemEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Source      string                 `json:"source"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   time.Time              `json:"timestamp"`
	Severity    string                 `json:"severity"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type TestScenario struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Actions     []TestAction           `json:"actions"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type TestAction struct {
	Type       string                 `json:"type"`
	Target     string                 `json:"target"`
	Parameters map[string]interface{} `json:"parameters"`
	Delay      int                    `json:"delay_seconds"`
}

type SystemState struct {
	Devices      []Device         `json:"devices"`
	Weather      WeatherData      `json:"weather"`
	Security     SecuritySystem   `json:"security"`
	Tasks        []ScheduledTask  `json:"tasks"`
	EnergyUsage  []EnergyUsage    `json:"energy_usage"`
	SystemEvents []SystemEvent    `json:"system_events"`
	Uptime       time.Duration    `json:"uptime"`
	Timestamp    time.Time        `json:"timestamp"`
}