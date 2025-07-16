# Smart Home Automation Hub + Multi-Agent Framework

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

## Multi-Agent Framework

This repository includes a comprehensive single-agent development framework with multi-agent coordination templates. While originally designed for multi-agent coordination, it provides valuable structured workflows, issue discovery, and project management capabilities for individual developers and teams. The framework is integrated with the "multi-agent-framework-testing" GitHub project board for automatic issue management and progress tracking.

### 🚨 **Important Note About Multi-Agent Capabilities**

**Current Limitation**: Claude Code cannot spawn separate Claude sessions. The "multi-agent" framework provides templates and coordination structures for single-agent workflows or team coordination, not true parallel AI agents.

**What It Actually Provides**:
- ✅ Comprehensive issue discovery and creation
- ✅ Structured implementation planning
- ✅ Conflict prevention logic for team coordination
- ✅ Progress tracking templates
- ✅ GitHub project board integration
- ✅ Workflow templates for structured development

### 🚀 Quick Start - Development Framework

#### Create Comprehensive Issues
```bash
/issue-discovery
```

#### Initialize Master Agent (Single Agent Mode)
```bash
/master-agent
```

#### Check System Status
```bash
/agent:status
```

#### Work on Issues Systematically
```bash
# After Master Agent initialization
I want to work on Issue #123. Please guide me through the sub-agent workflow.
```

#### Health Check
```bash
/project:health
```

### Master Agent Commands

#### Core Commands
- `status` - Show current system status and progress
- `issue-discovery` - Interactive process to create comprehensive GitHub issues
- `project-status` - Show GitHub project board status and available issues
- `available-issues` - List issues ready for assignment from project board
- `sync-project` - Synchronize system status with project board

#### Workflow Commands (Single Agent Mode)
- `assign <issue-number>` - Guide user through structured implementation of specific issue
- `approve-plans` - Review implementation plans (for team coordination)
- `conflicts` - Show potential conflicts and resolution strategies
- `branch-health` - Check coordination system health

#### Information Commands
- `help` - Show available commands and usage
- `coordination-files` - List current coordination files and their status
- `recent-activity` - Show recent completions and activities

### Framework Architecture

```
User Interface
    ↓
Master Agent (Central Coordinator)
    ↓
Structured Workflows & Templates
    ↓
GitHub Issues/Projects Integration
```

**Note**: While designed as a multi-agent system, the framework operates as a single-agent system with comprehensive coordination templates that can be used for team coordination.

## Smart Home API

### Running the Application

#### Prerequisites
- Go 1.21 or higher
- Git

#### Installation

```bash
# Clone the repository
git clone https://github.com/ramizarif/multi-agent-framework-testing.git
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

## Multi-Agent Framework Directory Structure

The `.ai` directory contains the complete framework for structured development workflows and team coordination:

```
.ai/
├── orchestrator/                    # Master Agent and coordination system
│   ├── master-agent.md              # Master Agent configuration and commands
│   ├── user-interface.md            # User command interface documentation
│   ├── master-agent-ready.md        # System ready status and overview
│   ├── issue-discovery-questions.md # Comprehensive issue discovery questionnaire
│   ├── issue-creation-template.md   # GitHub issue creation template
│   ├── issue-discovery-summary.md   # Issue discovery system documentation
│   └── project-board-integration.md # GitHub project board integration guide
├── agent-coordination/              # Agent coordination and tracking
│   └── {branch-name}.md             # Branch-specific coordination files
├── implementation-plans/            # Detailed implementation plans
│   └── {branch-name}/
│       └── {issue-title}-{number}.md # Agent implementation plans
└── agent-templates/                 # Templates and workflows
    ├── sub-agent-template.md        # Sub-agent workflow template
    ├── spawn-instructions.md        # Agent spawn procedures (legacy)
    ├── coordination-file-template.md # Coordination file structure
    ├── implementation-plan-template.md # Implementation plan structure
    ├── conflict-prevention-logic.md  # Conflict detection algorithms
    └── github-integration.md        # GitHub integration documentation

.claude/
└── commands/                        # Claude Code slash commands
    ├── master-agent.md              # Master Agent slash command
    ├── issue-discovery.md           # Issue discovery slash command
    ├── agent/
    │   ├── spawn.md                 # Agent spawn command (legacy)
    │   └── status.md                # Agent status command
    └── project/
        ├── health.md                # System health check
        ├── reset.md                 # System reset command
        └── board.md                 # Project board integration
```

### 📁 **Directory Breakdown**

#### `.ai/orchestrator/` - Core Framework
- **`master-agent.md`** - Central coordinator configuration with all commands
- **`issue-discovery-questions.md`** - 32-question framework for comprehensive issue creation
- **`issue-creation-template.md`** - GitHub issue template with labels and project board integration
- **`project-board-integration.md`** - GitHub project board integration documentation

#### `.ai/agent-coordination/` - Progress Tracking
- **`{branch-name}.md`** - Real-time coordination files tracking agent status, conflicts, and progress
- **Purpose**: Track multiple developers or structured workflow progress
- **Format**: Standardized status updates, conflict detection, and dependency management

#### `.ai/implementation-plans/` - Detailed Planning
- **`{branch-name}/{issue-title}-{number}.md`** - Comprehensive implementation plans
- **Includes**: Technical approach, conflict analysis, phase breakdown, testing requirements
- **Purpose**: Ensure no implementation details are missed

#### `.ai/agent-templates/` - Workflow Templates
- **`sub-agent-template.md`** - Structured workflow for implementing issues
- **`coordination-file-template.md`** - Template for tracking progress and conflicts
- **`implementation-plan-template.md`** - Template for creating detailed implementation plans
- **`conflict-prevention-logic.md`** - Algorithms for detecting and preventing conflicts

#### `.claude/commands/` - Claude Code Integration
- **`master-agent.md`** - Main framework activation via `/master-agent`
- **`issue-discovery.md`** - Issue creation workflow via `/issue-discovery`
- **`project/board.md`** - Project board management utilities

### 🎯 **How to Use the Framework**

#### **For Individual Developers**
1. **Issue Discovery**: Use `/issue-discovery` to create comprehensive GitHub issues
2. **Structured Implementation**: Follow sub-agent template for systematic development
3. **Progress Tracking**: Use coordination files to track implementation progress
4. **Conflict Prevention**: Use conflict detection logic to avoid implementation issues

#### **For Teams**
1. **Master Agent Coordination**: One person acts as coordinator using Master Agent
2. **Issue Assignment**: Use comprehensive issues created through discovery process
3. **Progress Monitoring**: Use coordination files to track team member progress
4. **Conflict Prevention**: Prevent developers from working on conflicting areas

#### **For Learning**
1. **Study Templates**: Examine the structured approach to issue creation and implementation
2. **Workflow Patterns**: Learn systematic development workflow patterns
3. **Conflict Detection**: Understand how to identify and prevent development conflicts
4. **Documentation Standards**: See examples of comprehensive technical documentation

### 🔧 **Framework Components Deep Dive**

#### **Issue Discovery System**
The framework includes a comprehensive 32-question system for creating detailed GitHub issues:

```markdown
# Example Issue Discovery Flow
1. Task Overview (What, Why, Type)
2. Technical Context (Files, Technologies, Dependencies)
3. Functional Requirements (Behavior, Interactions, I/O)
4. Technical Requirements (Performance, Security, Data)
5. Quality & Testing (Tests, Criteria, Edge Cases)
6. Implementation (Approach, Standards, Integration)
7. Dependencies (Prerequisites, Priority, Timeline)
8. Documentation (Updates, Communication, Deployment)
9. Definition of Done (Completion, Deliverables, Follow-up)
```

#### **Coordination System**
Real-time tracking of development progress with conflict prevention:

```markdown
# Example Coordination File
## Active Agents
### Agent-1704123456000 - Issue #123: Add Authentication
**Status**: 🔄 Implementation Phase 2
**Files**: auth.go, handlers.go, models.go
**Conflicts**: None detected
**Progress**: 60% complete
**Next**: Phase 3 - Testing (ETA: 30 min)
```

#### **Implementation Planning**
Detailed technical specifications with conflict analysis:

```markdown
# Example Implementation Plan
## Conflict Analysis
- File conflicts with Agent-X on models.go
- Service conflicts on authentication endpoints
- Resolution: Coordinate timing with Agent-X Phase 3

## Technical Approach
- JWT-based authentication
- Database schema updates
- API endpoint modifications
- Frontend login components
```

### 📚 **Learning Resources**

#### **Key Files to Study**
1. **`.ai/orchestrator/issue-discovery-questions.md`** - Learn comprehensive requirement gathering
2. **`.ai/agent-templates/sub-agent-template.md`** - Study structured development workflow
3. **`.ai/agent-templates/conflict-prevention-logic.md`** - Understand conflict detection algorithms
4. **`.ai/orchestrator/issue-creation-template.md`** - See comprehensive issue documentation

#### **Workflow Examples**
1. **Issue Creation**: `/issue-discovery` → Comprehensive GitHub issue
2. **Implementation**: Follow sub-agent template workflow
3. **Coordination**: Use coordination files for progress tracking
4. **Completion**: Structured testing and documentation

### 🎓 **Framework Benefits for Learning**

#### **Development Process**
- **Systematic Approach**: Structured workflow for any development task
- **Comprehensive Planning**: Detailed analysis before implementation
- **Conflict Prevention**: Proactive identification of potential issues
- **Progress Tracking**: Real-time monitoring of development progress

#### **Documentation Standards**
- **Issue Documentation**: Comprehensive GitHub issue templates
- **Technical Specifications**: Detailed implementation planning
- **Progress Reporting**: Structured status updates
- **Conflict Analysis**: Systematic conflict detection and resolution

#### **Team Coordination**
- **Role Definition**: Clear coordinator and implementer roles
- **Communication Protocols**: Structured approval and reporting processes
- **Resource Management**: Conflict prevention and resource allocation
- **Quality Assurance**: Built-in testing and validation requirements

This framework provides a complete system for structured development whether working alone or coordinating with a team, with comprehensive templates and workflows that can be adapted to any project.

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

This project is provided as-is for testing multi-agent frameworks and smart home automation development.