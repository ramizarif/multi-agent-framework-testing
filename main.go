package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"multi-agent-framework-testing/config"
	"multi-agent-framework-testing/handlers"
	"multi-agent-framework-testing/services"
	"multi-agent-framework-testing/storage"
	"multi-agent-framework-testing/workers"
)

func main() {
	cfg := config.Load()
	
	store := storage.NewMemoryStore()
	
	deviceService := services.NewDeviceService(store)
	weatherService := services.NewWeatherService()
	
	scheduler := workers.NewScheduler(store, deviceService, weatherService)
	
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	
	handler := handlers.NewHandler(store, deviceService, weatherService, scheduler, &upgrader)
	
	router := mux.NewRouter()
	
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/devices", handler.ListDevices).Methods("GET")
	router.HandleFunc("/devices", handler.AddDevice).Methods("POST")
	router.HandleFunc("/devices/{id}", handler.UpdateDevice).Methods("PUT")
	router.HandleFunc("/weather", handler.GetWeather).Methods("GET")
	router.HandleFunc("/energy/usage", handler.GetEnergyUsage).Methods("GET")
	router.HandleFunc("/security/arm", handler.ArmSecurity).Methods("POST")
	router.HandleFunc("/security/disarm", handler.DisarmSecurity).Methods("POST")
	router.HandleFunc("/analytics/summary", handler.GetAnalytics).Methods("GET")
	router.HandleFunc("/schedule/task", handler.CreateScheduledTask).Methods("POST")
	router.HandleFunc("/debug/state", handler.DebugState).Methods("GET")
	router.HandleFunc("/debug/reset", handler.ResetSystem).Methods("POST")
	router.HandleFunc("/debug/trigger/{scenario}", handler.TriggerScenario).Methods("POST")
	router.HandleFunc("/ws", handler.WebSocketHandler).Methods("GET")
	
	router.Use(handler.LoggingMiddleware)
	router.Use(handler.AuthMiddleware)
	router.Use(handler.RateLimitMiddleware)
	
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}
	
	go scheduler.Start()
	
	go func() {
		log.Printf("Smart Home Hub starting on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	scheduler.Stop()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	
	log.Println("Server exited")
}