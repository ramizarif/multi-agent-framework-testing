# Issue Discovery Process

Launch an interactive issue discovery session to create comprehensive GitHub issues.

## Process Overview

You are now conducting an issue discovery session to create a well-defined GitHub issue. You will guide the user through a structured questionnaire to gather all necessary information for creating a comprehensive issue that sub-agents can implement without guesswork.

## Instructions

### Step 1: Load Questions and Template
1. Read the questionnaire from `.ai/orchestrator/issue-discovery-questions.md`
2. Load the issue template from `.ai/orchestrator/issue-creation-template.md`
3. Prepare to conduct a thorough discovery session

### Step 2: Conduct Discovery Session
Follow the structured phases from the questionnaire:

1. **Task Overview** - Understand the basic scope and nature
2. **Technical Context** - Identify affected components and technologies
3. **Functional Requirements** - Define what needs to be built
4. **Technical Requirements** - Establish constraints and specifications
5. **Quality and Testing** - Define testing approach and acceptance criteria
6. **Implementation Details** - Gather implementation preferences
7. **Dependencies and Sequencing** - Understand task relationships
8. **Documentation and Communication** - Plan documentation needs
9. **Definition of Done** - Establish completion criteria

### Step 3: Interactive Guidelines
- **Ask one question at a time** - Don't overwhelm with multiple questions
- **Wait for complete answers** - Let the user fully respond before proceeding
- **Probe for clarity** - Ask follow-up questions when answers are vague
- **Skip irrelevant questions** - Based on task type, skip questions that don't apply
- **Validate understanding** - Summarize key points and confirm accuracy

### Step 4: Create Issue Structure
Using the responses, create a comprehensive issue following the template structure:
- Clear title with type and priority
- Detailed overview and requirements
- Technical specifications
- Acceptance criteria with checkboxes
- Testing requirements
- Dependencies and constraints
- Documentation needs
- Definition of done

### Step 5: Review and Confirm
- Present the complete issue draft to the user
- Highlight key requirements and acceptance criteria
- Ask for confirmation or adjustments
- Ensure all details are accurate and complete

### Step 6: Create GitHub Issue
- Use `mcp__GitHubProjects__create-issue` to create the issue
- Apply appropriate labels (type, priority, component, effort)
- Use `mcp__GitHubProjects__add-item-to-project` to add to project board
- Place in the appropriate column (To Do/Backlog)

## Success Criteria

The resulting issue should be:
- **Complete** - All necessary information provided
- **Clear** - No ambiguity in requirements
- **Actionable** - Sub-agent can implement without questions
- **Testable** - Clear acceptance criteria and testing requirements
- **Documented** - Includes documentation and communication needs

## Example Flow

```
Master Agent: Let's create a comprehensive issue for your task. I'll guide you through a series of questions to ensure we capture all necessary details.

First, what is the main task or feature you want to implement? Please provide a brief, one-sentence description.

User: [Response]

Master Agent: Thank you. Now, what type of task is this? Is it:
- New feature implementation
- Bug fix
- Code refactoring
- Documentation update
- Performance optimization
- Testing improvement
- Security enhancement
- Configuration change
- Other

[Continue through all relevant questions...]

Master Agent: Based on your responses, I've created a comprehensive issue. Here's the complete issue I'll create:

[Show formatted issue]

Does this accurately capture your requirements? Any adjustments needed?

User: [Confirmation or adjustments]

Master Agent: Perfect! I'm now creating the GitHub issue and adding it to your project board...

Issue created successfully: #[number]
Added to project board in "To Do" column
Labels applied: [list of labels]
Ready for sub-agent assignment!
```

## Important Notes

- **Be thorough but efficient** - Ask all necessary questions without being repetitive
- **Adapt to task type** - Skip questions that don't apply to the specific task
- **Maintain context** - Remember previous answers and reference them
- **Ensure completeness** - Don't create the issue until all critical information is gathered
- **Confirm before creating** - Always show the final issue to the user for approval

Begin the discovery session now by introducing the process and asking the first question.