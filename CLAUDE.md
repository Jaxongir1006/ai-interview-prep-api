# CLAUDE

## Skills

For any backend work (implementation, review, modification, planning), load first:

1. **`backend-guidelines`** — project architecture, conventions, and development rules

Do not improvise on patterns, naming, or structure — follow what the skill prescribes.

Before implementing a specific module's features, also read:

- `docs/specs/modules/{module}/overview.md` — module overview
- `docs/specs/modules/{module}/ERD.md` — module ERD
- `docs/specs/flows/` — relevant business flows

Don't read use case documents until you're actually working on them — token and context efficiency.

## Agents

Use specialized agents (`.claude/agents/`) instead of generic subagents. Each agent has focused instructions, restricted tools, and auto-loaded skills.

| Agent              | Role                            | Tools                             |
| ------------------ | ------------------------------- | --------------------------------- |
| **`architect`**    | Plans implementation            | Read, Glob, Grep *(read-only)*    |
| **`system-analyst`** | Writes UC docs, specs, ERDs   | Read, Write, Edit, Glob, Grep     |
| **`go-coder`**     | Implements code                 | Read, Write, Edit, Bash, Glob, Grep |
| **`go-tester`**    | Writes unit & system tests      | Read, Write, Edit, Bash, Glob, Grep |
| **`reviewer`**     | Reviews code against guidelines | Read, Glob, Grep, Bash *(no edit)* |

### How to delegate

```
Task tool → subagent_type: "architect" / "go-coder" / "go-tester" / "system-analyst" / "reviewer"
```

### Pipeline

For features and use cases, follow this pipeline:

1. **architect** — plan the implementation (read-only)
2. **system-analyst** — write/update UC docs (if needed)
3. **go-coder** + **go-tester** — implement code and tests (parallelizable)
4. **reviewer** — review against guidelines, run lint/tests (read-only)
5. Fix issues from reviewer, re-review if needed

## Rules

Always keep docs, code, and tests in sync.

## Verification

```bash
make lint          # Lint
make test          # Unit tests (pkg/ layer)
make test-system   # System tests (use case coverage)
```

- If lint fails, run `make fmt` first, then `make lint` again
- NEVER deliver code that hasn't passed all three

## Workflow

For large tasks (features, modules, multi-use-case work), break the work into manageable pieces using `TaskCreate` to plan and track them, and `Task` to parallelize independent pieces via specialized agents.

Use `/generate-module` slash command to scaffold a new module.
