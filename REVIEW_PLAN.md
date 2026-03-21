# Jiru Codebase Review — Implementation Plan

Based on the comprehensive codebase review (March 2026), this document tracks all recommendations and their implementation steps, ordered by priority.

---

## High Priority

### 1. Shell Completions (bash/zsh/fish)

Table stakes for CLI tools. Cobra has built-in support.

**Steps:**
1. Add `completion` subcommand to `main.go` via `cobra.GenBashCompletionV2`, `cobra.GenZshCompletion`, `cobra.GenFishCompletion`
2. Register completions for flags (`--profile` should complete from `profiles.yml` entries)
3. Register completions for positional args (issue keys where applicable)
4. Add `make completion` target to Makefile for generating static completion scripts
5. Document installation in README (brew, manual sourcing)
6. Add completion script to Homebrew formula via release workflow

---

### 2. Help Overlay (`?` keybinding)

Users must memorize 30+ keybindings with no in-app reference beyond the footer.

**Steps:**
1. Create `internal/ui/helpview/` package with a `Model` using `bubbles/viewport`
2. Render grouped shortcut table: Navigation, Issue Operations, Search & Filters, Confluence, General
3. Show context-sensitive content (highlight shortcuts relevant to the current view)
4. Wire `?` key in `app.go` — push `viewHelp` overlay, dismiss with `esc`/`q`/`?`
5. Suppress `?` when text input is active (`inputActive()` guard)
6. Update footer to show `? help` hint
7. Add tests for help view model

---

### 3. Consolidate Navigation Origin Variables

Five separate origin trackers (`searchOrigin`, `filterOrigin`, `tabOrigin`, `profileOrigin`, `issuePickOrigin`) are error-prone.

**Steps:**
1. Define a `navContext` struct holding: source view, cursor position, metadata (filter name, search query, etc.)
2. Replace the 5 origin variables with a `navStack []navContext` on `App`
3. Refactor `navigateBack()` to pop from the stack instead of checking each origin
4. Ensure all view transitions push onto the stack instead of setting individual origin fields
5. Add tests for navigation stack push/pop across all view transition paths
6. Verify no regressions in back-navigation from: search results, filter results, issue detail, Confluence pages, profile switch

---

### 4. Jira Data Center / On-Premise Support

Many enterprise users run Jira Server or Data Center, not Cloud.

**Steps:**
1. Add `JIRA_SERVER_TYPE` config field (`cloud` | `server` | `datacenter`, default `cloud`)
2. Add server type to setup wizard as a new step after domain
3. Branch API paths: Server uses `/rest/api/2/search` (not removed on Server), Cloud uses v3
4. Handle authentication differences: Server supports session auth + PAT, Cloud uses API tokens
5. Adjust pagination: Server uses `startAt` offset for all endpoints (no cursor pagination)
6. Skip Confluence v2 API on Server — fall back to v1 (`/rest/api/content`)
7. Test against a local Jira Server Docker instance (add to CI as optional)
8. Update README with Server/DC configuration notes
9. Add `validate.ServerURL()` that accepts self-hosted domains (not just `*.atlassian.net`)

---

## Medium Priority

### 5. CSV/JSON Output for CLI Subcommands

Enables scripting and integration with other tools.

**Steps:**
1. Add `--output` / `-o` flag to CLI subcommands (`get`, `search`, `list`): values `table` (default), `json`, `csv`
2. Create `internal/cli/output.go` with `FormatJSON()`, `FormatCSV()`, `FormatTable()` renderers
3. JSON output: marshal `jira.Issue` / `jira.Board` directly
4. CSV output: header row + data rows, configurable columns via `--fields` flag
5. Table output: current default behavior
6. Add `--no-header` flag for CSV (useful for piping)
7. Update `search`, `list`, `get`, `boards` subcommands
8. Add tests for each output format

---

### 6. Improve Status Message Persistence

Status messages clear on any keypress — users miss success feedback.

**Steps:**
1. Add `statusMsgTime time.Time` to `App` struct, set on status message creation
2. Change clearing logic: only clear when (a) a meaningful action is taken (view change, selection) OR (b) 5 seconds have elapsed
3. Add a `tea.Tick` command for auto-dismiss after timeout
4. Keep immediate clear on `esc` for explicit dismissal
5. Add tests for status message lifecycle

---

### 7. Non-Modal Toasts for Non-Critical Errors

Modal error dialogs block all interaction even for minor issues (e.g., failed user search suggestion).

**Steps:**
1. Add error severity enum: `errCritical` (modal), `errWarning` (toast)
2. Create toast rendering in `app.View()` — styled message in top-right corner with auto-dismiss
3. Route non-critical errors (user search failure, metadata fetch partial failure) to toast
4. Keep modal for: auth failure, network error, issue not found, permission denied
5. Add `tea.Tick` for toast auto-dismiss (3 seconds)
6. Add tests for error routing by severity

---

### 8. Issue Cloning

Missing feature that `jira-cli` supports.

**Steps:**
1. Add `CloneIssue()` to `internal/client/` — `POST /rest/api/2/issue` with fields copied from source
2. Add `y` keybinding in issue view for "clone" (mnemonic: yank/copy)
3. Show confirmation overlay with editable summary (pre-filled with "Clone of PROJ-123: <summary>")
4. Clone fields: project, type, priority, labels, description, parent. Skip: assignee, status, comments
5. Navigate to cloned issue on success
6. Add to footer in issue view
7. Add tests for clone API call and UI flow

---

### 9. Issue Watch/Unwatch

Missing feature that `jira-cli` supports.

**Steps:**
1. Add `WatchIssue()` / `UnwatchIssue()` to `internal/client/` — `POST/DELETE /rest/api/2/issue/{key}/watchers`
2. Add `IsWatching()` check to `GetIssue()` response parsing (check `watches.isWatching` field)
3. Display watch status in issue view metadata section
4. Add `w` keybinding in issue view to toggle watch
5. Show status message on toggle ("Watching PROJ-123" / "Unwatched PROJ-123")
6. Add tests

---

## Low Priority / Polish

### 10. Retry Button in Error Dialogs

After connection timeouts, users must re-navigate to retry.

**Steps:**
1. Store the last failed `tea.Cmd` on `App` as `retryCmd`
2. Add `r` keybinding in error dialog to re-execute the failed command
3. Show "r retry | esc dismiss" in error dialog footer
4. Clear `retryCmd` on successful retry or dismiss
5. Add tests for retry flow

---

### 11. Mouse Scroll Support

Bubble Tea supports mouse events; enabling scroll improves accessibility.

**Steps:**
1. Enable mouse support in Bubble Tea options: `tea.WithMouseCellMotion()`
2. Handle `tea.MouseMsg` in list views, board view, issue detail, search results
3. Map scroll wheel to cursor movement (scroll up = `k`, scroll down = `j`)
4. Map click on list items to selection
5. Test mouse support doesn't interfere with terminal copy/paste (may need `tea.WithMouseAllMotion()` investigation)
6. Add toggle: `--no-mouse` flag or config option for users who prefer pure keyboard

---

### 12. More Specific Loading Messages

"Connecting to Jira..." is too generic.

**Steps:**
1. Add `loadingMsg string` field to `App`
2. Set contextual messages: "Verifying credentials...", "Fetching boards...", "Loading sprint issues...", "Searching..."
3. Update `viewLoading` rendering to display `a.loadingMsg` instead of static text
4. Set appropriate message before each async command dispatch
5. Keep spinner animation

---

### 13. Client-Side Rate Limiting

No rate limiting — relies entirely on Jira server-side limits.

**Steps:**
1. Add `rate.Limiter` (from `golang.org/x/time/rate`) to `api.Client`
2. Default: 10 requests/second (conservative for Jira Cloud)
3. Make configurable via `JIRA_RATE_LIMIT` env var
4. Add `Wait()` call before each HTTP request in `Do()`
5. Add tests with mock time

---

### 14. mTLS Authentication

For on-premise deployments using client certificates.

**Steps:**
1. Add `JIRA_CLIENT_CERT` and `JIRA_CLIENT_KEY` config fields
2. Add `mtls` as a valid auth type in `validate.AuthType()`
3. Load certificate pair in `api.NewClient()` and configure `tls.Config`
4. Add cert/key path steps to setup wizard (conditional on auth type = mtls)
5. Add tests with `httptest.NewTLSServer()`

---

## Test Coverage Gaps

### 15. Fill Test Gaps

**Steps:**
1. `internal/theme/` — Add tests for `StatusStyle()`, `StatusCategory()`, `StatusSubPriority()`, `UserStyle()` with various inputs
2. `internal/cli/branch.go` — Add test file `branch_test.go` with mock HTTP backend
3. `internal/jql/` — Add tests for `Match()` completion logic and `RenderPopup()` output
4. `internal/api/` — Add tests for timeout handling and connection failure scenarios
5. `internal/filters/` — Add test for concurrent access and corrupted YAML recovery
6. Add E2E workflow test: setup → search → create issue → transition → comment (using mock HTTP server)

---

## UX Polish

### 16. Unsaved Changes Warning

Create/edit wizards and comment input lose all data on `esc` with no confirmation.

**Steps:**
1. Track `dirty` flag on createview, editview, commentview models (set true on any input change)
2. On `esc` when `dirty`, show confirmation: "Discard changes? (y/n)"
3. Only dismiss on explicit `y` or second `esc`
4. Add tests for dirty tracking and confirmation flow

---

### 17. Input Validation Visual Feedback

Validation errors appear as text but don't highlight the problematic field.

**Steps:**
1. Add `errField` index to setupview and createview models
2. When validation fails, set `errField` and render that field's border in `ColourError` (red)
3. Add character count display below text inputs (e.g., "45/100")
4. Add asterisk to required field labels
5. Add tests for visual feedback state

---

### 18. Confluence Discoverability

Confluence integration is hidden behind an undocumented `Tab` key.

**Steps:**
1. Add "Tab wiki" to footer hint on home/sprint/board views
2. On first `Tab` press (check via config flag `confluenceHintShown`), show a brief toast: "Switched to Confluence wiki. Press Tab to return to Jira."
3. Persist hint-shown flag to avoid repeated prompts
4. Add Confluence section to `?` help overlay

---

### 19. Breadcrumb Trail for Nested Navigation

When drilling into parent → parent → parent issues, there's no visual trail.

**Steps:**
1. Add `navBreadcrumb []string` to `App` (e.g., `["PROJ-456", "PROJ-100", "PROJ-50"]`)
2. Render breadcrumb bar below the title in issue view: `PROJ-456 > PROJ-100 > PROJ-50`
3. Update on push/pop of issue stack
4. Make breadcrumb items clickable (or numbered for keyboard access)
5. Add tests for breadcrumb state management

---

### 20. Home Shortcut

No quick way to return to the home view from deeply nested views.

**Steps:**
1. Add `H` (shift+h) keybinding for "go home" — return to `viewHome`
2. Clear navigation stack on home
3. Add to global keys (suppressed when text input active)
4. Show in footer when not on home view
5. Add tests for home navigation from various depths
