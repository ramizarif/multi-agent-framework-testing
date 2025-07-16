# Agent Spawn Instructions - Multi-Agent Framework

Instructions for the Master Agent to spawn new sub-agents for GitHub issues.

## Spawn Process

### 1. Issue Analysis

Before spawning a new agent, analyze the GitHub issue:

```bash
# Get current branch for coordination
BRANCH_NAME=$(git branch --show-current)

# Retrieve issue details using GitHub Projects API
# Extract: title, description, labels, assignees, linked PRs
```

### 2. Agent ID Generation

Create unique agent identifier:

```bash
# Generate timestamp-based ID
AGENT_ID=$(date +%s%N | cut -b1-13)

# Format: Agent-{timestamp}
# Example: Agent-1704123456000
```

### 3. Conflict Analysis

Analyze potential conflicts with existing agents:

```bash
# Check coordination file for current branch
COORDINATION_FILE=".ai/agent-coordination/${BRANCH_NAME}.md"

# Check implementation plans directory
PLANS_DIR=".ai/implementation-plans/${BRANCH_NAME}"

# Analyze:
# - Active agents and their files
# - Pending plans and their scope
# - Service overlaps
# - Database migration conflicts
```

### 4. Agent Instantiation

Create new agent from template:

```bash
# Copy template
cp .ai/agent-templates/sub-agent-template.md .ai/agent-templates/agent-${AGENT_ID}.md

# Replace template variables
sed -i "s/{AGENT_ID}/${AGENT_ID}/g" .ai/agent-templates/agent-${AGENT_ID}.md
sed -i "s/{ISSUE_NUMBER}/${ISSUE_NUMBER}/g" .ai/agent-templates/agent-${AGENT_ID}.md
sed -i "s/{ISSUE_TITLE}/${ISSUE_TITLE}/g" .ai/agent-templates/agent-${AGENT_ID}.md
sed -i "s/{BRANCH_NAME}/${BRANCH_NAME}/g" .ai/agent-templates/agent-${AGENT_ID}.md
```

### 5. Coordination File Setup

Initialize coordination file entry:

```bash
# Create coordination file if it doesn't exist
if [ ! -f "$COORDINATION_FILE" ]; then
    echo "# Multi-Agent Coordination - Branch: ${BRANCH_NAME}" > "$COORDINATION_FILE"
    echo "" >> "$COORDINATION_FILE"
    echo "## Active Agents" >> "$COORDINATION_FILE"
    echo "" >> "$COORDINATION_FILE"
fi

# Add new agent entry
cat >> "$COORDINATION_FILE" << EOF
### Agent-${AGENT_ID} - Issue #${ISSUE_NUMBER}: ${ISSUE_TITLE}
**Status**: 📋 Planning Phase
**Files being modified**: TBD
**Conflicts**: Analysis pending
**Master Agent approval**: ⏳ Pending
**Next milestone**: Implementation Plan Complete (ETA: 15 min)
**Dependencies**: None
**Created**: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
**Last updated**: $(date -u +"%Y-%m-%dT%H:%M:%SZ")

EOF
```

### 6. Implementation Plan Directory

Ensure implementation plan directory exists:

```bash
# Create directory for branch plans
mkdir -p ".ai/implementation-plans/${BRANCH_NAME}"

# Create slug for plan filename
ISSUE_SLUG=$(echo "$ISSUE_TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-\|-$//g')
PLAN_FILE="${ISSUE_SLUG}-${ISSUE_NUMBER}.md"
```

### 7. Spawn Command

Launch new agent with Task tool:

```bash
# Create agent prompt
AGENT_PROMPT="You are Agent-${AGENT_ID} assigned to GitHub issue #${ISSUE_NUMBER}: '${ISSUE_TITLE}' on branch '${BRANCH_NAME}'.

Follow the sub-agent template at .ai/agent-templates/agent-${AGENT_ID}.md

Your coordination file is: ${COORDINATION_FILE}
Your implementation plan should be: .ai/implementation-plans/${BRANCH_NAME}/${PLAN_FILE}

Begin by:
1. Reading your agent template
2. Analyzing the GitHub issue
3. Creating your implementation plan
4. Updating the coordination file
5. Requesting Master Agent approval

Remember: All approvals must go through the Master Agent. Do not proceed without approval."

# Launch agent
Task agent spawn: "$AGENT_PROMPT"
```

## Conflict-Aware Spawning

### Pre-Spawn Conflict Detection

Before spawning, check for conflicts:

```bash
# Check active agents
ACTIVE_AGENTS=$(grep -c "^### Agent-" "$COORDINATION_FILE" 2>/dev/null || echo "0")

# Check for file conflicts
if [ -f "$COORDINATION_FILE" ]; then
    # Extract files being modified by active agents
    ACTIVE_FILES=$(grep "**Files being modified**:" "$COORDINATION_FILE" | cut -d':' -f2 | tr ',' '\n' | sort -u)
fi

# Check for service conflicts
# Analyze issue labels and description for service indicators
SERVICE_INDICATORS=("auth" "database" "api" "frontend" "backend" "migration")
```

### Conflict Resolution Strategies

#### High Conflict Issues
```bash
# If conflict detected, provide specific spawn instructions
if [ "$CONFLICT_DETECTED" = "true" ]; then
    AGENT_PROMPT="$AGENT_PROMPT

⚠️  CONFLICT ALERT: This issue has potential conflicts with:
- Agent-${CONFLICTING_AGENT} (files: ${CONFLICTING_FILES})
- Recommended approach: ${RESOLUTION_STRATEGY}
- Timing: ${TIMING_INSTRUCTIONS}

You must coordinate with Master Agent before proceeding with any implementation."
fi
```

#### Sequential Execution
```bash
# If sequential execution required
if [ "$SEQUENTIAL_REQUIRED" = "true" ]; then
    AGENT_PROMPT="$AGENT_PROMPT

🔄 SEQUENTIAL EXECUTION REQUIRED:
- Wait for Agent-${DEPENDENCY_AGENT} to complete Phase ${DEPENDENCY_PHASE}
- Monitor coordination file for status updates
- Do not begin implementation until dependency is resolved"
fi
```

## Spawn Response Format

After spawning an agent, update user with:

```markdown
Analyzing issue #${ISSUE_NUMBER}: "${ISSUE_TITLE}"
├─ Conflict analysis: ${CONFLICT_STATUS}
├─ Dependencies: ${DEPENDENCY_STATUS}
├─ Agent-${AGENT_ID} spawned successfully
├─ Coordination file: ${COORDINATION_FILE}
├─ Implementation plan: .ai/implementation-plans/${BRANCH_NAME}/${PLAN_FILE}
└─ Status: Agent ready for planning phase
```

## Error Handling

### Common Spawn Failures

1. **GitHub API Errors**
   ```bash
   # Retry with exponential backoff
   # Fallback to manual issue description
   ```

2. **File System Errors**
   ```bash
   # Check permissions
   # Create missing directories
   # Validate file paths
   ```

3. **Conflict Detection Failures**
   ```bash
   # Default to high-conflict mode
   # Require manual Master Agent review
   # Provide conservative instructions
   ```

### Recovery Procedures

1. **Failed Agent Spawn**
   - Clean up partial files
   - Remove coordination file entry
   - Report error to user
   - Suggest manual intervention

2. **Coordination File Corruption**
   - Backup existing file
   - Reconstruct from git history
   - Validate agent entries
   - Resume normal operation

## Success Validation

After spawning, verify:

- [ ] Agent template created and customized
- [ ] Coordination file updated
- [ ] Implementation plan directory ready
- [ ] No file system errors
- [ ] Agent received proper instructions
- [ ] Conflict analysis completed
- [ ] User notified of spawn status

## Template Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `{AGENT_ID}` | Unique timestamp ID | `1704123456000` |
| `{ISSUE_NUMBER}` | GitHub issue number | `123` |
| `{ISSUE_TITLE}` | GitHub issue title | `Add user authentication` |
| `{BRANCH_NAME}` | Current git branch | `feature/user-auth` |
| `{ISSUE_TITLE_SLUG}` | URL-friendly title | `add-user-authentication` |
| `{CONFLICT_STATUS}` | Conflict analysis result | `No conflicts detected` |
| `{DEPENDENCY_STATUS}` | Dependencies identified | `None` |
| `{COORDINATION_FILE}` | Path to coordination file | `.ai/agent-coordination/feature-user-auth.md` |

This spawn system ensures each new agent is properly configured, conflict-aware, and ready for coordinated development work.