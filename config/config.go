package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port                 int    `json:"port"`
	AuthToken           string `json:"auth_token"`
	LogLevel            string `json:"log_level"`
	WeatherUpdateInterval int   `json:"weather_update_interval"`
	EnergyUpdateInterval  int   `json:"energy_update_interval"`
	SecurityTimeout      int    `json:"security_timeout"`
	RateLimitRPS         int    `json:"rate_limit_rps"`
	MaxDevices           int    `json:"max_devices"`
	EnableDebugMode      bool   `json:"enable_debug_mode"`
}

func Load() *Config {
	cfg := &Config{
		Port:                 8080,
		AuthToken:           "smarthome-secret-token",
		LogLevel:            "info",
		WeatherUpdateInterval: 30,
		EnergyUpdateInterval:  10,
		SecurityTimeout:      300,
		RateLimitRPS:         100,
		MaxDevices:           50,
		EnableDebugMode:      false,
	}
	
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		if err := loadFromFile(cfg, configFile); err != nil {
			log.Printf("Failed to load config from file: %v", err)
		}
	}
	
	loadFromEnv(cfg)
	
	return cfg
}

func loadFromFile(cfg *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, cfg)
}

func loadFromEnv(cfg *Config) {
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}
	
	if token := os.Getenv("AUTH_TOKEN"); token != "" {
		cfg.AuthToken = token
	}
	
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}
	
	if interval := os.Getenv("WEATHER_UPDATE_INTERVAL"); interval != "" {
		if i, err := strconv.Atoi(interval); err == nil {
			cfg.WeatherUpdateInterval = i
		}
	}
	
	if interval := os.Getenv("ENERGY_UPDATE_INTERVAL"); interval != "" {
		if i, err := strconv.Atoi(interval); err == nil {
			cfg.EnergyUpdateInterval = i
		}
	}
	
	if timeout := os.Getenv("SECURITY_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.SecurityTimeout = t
		}
	}
	
	if rps := os.Getenv("RATE_LIMIT_RPS"); rps != "" {
		if r, err := strconv.Atoi(rps); err == nil {
			cfg.RateLimitRPS = r
		}
	}
	
	if maxDevices := os.Getenv("MAX_DEVICES"); maxDevices != "" {
		if m, err := strconv.Atoi(maxDevices); err == nil {
			cfg.MaxDevices = m
		}
	}
	
	if debug := os.Getenv("ENABLE_DEBUG_MODE"); debug == "true" {
		cfg.EnableDebugMode = true
	}
}

func (c *Config) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}