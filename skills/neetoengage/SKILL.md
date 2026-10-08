---
name: neetoengage
description: >
  Manage NeetoEngage from the command line.
  Use when the user asks about operations exposed by the NeetoEngage CLI.
---

## Prerequisites

Run `neetoengage doctor` to check authentication and connectivity.
If not authenticated, run `neetoengage login`.

## Authentication & multi-subdomain

Credentials for every logged-in subdomain are stored together in
`~/.config/neetoengage/auth.json`. A command that talks to the API picks which
subdomain to use by these rules:

- 0 subdomains authenticated → every credential-using command errors with
  "Not authenticated. Run 'neetoengage login' to authenticate.".
- 1 subdomain authenticated → that one is the implicit default; `--subdomain`
  may be omitted.
- 2+ subdomains authenticated → **`--subdomain <name>` is required** on every
  credential-using command, including `doctor`. The error lists every
  authenticated subdomain so the agent can offer a choice.

`login` / `logout` / `whoami` have dedicated behavior:

| Command | Behavior |
|---|---|
| `neetoengage login --subdomain <name>` | Adds or refreshes the entry for `<name>`. No flag → prompts for the subdomain. |
| `neetoengage logout --subdomain <name>` | Removes that one entry. |
| `neetoengage logout --all` | Removes every entry. |
| `neetoengage logout` (no flag) | Removes the only entry if exactly one is logged in; errors if multiple. |
| `neetoengage whoami` | Lists every logged-in account. Marks the entry `(default)` when exactly one. |
| `neetoengage whoami --subdomain <name>` | Shows just that one. |

## Global flags (persistent on every command)

| Flag | Purpose |
|---|---|
| `--subdomain <name>` | Select which logged-in subdomain the command targets. Required when multiple are logged in. |
| `--json` | Force JSON envelope output even on a TTY. |
| `--quiet` | Emit only the raw payload — no envelope, no breadcrumbs. For action commands (create/update), emits just the resource identifier; `delete` emits `success`. Designed for scripting. |
| `--toon` | Emit TOON (Token Optimized Output Notation). Preferred for feeding list/show output back to an LLM; ~30–60% fewer tokens than JSON. |

Precedence if multiple are set: `--toon` > `--quiet` > `--json` > pretty.

## Output modes & response envelope

**Pretty (default on a TTY)** — tables for arrays, key-value for objects,
breadcrumbs appended. Not intended for machine consumption.

**JSON envelope** (non-TTY, or `--json`):
```json
{
  "data": <resource body>,
  "breadcrumbs": [{ "label": "List", "command": "neetoengage <resource> list" }],
  "pagination": {
    "current_page_number": 1,
    "total_pages": 10,
    "total_records": 250
  }
}
```
`breadcrumbs` is omitted when empty. `pagination` is present only for list
commands.

**Quiet** (`--quiet`) — `data` contents only, no envelope. For action
commands `PrintQuiet` unwraps a single-key wrapper and prints the first of
`sid` / `id` / `name`. For `delete` it prints `success`.

**TOON** (`--toon`) — same data as JSON, re-encoded into TOON. Shape is
equivalent but whitespace/keys are compressed. Parse by re-reading keys as
you would JSON.

### Pagination

List commands accept `--page` (1-indexed) and `--page-size` (max 100).
The envelope's `pagination` field always exposes:
`current_page_number`, `total_pages`, `total_records`. Agents should loop
by incrementing `--page` until `current_page_number == total_pages`.

## Discovery

The full, always-accurate command tree (including any flags added after
this skill was built) is available as JSON:

```bash
neetoengage commands
```

Each catalog entry has `command`, `description`, optional `flags` (with
`name`, `type`, `default`, `description`, `required`), and `subcommands`.
Use this whenever a user asks about a flag or command not covered below.

## Diagnostics & IDE setup

| Command | Purpose |
|---|---|
| `doctor` | Auth check + API reachability + version. Uses `--subdomain` when multiple are logged in. |
| `version` | Print CLI version / commit / build date. |
| `commands` | Emit the full command/flag catalog as JSON. |
| `setup claude` | Install NeetoEngage plugin into Claude Code (`plugin.json`, hooks, this SKILL.md). |
| `setup cursor` / `windsurf` / `copilot` / `gemini` / `codex` | Write NeetoEngage rule files into the current project directory; re-run after an upgrade to refresh them. |

## Environment variable override

Set `NEETOENGAGE_BASE_URL` to point the CLI at a staging or local server:

```bash
export NEETOENGAGE_BASE_URL=http://acme.lvh.me:8980
neetoengage login --subdomain acme
```

## Error surface

Every command exits non-zero on failure and writes a single-line message to
stderr. Common errors the agent should expect:

- `Not authenticated. Run 'neetoengage login' to authenticate.` — empty credential store.
- `Multiple subdomains authenticated (acme, beta); specify --subdomain.` — pick one.
- `Not authenticated for "foo". Authenticated subdomains: acme, beta.` — bad `--subdomain`.
- `required flag(s) "xxx" not set` (from cobra) — missing required flag.
- API errors come through with the server's message body; inspect the
  JSON envelope (or the `--quiet` payload) for `error` / `errors` / `notice`
  keys and any suggestions the API returns.

## Product-specific commands

| Command | Purpose |
|---|---|
| `neetoengage settings show` | Show the workspace's product name and website URL. |
| `neetoengage settings update --product-name <name>` | Change the product name. It cannot be blank. |
| `neetoengage settings update --website-url <url>` | Change the website URL. It must start with `http://`, `https://` or `www.`. |
| `neetoengage feature-requests search --query <words>` | Search public feature requests. Shows each match's `admin_url`, `track` and `votes_count`. |
| `neetoengage feature-requests create --title <t> --description <d> --customer-email <e> --customer-name <n> --note <note>` | Create a feature request. The customer becomes its first voter and the note is a private team note. |
| `neetoengage feature-requests add-voter <id> --customer-email <e> --customer-name <n> --note <note>` | Add a customer as a voter on an existing request and attach a private note. |
| `neetoengage feature-requests voters <id>` | List who voted for a request. |
| `neetoengage feature-requests move <id> --track-id <track-id>` | Move a request to another track (board column). Add `--notify-all-voters` to email voters, `--message <html>` to replace the default "now live" email, `--changelog-id <id>` to attach a changelog. |
| `neetoengage feature-requests comment <id> --content <text>` | Post a public comment. Add `--notify-all-voters` to email it to every voter. |
| `neetoengage votes list --email <email>` | List the requests a customer voted for, with admin links. |
| `neetoengage tracks list` | List the tracks (board columns) and their IDs. |
| `neetoengage changelogs list --query <title>` | List changelogs that are not archived, to find one to attach. |

`<id>` is a feature request's ID or the slug at the end of its admin URL.

### Filing a customer's request

1. Run `feature-requests search` with the gist of the request. Show the user the
   matches with their `admin_url`, `track` and `votes_count`, and ask whether one
   is the same request.
2. On a match, run `feature-requests add-voter` on it. Otherwise run
   `feature-requests create`.
3. The title, the description and comments are public on the roadmap. Never put
   the customer's name, email, company or ticket details in them. Those go only
   in `--note`, which only the team can see.
4. Write the note in this format, one item per line: the product and ticket or
   conversation it came from; the customer's name and email; what the customer
   asked for; a link to the ticket.
5. Show the user the request's `admin_url`.

To tell voters a feature shipped, run `tracks list` to find the done track, then
`feature-requests move <id> --track-id <done-track-id> --notify-all-voters`, with
`--changelog-id` from `changelogs list` when there is a changelog.

`settings update` takes either flag or both, and needs at least one. It prints
the product name and website URL after the change. The public NeetoEngage
pages link the product name and logo in their header to the website URL.

`neetoengage commands` always prints the true command surface as JSON, so use
it to check this section against the binary.
