# General Workflow

## GitHub Issues

Use [doc/ISSUE_WORKFLOW.md](doc/ISSUE_WORKFLOW.md) for the shared process for finding, implementing, and closing GitHub Issues. Treat it as the source of truth for Issue identifiers, labels, status, and links between Issues and pull requests.

For every task:

Inspect the relevant parts of the repository before proposing changes.
Read any nested AGENTS.md files that apply to directories you expect to modify.
Treat the repository and applicable AGENTS.md files as the source of truth. Do not assume the task description contains all relevant implementation details.


## Before Implementation

For any non-trivial change:

Restate the goal briefly to confirm your understanding.
Inspect the existing implementation and identify the files and components involved.
Propose a concise implementation plan.
Explain important architectural decisions, tradeoffs, or assumptions.
List the files you expect to modify.
Wait for explicit approval before modifying files.
Discussing features take a look into doc/ directory and check requirements
On any requirements conflict discuss it explicitly before implementation

Do not start implementation while discussing or refining the plan.

For trivial, explicitly requested edits where no architectural or design decision is involved, a plan is not required unless requested.

## During Implementation

After the plan is approved:

Implement only the approved scope.
Follow existing project patterns and applicable AGENTS.md instructions.
Prefer the smallest coherent change that satisfies the requirements.
Do not introduce unrelated refactoring, dependencies, abstractions, or infrastructure changes.
If implementation reveals a significant issue that requires changing the approved plan, stop and explain it before proceeding.
Add or update tests appropriate to the changed behavior.

Do not silently expand the scope of the task.

## Validation

After implementation:

Run the most relevant automated tests and checks available for the changed area.
Fix failures caused by the change.
Do not modify unrelated code merely to make unrelated existing failures disappear.
When practical, validate behavior at the lowest appropriate testing level before running broader tests.

Use existing Makefile targets when they provide the intended project workflow.

## Completion

When finished, provide a concise summary containing:

what changed;
important implementation decisions;
tests/checks that were run and their results;
any remaining limitations, risks, or follow-up work.

Do not claim that something was tested or verified if it was not actually run.

## Architecture and Scope

Do not make architectural decisions implicitly.

If a task appears to require changing an established architectural decision, dependency strategy, public interface, infrastructure pattern, or project convention:

identify the conflict;
explain why the existing approach may need to change;
present the proposed change and relevant tradeoffs;
wait for approval before proceeding.

Prefer adapting to the existing architecture over redesigning it unless redesign is part of the task.

## Uncertainty

Do not guess when repository inspection can answer the question.

If requirements remain materially ambiguous after inspecting the repository, ask for clarification before implementation.

For minor implementation details that do not affect behavior or architecture, use reasonable existing project conventions instead of asking unnecessary questions.

## Documentation

Documentation stored in the doc/ directory
Change documentation and even AGENTS.md files if the logic was changed during the implementation.
Keep comments and readme files up to date.

## Git and Pull Request Workflow

For implementation tasks, use a dedicated Git branch and a dedicated Git worktree unless explicitly instructed otherwise. This lets multiple tasks proceed in parallel without interfering with each other's working tree state, even when the user has not mentioned that other work is in progress.

For trivial, explicitly requested edits where no architectural or design decision is involved, working directly in the current checkout without a worktree is fine.

After the implementation plan is approved and before modifying files:

Confirm that the working tree does not contain unrelated uncommitted changes.
Use a short descriptive branch name appropriate to the task.

Examples:

feature/backend-config
fix/frontend-api-error
test/login-e2e

Check `git worktree list` first to avoid creating a duplicate worktree for a branch that already has one.
Create the branch and its worktree together in a single step, from the current intended base branch, so the main checkout stays on the base branch instead of switching to the new branch:

    git worktree add -b <branch-name> .agents/worktrees/<branch-name> <base-branch>

Example:

    git worktree add -b feature/backend-config .agents/worktrees/feature/backend-config main

Do all implementation work for the task inside that worktree directory, not in the main checkout.

Do not discard, overwrite, stash, or commit unrelated user changes without explicit approval.

After implementation:

Run the relevant tests and checks.
Review the resulting diff for unintended changes.
Commit only changes related to the task.

When committing changes you implemented, add yourself as a `Co-Authored-By`
trailer using your actual agent/model identity and your provider's official
attribution email, if available.

Never invent an attribution email or impersonate another agent.
Keep the user's Git identity as the primary author/committer and never modify
`git config user.name` or `git config user.email`.

Use a concise commit message describing the change.
Push the branch to the remote repository.
Open a pull request against the appropriate base branch.

The pull request should include:

a concise summary of the change;
important implementation or architectural decisions;
tests and checks that were run;
any known limitations or follow-up work.

Do not merge the pull request unless explicitly instructed to do so.

If branch creation, pushing, or pull request creation is not possible because of missing authentication, permissions, tools, or repository state, stop at the appropriate point and clearly report what remains to be done.

After the pull request is merged, or if the task is abandoned, remove the worktree instead of leaving it behind:

    git worktree remove .agents/worktrees/<branch-name>
