# Issue Discovery System Summary

The Master Agent now includes a comprehensive issue discovery system that creates well-defined GitHub issues through an interactive questionnaire process.

## 🎯 System Overview

### Purpose
Transform user ideas into comprehensive GitHub issues that sub-agents can implement without guesswork.

### Command
`issue-discovery` - Available as both a Master Agent command and standalone slash command

### Process Flow
1. **Interactive Questionnaire** - Guided question session
2. **Comprehensive Issue Creation** - Structured GitHub issue
3. **Project Board Integration** - Automatic placement in project board
4. **Sub-Agent Ready** - Complete specifications for implementation

## 📋 Discovery Components

### 1. Question Framework (`.ai/orchestrator/issue-discovery-questions.md`)
- **9 structured phases** covering all aspects of task definition
- **32 core questions** with conditional follow-ups
- **Validation checks** to ensure completeness
- **Adaptive flow** based on task type

### 2. Issue Template (`.ai/orchestrator/issue-creation-template.md`)
- **Comprehensive structure** with all necessary sections
- **Clear formatting** for easy sub-agent consumption
- **Label assignment rules** for categorization
- **Project board placement** logic

### 3. Discovery Process (`.claude/commands/issue-discovery.md`)
- **Step-by-step guidance** for conducting sessions
- **Interactive guidelines** for natural conversation
- **Success criteria** for quality assurance
- **Example flow** for consistent execution

## 🔄 Discovery Phases

### Phase 1: Task Overview
- Main task description
- Task type classification
- Primary goal identification

### Phase 2: Technical Context
- Affected codebase areas
- Technology stack involved
- External dependencies

### Phase 3: Functional Requirements
- Specific functionality needed
- User interaction patterns
- Input/output specifications

### Phase 4: Technical Requirements
- Performance expectations
- Security considerations
- Data requirements

### Phase 5: Quality and Testing
- Testing requirements
- Acceptance criteria
- Edge case handling

### Phase 6: Implementation Details
- Preferred approaches
- Code standards
- Integration requirements

### Phase 7: Dependencies and Sequencing
- Task dependencies
- Priority level
- Timeline constraints

### Phase 8: Documentation and Communication
- Documentation needs
- Stakeholder notification
- Deployment considerations

### Phase 9: Definition of Done
- Completion criteria
- Final deliverables
- Post-completion steps

## 🏷️ Label System

### Task Type Labels
- `type:feature` - New functionality
- `type:bug` - Bug fixes
- `type:refactor` - Code improvements
- `type:documentation` - Documentation updates
- `type:performance` - Performance optimizations
- `type:security` - Security enhancements
- `type:test` - Testing improvements

### Priority Labels
- `priority:critical` - Blocking other work
- `priority:high` - Important for upcoming milestone
- `priority:medium` - Normal priority
- `priority:low` - Nice to have

### Component Labels
- `component:frontend` - UI/Frontend changes
- `component:backend` - Backend/API changes
- `component:database` - Database modifications
- `component:config` - Configuration changes
- `component:docs` - Documentation
- `component:tests` - Testing

### Effort Labels
- `effort:small` - 1-2 hours
- `effort:medium` - 3-8 hours
- `effort:large` - 1-2 days
- `effort:xl` - 3+ days

## 📊 Project Board Integration

### Column Placement
- **To Do** - New issues ready for development
- **Backlog** - Issues needing refinement or low priority
- **In Progress** - Issues assigned to agents
- **Review** - Issues awaiting code review
- **Done** - Completed issues

### Automatic Placement Rules
- Critical bugs → "To Do" column
- High priority features → "To Do" column
- Medium/Low priority → "Backlog" column
- Documentation → "To Do" column

## 🎯 Quality Assurance

### Validation Checks
- **Completeness** - All required information provided
- **Consistency** - No contradictory requirements
- **Clarity** - Sub-agent can understand without questions
- **Feasibility** - Technical approach is sound

### Success Metrics
- Issue is comprehensive enough for sub-agent implementation
- All technical requirements are clearly specified
- Acceptance criteria are measurable and specific
- Testing requirements are defined
- Dependencies and constraints are identified

## 🚀 Usage Examples

### Basic Usage
```bash
# Initialize Master Agent
/master-agent

# Start issue discovery
Master Agent: issue-discovery

# Follow interactive questionnaire
# Issue created and added to project board
```

### Direct Command
```bash
# Use standalone command
/issue-discovery

# Complete guided session
# Issue automatically created in project board
```

## 💡 Benefits

### For Users
- **Guided process** - No need to remember all requirements
- **Comprehensive coverage** - Nothing important gets missed
- **Consistent quality** - All issues follow same structure
- **Time savings** - Structured approach is faster than ad-hoc

### For Sub-Agents
- **Clear requirements** - No guesswork needed
- **Complete specifications** - All technical details provided
- **Defined success criteria** - Clear acceptance criteria
- **Testing guidance** - Specific testing requirements

### For Project Management
- **Consistent formatting** - All issues follow same structure
- **Proper categorization** - Automatic labeling and placement
- **Dependency tracking** - Clear prerequisite identification
- **Progress visibility** - Issues properly placed in project board

## 🔧 Technical Implementation

### MCP Tools Used
- `mcp__GitHubProjects__create-issue` - Create the GitHub issue
- `mcp__GitHubProjects__add-item-to-project` - Add to project board
- `mcp__GitHubProjects__update-project-item-field` - Set project board fields

### File Structure
```
.ai/orchestrator/
├── issue-discovery-questions.md    # Question framework
├── issue-creation-template.md      # Issue template
└── issue-discovery-summary.md      # This file

.claude/commands/
└── issue-discovery.md              # Standalone command
```

## 📈 Future Enhancements

- **Template customization** - Different templates for different project types
- **Question personalization** - Adapt questions based on user preferences
- **Integration with estimations** - Automatic effort estimation
- **Dependency analysis** - Automatic dependency detection
- **Similar issue detection** - Identify related existing issues

The issue discovery system transforms the process of creating GitHub issues from ad-hoc descriptions into comprehensive, actionable specifications that enable efficient sub-agent implementation.