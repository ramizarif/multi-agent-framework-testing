# Quick Agent Spawn

Spawn a new sub-agent for a GitHub issue.

Usage: `/agent:spawn <issue-number>`

You are a sub-agent spawned to work on GitHub issue #$ARGUMENTS. Follow the sub-agent template workflow:

1. **Retrieve the assigned GitHub issue details** using mcp__GitHubProjects__get-issue
2. **Create implementation plan** in `.ai/implementation-plans/{branch}/issue-$ARGUMENTS.md`
3. **Analyze conflicts** with other active agents by reading existing coordination files
4. **Request Master Agent approval** via coordination file update
5. **Wait for approval** before beginning implementation
6. **Follow the parallel task runner workflow** once approved

Your issue number: #$ARGUMENTS

Begin by retrieving the issue details and creating your implementation plan. Remember to coordinate with the Master Agent for all approvals.