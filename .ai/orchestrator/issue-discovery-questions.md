# Issue Discovery Questions

This file defines the comprehensive questionnaire for the Master Agent's `issue-discovery` command. The questions are designed to gather all necessary information to create a complete GitHub issue that sub-agents can implement without guesswork.

## Discovery Phase Structure

### Phase 1: Task Overview
**Purpose**: Establish the basic scope and nature of the task

1. **What is the main task or feature you want to implement?**
   - Provide a brief, one-sentence description of what you want to accomplish

2. **What type of task is this?**
   - [ ] New feature implementation
   - [ ] Bug fix
   - [ ] Code refactoring
   - [ ] Documentation update
   - [ ] Performance optimization
   - [ ] Testing improvement
   - [ ] Security enhancement
   - [ ] Configuration change
   - [ ] Other (please specify)

3. **What is the primary goal or outcome you want to achieve?**
   - Describe the end result and why this task is important

### Phase 2: Technical Context
**Purpose**: Understand the technical scope and requirements

4. **Which parts of the codebase will this task affect?**
   - [ ] Frontend/UI components
   - [ ] Backend API endpoints
   - [ ] Database schema/models
   - [ ] Configuration files
   - [ ] Documentation
   - [ ] Tests
   - [ ] Build/deployment scripts
   - [ ] Other (please specify)

5. **Are there specific files, functions, or modules you know will need to be modified?**
   - List any specific files or code areas you're aware of

6. **What programming languages, frameworks, or technologies are involved?**
   - List the technical stack components relevant to this task

7. **Are there any external dependencies or integrations involved?**
   - APIs, libraries, services, or third-party components

### Phase 3: Functional Requirements
**Purpose**: Define what exactly needs to be built or changed

8. **What specific functionality needs to be implemented?**
   - Describe the detailed behavior and features required

9. **What are the key user interactions or use cases?**
   - How will users interact with this feature or change?

10. **What are the expected inputs and outputs?**
    - Data formats, parameters, return values, etc.

11. **Are there any specific business rules or logic that must be followed?**
    - Validation rules, calculations, workflows, etc.

### Phase 4: Technical Requirements
**Purpose**: Establish technical constraints and specifications

12. **What are the performance requirements?**
    - Speed, scalability, resource usage expectations

13. **Are there any security considerations?**
    - Authentication, authorization, data protection, etc.

14. **What are the data requirements?**
    - Data structures, storage needs, migration requirements

15. **Are there any compatibility requirements?**
    - Browser support, version compatibility, API versions

### Phase 5: Quality and Testing
**Purpose**: Define quality standards and testing approach

16. **What testing is required?**
    - [ ] Unit tests
    - [ ] Integration tests
    - [ ] End-to-end tests
    - [ ] Performance tests
    - [ ] Security tests
    - [ ] Manual testing
    - [ ] Other (please specify)

17. **What are the acceptance criteria?**
    - List specific, measurable criteria for task completion

18. **Are there any edge cases or error conditions to handle?**
    - Unusual inputs, failure scenarios, boundary conditions

### Phase 6: Implementation Details
**Purpose**: Gather implementation-specific information

19. **Do you have any preferred implementation approach?**
    - Architectural patterns, design preferences, etc.

20. **Are there any existing code patterns or standards to follow?**
    - Coding conventions, architectural patterns used in the project

21. **Are there any constraints or limitations to consider?**
    - Technical limitations, resource constraints, time constraints

22. **Should this work integrate with any existing features?**
    - How does this connect to current functionality?

### Phase 7: Dependencies and Sequencing
**Purpose**: Understand task dependencies and prioritization

23. **Does this task depend on any other tasks or issues?**
    - List prerequisite work that must be completed first

24. **Will other tasks depend on this one?**
    - Identify work that will be blocked until this is done

25. **What is the priority level of this task?**
    - [ ] Critical (blocks other work)
    - [ ] High (important for upcoming milestone)
    - [ ] Medium (normal priority)
    - [ ] Low (nice to have)

26. **Is there a specific deadline or timeline?**
    - When does this need to be completed?

### Phase 8: Documentation and Communication
**Purpose**: Ensure proper documentation and stakeholder communication

27. **What documentation needs to be created or updated?**
    - README updates, API docs, user guides, etc.

28. **Who should be notified when this task is completed?**
    - Stakeholders, team members, users

29. **Are there any special deployment or rollout considerations?**
    - Staging requirements, gradual rollout, feature flags

### Phase 9: Definition of Done
**Purpose**: Establish clear completion criteria

30. **How will you know this task is completely finished?**
    - Specific, measurable completion criteria

31. **What should the final deliverable include?**
    - Code, tests, documentation, configuration, etc.

32. **Are there any post-completion steps required?**
    - Monitoring, user communication, follow-up tasks

## Question Flow Logic

### Conditional Questions
- **If bug fix selected**: Add questions about reproduction steps, affected versions, error messages
- **If new feature**: Add questions about user stories, mockups, API specifications
- **If refactoring**: Add questions about current problems, performance goals, backward compatibility
- **If documentation**: Add questions about target audience, format preferences, examples needed

### Follow-up Prompts
- **For vague answers**: "Can you provide more specific details about..."
- **For technical terms**: "Please explain what you mean by..."
- **For missing information**: "You mentioned X, but could you also clarify..."

### Validation Checks
- Ensure all critical information is provided
- Check for consistency between answers
- Identify potential conflicts or ambiguities
- Verify technical feasibility

## Output Requirements

After completing the questionnaire, the Master Agent should have sufficient information to:

1. **Create a comprehensive GitHub issue** with:
   - Clear title and description
   - Detailed requirements and acceptance criteria
   - Technical specifications
   - Implementation guidelines
   - Testing requirements

2. **Assign appropriate labels** such as:
   - Task type (feature, bug, refactor, etc.)
   - Priority level
   - Affected components
   - Estimated effort

3. **Add to project board** in the appropriate column (usually "To Do" or "Backlog")

4. **Include all necessary context** for a sub-agent to:
   - Understand the requirements completely
   - Implement the solution without additional questions
   - Know when the task is complete
   - Follow proper testing and documentation procedures

## Master Agent Instructions

When using this questionnaire:

1. **Ask questions in order** but skip irrelevant ones based on task type
2. **Probe for clarity** when answers are vague or incomplete
3. **Validate consistency** between different answers
4. **Summarize understanding** before creating the issue
5. **Confirm final issue** content with the user before submission