# Smart Home Automation Hub

A comprehensive Go-based microservice for testing multi-agent frameworks. This application simulates a smart home automation system with interconnected components, REST APIs, background workers, and real-time updates.

## Features

- **Device Management**: Control various smart devices (lights, thermostats, cameras, sensors, locks)
- **Weather Simulation**: Dynamic weather data with impact on device behavior
- **Energy Monitoring**: Track and analyze energy consumption across devices
- **Security System**: Armed/disarmed states with sensor integration
- **Task Scheduling**: Automate device actions with scheduled tasks
- **Real-time Updates**: WebSocket support for live device state changes
- **Analytics Engine**: Process usage patterns and generate reports
- **Testing Scenarios**: Built-in test scenarios for multi-agent testing

## Architecture

The application consists of 9 main files:
- `main.go` - Entry point and server setup
- `config/config.go` - Configuration management
- `models/models.go` - Data structures and types
- `handlers/handlers.go` - HTTP handlers and routing
- `services/device_service.go` - Device management logic
- `services/weather_service.go` - Weather simulation
- `workers/scheduler.go` - Background task processing
- `utils/utils.go` - Helper functions and utilities
- `storage/memory_store.go` - In-memory data persistence

## Running the Application

### Prerequisites
- Go 1.21 or higher
- Git

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd multi-agent-framework-testing

# Download dependencies
go mod download

# Run the application
go run main.go
```

### Configuration

The application can be configured via:
1. Environment variables
2. Configuration file (`config.json`)

#### Environment Variables
- `PORT` - Server port (default: 8080)
- `AUTH_TOKEN` - Bearer authentication token
- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- `CONFIG_FILE` - Path to configuration file
- `WEATHER_UPDATE_INTERVAL` - Weather update frequency in seconds
- `ENERGY_UPDATE_INTERVAL` - Energy monitoring frequency in seconds
- `SECURITY_TIMEOUT` - Security alarm auto-reset timeout in seconds
- `RATE_LIMIT_RPS` - Rate limit requests per second
- `MAX_DEVICES` - Maximum number of devices allowed
- `ENABLE_DEBUG_MODE` - Enable debug endpoints

## API Endpoints

### Device Management
- `GET /devices` - List all devices
- `POST /devices` - Add a new device
- `PUT /devices/{id}` - Update device state

### Weather
- `GET /weather` - Get current weather and forecast

### Energy
- `GET /energy/usage` - Get energy consumption data

### Security
- `POST /security/arm` - Arm security system
- `POST /security/disarm` - Disarm security system

### Analytics
- `GET /analytics/summary` - Get system analytics

### Scheduling
- `POST /schedule/task` - Create scheduled automation

### Debug & Testing
- `GET /debug/state` - Get complete system state
- `POST /debug/reset` - Reset system to initial state
- `POST /debug/trigger/{scenario}` - Trigger test scenarios

### Health & Monitoring
- `GET /health` - Health check endpoint

### WebSocket
- `GET /ws` - WebSocket connection for real-time updates

## Authentication

All endpoints (except `/health` and `/ws`) require Bearer token authentication:
```
Authorization: Bearer smarthome-secret-token
```

## Test Scenarios

Trigger test scenarios via `/debug/trigger/{scenario}`:

- `weather` - Simulate extreme weather conditions
- `storm` - Trigger storm weather
- `heatwave` - Trigger extreme heat
- `cold_snap` - Trigger freezing temperatures
- `device_failure` - Simulate device going offline
- `power_surge` - Simulate high energy usage
- `morning_routine` - Execute morning automation
- `evening_routine` - Execute evening automation
- `away_mode` - Activate away mode
- `sleep_mode` - Activate sleep mode
- `security_breach` - Trigger security alarm

## Device Types

Supported device types:
- `light` - Smart lights with brightness control
- `thermostat` - Temperature control with heating/cooling
- `camera` - Security cameras with recording
- `sensor` - Motion sensors with security integration
- `lock` - Smart locks with remote control
- `alarm` - Security alarm system

## WebSocket Events

Real-time events broadcast via WebSocket:
- `initial_state` - Complete system state on connection
- `device_added` - New device added
- `device_updated` - Device state changed
- `security_armed` - Security system armed
- `security_disarmed` - Security system disarmed
- `state_update` - Periodic state updates

## Example Usage

### Add a Device
```bash
curl -X POST http://localhost:8080/devices \
  -H "Authorization: Bearer smarthome-secret-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kitchen Light",
    "type": "light",
    "location": "Kitchen",
    "properties": {
      "brightness": 75,
      "power": true
    }
  }'
```

### Create Scheduled Task
```bash
curl -X POST http://localhost:8080/schedule/task \
  -H "Authorization: Bearer smarthome-secret-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Turn off lights at night",
    "device_id": "light_001",
    "action": "turn_off",
    "schedule": "daily",
    "parameters": {}
  }'
```

### WebSocket Connection
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message.type, message.data);
};
```

## Testing Multi-Agent Interactions

This application is designed to test multi-agent frameworks with:
- Complex interdependencies between components
- Concurrent operations requiring synchronization
- Event-driven updates between services
- Data validation and error propagation
- Configurable delays and failure modes
- State management across multiple services

## Development

### Adding New Features

1. Define models in `models/models.go`
2. Add storage methods in `storage/memory_store.go`
3. Implement business logic in services
4. Create HTTP handlers in `handlers/handlers.go`
5. Add routes in `main.go`

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o smart-home-hub
```

## Performance Considerations

- Rate limiting prevents API abuse
- Background workers use goroutines for concurrent processing
- In-memory storage limits data retention (1000 energy records, 500 system events)
- WebSocket broadcasting is throttled to 5-second intervals
- Graceful shutdown ensures clean termination

## License

This project is provided as-is for testing multi-agent frameworks.