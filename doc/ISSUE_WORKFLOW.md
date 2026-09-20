# GitHub Issue Workflow

GitHub Issues are the delivery backlog and the concise record of implementation status. The PRD describes the intended final product; an Issue shows whether a bounded piece of that product is ready, underway, blocked, or complete in `main`.

Use an Issue for planned backlog work. An Issue is optional for an isolated ad hoc change, such as a small documentation or repository correction that does not need backlog tracking.

## Finding work and reading status

Use the repository's **Issues** tab and these filters:

- `is:issue is:open label:status:ready` — work that can be started.
- `is:issue is:open label:status:in-progress` — work with an active branch or pull request.
- `is:issue is:open label:status:blocked` — work awaiting an explicitly recorded dependency or decision.
- `is:issue is:closed` — work completed in `main`.

An open Issue without `status:ready` must not be assumed ready for implementation. Closed means that its final pull request was merged into `main`; it does not by itself claim that the change is deployed to production.

## Issue identifiers and type labels

Start each Issue title with one identifier. Use the same identifier in related pull requests and discussion.

| Identifier | Type label | Use |
| --- | --- | --- |
| `FR-X` | `type:feature` | A functional requirement from the PRD. `X` matches the PRD requirement number. |
| `BUG-X` | `type:bug` | A reproducible defect in existing behavior. |
| `NFR-X` | `type:nfr` | A non-functional requirement, such as accessibility, performance, security, or reliability. |
| `TASK-X` | `type:task` | A bounded documentation, design, maintenance, or repository task that is not an FR, bug, or NFR. |
| `DEBT-X` | `type:debt` | An accepted temporary shortcut that must be replaced later. |
| `DEC-X` | `type:decision` | A decision required before dependent work can proceed. |

Use the next unused number within `BUG`, `NFR`, `TASK`, `DEBT`, and `DEC`. Do not introduce priorities or area labels until the backlog needs them.

## Status labels

Every open Issue has exactly one status label:

- `status:ready` — scope and acceptance conditions are sufficient to implement; there is no unresolved blocker.
- `status:in-progress` — implementation has started; the Issue links to its branch or pull request.
- `status:blocked` — the Issue records what is blocking it and where the required decision or dependency is tracked.

Do not use a `done` label. Closing the Issue is the status for delivered work.

## Issue contents

Every Issue states its outcome, scope, and completion conditions.

- An `FR-X` Issue links to its PRD section and lists every applicable acceptance criterion as a checklist. It closes only after all listed criteria are complete.
- A `BUG-X` Issue includes reproduction steps, expected result, actual result, and relevant environment details.
- An `NFR-X` Issue states the measurable requirement and how it will be verified.
- A `TASK-X` Issue states the bounded result and why it is needed.
- A `DEBT-X` Issue states the temporary shortcut, why it is acceptable now, the intended end state, and the condition that should trigger returning to it.
- A `DEC-X` Issue states the decision needed, the options or constraints known so far, and which work it blocks. Close it when the decision is recorded; create or unblock the resulting implementation Issues separately.

## Pull requests and completion

One Issue may have more than one pull request. When a pull request is linked to an Issue, each active pull request uses `Refs #<issue-number>` in its description. The final pull request uses `Closes #<issue-number>` only when it completes the entire Issue.

Keep the Issue `status:in-progress` until every linked pull request is merged and every completion condition is satisfied. A merged partial pull request does not close the Issue.

Before opening a final pull request, confirm that the Issue scope still matches the PRD and any linked decisions. If the work uncovers a new independent concern, create a separate Issue rather than silently expanding the current one.
