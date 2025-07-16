# User Interface - Multi-Agent Framework Commands

This document defines the user interface for interacting with the multi-agent framework through the Master Agent.

## Core Commands

### Basic Status and Control

#### `status`
**Description**: Show all agent activity and progress
**Usage**: `status`
**Returns**: Complete overview of all active agents, their current phases, progress, and any conflicts

#### `assign <issue-number>`
**Description**: Assign a GitHub issue to a new agent
**Usage**: `assign 123`
**Process**: 
- Retrieves GitHub issue details
- Analyzes potential conflicts with existing agents
- Spawns new agent with unique ID
- Initializes coordination tracking

#### `agents`
**Description**: List all active agents with detailed status
**Usage**: `agents`
**Returns**: Detailed list of all agents including IDs, assigned issues, current status, and files being worked on

### Approval Management

#### `approve-plans`
**Description**: Review and approve pending implementation plans
**Usage**: `approve-plans`
**Process**:
- Shows all pending plans awaiting approval
- Displays conflict analysis for each plan
- Provides recommendation for approval timing
- Allows selective approval

#### `approve-all`
**Description**: Batch approve all pending plans with automatic coordination
**Usage**: `approve-all`
**Process**:
- Approves all pending plans
- Automatically sequences conflicting work
- Updates coordination files
- Notifies agents of execution timing

#### `approve <issue-numbers>`
**Description**: Selectively approve specific implementation plans
**Usage**: `approve 116,123`
**Format**: Comma-separated list of issue numbers

#### `phase-approvals`
**Description**: Review agents waiting for phase progression approval
**Usage**: `phase-approvals`
**Returns**: List of agents waiting to progress to next phase with completion status

#### `approve-phases <agent-ids>`
**Description**: Approve phase progression for specific agents
**Usage**: `approve-phases 1704123456000,1704123789000`
**Format**: Comma-separated list of agent IDs

### Conflict Management

#### `conflicts`
**Description**: Show current and predicted conflicts with resolution suggestions
**Usage**: `conflicts`
**Returns**: 
- Active conflicts with resolution status
- Predicted conflicts with prevention strategies
- Recommended actions for conflict resolution

#### `branch-health`
**Description**: Check coordination system health and integrity
**Usage**: `branch-health`
**Returns**: System health status, file integrity, and coordination system status

### Agent Management

#### `terminate <agent-id>`
**Description**: Stop a problematic agent with cleanup
**Usage**: `terminate 1704123456000`
**Process**:
- Gracefully stops agent execution
- Updates coordination files
- Releases file locks
- Notifies dependent agents

## Command Flow Examples

### Starting New Work

```bash
# Check current status
status

# Assign new issue to agent
assign 123

# Agent creates plan and requests approval
approve-plans

# Approve the plan
approve 123

# Monitor progress
status
```

### Managing Multiple Agents

```bash
# Check all active agents
agents

# Review conflicts
conflicts

# Approve phase progressions
phase-approvals
approve-phases 1704123456000,1704123789000

# Batch approve all pending plans
approve-all
```

### Conflict Resolution

```bash
# Identify conflicts
conflicts

# Review specific plans
approve-plans

# Selective approval to sequence work
approve 116
# Wait for Agent-116 to complete Phase 3, then:
approve 123
```

### System Monitoring

```bash
# Overall system health
branch-health

# Current activity
status

# Detailed agent information
agents
```

## Response Formats

### Status Command Response
```
Multi-Agent Status Report
========================

Active Agents: 3
Branch: feature/user-auth

Agent-1704123456000 - Issue #113: Add user authentication
├─ Status: 🔄 Implementation Phase 2
├─ Files: auth.go, handlers.go, models.go
├─ Progress: 60% complete
├─ Next milestone: Phase 3 - Testing (ETA: 30 min)
└─ Dependencies: None

Pending Approvals: 1
Conflicts Detected: 1 (managed)
```

### Approval Review Response
```
Pending Implementation Plans Review
==================================

Plan 1: Agent-1704124000000 - Issue #116: Update user schema
├─ Plan: .ai/implementation-plans/feature-user-auth/update-user-schema-116.md
├─ Conflicts: models.go overlap with Agent-1704123456000
├─ Recommendation: Approve after Agent-1704123456000 Phase 3
└─ Risk level: Medium

Actions:
[a] Approve all with coordination
[s] Selective approval
[r] Reject and request modifications
```

### Conflict Report Response
```
Conflict Analysis Report
=======================

Active Conflicts: 1
Predicted Conflicts: 2

🔥 Active Conflict: File Access
├─ Agents: Agent-1704123456000, Agent-1704124000000
├─ File: models/user.go
├─ Resolution: Agent-1704124000000 waiting for Agent-1704123456000 Phase 3
└─ Status: Managed - no action needed

⚠️  Predicted Conflict: Database Migration
├─ Agents: Agent-1704124000000, Agent-1704124567000
├─ Issue: Both plan to modify user table schema
├─ Resolution: Coordinate migration timestamps
└─ Action: Review plans before approval
```

## Error Handling

### Invalid Commands
- Unknown commands receive help text
- Invalid parameters show usage examples
- Missing required arguments prompt for input

### System Errors
- Agent timeouts reported with recovery options
- File access errors with resolution suggestions
- GitHub API errors with fallback procedures

## Best Practices

### Regular Monitoring
- Run `status` frequently to monitor progress
- Check `conflicts` before approving new plans
- Use `branch-health` to ensure system integrity

### Approval Strategy
- Review `approve-plans` before batch operations
- Use selective approval for complex conflict scenarios
- Monitor `phase-approvals` for stuck agents

### Conflict Prevention
- Always check `conflicts` before assigning new issues
- Use `agents` to understand current workload
- Coordinate timing through selective approvals

## Integration Notes

### GitHub Integration
- Issue numbers correspond to GitHub issues
- Agent progress updates GitHub issue status
- Commits are automatically linked to issues

### File System Integration
- Coordination files track all activity
- Implementation plans provide detailed blueprints
- System maintains audit trail of all decisions

This interface provides complete control over the multi-agent framework while maintaining simplicity and preventing conflicts through intelligent coordination.