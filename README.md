# NeetoEngage CLI

NeetoEngage CLI

<!-- neeto-cli-commons:installation:start -->
## Installation

### macOS / Linux

**Homebrew (recommended on macOS):**

```bash
brew install neetozone/tap/neetoengage
```

**Shell script:**

```bash
curl -fsSL https://neeto-downloads.s3.amazonaws.com/cli/NeetoEngage/latest/install.sh | sh
```

This verifies the download's SHA-256 checksum against the published `SHA256SUMS`,
then installs to `/usr/local/bin` (may prompt for sudo). Set `NEETOENGAGE_INSTALL_DIR`
to a directory you own to install without sudo.

### Windows

**PowerShell:**

```powershell
irm https://neeto-downloads.s3.amazonaws.com/cli/NeetoEngage/latest/install.ps1 | iex
```

**Command Prompt (CMD):**

```cmd
curl -fsSL https://neeto-downloads.s3.amazonaws.com/cli/NeetoEngage/latest/install.cmd -o install.cmd && install.cmd
```

Both verify the download's SHA-256 checksum before installing to
`%LOCALAPPDATA%\Programs\neetoengage` and adding it to your user PATH. Set
`NEETOENGAGE_INSTALL_DIR` to install somewhere else.
<!-- neeto-cli-commons:installation:end -->

<!-- neeto-cli-commons:verify-installation:start -->
### Verify installation

```bash
neetoengage --help
```
<!-- neeto-cli-commons:verify-installation:end -->

<!-- neeto-cli-commons:ai-coding-assistants:start -->
## AI coding assistants

```bash
neetoengage setup claude      # Register plugin with Claude Code
neetoengage setup cursor      # Write .cursor/rules/neetoengage.mdc
neetoengage setup windsurf    # Write .windsurf/rules/neetoengage.md
neetoengage setup copilot     # Add a NeetoEngage section to .github/copilot-instructions.md
neetoengage setup gemini      # Add a NeetoEngage section to GEMINI.md
neetoengage setup codex       # Add a NeetoEngage section to AGENTS.md
```

Every command except `setup claude` writes into the current project directory, so
run these commands from the root of the project the assistant works in. Re-run
them after every upgrade: `setup cursor` and `setup windsurf` overwrite their rule
file, while `setup copilot`, `setup gemini` and `setup codex` keep the existing
content of their file and replace only the NeetoEngage section instead of adding a
duplicate.
<!-- neeto-cli-commons:ai-coding-assistants:end -->

<!-- neeto-cli-commons:prerequisites:start -->
## Prerequisites (development)

- [Go](https://go.dev/dl/) 1.26.1+
- Access to a NeetoEngage organization
<!-- neeto-cli-commons:prerequisites:end -->

## Development

### Setup

```bash
git clone https://github.com/neetozone/neeto-engage-cli.git
cd neeto-engage-cli
bin/setup
```

This installs Go dependencies, golangci-lint, configures git hooks, and builds the binary.

### Build and run

```bash
make build
./neetoengage --help
```

<!-- neeto-cli-commons:make-targets:start -->
### Make targets

```bash
make build          # Builds ./neetoengage
make test           # Run tests
make lint           # golangci-lint
make fmt            # gofmt -w
make vet            # go vet
make check          # fmt + vet + test
make install        # Installs to /usr/local/bin
make clean          # Remove built binary
```
<!-- neeto-cli-commons:make-targets:end -->

### Adding a command

`internal/commands/example.go` is a working two-command resource — a `list` and a
`show`. Copy it, rename the resource, point it at the right API path, and delete the
example. Registration happens in `init()`, so a new file needs no wiring anywhere else.

### Pointing to a local or staging server

By default the CLI targets `https://{subdomain}.neetoengage.com`. Set `NEETOENGAGE_BASE_URL`
to point somewhere else:

```bash
export NEETOENGAGE_BASE_URL=http://acme.lvh.me:3000
neetoengage login --subdomain acme
neetoengage doctor
```

<!-- neeto-cli-commons:global-flags:start -->
## Global flags

Every command accepts:

| Flag | Description |
|---|---|
| `--subdomain <name>` | Which logged-in subdomain to use (required when multiple are logged in). |
| `--json` | Force JSON envelope output. |
| `--quiet` | Emit raw data only. Action commands print just the identifier; `delete` prints `success`. |
| `--toon` | TOON (Token-Optimized Output Notation) — compact format for LLMs. |
| `--verbose` | Expand every field of a record instead of a table. |
<!-- neeto-cli-commons:global-flags:end -->

<!-- neeto-cli-commons:release:start -->
## Release

Releases are cut by the CI pipeline defined in `.neetoci/release.yml`.

Merging a PR with a `major`, `minor`, or `patch` label to `main` triggers the shared release script from `neeto-cli-commons`. The script:

* Bumps and tags `VERSION`
* Runs GoReleaser
* Uploads artifacts to `s3://neeto-downloads/cli/NeetoEngage/`
* Updates the Homebrew tap (`neetozone/tap`)
* Pushes the version bump commit to `main`
<!-- neeto-cli-commons:release:end -->

## Support

Questions or problems: support@neetoengage.com
