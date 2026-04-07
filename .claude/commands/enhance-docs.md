# Enhance Documentation

You are enhancing project documentation. The user will describe a concept, rule, or guideline to add.

## Process

1. **Identify the target file.** Determine which file the concept belongs to:
   - `docs/guidelines/` — development rules and coding standards
   - `docs/architecture/` — system architecture and codebase structure
   - `docs/specs/` — module specs, use cases, flows
   - `CLAUDE.md` — agent-facing instructions

2. **Read the entire target file.** Understand the current structure, sections, and flow.

3. **Decide where the concept fits.** Don't just append to the bottom. Think about:
   - Which existing section does it belong to?
   - Does it need a new section? If so, where does that section logically fit among others?
   - Does adding this concept make an existing section too long or unfocused? If so, split it.

4. **Integrate the concept.** Write it in the same style as the surrounding content:
   - Match the tone (concise, direct, no fluff)
   - Match the formatting (bullet points, tables, code blocks — whatever the file uses)
   - If the concept relates to existing points, group them together

5. **Restructure if needed.** After integration, re-read the file and check:
   - Does the section ordering still make sense?
   - Are related concepts grouped together?
   - Is anything redundant now? Remove duplication.
   - Does the file still flow well from top to bottom?

6. **Sync the backend-guidelines skill.** After modifying any file in `docs/guidelines/`, `docs/architecture/codebase.md`, or `docs/specs/api/general.md`, you MUST update the corresponding section in `.claude/skills/backend-guidelines/SKILL.md`:

   **Mapping (all sections are in SKILL.md):**
   | Source docs | SKILL.md section |
   |---|---|
   | `docs/architecture/codebase.md` | Architecture |
   | `docs/guidelines/01_code-style.md` | Code Style |
   | `docs/guidelines/02_controllers.md` | Controllers |
   | `docs/guidelines/03_use-cases.md` | Use Cases |
   | `docs/guidelines/04_pblc.md` | PBLC |
   | `docs/guidelines/05_infra.md` | Infrastructure |
   | `docs/guidelines/06_di-containers.md` | DI Containers |
   | `docs/guidelines/07_tx-management.md` | Transaction Management (UOW) |
   | `docs/guidelines/08_error-handling.md` | Error Handling |
   | `docs/guidelines/09_validation.md` | Validation |
   | `docs/guidelines/10_list-manipulations.md` | List Manipulations |
   | `docs/guidelines/11_testing.md` | Testing |
   | `docs/guidelines/12_observability.md` | Observability |
   | `docs/guidelines/13_db_migrations.md` | DB Migrations |
   | `docs/guidelines/14_documentation.md` | Documentation |
   | `docs/guidelines/15_dev-workflow.md` | Development Workflow |
   | `docs/guidelines/99_common_mistakes.md` | Common Mistakes |
   | `docs/specs/api/general.md` | API Design |

   **Sync rules:**
   - The skill is a condensed, agent-optimized version — not a copy
   - Maintain the concise style (no explanations Claude already knows)
   - Include all project-specific patterns and conventions

## Rules

- **Never lose existing concepts.** Every idea currently in the file must survive the edit.
- **Be concise.** Add what's needed, nothing more. Don't pad with explanations unless the concept is non-obvious.
- **Preserve formatting conventions.** If the file uses `##` for sections and `-` for bullets, keep doing that.
- **Show the user what changed.** After editing, briefly summarize what was added and any restructuring done.
